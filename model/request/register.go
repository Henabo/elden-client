package request

type UserRegister struct {
	ID        string `json:"id"`
	MacAddr   string `json:"macAddr"`
	PublicKey string `json:"publicKey"`
}
