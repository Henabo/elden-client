package initialize

import (
	"github.com/hiro942/elden-client/global"
	"github.com/hiro942/elden-client/mock"
)

func MockInit() {
	mock.MyHashedId = "hashed-UUU"
	mock.SatelliteId = "satellite-AAA"
	mock.NewSatelliteId = "satellite-BBB"
	global.SatelliteSockets[mock.SatelliteId] = "localhost:20000"
	global.SatelliteSockets[mock.NewSatelliteId] = "localhost:20001"
}
