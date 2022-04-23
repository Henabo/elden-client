package utils

import (
	"encoding/json"
	"fmt"
	"github.com/hiro942/elden-client/global"
	"github.com/hiro942/elden-client/model"
	"github.com/tjfoc/gmsm/x509"
	"log"
	"os"
)

func FileExist(path string) bool {
	_, err := os.Lstat(path)
	return os.IsExist(err)
}

func ReadKeyPair() {
	// 读pem格式私钥和公钥
	privateKeyPem := ReadFile(global.PrivateKeyPath)
	publicKeyPem := ReadFile(global.PublicKeyPath)

	// 转化私钥
	privateKey, err := x509.ReadPrivateKeyFromPem(privateKeyPem, global.PrivateKeyPwd)
	if err != nil {
		log.Panic(fmt.Printf("failed to convert pem to sm2 private key: %+v", err))
	}

	// 转化公钥
	publicKey, err := x509.ReadPublicKeyFromPem(publicKeyPem)
	if err != nil {
		log.Panic(fmt.Printf("failed to convert pem to sm2 public key: %+v", err))
	}

	global.PrivateKey = privateKey
	global.PublicKey = publicKey
}

func WriteNewSessionRecord(newRecord model.SessionRecord) {
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

	// 反序列化
	records := JsonUnmarshal[[]model.SessionRecord](recordsBytes)

	return records
}

func WriteFile(path string, data []byte) {
	err := os.WriteFile(path, data, global.DefaultFilePerm)
	if err != nil {
		log.Panic(fmt.Printf("failed to write file: %+v", err))
	}
}

func ReadFile(path string) []byte {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Panic(fmt.Printf("failed to read file: %+v", err))
	}
	return data
}

func JsonMarshal(v any) []byte {
	result, err := json.Marshal(v)
	if err != nil {
		panic("json marshal error")
	}
	return result
}

func JsonUnmarshal[T any](data []byte) T {
	var result T
	err := json.Unmarshal(data, &result)
	if err != nil {
		panic("json unmarshal error")
	}
	return result
}
