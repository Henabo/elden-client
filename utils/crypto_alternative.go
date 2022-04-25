package utils

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
)

// Sha256Hash 获取数据的16进制SHA256哈希值
func Sha256Hash(data string) string {
	hashedBytes := sha256.Sum256([]byte(data))
	hashedHex := hex.EncodeToString(hashedBytes[:])
	return hashedHex
}

// GenerateRSAKeyPair 生成RSA密钥对
func GenerateRSAKeyPair() (*rsa.PrivateKey, *rsa.PublicKey) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)  //1024bit
	if err != nil {
		log.Panicln(fmt.Errorf("failed to generate the RSA Key: %w", err))
	}
	publicKey := &privateKey.PublicKey
	return privateKey, publicKey
}

// RSAEncrypt 使用RSA算法加密明文，得到字节流
func RSAEncrypt(src []byte, key *rsa.PublicKey) []byte {
	encrypted, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, key, src, nil)
	if err != nil {
		log.Panicln(fmt.Errorf("failed to encrypt the message: %w", err))
	}
	return encrypted
}

// RSADecrypt 使用RSA算法解密密文，得到字符串
func RSADecrypt(src []byte, key *rsa.PrivateKey) []byte {
	decrypted, err := key.Decrypt(nil, src, &rsa.OAEPOptions{Hash: crypto.SHA256})
	if err != nil {
		log.Panicln(fmt.Errorf("failed to decrypt the cipher: %w", err))
	}
	return decrypted
}

// RSASign 给消息摘要签名
func RSASign(key *rsa.PrivateKey, hashed []byte) []byte {
	sig, err := rsa.SignPKCS1v15(nil, key, crypto.SHA256, hashed)
	if err != nil {
		log.Panicln(fmt.Errorf("failed to sign: %w", err))
	}
	return sig
}

// RSAVerify 验证签名
func RSAVerify(key *rsa.PublicKey, hashed []byte, sig []byte) {
	if err := rsa.VerifyPKCS1v15(key, crypto.SHA256, hashed[:], sig); err != nil {
		log.Panicln(fmt.Errorf("failed to verify the signature: %w", err))
	}
}