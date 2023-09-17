package request

// FAR first access request
type FAR struct {
	HashedIMSI  string `json:"hashedIMSI"`
	MacAddr     string `json:"macAddr"`
	SatelliteID string `json:"satelliteID"`
}

type FARWithRand struct {
	Rand int `json:"rand"`
}
