package utils

import (
	"encoding/json"
	"fmt"
	"github.com/hiro942/elden-client/global"
	"github.com/hiro942/elden-client/model"
	"github.com/tjfoc/gmsm/x509"
	"log"
	"net"
	"os"
)

func GetMacAddress() string {
	netInterface, err := net.InterfaceByName("en0")
	if err != nil {
		log.Panicln(fmt.Errorf("getting net interfact error: %v", err))
	}
	return netInterface.HardwareAddr.String()
}

func FileExist(path string) bool {
	_, err := os.Lstat(path)
	return !os.IsNotExist(err)
}

func ReadKeyPair() {
	// 读pem格式私钥和公钥
	privateKeyPem := ReadFile(global.PrivateKeyPath)
	publicKeyPem := ReadFile(global.PublicKeyPath)

	// 转化私钥
	privateKey, err := x509.ReadPrivateKeyFromPem(privateKeyPem, global.PrivateKeyPwd)
	if err != nil {
		log.Panicln(fmt.Printf("failed to convert pem to sm2 private key: %+v", err))
	}

	// 转化公钥
	publicKey, err := x509.ReadPublicKeyFromPem(publicKeyPem)
	if err != nil {
		log.Panicln(fmt.Printf("failed to convert pem to sm2 public key: %+v", err))
	}

	global.PrivateKey = privateKey
	global.PublicKey = publicKey
}

func WriteNewSessionRecord(newRecord model.SessionRecord) {
	// 若不存在记录文件则先创建再读
	if !FileExist(global.SessionRecordsFilePath) {
		WriteFile(global.SessionRecordsFilePath, nil)
	}

	// 读出原内容
	records := ReadSessionRecords()

	// 添加新数据
	records = append(records, newRecord)
	recordsBytes := JsonMarshal(records)

	// 写入更新后的内容
	WriteFile(global.SessionRecordsFilePath, recordsBytes)
}

func ReadSessionRecords() []model.SessionRecord {
	// 读文件
	recordsBytes := ReadFile(global.SessionRecordsFilePath)

	// 若文件本身为空，则不会反序列化成功，直接返回空记录切片即可
	if len(recordsBytes) == 0 {
		return []model.SessionRecord{}
	}

	// 反序列化
	records := JsonUnmarshal[[]model.SessionRecord](recordsBytes)

	return records
}

func WriteFile(path string, data []byte) {
	err := os.WriteFile(path, data, global.DefaultFilePerm)
	if err != nil {
		log.Panicln(fmt.Printf("failed to write file: %+v", err))
	}
}

func ReadFile(path string) []byte {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Panicln(fmt.Printf("failed to read file: %+v", err))
	}
	return data
}

func JsonMarshal(v any) []byte {
	result, err := json.Marshal(v)
	if err != nil {
		log.Panicln("json marshal error")
	}
	return result
}

func JsonUnmarshal[T any](data []byte) T {
	var result T
	err := json.Unmarshal(data, &result)
	if err != nil {
		log.Panicln("json unmarshal error: ", err)
	}
	return result
}
