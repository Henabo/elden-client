package request

// NAR normal access request
type NAR[T ~string | SessionKeyKeyWithExpDate] struct {
	HashedIMSI     string `json:"hashedIMSI"`
	MacAddr        string `json:"macAddr"`
	SatelliteId    string `json:"satelliteId"`
	SessionKeyInfo T      `json:"sessionKey"`
	TimeStamp      int64  `json:"timeStamp"`
}

type HashedSessionKey string

type SessionKeyKeyWithExpDate struct {
	EncryptedSessionKey string `json:"encryptedSessionKey"`
	ExpirationDate      int64  `json:"expirationDate"`
}
