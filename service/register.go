package service

import (
	"fmt"
	"github.com/hiro942/elden-client/global"
	"github.com/hiro942/elden-client/model/request"
	"github.com/hiro942/elden-client/model/response"
	"github.com/hiro942/elden-client/utils"
	"github.com/hiro942/elden-client/utils/gxios"
	"github.com/pkg/errors"
	"github.com/tjfoc/gmsm/x509"
	"log"
	"os"
)

func Register() error {
	err := os.MkdirAll(fmt.Sprintf("./.crypto/%s", global.MyHashedIMSI), global.DefaultFilePerm)
	if err != nil {
		log.Panicln(fmt.Printf("failed to make directory: %+v", err))
	}

	// 生成公私钥
	global.PrivateKey, global.PublicKey = utils.GenerateSm2KeyPair()

	// 公私钥转为pem格式
	privateKeyPem := utils.WritePrivateKeyToPem(global.PrivateKey)
	publicKeyPem := utils.WritePublicKeyToPem(global.PublicKey)

	// 公私钥写入文件
	utils.WriteFile(global.PrivateKeyPath, privateKeyPem)
	utils.WriteFile(global.PublicKeyPath, publicKeyPem)

	// HTTP[POST] 添加用户公钥至区块链
	responseBytes := gxios.POST(
		global.FabricAppBaseUrl+"/node/user/register",
		request.UserRegister{
			Id:        global.MyHashedIMSI,
			MacAddr:   global.MyMacAddr,
			PublicKey: x509.WritePublicKeyToHex(global.PublicKey),
		},
	)

	// 解析http响应
	if res := utils.JsonUnmarshal[response.Response[any]](responseBytes); res.Code != 0 {
		return errors.Errorf("message: %s, decription: %s",
			res.Message, res.Description)
	}

	// 注册成功
	return nil
}
