package initialize

import "github.com/hiro942/elden-client/global"

func MockInit() {
	global.SIMCardExist = true

	global.SatelliteIPAddrs["satellite-3333"] = "localhost"

	global.MyHashedIMSI = "hashed-3333"
}
