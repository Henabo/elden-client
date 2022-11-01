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
)

func FirstAccess(satelliteId string) error {
	log.Println("[First-Access] Go First Access")

	// ********************** Step1 ****************************
	log.Println("[First-Access] First Access Step1: Send basic information.")

	// HTTP[POST] 发送用户基本信息
	resBytes := gxios.POST(
		url.FirstAccessStep1(satelliteId),
		utils.GetMessageWithSig[request.FAR](request.FAR{
			HashedIMSI:  global.MyHashedIMSI,
			MacAddr:     global.MyMacAddr,
			SatelliteId: satelliteId,
		}))

	res := utils.JsonUnmarshal[response.Response[[]byte]](resBytes)
	if res.Code != 0 {
		return errors.Errorf("message: %s, decription: %s",
			res.Message, res.Description)
	}
	log.Println("[First-Access] First Access Step1: Response OK, get session key and random number from the satellite.")

	// 解密data并得到卫星签名消息
	dataWithSigBytes := utils.Sm2Decrypt(global.PrivateKey, res.Data)
	dataWithSig := utils.JsonUnmarshal[response.FARWithSig](dataWithSigBytes)

	// HTTP[GET] 获取卫星公钥
	satellitePublicKeyHex, _ := gxios.QuerySatellitePublicKey(satelliteId)
	satellitePublicKey := utils.ReadPublicKeyFromHex(satellitePublicKeyHex)

	// 保存卫星公钥
	global.SatellitePubKeys[satelliteId] = satellitePublicKey

	// 验证签名
	if !utils.Sm2Verify(satellitePublicKey, dataWithSig.Plain, dataWithSig.Signature) {
		return errors.New("failed to verify signature")
	}

	// 解析出卫星响应明文
	data := utils.JsonUnmarshal[response.FAR](dataWithSig.Plain)

	// ********************** Step2 ****************************
	log.Println("[First-Access] First Access Step2: Return the random number received from the satellite.")

	// HTTP[POST] 发送加密后的随机数
	res2Bytes := gxios.POST(
		url.FirstAccessStep2(satelliteId),
		utils.GetMessageCipherWithSm4[request.FARWithRand](
			request.FARWithRand{Rand: data.Rand},
			[]byte(data.SessionKey),
		))

	if res2 := utils.JsonUnmarshal[response.Response[any]](res2Bytes); res2.Code != 0 {
		return errors.Errorf("message: %s, decription: %s",
			res.Message, res.Description)
	}
	log.Println("[First-Access] First Access Step2: Response OK.")

	// 认证完成后，保存密钥至本地
	utils.WriteNewSessionRecord(model.SessionRecord{
		SatelliteId:    satelliteId,
		SessionKey:     data.SessionKey,
		ExpirationDate: data.ExpirationDate,
	})

	log.Println("[First-Access] First Access Success! Connecting To:", satelliteId)

	// 记录当前会话
	global.CurrentSession = model.Session{
		AuthStatus:  true,
		SatelliteId: satelliteId,
		Socket:      global.SatelliteSockets[satelliteId],
		AccessType:  "first",
		SessionKey:  []byte(data.SessionKey),
	}
	return nil
}
