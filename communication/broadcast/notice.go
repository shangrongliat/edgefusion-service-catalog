package broadcast

import (
	"encoding/json"
	"log"
	"net"
	"time"

	"edgefusion-service-catalog/cache"
)

// NewNotice 广播通知
func NewNotice(cache *cache.Cache) {
	// 创建UDP地址，使用广播地址和端口
	addr, err := net.ResolveUDPAddr("udp4", "255.255.255.255:9999")
	if err != nil {
		log.Println("Failed to resolve UDP address", err)
	}
	// 创建UDP连接
	conn, err := net.DialUDP("udp4", nil, addr)
	if err != nil {
		log.Println("Failed to dial UDP", err)
	}
	defer conn.Close()
	for {
		br := cache.GetBroadcastInfo()
		if br.ID != "" && br.Version != "" {
			message, err := json.Marshal(br)
			if err != nil {
				log.Println("广播信息Json化失败.", err)
			}
			// 发送数据
			_, err = conn.Write(message)
			if err != nil {
				log.Println("Failed to send data.", err)
			}
		}
		// 每隔10秒发送一次
		time.Sleep(10 * time.Second)
	}
}
