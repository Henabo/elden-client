package initialize

import (
	"github.com/hiro942/elden-client/global"
	"github.com/hiro942/elden-client/router"
	"github.com/hiro942/elden-client/service"
	"github.com/hiro942/elden-client/utils"
	"github.com/hiro942/elden-client/utils/gxios"
	"log"
	"time"
)

func SysInit() {
	//go ListenSIMCardStatus()                 // 监听SIM卡的插拔

	global.MyMacAddr = utils.GetMacAddress() //读取设备MAC地址

	// SIM卡存在则启动认证服务
	simCardExistsForSeconds := 0
	for {
		if global.SIMCardExist {
			if simCardExistsForSeconds == 0 { // SIM卡插入时开启服务
				//go ListenBroadcast()                     // 开始监听卫星广播
				AuthenticationInit()                                            // 认证加密材料初始化
				go router.Routers().Run(":" + global.DefaultAuthenticationPort) // 启动路由
				go service.Authentication()                                     // 启动认证服务
				go ListenHandover()                                             // 监听卫星切换
			}
			simCardExistsForSeconds += 3
		} else {
			// 若SIM卡拔出，记录服务时间，停止认证服务
			if simCardExistsForSeconds > 0 {
				log.Printf("SIM card with hashed IMSI <%s> exists for %d seconds",
					global.MyHashedIMSI, simCardExistsForSeconds)
				gxios.Disconnect(global.CurrentSession.SatelliteId, false) // 断开与卫星的会话
				simCardExistsForSeconds = 0                                // 服务时间清零
			}

		}
		time.Sleep(time.Second * 3)
	}

}
