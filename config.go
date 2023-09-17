package main

import (
	"github.com/spf13/viper"
	"log"
)

var config = new(Config)

type Config struct {
	ClientID          string // 客户端 id
	FabricAppHostPath string // fabric-app 地址
	HttpServerPort    string // http端口
	PrivateKeyPwd     string // 私钥加密密钥
}

func LoadConfig() {
	// 指定配置文件路径
	viper.SetConfigName("app")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")

	// 读取配置
	if err := viper.ReadInConfig(); err != nil {
		log.Panicln(err)
	}

	// 配置绑定
	if err := viper.Unmarshal(config); err != nil {
		log.Panicln(err)
	}
}
