package request

// MessageWithSig indicates the messages with signature
type MessageWithSig struct {
	Plain     []byte `json:"plain"`
	Signature []byte `json:"signature"`
}

// MessageCipher indicates the encrypted message
type MessageCipher struct {
	Cipher string `json:"cipher"`
}