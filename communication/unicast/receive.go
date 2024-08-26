package unicast

import (
	"log"
	"net"
)

func NewReceive() {
	// 创建UDP地址
	addr, err := net.ResolveUDPAddr("udp4", ":9999")
	if err != nil {
		log.Fatalf("Failed to resolve UDP address: %v", err)
	}

	// 创建UDP连接
	conn, err := net.ListenUDP("udp4", addr)
	if err != nil {
		log.Fatalf("Failed to listen on UDP: %v", err)
	}
	defer conn.Close()

	buffer := make([]byte, 1024)
	for {
		// 接收数据
		n, srcAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			log.Fatalf("Failed to read from UDP: %v", err)
		}
		log.Printf("Received message from %s: %s\n", srcAddr.IP, string(buffer[:n]))
	}
}
