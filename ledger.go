package main

import (
	"github.com/hiro942/elden-client/model"
	"github.com/hiro942/elden-client/model/request"
	"github.com/hiro942/elden-client/utils/ghttp"
	"github.com/tjfoc/gmsm/sm2"
	"github.com/tjfoc/gmsm/x509"
)

type Ledger struct {
	URL *URLService
}

func NewLedger(url *URLService) *Ledger {
	return &Ledger{
		URL: url,
	}
}

func (ledger *Ledger) QuerySatellitePublicKeyHex(id string) (keyHex string, err error) {
	key, err := ghttp.GET[string](ledger.URL.QuerySatellitePublicKeyHex(id))
	if err != nil {
		return "", err
	}
	return key, nil
}

func (ledger *Ledger) QueryHasAccessedSatellite(clientID, clientMacAddr, satelliteID string) (bool, error) {
	user, err := ledger.QueryNodeByID(clientID)
	if err != nil {
		return false, err
	}

	for _, record := range user.AccessRecord[clientMacAddr] {
		if record.SatelliteID == satelliteID {
			return true, nil
		}
	}
	return false, nil
}

func (ledger *Ledger) QueryNodeByID(id string) (model.Node, error) {
	node, err := ghttp.GET[model.Node](ledger.URL.QueryNodeByID(id))
	if err != nil {
		return model.Node{}, err
	}
	return node, nil
}

func (ledger *Ledger) Register(clientID, clientMadAddr string, publicKey *sm2.PublicKey) error {
	// HTTP[POST] 添加用户公钥至区块链
	_, err := ghttp.POST[any](
		config.FabricAppHostPath+"/node/user/register",
		request.UserRegister{
			ID:        clientID,
			MacAddr:   clientMadAddr,
			PublicKey: x509.WritePublicKeyToHex(publicKey),
		},
	)
	return err
}
