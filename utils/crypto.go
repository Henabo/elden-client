package utils

import (
	"encoding/hex"
	"fmt"
	"github.com/tjfoc/gmsm/sm2"
	"github.com/tjfoc/gmsm/sm3"
	"github.com/tjfoc/gmsm/sm4"
	"math/rand"
	"time"
)

// Sm3Hash 使用sm3算法生成消息摘要
func Sm3Hash(src []byte) string {
	hashedBytes := sm3.Sm3Sum(src)
	hashedHex := hex.EncodeToString(hashedBytes)
	return hashedHex
}

// GenerateSm2KeyPair 生成sm2密钥对
func GenerateSm2KeyPair() (*sm2.PrivateKey, *sm2.PublicKey) {
	privateKey, err := sm2.GenerateKey(nil)
	if err != nil {
		panic(fmt.Errorf("failed to generate the Sm2 private key: %w", err))
	}
	publicKey := &privateKey.PublicKey
	return privateKey, publicKey
}

// Sm2Encrypt 使用sm2算法加密明文，得到字节流
func Sm2Encrypt(key *sm2.PublicKey, src []byte)  []byte {
	encrypted, err := sm2.Encrypt(key, src, nil, 0)
	if err != nil {
		panic(fmt.Errorf("failed to encrypt with sm2: %w", err))
	}
	return encrypted
}

// Sm2Decrypt 使用RSA算法解密密文，得到字符串
func Sm2Decrypt(key *sm2.PrivateKey, encrypted []byte) []byte {
	decrypted, err := sm2.Decrypt(key, encrypted, 0)
	if err != nil {
		panic(fmt.Errorf("failed to decrypt with sm2: %w", err))
	}
	return decrypted
}

// Sm2Sign 给消息摘要签名（会先sm3哈希，再对该哈希值签名）
func Sm2Sign(key *sm2.PrivateKey, src []byte) []byte {
	sign, err := key.Sign(nil, src, nil)
	if err != nil {
		panic(fmt.Errorf("failed to sign with sm2: %w", err))
	}
	return sign
}

// Sm2Verify 验证签名
func Sm2Verify(key *sm2.PublicKey, src []byte, sign []byte) bool {
	ok := key.Verify(src, sign)
	return ok
}

// GenerateSm4Key 随机生成16字节的sm4对称密钥
func GenerateSm4Key() []byte {
	// 使用系统时间作为随机种子
	rand.Seed(time.Now().Unix())

	// 从以下62字符随机选择16个
	var letters = []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	key := make([]byte, 16)
	arc := uint8(0)
	if _, err := rand.Read(key); err != nil {
		panic(fmt.Errorf("failed to generate sm4 key: %w", err))
	}

	for i, val := range key {
		arc = val & 61
		key[i] = letters[arc]
	}
	return key
}

// Sm4Encrypt 使用sm4算法对消息进行对称加密
func Sm4Encrypt(key []byte, src []byte) []byte {
	encrypted, err := sm4.Sm4Ecb(key, src, true)
	if err != nil {
		panic(fmt.Errorf("failed to encrypt with sm4: %w", err))
	}
	//fmt.Println(hex.EncodeToString(encrypted))
	return encrypted
}

// Sm4Decrypt 使用sm4算法对密文进行对称解密
func Sm4Decrypt(key []byte, encrypted []byte) []byte {
	decrypted, _ := sm4.Sm4Ecb(key, encrypted, false)
	fmt.Println(string(decrypted))
	return decrypted
}
