package initialize

import (
	"fmt"
	"github.com/hiro942/elden-client/global"
	"github.com/hiro942/elden-client/router"
	"github.com/hiro942/elden-client/service"
	"github.com/hiro942/elden-client/utils"
	"net"
	"time"
)

// ListenBroadcast 接收卫星广播UDP报文
func ListenBroadcast() {
	packageConn, err := net.ListenPacket("udp4", ":8829")
	if err != nil {
		panic(err)
	}
	defer packageConn.Close()

	for {
		buf := make([]byte, 1024)
		n, addr, err := packageConn.ReadFrom(buf)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s sent this: %s\n", addr, buf[:n])

		time.Sleep(time.Second * 3)
	}
}

func ListenSIMCardStatus() {
	// 通过系统提供的接口取得IMSI
	for {
		/*
			if 获取IMSI成功(sim卡已插入) && global.SIMCardExist == false{
				global.SIMCardExist = true
			}
			if 获取IMSI失败 && global.SIMCardExist == true
				global.SIMCardExist = false
				停止服务
			}
		*/
		time.Sleep(time.Second * 5)
	}
}

func SysInit() {
	go ListenSIMCardStatus() // 监听SIM卡的插拔
	go ListenBroadcast()     // 监听卫星广播

	// 读取mac地址
	netInterface, err := net.InterfaceByName("en0")
	if err != nil {
		panic(fmt.Errorf("getting net interfact error: %v", err))
	}
	global.MyMacAddr = netInterface.HardwareAddr.String()

	global.PrivateKeyPath = global.CryptoPath + global.MyHashedIMSI + "/" + global.PrivateKeyPemFileName
	global.PublicKeyPath = global.CryptoPath + global.MyHashedIMSI + "/" + global.PublicKeyPemFileName
	global.SessionRecordsFilePath = global.CryptoPath + global.MyHashedIMSI + "/" + global.SessionRecordsFileName

	// 读取公私钥，或
	// 生成公私钥，并将公钥注册上链
	privateKeyPath := global.CryptoPath + global.MyHashedIMSI + "/" + global.PrivateKeyPemFileName
	publicKeyPath := global.CryptoPath + global.MyHashedIMSI + "/" + global.PublicKeyPemFileName
	if !utils.FileExist(privateKeyPath) || !utils.FileExist(publicKeyPath) {
		if err := service.Register(); err != nil {
			panic(fmt.Errorf("register error: %v", err))
		}
	} else {
		utils.ReadKeyPair()
	}

	r := router.Routers()
	r.Run(global.DefaultAuthenticationPort)
}
