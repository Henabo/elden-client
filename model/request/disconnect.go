package request

type Disconnect struct {
	Id         string `json:"id"`
	MacAddr    string `json:"macAddr"`
	IsHandover bool   `json:"isHandover"`
}
