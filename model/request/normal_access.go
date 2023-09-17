package request

// NARHashed normal access request with hashed session key
type NARHashed struct {
	HashedIMSI       string `json:"hashedIMSI"`
	MacAddr          string `json:"macAddr"`
	SatelliteID      string `json:"satelliteID"`
	HashedSessionKey string `json:"hashedSessionKey"`
	TimeStamp        int64  `json:"timeStamp"`
}

// NAREncrypted normal access request with encrypted session key
type NAREncrypted struct {
	HashedIMSI  string `json:"hashedIMSI"`
	MacAddr     string `json:"macAddr"`
	SatelliteID string `json:"satelliteID"`
	EncryptedSessionKeyWithExpDate
	TimeStamp int64 `json:"timeStamp"`
}

type EncryptedSessionKeyWithExpDate struct {
	EncryptedSessionKey []byte `json:"encryptedSessionKey"`
	ExpirationDate      int64  `json:"expirationDate"`
}
