package broadcast

import (
	"log"
	"net"
	"time"
)

func NewNotice() {
	// 创建UDP地址，使用广播地址和端口
	addr, err := net.ResolveUDPAddr("udp4", "255.255.255.255:9999")
	if err != nil {
		log.Fatalf("Failed to resolve UDP address: %v", err)
	}

	// 创建UDP连接
	conn, err := net.DialUDP("udp4", nil, addr)
	if err != nil {
		log.Fatalf("Failed to dial UDP: %v", err)
	}
	defer conn.Close()

	for {
		// 发送数据
		message := []byte("Hello, UDP broadcast!")
		_, err := conn.Write(message)
		if err != nil {
			log.Fatalf("Failed to send data: %v", err)
		}
		log.Println("Broadcast message sent")

		// 每隔1秒发送一次
		time.Sleep(1 * time.Second)
	}
}
