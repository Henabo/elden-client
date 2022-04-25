package gxios

import (
	"fmt"
	"github.com/hiro942/elden-client/global"
	"github.com/hiro942/elden-client/model/request"
	"github.com/hiro942/elden-client/model/response"
	"github.com/hiro942/elden-client/utils"
	"github.com/pkg/errors"
)

func GetBaseURL(satelliteId string) string {
	url := fmt.Sprintf("http://%s/auth", global.SatelliteSocket[satelliteId])
	return url
}

func GetFormatResponse[DataType any](resBytes []byte) response.Response[DataType] {
	res := utils.JsonUnmarshal[response.Response[DataType]](resBytes)
	return res
}

func QuerySatellitePublicKey(id string) (keyHex string, err error) {
	url := fmt.Sprintf("%s/node/satellite/publicKey?id=%s", global.FabricAppBaseUrl, id)
	resBytes := GET(url)
	res := utils.JsonUnmarshal[response.Response[string]](resBytes)
	if res.Code != 0 {
		return "", errors.Errorf("message: %s, decription: %s", res.Message, res.Description)
	}
	return res.Data, nil
}

func Disconnect(targetSatelliteId string, isHandover bool) error {
	disconnectRequest := request.Disconnect{
		Id:         global.MyHashedIMSI,
		MacAddr:    global.MyMacAddr,
		IsHandover: isHandover,
	}
	resBytes := POST(
		fmt.Sprintf("%s/disconnect", GetBaseURL(targetSatelliteId)),
		disconnectRequest,
	)
	if res := utils.JsonUnmarshal[response.Response[any]](resBytes); res.Code != 0 {
		return errors.Errorf("message: %s, decription: %s", res.Message, res.Description)
	}
	return nil
}
