package request

type Disconnect struct {
	ID         string `json:"id"`
	MacAddr    string `json:"macAddr"`
	IsHandover bool   `json:"isHandover"`
}
