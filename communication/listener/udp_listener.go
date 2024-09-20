package listener

import (
	"log"
	"net"

	"edgefusion-service-catalog/cache"
)

type Listener struct {
	conn net.PacketConn
}

// 监听本地端口
func NewLister() *Listener {
	conn, err := net.ListenPacket("udp", "0.0.0.0:64505")
	if err != nil {
		log.Printf("Failed to bind to address . %v \r\n", err)
	}
	log.Printf("udp lister start \r\n")
	return &Listener{conn: conn}
}

func (l *Listener) Lister(cache *cache.Cache) {
	defer func(conn net.PacketConn) {
		if err := conn.Close(); err != nil {
			log.Printf("Failed to close from UDP: %v \r\n", err)
		}
	}(l.conn)
	buf := make([]byte, 4050)
	for {
		n, addr, err := l.conn.ReadFrom(buf)
		if err != nil {
			log.Printf("Failed to read from UDP: %v \r\n", err)
			continue
		}
		// 跳过本地请求
		if addr.String() == cache.GetLocalCache().IP {
			continue
		}
		//TODO 这里收到的消息会有2种：
		dataType := buf[0] // 获取数据类型
		switch dataType {
		case 0:
			// 0 是本机发出询问后收到的反馈
			// 收到其他节点发送的缓存信息
			cache.AddECacheBinary(buf[1:n]) // 将收到的message添加到缓存中，由添加方法进行解析处理
		case 1:
			// 2是由其他设备发过来的询问请求
			catalog := cache.GetCacheBinary(buf[1:n])
			if catalog == nil {
				continue
			}
			data := []byte{0}
			data = append(data, catalog...)
			l.Transmit(data, addr.String())
		default:
			// 未知类型不处理
		}
	}
}

func (t *Listener) Transmit(data []byte, remoteAddr string) {
	remote, err := net.ResolveUDPAddr("udp", remoteAddr)
	if err != nil {
		log.Printf("节点同步信息询问远程地址[ %s ]初始化失败.%v \r\n", remoteAddr, err)
	}
	if _, err := t.conn.WriteTo(data, remote); err != nil {
		log.Printf("Error sending UDP packet: %v \r\n", err)
		return
	}
}
