package request

type UserRegister struct {
	Id        string `json:"id"`
	MacAddr   string `json:"macAddr"`
	PublicKey string `json:"publicKey"`
}
