package initialize

import (
	"github.com/hiro942/elden-client/global"
	"github.com/hiro942/elden-client/mock"
	"github.com/hiro942/elden-client/utils/gxios"
	"log"
	"os"
	"time"
)

func MyTest() {
	time.Sleep(time.Second * 1)
	log.Println("[Test] 测试开始：监听SIM卡状态")
	time.Sleep(time.Second * 3)

	log.Println("[Test] 检测到SIM卡插入")
	global.MyHashedIMSI = mock.MyHashedId // 获取到SIM卡标识
	global.SIMCardExist = true            // SIM卡存在

	log.Println("[Test] 监听卫星广播")
	global.SatelliteId = mock.SatelliteId
	log.Println("[Test] 接收到卫星ID: ", global.SatelliteId)
	log.Println("[Test] 准备接入...")
	for {
		if global.CurrentSession.AuthStatus == true {
			log.Println("[Test] 首次接入成功，卫星ID:", global.CurrentSession.SatelliteId)
			break
		}
	}

	log.Printf("[Test] 与卫星 %s 断开", global.SatelliteId)
	gxios.Disconnect(global.SatelliteId, false)
	global.SatelliteId = ""

	log.Println("[Test] 等待三秒...")
	time.Sleep(time.Second * 3)

	log.Println("[Test] 测试快速接入，再次接入卫星:", global.SatelliteId)
	global.SatelliteId = mock.SatelliteId

	for {
		if global.CurrentSession.AuthStatus == true {
			log.Println("[Test] 快速接入成功，卫星ID:", global.CurrentSession.SatelliteId)
			break
		}
		time.Sleep(time.Microsecond * 200)
	}

	log.Println("[Test] 测试交接认证, 当前卫星会在3s后发送切换消息")
	for {
		// 假设接收到新卫星ID后3秒，收到了新卫星的广播
		if global.HandoverSatellite != "" {
			time.Sleep(time.Second * 3)
			global.SatelliteId = mock.NewSatelliteId
			break
		}
		time.Sleep(time.Microsecond * 200)
	}

	for {
		if global.CurrentSession.AuthStatus == true && global.CurrentSession.SatelliteId == mock.NewSatelliteId {
			log.Printf("[Test] 已连接至新卫星 %s", mock.NewSatelliteId)
			log.Printf("[Test] 断开与原卫星 %s 的连接", mock.SatelliteId)
			gxios.Disconnect(mock.SatelliteId, true)
			break
		}
		time.Sleep(time.Microsecond * 200)
	}

	log.Printf("[Test] 保持与新卫星 %s 会话 ...", mock.NewSatelliteId)
	time.Sleep(time.Second * 2)

	log.Printf("[Test] SIM卡拔出后，断开当前与卫星 %s 的会话", mock.NewSatelliteId)
	global.SIMCardExist = false

	log.Println("[Test] 测试结束, 退出程序")
	os.Exit(0)
}
