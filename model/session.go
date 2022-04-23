package model

type SessionRecord struct {
	SatelliteId    string `json:"satelliteId"`
	SessionKey     string `json:"sessionKey"`
	ExpirationDate int64  `json:"expirationDate"`
}
