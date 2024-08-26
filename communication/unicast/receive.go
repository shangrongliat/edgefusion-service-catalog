package unicast

import (
	"edgefusion-service-catalog/cache"
	"edgefusion-service-catalog/communication/listener"
	"edgefusion-service-catalog/model"
	"encoding/json"
	"fmt"
	"log"
	"net"
)

func NewReceive(listener *listener.Listener, cache *cache.Cache) {
	// 创建UDP地址
	addr, err := net.ResolveUDPAddr("udp4", ":9999")
	if err != nil {
		log.Printf("Failed to resolve UDP address: %v \n", err)
	}
	// 创建UDP连接
	conn, err := net.ListenUDP("udp4", addr)
	if err != nil {
		log.Printf("Failed to listen on UDP: %v \n", err)
	}
	defer conn.Close()

	buf := make([]byte, 1024)
	for {
		// 接收数据
		n, srcAddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			log.Printf("Failed to read from UDP: %v \n", err)
		}
		// 收到各个节点发送的同步信息，
		if broadcast := dataUnmarshal(buf[:n]); broadcast != nil {
			// 1. 判断该节点是否在内存中，如果不存在则发起询问信息，将获取到的信息写入内存中
			// 2. 如果内存中存在对应节点，则判断版本是否一致，如果不一致则发起询问获取新的节点信息
			if _, exists := cache.GetECache(broadcast); exists {
				continue
			}
			data := []byte{0}
			data = append(data, buf[:n]...)
			// 发起询问
			listener.Transmit(data, fmt.Sprintf("%s:64505", srcAddr.IP))
		}
	}
}

func dataUnmarshal(data []byte) (broadcast *model.Broadcast) {
	if err := json.Unmarshal(data, &broadcast); err != nil {
		log.Printf("广播数据解析失败,data: [%s].%v \n", string(data), err)
	}
	return
}
