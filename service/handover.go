package service

import (
	"fmt"
	"github.com/hiro942/elden-client/global"
	"github.com/hiro942/elden-client/model"
	"github.com/hiro942/elden-client/model/request"
	"github.com/hiro942/elden-client/model/response"
	"github.com/hiro942/elden-client/utils"
	"github.com/hiro942/elden-client/utils/gxios"
	"github.com/pkg/errors"
	"github.com/tjfoc/gmsm/x509"
	"log"
	"time"
)

func HandoverAccess(satelliteId string) error {
	log.Println("Go Handover Access")

	// 是否存在该卫星公钥
	if _, ok := global.SatellitePubKeys[satelliteId]; !ok {
		// 查询卫星公钥
		satellitePublicKeyHex, _ := gxios.QuerySatellitePublicKey(satelliteId)
		satellitePublicKey, err := x509.ReadPublicKeyFromHex(satellitePublicKeyHex)
		if err != nil {
			log.Panicln(fmt.Printf("failed to resolve satellite's public key: %+v", err))
		}
		global.SatellitePubKeys[satelliteId] = satellitePublicKey
	}

	var isExpired = false
	var sessionKeyInfo any

	log.Println("Handover Access: Read local file of session record and check whether the session key has expired.")

	// 读会话记录文件，判断有无密钥记录
	hasRecord := false
	records := utils.ReadSessionRecords()
	for _, record := range records {
		if record.SatelliteId == satelliteId {
			hasRecord = true
			break
		}
	}

	// 没有则创建，并写入本地
	// 有则读取记录
	if !hasRecord {

		newSessionKey := utils.GenerateSm4Key()
		expDate := time.Now().Unix() + global.DefaultSessionKeyAge
		sessionKeyInfo = request.EncryptedSessionKeyWithExpDate{
			EncryptedSessionKey: utils.Sm2Encrypt(global.SatellitePubKeys[satelliteId], newSessionKey),
			ExpirationDate:      expDate,
		}

		utils.WriteNewSessionRecord(model.SessionRecord{
			SatelliteId:    satelliteId,
			SessionKey:     string(newSessionKey),
			ExpirationDate: expDate,
		})

	} else {

		for _, record := range records {
			if record.SatelliteId == satelliteId {
				hasRecord = true
				// 判断密钥是否过期。若密钥已经过期，则生成一个新密钥，并更新会话记录文件。
				if record.ExpirationDate < time.Now().Unix()+60 {
					isExpired = true
					newSessionKey := utils.GenerateSm4Key()
					log.Println("Generated Session Key:", string(newSessionKey))
					record.SessionKey = string(newSessionKey) // 更新会话密钥
					recordsBytes := utils.JsonMarshal(records)
					utils.WriteFile(global.SessionRecordsFilePath, recordsBytes) // 写回
					sessionKeyInfo = request.EncryptedSessionKeyWithExpDate{
						EncryptedSessionKey: utils.Sm2Encrypt(global.SatellitePubKeys[satelliteId], newSessionKey),
						ExpirationDate:      time.Now().Unix() + global.DefaultSessionKeyAge,
					}
				} else {
					sessionKeyInfo = utils.Sm3Hash([]byte(record.SessionKey))
				}
				break
			}
		}
	}

	if hasRecord && !isExpired {

		log.Println("Handover Access (hashed): The session key has not expired so that user send the hashed form to the satellite.")

		NARWithSig := utils.GetMessageWithSig[request.NARHashed](request.NARHashed{
			HashedIMSI:       global.MyHashedIMSI,
			MacAddr:          global.MyMacAddr,
			SatelliteId:      satelliteId,
			HashedSessionKey: sessionKeyInfo.(string),
			TimeStamp:        time.Now().Unix(),
		})

		url := fmt.Sprintf("http://%s/auth/handover?type=hashed", global.SatelliteSocket[satelliteId])
		resBytes := gxios.POST(url, NARWithSig)

		res := utils.JsonUnmarshal[response.Response[any]](resBytes)

		if res.Code != 0 {
			return errors.Errorf("message: %s, decription: %s",
				res.Message, res.Description)
		}

		log.Println("Handover Access (hashed): Response OK.")

	} else {

		log.Println("Handover Access (encrypted): The session key has expired so that user generate a new key and send the encrypted form to the satellite.")

		NARWithSig := utils.GetMessageWithSig[request.NAREncrypted](request.NAREncrypted{
			HashedIMSI:                     global.MyHashedIMSI,
			MacAddr:                        global.MyMacAddr,
			SatelliteId:                    satelliteId,
			EncryptedSessionKeyWithExpDate: sessionKeyInfo.(request.EncryptedSessionKeyWithExpDate),
			TimeStamp:                      time.Now().Unix(),
		})

		url := fmt.Sprintf("http://%s/auth/handover?type=encrypted", global.SatelliteSocket[satelliteId])
		resBytes := gxios.POST(url, NARWithSig)

		res := utils.JsonUnmarshal[response.Response[any]](resBytes)

		if res.Code != 0 {
			return errors.Errorf("message: %s, decription: %s",
				res.Message, res.Description)
		}

	}

	log.Println("Handover Access Success! Connecting To:", satelliteId)

	return nil
}
