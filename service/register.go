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

	// 私钥转为pem格式
	privateKeyPem, err := x509.WritePrivateKeyToPem(global.PrivateKey, global.PrivateKeyPwd)
	if err != nil {
		return errors.Wrap(err, "failed to convert private key to pem")
	}

	// 写私钥
	utils.WriteFile(global.PrivateKeyPath, privateKeyPem)

	// 公钥转为pem格式
	publicKeyPem, err := x509.WritePublicKeyToPem(global.PublicKey)
	if err != nil {
		return errors.Wrap(err, "failed to convert public key to pem")
	}

	// 写公钥
	utils.WriteFile(global.PublicKeyPath, publicKeyPem)

	// 获取公钥的十六进制字符串形式
	publicKeyHex := x509.WritePublicKeyToHex(global.PublicKey)

	// HTTP[POST] 添加用户公钥至区块链
	responseBytes := gxios.POST(
		global.FabricAppBaseUrl+"/node/user/register",
		request.UserRegister{
			Id:        global.MyHashedIMSI,
			MacAddr:   global.MyMacAddr,
			PublicKey: publicKeyHex,
		},
	)

	// 解析http响应
	res := utils.JsonUnmarshal[response.Response[any]](responseBytes)

	// 服务端返回错误
	if res.Code != 0 {
		return errors.Errorf("message: %s, decription: %s",
			res.Message, res.Description)
	}

	// 注册成功
	return nil
}
