package request

// FAR first access request
type FAR struct {
	HashedIMSI  string `json:"hashedIMSI"`
	MacAddr     string `json:"macAddr"`
	SatelliteId string `json:"satelliteId"`
}

type FARWithRand struct {
	Rand int `json:"rand"`
}
