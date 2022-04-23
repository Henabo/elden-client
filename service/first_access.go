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
)

func FirstAccess(satelliteId string) error {
	// ********************** Step1 ****************************

	// HTTP[POST] 发送用户基本信息
	url := gxios.GetBaseURL(satelliteId) + "/first/step1"
	resBytes := gxios.POST(url, utils.GetMessageWithSig[request.FAR](request.FAR{
		HashedIMSI:  global.MyHashedIMSI,
		MacAddr:     global.MyMacAddr,
		SatelliteId: satelliteId,
	}))

	res := gxios.GetFormatResponse[string](resBytes)
	if res.Code != 0 {
		return errors.Errorf("message: %s, decription: %s",
			res.Message, res.Description)
	}

	// 解密data并得到卫星签名消息
	dataWithSigBytes := utils.Sm2Decrypt(global.PrivateKey, []byte(res.Data))
	dataWithSig := utils.JsonUnmarshal[response.FARWithSig](dataWithSigBytes)

	// HTTP[GET] 获取卫星公钥
	satellitePublicKeyHex := gxios.QuerySatellitePublicKey(satelliteId)
	satellitePublicKey, err := x509.ReadPublicKeyFromHex(satellitePublicKeyHex)
	if err != nil {
		log.Panic(fmt.Printf("failed to resolve satellite's public key: %+v", err))
	}
	global.SatellitePubKeys[satelliteId] = satellitePublicKey

	// 验证签名
	if !utils.Sm2Verify(satellitePublicKey, dataWithSig.Plain, dataWithSig.Signature) {
		return errors.New("failed to verify signature")
	}

	// 解析出卫星响应明文
	data := utils.JsonUnmarshal[response.FAR](dataWithSig.Plain)

	// ********************** Step2 ****************************

	// HTTP[POST] 发送加密后的随机数
	url = fmt.Sprintf("%s/first/step2?id=%s&mac=%s", gxios.GetBaseURL(satelliteId), global.MyHashedIMSI, global.MyMacAddr)
	res2Bytes := gxios.POST(url, utils.GetMessageCipherWithSm4[request.FARWithRand](
		request.FARWithRand{Rand: data.Rand},
		[]byte(data.SessionKey),
	))

	if res = utils.JsonUnmarshal[response.Response[string]](res2Bytes); res.Code != 0 {
		return errors.Errorf("message: %s, decription: %s",
			res.Message, res.Description)
	}

	log.Println("First authentication success!")

	// 认证完成后，保存密钥至本地
	utils.WriteNewSessionRecord(model.SessionRecord{
		SatelliteId:    satelliteId,
		SessionKey:     data.SessionKey,
		ExpirationDate: data.ExpirationDate,
	})

	return nil
}
