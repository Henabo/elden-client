package model

type Session struct {
	AuthStatus  bool   // 是否已经接入卫星
	SatelliteId string // 卫星ID
	Socket      string // 卫星套接字
	AccessType  string
	SessionKey  []byte
}

type SessionRecord struct {
	SatelliteId    string `json:"satelliteId"`
	SessionKey     string `json:"sessionKey"`
	ExpirationDate int64  `json:"expirationDate"`
}
