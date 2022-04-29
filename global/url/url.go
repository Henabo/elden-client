/*
	HTTP请求地址
*/

package url

import (
	"fmt"
	"github.com/hiro942/elden-client/global"
)

func QuerySatellitePublicKey(id string) string {
	return global.FabricAppBaseUrl + fmt.Sprintf("/node/satellite/publicKey?id=%s", id)
}

func QueryNodeById(id string) string {
	return global.FabricAppBaseUrl + fmt.Sprintf("/node/search?id=%s", id)
}

func Disconnect(satelliteId string) string {
	return fmt.Sprintf("http://%s/auth/disconnect", global.SatelliteSockets[satelliteId])
}

func FirstAccessStep1(satelliteId string) string {
	return fmt.Sprintf("http://%s/auth/first/step1", global.SatelliteSockets[satelliteId])
}

func FirstAccessStep2(satelliteId string) string {
	return fmt.Sprintf("http://%s/auth/first/step2?id=%s&mac=%s",
		global.SatelliteSockets[satelliteId], global.MyHashedIMSI, global.MyMacAddr)
}

func NormalAccess(satelliteId string, isHandover bool, isKeyTypeHashed bool) string {
	var accessType = "normal"
	var keyType = "hashed"
	if !isKeyTypeHashed {
		keyType = "encrypted"
	}
	if isHandover {
		accessType = "handover"
	}
	return fmt.Sprintf("http://%s/auth/%s?type=%s", global.SatelliteSockets[satelliteId], accessType, keyType)
}
