package enums

type (
	// ClientStatus 客户端状态
	ClientStatus int8

	// SessionStatus 会话状态
	SessionStatus int8

	// AccessType 接入类型
	AccessType int8

	// AccessKeyMode 快速接入模式
	AccessKeyMode int8

	// Cipher 密文
	Cipher []byte
)

const (
	ClientStatusWithoutSIM    ClientStatus = 0 // 未开始（未插sim卡）
	ClientStatusWaitVerify    ClientStatus = 1 // 未认证
	ClientStatusVerifying     ClientStatus = 2 // 认证中
	ClientStatusVerifySuccess ClientStatus = 3 // 认证成功
	ClientStatusVerifyFailed  ClientStatus = 4 // 认证失败

	SessionStatusNull                     = 0 // 无会话
	SessionStatusProcessing SessionStatus = 1 // 进行中

	AccessTypeStrict   AccessType = 1 // 首次认证
	AccessTypeNormal   AccessType = 2 // 快速认证
	AccessTypeHandover AccessType = 3 // 交接认证

	AccessKeyModeHashed    AccessKeyMode = 1 // 发送哈希
	AccessKeyModeEncrypted AccessKeyMode = 2 // 发送新会话密钥
)

func (accessType AccessType) Format() string {
	switch accessType {
	case AccessTypeStrict:
		return "strict"
	case AccessTypeNormal:
		return "normal"
	case AccessTypeHandover:
		return "handover"
	}
	return ""
}

func (mode AccessKeyMode) Format() string {
	switch mode {
	case AccessKeyModeHashed:
		return "hashed"
	case AccessKeyModeEncrypted:
		return "encrypted"
	}
	return ""
}
