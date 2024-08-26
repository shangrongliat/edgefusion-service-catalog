package httplink

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"edgefusion-service-catalog/model"
	"github.com/pebbe/zmq4"
)

func Subscribe() {
	// 创建一个 ZMQ SUB 套接字
	subscriber, err := zmq4.NewSocket(zmq4.SUB)
	if err != nil {
		log.Fatalf("Failed to create ZMQ socket: %v", err)
	}
	defer subscriber.Close()
	// 连接到发布者（Publisher）
	err = subscriber.Connect("tcp://127.0.0.1:19400")
	if err != nil {
		log.Fatalf("Failed to connect to publisher: %v", err)
	}

	// 设置订阅的过滤器，可以设置为空字符串 "" 来接收所有消息
	err = subscriber.SetSubscribe("")
	if err != nil {
		log.Fatalf("Failed to set subscribe filter: %v", err)
	}
	fmt.Println("Waiting for messages...")
	for {
		// 接收消息
		message, err := subscriber.Recv(0)
		if err != nil {
			log.Fatalf("Failed to receive message: %v", err)
		}
		topic, data := parseZmqData(message)
		if topic == "node/np" {
			var node model.Node
			if err := json.Unmarshal([]byte(data), &node); err != nil {
				log.Fatalf("Failed to unmarshal node: %v", err)
			}
		} else if topic == "app/instances" {
			var service []model.Service
			if err := json.Unmarshal([]byte(data), &service); err != nil {
				log.Fatalf("Failed to unmarshal service: %v", err)
			}
		}
	}
}

func parseZmqData(data string) (string, string) {
	i := strings.Index(data, "@")
	return data[:i], data[i+1:]
}

// 解析zmq消息，使用json将content解析为泛型指定的类型
// 泛型不要传入指针
func parseZmqMsg[T any](msg string) (string, T, error) {
	topic, content, err := parseZmgTopicAndContent(msg)
	if err != nil {
		return topic, *new(T), err
	}
	var t T
	if err = json.Unmarshal([]byte(content), &t); err != nil {
		return topic, *new(T), err
	}
	return topic, t, nil
}

func parseZmgTopicAndContent(msg string) (string, string, error) {
	i := strings.Index(msg, "@")
	//必须包含有效的分隔符@，且topic不能为""
	if i == -1 || i == 0 {
		return "", "", fmt.Errorf("error to parse zmg msg, invalid msg: %s", msg)

	}
	return msg[:i], msg[i+1:], nil
}
