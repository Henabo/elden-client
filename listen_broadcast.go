package main

import (
	"fmt"
	"log"
	"net"
	"time"
)

// ListenBroadcast 接收卫星广播UDP报文
func ListenBroadcast() {
	packageConn, err := net.ListenPacket("udp4", ":8829")
	if err != nil {
		log.Panicln("listen broadcast failed, error:", err)
	}
	defer packageConn.Close()

	for {
		buf := make([]byte, 1024)
		n, addr, err := packageConn.ReadFrom(buf)
		if err != nil {
			log.Panicln(err)
		}
		fmt.Printf("%s sent this: %s\n", addr, buf[:n])

		time.Sleep(time.Second * 3)
	}
}
