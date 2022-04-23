package initialize

import "github.com/hiro942/elden-client/global"

func MockInit() {
	global.SIMCardExist = true

	global.SatelliteIPAddr["satellite-1"] = "localhost"

	global.HashedIMSI = "400100123456788"
}
