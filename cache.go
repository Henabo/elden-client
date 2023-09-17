package main

import (
	"github.com/tjfoc/gmsm/sm2"
	"sync"
)

type Cache struct {
	SatellitePubKeys map[string]*sm2.PublicKey // 卫星公钥
	SatelliteSockets map[string]string         // 卫星套接字
	mu               *sync.Mutex
}

func NewCache() *Cache {
	c := &Cache{
		SatellitePubKeys: make(map[string]*sm2.PublicKey),
		SatelliteSockets: make(map[string]string),
		mu:               &sync.Mutex{},
	}
	c.SetSatelliteSocket("s0", "localhost:19999")
	c.SetSatelliteSocket("s1", "localhost:19998")
	return c
}

func (c *Cache) GetSatellitePublicKey(sid string) *sm2.PublicKey {
	return c.SatellitePubKeys[sid]
}

func (c *Cache) SetSatellitePublicKey(sid string, key *sm2.PublicKey) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.SatellitePubKeys[sid] = key
}

func (c *Cache) GetSatelliteSocket(sid string) string {
	return c.SatelliteSockets[sid]
}

func (c *Cache) SetSatelliteSocket(sid string, socket string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.SatelliteSockets[sid] = socket
}
