package gxios

import (
	"fmt"
	"github.com/hiro942/elden-client/global"
	"github.com/hiro942/elden-client/model/response"
	"github.com/hiro942/elden-client/utils"
)

func GetBaseURL(satelliteId string) string {
	url := fmt.Sprintf("http://%s:%s/auth",
		global.SatelliteIPAddrs[satelliteId], global.DefaultAuthenticationPort)
	return url
}

func GetFormatResponse[DataType any](resBytes []byte) response.Response[DataType] {
	res := utils.JsonUnmarshal[response.Response[DataType]](resBytes)
	return res
}

func QuerySatellitePublicKey(id string) (keyHex string) {
	url := fmt.Sprintf("%s/node/satellite/publicKey?id=%s", global.FabricAppBaseUrl, id)
	resBytes := GET(url)
	res := GetFormatResponse[string](resBytes)
	return res.Data
}
