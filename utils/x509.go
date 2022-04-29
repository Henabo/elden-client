package utils

import (
	"github.com/hiro942/elden-client/global"
	"github.com/tjfoc/gmsm/sm2"
	"github.com/tjfoc/gmsm/x509"
	"log"
)

func ReadPublicKeyFromHex(publicKeyHex string) *sm2.PublicKey {
	publicKey, err := x509.ReadPublicKeyFromHex(publicKeyHex)
	if err != nil {
		log.Panicf("failed to resolve public key: %+v\n", err)
	}
	return publicKey
}

func ReadPrivateKeyFromPem(privateKeyPem []byte) *sm2.PrivateKey {
	privateKey, err := x509.ReadPrivateKeyFromPem(privateKeyPem, global.PrivateKeyPwd)
	if err != nil {
		log.Panicln("failed to convert pem to sm2 private key:", err)
	}
	return privateKey
}

func ReadPublicKeyFromPem(publicKeyPem []byte) *sm2.PublicKey {
	publicKey, err := x509.ReadPublicKeyFromPem(publicKeyPem)
	if err != nil {
		log.Panicln("failed to convert pem to sm2 public key:", err)
	}
	return publicKey
}

func WritePublicKeyToPem(publicKey *sm2.PublicKey) []byte {
	publicKeyPem, err := x509.WritePublicKeyToPem(publicKey)
	if err != nil {
		log.Panicln("failed to convert public key to pem:", err)
	}
	return publicKeyPem
}

func WritePrivateKeyToPem(privateKey *sm2.PrivateKey) []byte {
	privateKeyPem, err := x509.WritePrivateKeyToPem(privateKey, global.PrivateKeyPwd)
	if err != nil {
		log.Panicln("failed to convert private key to pem:", err)
	}
	return privateKeyPem
}
