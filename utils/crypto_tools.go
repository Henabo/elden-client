package utils

import (
	"github.com/hiro942/elden-client/global"
	"github.com/hiro942/elden-client/model/request"
	"github.com/tjfoc/gmsm/sm2"
)

func GetMessageWithSig[T any](message T) request.MessageWithSig {
	messageBytes := JsonMarshal(message)
	messageSig := Sm2Sign(global.PrivateKey, messageBytes)

	messageWithSig := request.MessageWithSig{
		Plain:     messageBytes,
		Signature: messageSig,
	}

	return messageWithSig
}

func GetMessageCipherWithSm2[T any](message T, pubKey *sm2.PublicKey) request.MessageCipher {
	messageBytes := JsonMarshal(message)
	return request.MessageCipher{Cipher: string(Sm2Encrypt(pubKey, messageBytes))}
}

func GetMessageCipherWithSm4[T any](message T, key []byte) request.MessageCipher {
	messageBytes := JsonMarshal(message)
	return request.MessageCipher{Cipher: string(Sm4Encrypt(key, messageBytes))}
}

func GetMessageWithSigCipher[T any](message T, pubKey *sm2.PublicKey) request.MessageCipher {
	messageWithSig := GetMessageWithSig[T](message)
	messageWithSigBytes := JsonMarshal(messageWithSig)
	return request.MessageCipher{
		Cipher: string(Sm2Encrypt(pubKey, messageWithSigBytes)),
	}
}
