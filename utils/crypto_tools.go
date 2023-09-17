package utils

import (
	"github.com/hiro942/elden-client/model/request"
	"github.com/tjfoc/gmsm/sm2"
)

func GetMessageWithSig[T any](message T, privateKey *sm2.PrivateKey) request.MessageWithSig {
	messageBytes := JsonMarshal(message)
	messageSig := Sm2Sign(privateKey, messageBytes)

	messageWithSig := request.MessageWithSig{
		Plain:     messageBytes,
		Signature: messageSig,
	}

	return messageWithSig
}

func GetMessageCipherWithSm2[T any](message T, pubKey *sm2.PublicKey) request.MessageCipher {
	messageBytes := JsonMarshal(message)
	return request.MessageCipher{Cipher: Sm2Encrypt(pubKey, messageBytes)}
}

func GetMessageCipherWithSm4[T any](message T, sm4Key []byte) request.MessageCipher {
	messageBytes := JsonMarshal(message)
	return request.MessageCipher{Cipher: Sm4Encrypt(sm4Key, messageBytes)}
}

func GetMessageWithSigCipher[T any](message T, publicKey *sm2.PublicKey, privateKey *sm2.PrivateKey) request.MessageCipher {
	messageWithSig := GetMessageWithSig[T](message, privateKey)
	messageWithSigBytes := JsonMarshal(messageWithSig)
	return request.MessageCipher{
		Cipher: Sm2Encrypt(publicKey, messageWithSigBytes),
	}
}
