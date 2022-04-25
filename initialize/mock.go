package initialize

import "github.com/hiro942/elden-client/global"

func MockInit() {
	global.SIMCardExist = true
	global.MyHashedIMSI = "hashed-9"

	global.MockSatelliteId = "satellite-99"
	global.MockNewSatelliteId = "satellite-999"
	global.SatelliteSocket[global.MockSatelliteId] = "localhost:20000"
	global.SatelliteSocket[global.MockNewSatelliteId] = "localhost:20001"
}
