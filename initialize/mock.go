package initialize

import (
	"github.com/hiro942/elden-client/global"
	"github.com/hiro942/elden-client/mock"
)

func MockInit() {
	mock.MyHashedId = "hashed-19999-test"
	mock.SatelliteId = "satellite-20000-test"
	mock.NewSatelliteId = "satellite-20001-test"
	global.SatelliteSockets[mock.SatelliteId] = "localhost:20000"
	global.SatelliteSockets[mock.NewSatelliteId] = "localhost:20001"
}
