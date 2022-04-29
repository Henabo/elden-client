package initialize

import (
	"github.com/hiro942/elden-client/global"
	"github.com/hiro942/elden-client/service"
	"log"
	"time"
)

func ListenHandover() {
	for {
		if global.HandoverSatellite != "" && global.SatelliteId == global.HandoverSatellite {
			log.Println("Getting handover satellite id.")
			if err := service.HandoverAccess(global.HandoverSatellite); err != nil {
				log.Printf("handover error: %+v", err)
			}
			// 清除该交接卫星
			global.HandoverSatellite = ""
		}

		time.Sleep(time.Microsecond * 500)
	}
}
