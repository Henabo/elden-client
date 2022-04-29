package gxios

import (
	"github.com/hiro942/elden-client/global"
	"github.com/hiro942/elden-client/global/url"
	"github.com/hiro942/elden-client/model"
	"github.com/hiro942/elden-client/model/request"
	"github.com/hiro942/elden-client/model/response"
	"github.com/hiro942/elden-client/utils"
	"github.com/pkg/errors"
)

func QuerySatellitePublicKey(id string) (keyHex string, err error) {
	resBytes := GET(url.QuerySatellitePublicKey(id))
	res := utils.JsonUnmarshal[response.Response[string]](resBytes)
	if res.Code != 0 {
		return "", errors.Errorf("message: %s, decription: %s", res.Message, res.Description)
	}
	return res.Data, nil
}

func QueryHasAccessedSatellite(userId string, macAddr string, satelliteId string) bool {
	user, _ := QueryNodeById(userId)

	for _, record := range user.AccessRecord[macAddr] {
		if record.SatelliteId == satelliteId {
			return true
		}
	}
	return false
}

func QueryNodeById(id string) (model.Node, error) {
	resBytes := GET(url.QueryNodeById(id))
	res := utils.JsonUnmarshal[response.Response[model.Node]](resBytes)
	if res.Code != 0 {
		return model.Node{}, errors.Errorf("message: %s, decription: %s", res.Message, res.Description)
	}
	return res.Data, nil
}

func Disconnect(targetSatelliteId string, isHandover bool) error {
	disconnectRequest := request.Disconnect{
		Id:         global.MyHashedIMSI,
		MacAddr:    global.MyMacAddr,
		IsHandover: isHandover,
	}
	resBytes := POST(url.Disconnect(targetSatelliteId), disconnectRequest)
	if res := utils.JsonUnmarshal[response.Response[any]](resBytes); res.Code != 0 {
		return errors.Errorf("message: %s, decription: %s", res.Message, res.Description)
	}

	// 更新认证态
	if !isHandover {
		global.CurrentSession.AuthStatus = false
	}
	return nil
}
