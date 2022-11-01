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

func NormalAccess(satelliteId string) error {
	log.Println("[Normal-Access] Go Normal Access")

	// 是否存在该卫星公钥
	if _, ok := global.SatellitePubKeys[satelliteId]; !ok {
		// 查询卫星公钥
		satellitePublicKeyHex, _ := gxios.QuerySatellitePublicKey(satelliteId)
		satellitePublicKey := utils.ReadPublicKeyFromHex(satellitePublicKeyHex)
		global.SatellitePubKeys[satelliteId] = satellitePublicKey
	}

	log.Println("[Normal-Access] Read local file of session record and check whether the session key has expired.")
	// 读会话记录文件，判断会话密钥是否过期
	records := utils.ReadSessionRecords()
	var sessionKey []byte
	var expDate int64
	var isExpired = false
	var sessionKeyInfo any
	for _, record := range records {
		if record.SatelliteId == satelliteId {
			// 判断密钥是否过期。若密钥已经过期，则生成一个新密钥，并更新会话记录文件。
			if record.ExpirationDate < time.Now().Unix()+60 {
				isExpired = true
				sessionKey = utils.GenerateSm4Key()
				expDate = time.Now().Unix() + global.DefaultSessionKeyAge
				log.Println("Generated Session Key:", string(sessionKey))
				record.SessionKey = string(sessionKey) // 更新会话密钥
				recordsBytes := utils.JsonMarshal(records)
				utils.WriteFile(global.SessionRecordsFilePath, recordsBytes) // 写回
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

	if isExpired {
		log.Println("[Normal-Access](encrypted): The session key has expired so that user generate a new key and send the encrypted form to the satellite.")

		resBytes := gxios.POST(
			url.NormalAccess(satelliteId, false, false),
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

		log.Println("[Normal-Access] (encrypted): Response OK.")

	} else {

		log.Println("[Normal-Access](hashed): The session key has not expired so that user send the hashed form to the satellite.")

		resBytes := gxios.POST(
			url.NormalAccess(satelliteId, false, true),
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

		log.Println("[Normal-Access](hashed): Response OK.")
	}

	log.Println("[Normal-Access] Normal Access Success!")

	// 记录当前会话
	global.CurrentSession = model.Session{
		AuthStatus:  true,
		SatelliteId: satelliteId,
		Socket:      global.SatelliteSockets[satelliteId],
		AccessType:  "first",
		SessionKey:  sessionKey,
	}
	return nil
}
