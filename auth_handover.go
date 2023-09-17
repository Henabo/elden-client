package main

import (
	"github.com/hiro942/elden-client/constant"
	"github.com/hiro942/elden-client/model/enums"
	"github.com/hiro942/elden-client/model/request"
	"github.com/hiro942/elden-client/utils"
	"github.com/hiro942/elden-client/utils/ghttp"
	"log"
	"time"
)

func (auth *AuthenticationService) HandoverAccess(sid string) error {
	session := auth.Session
	client := session.Client

	client.Status = enums.ClientStatusVerifying

	// 是否存在该卫星公钥
	satellitePublicKey := client.Cache.GetSatellitePublicKey(sid)
	if satellitePublicKey == nil {
		// 查询卫星公钥
		satellitePublicKeyHex, err := client.Ledger.QuerySatellitePublicKeyHex(sid)
		if err != nil {
			client.FailedReason = err.Error()
			return err
		}
		// 保存卫星公钥
		client.Cache.SetSatellitePublicKey(sid, utils.ReadPublicKeyFromHex(satellitePublicKeyHex))
	}

	log.Println("【交接认证】判断本地保存的会话密钥是否过期")
	records := client.SessionRecords
	var sessionKey []byte
	var expDate int64
	var isExpired = false
	var sessionKeyInfo any

	// 读会话记录文件，判断有无密钥记录
	hasRecord := false
	for _, record := range records {
		if record.SatelliteID == sid {
			hasRecord = true
			break
		}
	}

	// 没有则创建，并写入本地
	// 有则读取记录
	if !hasRecord {
		sessionKey = utils.GenerateSm4Key()
		expDate = time.Now().Unix() + constant.DefaultSessionKeyAge
		sessionKeyInfo = request.EncryptedSessionKeyWithExpDate{
			EncryptedSessionKey: utils.Sm2Encrypt(client.Cache.GetSatellitePublicKey(sid), sessionKey),
			ExpirationDate:      expDate,
		}

		client.WriteNewSessionRecord(SessionRecord{
			SatelliteID:    sid,
			SessionKey:     string(sessionKey),
			ExpirationDate: expDate,
		})

	} else {

		for _, record := range records {
			if record.SatelliteID == sid {
				hasRecord = true
				// 判断密钥是否过期。若密钥已经过期，则生成一个新密钥，并更新会话记录文件。
				if record.ExpirationDate < time.Now().Unix()+60 {
					isExpired = true
					sessionKey = utils.GenerateSm4Key()
					expDate = time.Now().Unix() + constant.DefaultSessionKeyAge
					log.Println("Generated Session Key:", string(sessionKey))
					record.SessionKey = string(sessionKey)                                         // 更新会话密钥
					utils.WriteFile(client.GetSessionRecordFilePath(), utils.JsonMarshal(records)) // 会话记录写回
					sessionKeyInfo = request.EncryptedSessionKeyWithExpDate{
						EncryptedSessionKey: utils.Sm2Encrypt(client.Cache.GetSatellitePublicKey(sid), sessionKey),
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
		log.Printf("【交接认证】第一次握手：密钥处于有效期内，终端「%s」发送哈希后的会话密钥给卫星「%s」。\n", client.ID, sid)

		_, err := ghttp.POST[any](
			auth.Session.Client.Ledger.URL.NormalAccess(sid, enums.AccessTypeHandover, enums.AccessKeyModeHashed),
			utils.GetMessageWithSig[request.NARHashed](request.NARHashed{
				HashedIMSI:       client.ID,
				MacAddr:          client.MacAddr,
				SatelliteID:      sid,
				HashedSessionKey: sessionKeyInfo.(string),
				TimeStamp:        time.Now().Unix(),
			}, client.PrivateKey))
		if err != nil {
			client.FailedReason = err.Error()
			return err
		}

		log.Printf("【交接认证】第二次握手：卫星「%s」返回成功信息。\n", sid)

	} else {
		log.Printf("【交接认证】第一次握手：密钥已过期，终端「%s」生成新的会话密钥加密后发送给卫星「%s」。\n", client.ID, sid)

		_, err := ghttp.POST[any](
			auth.Session.Client.Ledger.URL.NormalAccess(sid, enums.AccessTypeHandover, enums.AccessKeyModeEncrypted),
			utils.GetMessageWithSig[request.NAREncrypted](request.NAREncrypted{
				HashedIMSI:                     client.ID,
				MacAddr:                        client.MacAddr,
				SatelliteID:                    sid,
				EncryptedSessionKeyWithExpDate: sessionKeyInfo.(request.EncryptedSessionKeyWithExpDate),
				TimeStamp:                      time.Now().Unix(),
			}, client.PrivateKey))
		if err != nil {
			client.FailedReason = err.Error()
			return err
		}
	}

	log.Println("【交接认证】成功。")
	if err := auth.Session.Disconnect(true); err != nil {
		return err
	}

	// 记录当前会话
	client.Status = enums.ClientStatusVerifySuccess
	session.Status = enums.SessionStatusProcessing
	session.AccessType = enums.AccessTypeHandover
	session.SatelliteID = sid
	session.SatelliteSocket = client.Cache.GetSatelliteSocket(sid)
	session.SessionKey = sessionKey

	return nil
}
