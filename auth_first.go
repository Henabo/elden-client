package main

import (
	"fmt"
	"github.com/hiro942/elden-client/model/enums"
	"github.com/hiro942/elden-client/model/request"
	"github.com/hiro942/elden-client/model/response"
	"github.com/hiro942/elden-client/utils"
	"github.com/hiro942/elden-client/utils/ghttp"
	"github.com/pkg/errors"
	"log"
)

func (auth *AuthenticationService) FirstAccess(sid string) error {
	session := auth.Session
	client := session.Client

	client.FailedReason = ""
	client.Status = enums.ClientStatusVerifying

	// 发送用户基本信息
	log.Printf("【首次认证】第一次握手: 用户发送基本信息，包括终端ID「%s」，MAC地址：「%s」\n", client.ID, client.MacAddr)
	cipher, err := ghttp.POST[enums.Cipher](
		auth.Session.Client.Ledger.URL.FirstAccessStep1(sid),
		utils.GetMessageWithSig[request.FAR](request.FAR{
			HashedIMSI:  client.ID,
			MacAddr:     client.MacAddr,
			SatelliteID: sid,
		}, client.PrivateKey))
	if err != nil {
		client.FailedReason = err.Error()
		return err
	}

	// 解密并得到卫星签名消息
	farWithSignBytes := utils.Sm2Decrypt(client.PrivateKey, cipher)
	farWithSign := utils.JsonUnmarshal[response.FARWithSign](farWithSignBytes)

	// HTTP[GET] 获取卫星公钥
	satellitePublicKeyHex, err := client.Ledger.QuerySatellitePublicKeyHex(sid)
	if err != nil {
		client.FailedReason = err.Error()
		return errors.Wrap(err, fmt.Sprintf("账本查询卫星「%s」公钥失败", sid))
	}
	// 缓存卫星公钥
	client.Cache.SetSatellitePublicKey(sid, utils.ReadPublicKeyFromHex(satellitePublicKeyHex))

	// 验证签名
	signVerified := utils.Sm2Verify(client.Cache.GetSatellitePublicKey(sid), farWithSign.Plain, farWithSign.Signature)
	if !signVerified {
		client.FailedReason = err.Error()
		return errors.Errorf("无法验证卫星「%s」的消息签名", sid)
	}

	// 解析出卫星响应明文
	far := utils.JsonUnmarshal[response.FAR](farWithSign.Plain)

	log.Printf("【首次认证】第二次握手: 卫星发送了会话密钥「%s」 和随机数 「%d」.\n", far.SessionKey, far.Rand)

	// HTTP[POST] 发送加密后的随机数
	log.Printf("【首次认证】第三次握手: 终端返回从卫星「%s」处接收到的随机数「%d」\n", sid, far.Rand)
	_, err = ghttp.POST[any](
		auth.Session.Client.Ledger.URL.FirstAccessStep2(client.ID, client.MacAddr, sid),
		utils.GetMessageCipherWithSm4[request.FARWithRand](
			request.FARWithRand{Rand: far.Rand},
			[]byte(far.SessionKey),
		))
	if err != nil {
		client.FailedReason = err.Error()
		return err
	}

	log.Printf("【首次认证】第四次握手: 卫星「%s」返回接入成功信息\n", sid)

	// 认证完成后，保存密钥至本地
	client.WriteNewSessionRecord(SessionRecord{
		SatelliteID:    sid,
		SessionKey:     far.SessionKey,
		ExpirationDate: far.ExpirationDate,
	})

	log.Println("【首次认证】成功")

	client.Status = enums.ClientStatusVerifySuccess
	session.Status = enums.SessionStatusProcessing
	session.AccessType = enums.AccessTypeStrict
	session.SatelliteID = sid
	session.SatelliteSocket = client.Cache.GetSatelliteSocket(sid)
	session.SessionKey = []byte(far.SessionKey)

	return nil
}
