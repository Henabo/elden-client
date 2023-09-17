package main

import (
	"github.com/hiro942/elden-client/model/enums"
	"github.com/hiro942/elden-client/model/request"
	"github.com/hiro942/elden-client/utils/ghttp"
	"log"
)

type Session struct {
	SatelliteID     string // 卫星ID
	Status          enums.SessionStatus
	SatelliteSocket string // 卫星套接字
	AccessType      enums.AccessType
	SessionKey      []byte

	Client *Client
}

func NewSession(client *Client) *Session {
	return &Session{
		SatelliteID:     "",
		Status:          enums.SessionStatusNull,
		SatelliteSocket: "",
		AccessType:      0,
		SessionKey:      nil,
		Client:          client,
	}
}

func (s *Session) Disconnect(isHandover bool) error {
	_, err := ghttp.POST[any](
		s.Client.Ledger.URL.Disconnect(s.SatelliteID),
		request.Disconnect{
			ID:         s.Client.ID,
			MacAddr:    s.Client.MacAddr,
			IsHandover: isHandover,
		})
	if err != nil {
		return err
	}

	// 更新认证态
	if !isHandover {
		s.Client.Status = enums.ClientStatusWaitVerify
		s.Client.Status = enums.SessionStatusNull
	}

	log.Printf("已断开与卫星「%s」的会话", s.SatelliteID)
	return nil
}
