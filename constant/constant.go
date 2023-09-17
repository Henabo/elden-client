package constant

const (
	// CryptoPath 加密材料存储路径
	CryptoPath             = "./.crypto/"
	PrivateKeyPemFileName  = "id_sm2"
	PublicKeyPemFileName   = "id_sm2.pub"
	SessionRecordsFileName = "session_records.json"

	// DefaultFilePerm 文件权限
	DefaultFilePerm = 0777

	// DefaultSessionKeyAge 会话密钥默认寿命
	DefaultSessionKeyAge = 3600 * 24
)
