package initialize

import (
	"fmt"
	"github.com/hiro942/elden-client/global"
	"github.com/hiro942/elden-client/service"
	"github.com/hiro942/elden-client/utils"
	"log"
)

func AuthenticationInit() {
	// 获取公私钥文件路径
	global.PrivateKeyPath = global.CryptoPath + global.MyHashedIMSI + "/" + global.PrivateKeyPemFileName
	global.PublicKeyPath = global.CryptoPath + global.MyHashedIMSI + "/" + global.PublicKeyPemFileName
	global.SessionRecordsFilePath = global.CryptoPath + global.MyHashedIMSI + "/" + global.SessionRecordsFileName

	// 若存在公私钥，则读取公私钥
	// 若不存在公私钥，则生成公私钥，并将公钥注册上链
	privateKeyPath := global.CryptoPath + global.MyHashedIMSI + "/" + global.PrivateKeyPemFileName
	publicKeyPath := global.CryptoPath + global.MyHashedIMSI + "/" + global.PublicKeyPemFileName
	if !utils.FileExist(privateKeyPath) || !utils.FileExist(publicKeyPath) {
		if err := service.Register(); err != nil {
			log.Panicln(fmt.Errorf("register error: %+v", err))
		}
	} else {
		utils.ReadKeyPair()
	}
}
