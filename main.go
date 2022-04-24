package main

import "github.com/hiro942/elden-client/initialize"

func main() {
	//port := *flag.String("p", "20000", "application running port")
	//flag.Parse()
	//fmt.Println("running port:", port)

	initialize.MockInit() // mock模式
	initialize.SysInit()  // 系统初始化

}
