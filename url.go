package main

import (
	"fmt"
	"github.com/hiro942/elden-client/model/enums"
)

type URLService struct {
	Cache *Cache
}

func NewURLService(cache *Cache) *URLService {
	return &URLService{Cache: cache}
}

func (url *URLService) QuerySatellitePublicKeyHex(satelliteID string) string {
	return config.FabricAppHostPath + fmt.Sprintf("/node/satellite/publicKey?id=%s", satelliteID)
}

func (url *URLService) QueryNodeByID(id string) string {
	return config.FabricAppHostPath + fmt.Sprintf("/node/search?id=%s", id)
}

func (url *URLService) Disconnect(satelliteID string) string {
	socket := url.Cache.GetSatelliteSocket(satelliteID)
	return fmt.Sprintf("http://%s/auth/disconnect", socket)
}

func (url *URLService) FirstAccessStep1(satelliteID string) string {
	socket := url.Cache.GetSatelliteSocket(satelliteID)
	return fmt.Sprintf("http://%s/auth/first/step1", socket)
}

func (url *URLService) FirstAccessStep2(clientID, clientMacAddr, satelliteID string) string {
	socket := url.Cache.GetSatelliteSocket(satelliteID)
	return fmt.Sprintf("http://%s/auth/first/step2?id=%s&mac=%s",
		socket, clientID, clientMacAddr)
}

func (url *URLService) NormalAccess(satelliteID string, accessType enums.AccessType, accessKeyMode enums.AccessKeyMode) string {
	socket := url.Cache.GetSatelliteSocket(satelliteID)
	return fmt.Sprintf("http://%s/auth/%s?type=%s", socket, accessType.Format(), accessKeyMode.Format())
}
