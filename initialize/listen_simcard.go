package initialize

import "time"

func ListenSIMCardStatus() {
	// 通过系统提供的接口取得IMSI
	for {
		/*
			if 获取IMSI成功(sim卡已插入) && global.SIMCardExist == false{
				global.SIMCardExist = true
				global.MyHashedIMSI = "xxx"
			}
			if 获取IMSI失败 && global.SIMCardExist == true
				global.SIMCardExist = false
				global.MyHashedIMSI = ""
				停止服务
			}
		*/
		time.Sleep(time.Second * 5)
	}
}
