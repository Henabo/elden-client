package service

import (
	"github.com/hiro942/elden-client/global"
	"github.com/hiro942/elden-client/model/request"
	"github.com/hiro942/elden-client/model/response"
	"github.com/hiro942/elden-client/utils"
	"github.com/hiro942/elden-client/utils/gxios"
	"github.com/pkg/errors"
	"log"
	"time"
)

func NormalAccess(satelliteId string) error {
	// 读会话记录文件
	records := utils.ReadSessionRecords()

	// 判断会话密钥是否过期
	var isExpired = false
	var sessionKeyInfo any
	for _, record := range records {
		if record.SatelliteId == satelliteId {
			// 判断密钥是否过期。若密钥已经过期，则生成一个新密钥，并更新会话记录文件。
			if record.ExpirationDate < time.Now().Unix()+60 {
				isExpired = true
				newSessionKey := utils.GenerateSm4Key()
				record.SessionKey = string(newSessionKey) // 更新会话密钥
				recordsBytes := utils.JsonMarshal(records)
				utils.WriteFile(global.SessionRecordsFilePath, recordsBytes) // 写回
				sessionKeyInfo = request.SessionKeyKeyWithExpDate{
					EncryptedSessionKey: string(utils.Sm2Encrypt(global.SatellitePubKeys[satelliteId], newSessionKey)),
					ExpirationDate:      time.Now().Unix() + global.DefaultSessionKeyAge,
				}
			} else {
				sessionKeyInfo = request.HashedSessionKey(record.SessionKey)
			}
			break
		}
	}

	if isExpired {
		NARWithSig := utils.GetMessageWithSig[request.NAR[request.SessionKeyKeyWithExpDate]](request.NAR[request.SessionKeyKeyWithExpDate]{
			HashedIMSI:     global.MyHashedIMSI,
			MacAddr:        global.MyMacAddr,
			SatelliteId:    satelliteId,
			SessionKeyInfo: sessionKeyInfo.(request.SessionKeyKeyWithExpDate),
			TimeStamp:      time.Now().Unix(),
		})

		url := gxios.GetBaseURL(satelliteId) + "/normal?type=encrypted"
		resBytes := gxios.POST(url, NARWithSig)

		res := utils.JsonUnmarshal[response.Response[any]](resBytes)

		if res.Code != 0 {
			return errors.Errorf("message: %s, decription: %s",
				res.Message, res.Description)
		}

		log.Println(res.Message)
	} else {
		NARWithSig := utils.GetMessageWithSig[request.NAR[request.HashedSessionKey]](request.NAR[request.HashedSessionKey]{
			HashedIMSI:     global.MyHashedIMSI,
			MacAddr:        global.MyMacAddr,
			SatelliteId:    satelliteId,
			SessionKeyInfo: sessionKeyInfo.(request.HashedSessionKey),
		})

		url := gxios.GetBaseURL(satelliteId) + "/normal?type=hashed"
		resBytes := gxios.POST(url, NARWithSig)

		res := utils.JsonUnmarshal[response.Response[any]](resBytes)

		if res.Code != 0 {
			return errors.Errorf("message: %s, decription: %s",
				res.Message, res.Description)
		}

		log.Println(res.Message)
	}

	return nil
}
