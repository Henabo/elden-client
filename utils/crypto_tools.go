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
	return request.MessageCipher{Cipher: Sm2Encrypt(pubKey, messageBytes)}
}

func GetMessageCipherWithSm4[T any](message T, sm4Key []byte) request.MessageCipher {
	messageBytes := JsonMarshal(message)
	return request.MessageCipher{Cipher: Sm4Encrypt(sm4Key, messageBytes)}
}

func GetMessageWithSigCipher[T any](message T, pubKey *sm2.PublicKey) request.MessageCipher {
	messageWithSig := GetMessageWithSig[T](message)
	messageWithSigBytes := JsonMarshal(messageWithSig)
	return request.MessageCipher{
		Cipher: Sm2Encrypt(pubKey, messageWithSigBytes),
	}
}
