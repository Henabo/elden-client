package initialize

import (
	"fmt"
	"github.com/hiro942/elden-client/global"
	"github.com/hiro942/elden-client/router"
	"github.com/hiro942/elden-client/service"
	"github.com/hiro942/elden-client/utils"
	"github.com/hiro942/elden-client/utils/gxios"
	"log"
	"net"
	"time"
)

// ListenBroadcast 接收卫星广播UDP报文
func ListenBroadcast() {
	packageConn, err := net.ListenPacket("udp4", ":8829")
	if err != nil {
		log.Panicln(err)
	}
	defer packageConn.Close()

	for {
		buf := make([]byte, 1024)
		n, addr, err := packageConn.ReadFrom(buf)
		if err != nil {
			log.Panicln(err)
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
	//go ListenSIMCardStatus()                 // 监听SIM卡的插拔
	//go ListenBroadcast()                     // 监听卫星广播
	global.MyMacAddr = utils.GetMacAddress() //读取设备MAC地址

	// SIM卡存在则启动认证服务
	simCardExistsForSeconds := 0
	for {
		if global.SIMCardExist { // SIM卡存在则启动
			if simCardExistsForSeconds == 0 {
				go router.Routers().Run(":19999") // 启动路由
				AuthenticationInit()              // 认证服务初始化
				err := service.FirstAccess(global.MockSatelliteId)
				if err != nil {
					log.Panicln(fmt.Errorf("register error: %+v", err))
				}

				log.Printf("Disconnect With %s\n", global.MockSatelliteId)
				// 断开连接
				gxios.Disconnect(global.MockSatelliteId, false)

				log.Println("Sleep for 3 seconds ...")
				time.Sleep(time.Second * 3)

				err = service.NormalAccess(global.MockSatelliteId)
				if err != nil {
					log.Panicln(fmt.Errorf("register error: %v", err))
				}

				//go router.Routers().Run(global.DefaultAuthenticationPort) // 启动路由
			}
			simCardExistsForSeconds += 3
		} else {
			// 若SIM卡拔出，记录服务时间，停止认证服务
			if simCardExistsForSeconds > 0 {
				log.Printf("SIMCARD with hashed IMSI <%s> exists for %d seconds",
					global.MyHashedIMSI, simCardExistsForSeconds)
			}
			// todo 停止认证服务
		}
		time.Sleep(time.Second * 3)
	}

}

func AuthenticationInit() {
	// 更新（切换）文件路径
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

func Authentication() {
	/*
		若没接入任何卫星
			根据接收到的卫星ID，读本地文件，判断是否接入过该卫星（是否有会话密钥记录）
			若无记录，那么走First Access
			若有记录，那么走Normal Access
		若已接入某一卫星

	*/
}
