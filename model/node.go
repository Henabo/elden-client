package model

// PublicKeys indicates the structure how public keys are saved
type PublicKeys map[string]string

// UserAccessRecord is single access log
type UserAccessRecord struct {
	AccessType          string `json:"accessType"`          // "first" || "normal" || "handover"
	SatelliteId         string `json:"satelliteId"`         // current satellite
	PreviousSatelliteId string `json:"previousSatelliteId"` // previous satellite in handover
	StartAt             string `json:"startAt"`             // when to start
	EndAt               string `json:"endAt"`               // when to end
}

// UserAccessRecords indicates access records for a specific device
type UserAccessRecords map[string][]UserAccessRecord

type Node struct {
	Id           string            `json:"id"`
	NodeType     string            `json:"nodeType"`
	PublicKey    PublicKeys        `json:"publicKey"`
	AccessRecord UserAccessRecords `json:"accessRecord"`
	CreatedAt    string            `json:"createdAt"`
	UpdatedAt    string            `json:"updatedAt"`
}

type CreateAccessRecord struct {
	Id           string           `json:"id"`
	MacAddr      string           `json:"macAddr"`
	AccessRecord UserAccessRecord `json:"accessRecord"`
}
