package main

func main() {
	// 读取配置
	LoadConfig()

	// 实例初始化
	cache := NewCache()
	url := NewURLService(cache)
	ledger := NewLedger(url)
	client := NewClient(config.ClientID, ledger, cache)
	session := NewSession(client)
	authService := NewAuthenticationService(session)

	// 监听SIM卡的插拔
	go authService.Watcher()

	// 启动 router
	authService.RunRouter()
}
