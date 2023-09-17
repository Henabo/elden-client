package main

import (
	"fmt"
	"github.com/eiannone/keyboard"
	"github.com/hiro942/elden-client/model/enums"
	"log"
	"os"
	"time"
)

func (auth *AuthenticationService) Watcher() {
	session := auth.Session
	client := session.Client

	for {
		char, _, err := keyboard.GetSingleKey()
		if err != nil {
			log.Panicf("watcher error：%+v", err)
		}
		fmt.Printf("You pressed: %q\r\n", char)

		switch string(char) {
		case "w": // 插入sim卡，自动开始认证
			if client.Status != enums.ClientStatusWithoutSIM {
				log.Println("检测到SIM卡已插入，请勿重复操作")
				break
			}
			log.Println("【事件】SIM卡插入成功！")

			// 状态流转
			if client.Status != enums.ClientStatusWithoutSIM {
				break
			}
			client.Status = enums.ClientStatusWaitVerify

			// 发起认证
			if err = auth.LaunchAuthentication("s0"); err != nil {
				log.Println(err.Error())
			}

		case "s": // 拔出sim卡
			if client.Status == enums.ClientStatusWithoutSIM {
				log.Println("检测到SIM卡已拔出，请勿重复操作")
				break
			}
			log.Println("【事件】SIM卡移出成功！")

			// 发断开连接消息
			if err := session.Disconnect(false); err != nil {
				log.Printf("发送断开消息失败: %+v\n", err)
			}
			// TODO 可以加一个心跳机制，断开消息失败也会能正常让卫星更新账本

			// 状态流转
			client.Status = enums.ClientStatusWithoutSIM
			session.Status = enums.SessionStatusNull

		case "x": // 手动开启认证
			log.Println("【事件】手动发起认证")
			if client.Status == enums.ClientStatusWithoutSIM {
				log.Println("未检测到SIM卡")
				break
			}
			if client.Status == enums.ClientStatusVerifySuccess {
				log.Println("已经认证成功，请勿重复操作")
				break
			}
			if err = auth.LaunchAuthentication("s0"); err != nil {
				log.Println(err.Error())
			}

		case "y": // 手动结束当前会话
			log.Println("【事件】手动断开当前会话。")
			if session.Status == enums.SessionStatusNull {
				log.Println("当前暂无会话")
				break
			}
			if err = session.Disconnect(false); err != nil {
				log.Println(err.Error())
			}

		case "q":
			os.Exit(1)

		}

		time.Sleep(time.Second)
	}
}
