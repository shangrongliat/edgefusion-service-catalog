package main

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"

	"edgefusion-service-catalog/httplink"
)

// TODO ./consul agent -data-dir=/home/work/cata/ds -bind=172.16.100.81 -client=0.0.0.0 -server -ui -bootstrap-expect=1
func main() {
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
