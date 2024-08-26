package broadcast

import (
	"edgefusion-service-catalog/cache"
	"encoding/json"
	"log"
	"net"
	"time"
)

func NewNotice(cache *cache.Cache) {
	// 创建UDP地址，使用广播地址和端口
	addr, err := net.ResolveUDPAddr("udp4", "255.255.255.255:9999")
	if err != nil {
		log.Printf("Failed to resolve UDP address: %v \n", err)
	}
	// 创建UDP连接
	conn, err := net.DialUDP("udp4", nil, addr)
	if err != nil {
		log.Printf("Failed to dial UDP: %v \n", err)
	}
	defer conn.Close()
	for {
		message, err := json.Marshal(cache.GetBroadcastInfo())
		if err != nil {
			log.Printf("广播信息Json化失败.%v \n", err)
		}
		// 发送数据
		_, err = conn.Write(message)
		if err != nil {
			log.Printf("Failed to send data: %v \n", err)
		}
		// 每隔1秒发送一次
		time.Sleep(10 * time.Second)
	}
}
