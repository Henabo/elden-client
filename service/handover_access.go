package service

import (
	"github.com/hiro942/elden-client/global"
	"github.com/hiro942/elden-client/global/url"
	"github.com/hiro942/elden-client/model"
	"github.com/hiro942/elden-client/model/request"
	"github.com/hiro942/elden-client/model/response"
	"github.com/hiro942/elden-client/utils"
	"github.com/hiro942/elden-client/utils/gxios"
	"github.com/pkg/errors"
	"log"
	"time"
)

func HandoverAccess(satelliteId string) error {
	log.Println("Go Handover Access")

	// 是否存在该卫星公钥
	if _, ok := global.SatellitePubKeys[satelliteId]; !ok {
		// 查询卫星公钥
		satellitePublicKeyHex, _ := gxios.QuerySatellitePublicKey(satelliteId)
		satellitePublicKey := utils.ReadPublicKeyFromHex(satellitePublicKeyHex)

		// 保存卫星公钥
		global.SatellitePubKeys[satelliteId] = satellitePublicKey
	}

	var sessionKey []byte
	var expDate int64
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
		sessionKey = utils.GenerateSm4Key()
		expDate = time.Now().Unix() + global.DefaultSessionKeyAge
		sessionKeyInfo = request.EncryptedSessionKeyWithExpDate{
			EncryptedSessionKey: utils.Sm2Encrypt(global.SatellitePubKeys[satelliteId], sessionKey),
			ExpirationDate:      expDate,
		}

		utils.WriteNewSessionRecord(model.SessionRecord{
			SatelliteId:    satelliteId,
			SessionKey:     string(sessionKey),
			ExpirationDate: expDate,
		})

	} else {

		for _, record := range records {
			if record.SatelliteId == satelliteId {
				hasRecord = true
				// 判断密钥是否过期。若密钥已经过期，则生成一个新密钥，并更新会话记录文件。
				if record.ExpirationDate < time.Now().Unix()+60 {
					isExpired = true
					sessionKey = utils.GenerateSm4Key()
					expDate = time.Now().Unix() + global.DefaultSessionKeyAge
					log.Println("Generated Session Key:", string(sessionKey))
					record.SessionKey = string(sessionKey)                                     // 更新会话密钥
					utils.WriteFile(global.SessionRecordsFilePath, utils.JsonMarshal(records)) // 会话记录写回
					sessionKeyInfo = request.EncryptedSessionKeyWithExpDate{
						EncryptedSessionKey: utils.Sm2Encrypt(global.SatellitePubKeys[satelliteId], sessionKey),
						ExpirationDate:      expDate,
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

		resBytes := gxios.POST(
			url.NormalAccess(satelliteId, true, true),
			utils.GetMessageWithSig[request.NARHashed](request.NARHashed{
				HashedIMSI:       global.MyHashedIMSI,
				MacAddr:          global.MyMacAddr,
				SatelliteId:      satelliteId,
				HashedSessionKey: sessionKeyInfo.(string),
				TimeStamp:        time.Now().Unix(),
			}))
		if res := utils.JsonUnmarshal[response.Response[any]](resBytes); res.Code != 0 {
			return errors.Errorf("message: %s, decription: %s",
				res.Message, res.Description)
		}

		log.Println("Handover Access (hashed): Response OK.")

	} else {

		log.Println("Handover Access (encrypted): The session key has expired so that user generate a new key and send the encrypted form to the satellite.")

		resBytes := gxios.POST(
			url.NormalAccess(satelliteId, true, false),
			utils.GetMessageWithSig[request.NAREncrypted](request.NAREncrypted{
				HashedIMSI:                     global.MyHashedIMSI,
				MacAddr:                        global.MyMacAddr,
				SatelliteId:                    satelliteId,
				EncryptedSessionKeyWithExpDate: sessionKeyInfo.(request.EncryptedSessionKeyWithExpDate),
				TimeStamp:                      time.Now().Unix(),
			}))
		if res := utils.JsonUnmarshal[response.Response[any]](resBytes); res.Code != 0 {
			return errors.Errorf("message: %s, decription: %s",
				res.Message, res.Description)
		}
	}

	log.Println("Handover Access Success! Connecting To:", satelliteId)

	// 记录当前会话
	global.CurrentSession = model.Session{
		AuthStatus:  true,
		SatelliteId: satelliteId,
		Socket:      global.SatelliteSockets[satelliteId],
		AccessType:  "handover",
		SessionKey:  sessionKey,
	}
	return nil
}
