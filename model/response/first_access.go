package response

type FAR struct {
	SessionKey     string `json:"sessionKey"`
	ExpirationDate int64  `json:"expirationDate"`
	Timestamp      int64  `json:"timestamp"`
	Rand           int    `json:"rand"`
}

type FARWithSign struct {
	Plain     []byte `json:"plain"`
	Signature []byte `json:"signature"`
}
