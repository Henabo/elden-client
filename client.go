package main

import (
	"fmt"
	"github.com/hiro942/elden-client/constant"
	"github.com/hiro942/elden-client/model/enums"
	"github.com/hiro942/elden-client/utils"
	"github.com/tjfoc/gmsm/sm2"
	"log"
	"net"
	"os"
)

type Client struct {
	ID           string
	Status       enums.ClientStatus
	FailedReason string // 认证失败的原因
	MacAddr      string
	PublicKey    *sm2.PublicKey
	PrivateKey   *sm2.PrivateKey

	Ledger         *Ledger
	Cache          *Cache
	SessionRecords []SessionRecord
}

type SessionRecord struct {
	SatelliteID    string `json:"satelliteID"`
	SessionKey     string `json:"sessionKey"`
	ExpirationDate int64  `json:"expirationDate"`
}

func NewClient(id string, ledger *Ledger, cache *Cache) *Client {
	c := &Client{
		ID:     id,
		Status: enums.ClientStatusWithoutSIM,
		Ledger: ledger,
		Cache:  cache,
	}
	c.SetMacAddress()
	c.LoadKeyPair()
	c.LoadSessionRecords()
	return c
}

func (c *Client) LoadSessionRecords() {
	if !utils.FileExist(c.GetSessionRecordFilePath()) {
		c.SessionRecords = make([]SessionRecord, 0)
		return
	}

	// 读文件
	recordsBytes := utils.ReadFile(c.GetSessionRecordFilePath())

	// 若文件本身为空，则不会反序列化成功，直接返回空记录切片即可
	if len(recordsBytes) == 0 {
		c.SessionRecords = make([]SessionRecord, 0)
	}

	// 反序列化
	records := utils.JsonUnmarshal[[]SessionRecord](recordsBytes)

	c.SessionRecords = records
}

func (c *Client) LoadKeyPair() {
	privateKeyPath := c.GetPrivateKeyPath()
	publicKeyPath := c.GetPublicKeyPath()

	// 没有则生成-注册，有则load出来
	if !utils.FileExist(privateKeyPath) || !utils.FileExist(publicKeyPath) {
		// 生成公私钥
		c.PrivateKey, c.PublicKey = utils.GenerateSm2KeyPair()

		// 公私钥转为pem格式
		privateKeyPem := utils.WritePrivateKeyToPem(c.PrivateKey, []byte(config.PrivateKeyPwd))
		publicKeyPem := utils.WritePublicKeyToPem(c.PublicKey)

		// 生成本地目录
		err := os.MkdirAll(fmt.Sprintf("./.crypto/%s", c.ID), constant.DefaultFilePerm)
		if err != nil {
			log.Panicln(fmt.Printf("failed to make directory: %+v", err))
		}

		// 公私钥写入文件
		utils.WriteFile(c.GetPrivateKeyPath(), privateKeyPem)
		utils.WriteFile(c.GetPublicKeyPath(), publicKeyPem)

		// 注册进区块链
		if err = c.Ledger.Register(c.ID, c.MacAddr, c.PublicKey); err != nil {
			log.Panicln(fmt.Printf("failed to register: %+v", err))
		}

	} else {
		// 读pem格式私钥和公钥
		privateKeyPem := utils.ReadFile(privateKeyPath)
		publicKeyPem := utils.ReadFile(publicKeyPath)

		// 公私钥转化
		privateKey := utils.ReadPrivateKeyFromPem(privateKeyPem, []byte(config.PrivateKeyPwd))
		publicKey := utils.ReadPublicKeyFromPem(publicKeyPem)

		// 更新
		c.PrivateKey = privateKey
		c.PublicKey = publicKey
	}
}

func (c *Client) WriteNewSessionRecord(newRecord SessionRecord) {
	// 若不存在记录文件则先创建再读
	if !utils.FileExist(c.GetSessionRecordFilePath()) {
		utils.WriteFile(c.GetSessionRecordFilePath(), nil)
	}

	// 添加新数据
	c.SessionRecords = append(c.SessionRecords, newRecord)
	recordsBytes := utils.JsonMarshal(c.SessionRecords)

	// 写入更新后的内容
	utils.WriteFile(c.GetSessionRecordFilePath(), recordsBytes)
}

func (c *Client) SetMacAddress() {
	netInterface, err := net.InterfaceByName("en0")
	if err != nil {
		log.Panicln("getting net interface error:", err)
	}
	c.MacAddr = netInterface.HardwareAddr.String()
}

func (c *Client) GetPrivateKeyPath() string {
	return constant.CryptoPath + c.ID + "/" + constant.PrivateKeyPemFileName
}

func (c *Client) GetPublicKeyPath() string {
	return constant.CryptoPath + c.ID + "/" + constant.PublicKeyPemFileName
}

func (c *Client) GetSessionRecordFilePath() string {
	return constant.CryptoPath + c.ID + "/" + constant.SessionRecordsFileName
}
