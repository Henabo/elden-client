package service

import (
	"github.com/hiro942/elden-client/global"
	"github.com/hiro942/elden-client/utils/gxios"
	"time"
)

func Authentication() {
	for {
		// SIM卡存在且处于未认证状态下，开始接入
		if global.SIMCardExist && global.CurrentSession.AuthStatus == false && global.SatelliteId != "" {
			// 判断账本中有无该卫星的接入记录，选择首次认证或常规认证
			if !gxios.QueryHasAccessedSatellite(global.MyHashedIMSI, global.MyMacAddr, global.SatelliteId) {
				// 没有记录：首次认证
				if err := FirstAccess(global.SatelliteId); err != nil {
					panic(err)
				}
			} else {
				// 有记录：常规认证
				if err := NormalAccess(global.SatelliteId); err != nil {
					panic(err)
				}
			}
		}

		time.Sleep(time.Microsecond * 500)
	}
}
