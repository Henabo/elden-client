package initialize

import (
	"github.com/hiro942/elden-client/global"
	"github.com/hiro942/elden-client/mock"
	"github.com/hiro942/elden-client/utils/gxios"
	"log"
	"time"
)

func MyTest() {
	time.Sleep(time.Second * 1)
	log.Println("[Test] Begin")
	log.Println("[Test] Watching SIM card status")
	time.Sleep(time.Second * 3)

	log.Println("[Test] SIM card is ready")
	global.MyHashedIMSI = mock.MyHashedId // 获取到SIM卡标识
	global.SIMCardExist = true            // 修改SIM卡状态

	log.Println("[Test] Watching broadcast from satellites ...")
	global.SatelliteId = mock.SatelliteId
	log.Println("[Test] Received satellite ID: ", global.SatelliteId)

	log.Println("[Test] Testing the first access process ...")

	for {
		if global.CurrentSession.AuthStatus == true {
			log.Printf("[Test] First access is successful. Current session with %s", global.CurrentSession.SatelliteId)
			break
		}
	}

	log.Printf("[Test] Disconnect from %s", global.SatelliteId)
	gxios.Disconnect(global.SatelliteId, false)
	global.SatelliteId = ""

	log.Println("[Test] Sleep for 90 seconds ...")
	time.Sleep(time.Second * 90)

	global.SatelliteId = mock.SatelliteId
	log.Printf("[Test] To test the normal access process, reconnecting to %s", global.SatelliteId)

	for {
		if global.CurrentSession.AuthStatus == true {
			log.Printf("[Test] Normal access is successful. Current session with %s:", global.CurrentSession.SatelliteId)
			break
		}
		time.Sleep(time.Microsecond * 200)
	}

	log.Println("[Test] To test the handover access process, current satellite will send handover message in 90 seconds ...")
	for {
		// 假设接收到新卫星ID后3秒，收到了新卫星的广播
		if global.HandoverSatellite != "" {
			time.Sleep(time.Second * 90)
			global.SatelliteId = mock.NewSatelliteId
			break
		}
		time.Sleep(time.Microsecond * 200)
	}

	for {
		if global.CurrentSession.AuthStatus == true && global.CurrentSession.SatelliteId == mock.NewSatelliteId {
			log.Printf("[Test] Handover access is successful. Current session with %s", mock.NewSatelliteId)
			log.Printf("[Test] Disconnect from %s", mock.SatelliteId)
			gxios.Disconnect(mock.SatelliteId, true)
			break
		}
		time.Sleep(time.Microsecond * 200)
	}

	log.Printf("[Test] Keep session with %s ...", mock.NewSatelliteId)
	time.Sleep(time.Second * 10)

	log.Printf("[Test] SIM card out of device，disconnect the session from %s", mock.NewSatelliteId)
	global.SIMCardExist = false

	time.Sleep(time.Second * 10)

	log.Println("[Test] End")
	//os.Exit(0)
}
