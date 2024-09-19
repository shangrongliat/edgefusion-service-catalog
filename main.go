package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"sync"
	"time"

	"edgefusion-service-catalog/httplink"
)

// TODO ./consul agent -data-dir=/home/work/cata/ds -bind=172.16.100.81 -client=0.0.0.0 -server -ui -bootstrap-expect=1
func main2() {
	//initLog(false)
	// 设置 log 包的日志输出
	group := sync.WaitGroup{}
	group.Add(1)
	defer group.Done()
	// 加载配置文件
	//yamlFile, err := ioutil.ReadFile("./config.yml")
	//if err != nil {
	//	log.Fatalf("Error reading YAML file: %v", err)
	//}
	//// 解析 YAML 文件
	//var cfg config.Config
	//if err = yaml.Unmarshal(yamlFile, &cfg); err != nil {
	//	log.Fatalf("Error unmarshalling YAML data: %v", err)
	//}

	//go broadcast.NewNotice()
	//go unicast.NewReceive()
	go httplink.Subscribe(nil)
	group.Wait()
}

func initLog(terminal bool) {
	// 构建日志文件的完整路径
	logFilePath := filepath.Join("/etc/edgefusion/video/push/", "logs", "app.log")
	// 创建文件夹 "logs" 如果它不存在
	err := os.MkdirAll(filepath.Dir(logFilePath), 0755)
	if err != nil {
		log.Fatalf("Error creating logs folder: %v", err)
	}
	// 打开一个文件用于写入日志
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	// 设置 log 包的日志输出
	log.SetOutput(logFile)
	if terminal {
		// 创建一个 io.MultiWriter 实例，它允许我们将日志输出到多个地方
		multiWriter := io.MultiWriter(os.Stdout, logFile)
		// 设置 log 包的日志输出
		log.SetOutput(multiWriter)
	}
}

func main() {
	// Consul DNS 服务器地址和端口
	consulDNSServer := "127.0.0.1:8600"

	// 要解析的域名
	domain := "montage01.service.consul"

	// 创建一个自定义的 DNS 解析器
	resolver := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{
				Timeout:   10 * time.Second,
				KeepAlive: 30 * time.Second,
			}
			return d.Dial("udp", consulDNSServer)
		},
	}

	// 解析域名
	addrs, err := resolver.LookupHost(context.Background(), domain)
	if err != nil {
		fmt.Printf("Error resolving %s: %v\n", domain, err)
		return
	}

	// 打印解析结果
	for _, addr := range addrs {
		fmt.Println(addr)
	}
}
