package main

import (
	"github.com/pkg/errors"
	"log"
)

type AuthenticationService struct {
	Session *Session
}

func NewAuthenticationService(session *Session) *AuthenticationService {
	return &AuthenticationService{
		Session: session,
	}
}

func (auth *AuthenticationService) LaunchAuthentication(sid string) error {
	client := auth.Session.Client

	// 判断接入类型
	accessed, err := auth.Session.Client.Ledger.QueryHasAccessedSatellite(client.ID, client.MacAddr, sid)
	if err != nil {
		return errors.Wrap(err, "查账本接入记录失败")
	}

	if !accessed {
		log.Println("")
		err = auth.FirstAccess(sid)
	} else {
		err = auth.NormalAccess(sid)
	}

	return errors.Wrap(err, "接入失败")
}
