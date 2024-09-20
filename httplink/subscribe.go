package httplink

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"edgefusion-service-catalog/cache"
	"edgefusion-service-catalog/model"
	"github.com/go-zeromq/zmq4"
)

func Subscribe(cache *cache.Cache) {
	log.Println("ZMQ Subscribing start")
	// 创建一个 ZMQ SUB 套接字
	sub := zmq4.NewSub(context.Background())
	//连接socket
	if err := sub.Dial("tcp://127.0.0.1:19400"); err != nil {
		return
	}
	// 订阅所有主题 (空字符串表示订阅所有消息)
	if err := sub.SetOption(zmq4.OptionSubscribe, "node/np"); err != nil {
		return
	}
	if err := sub.SetOption(zmq4.OptionSubscribe, "app/instances"); err != nil {
		return
	}
	for {
		message, err := sub.Recv()
		if err != nil {
			log.Println("error receive zmq msg", "err", err)
			continue
		}
		topic, data := parseZmqData(string(message.Bytes()))
		if topic == "node/np" {
			var node model.Node
			if err := json.Unmarshal([]byte(data), &node); err != nil {
				log.Printf("Failed to unmarshal node: %v \n", err)
			}
			// 根据上报的信息，判断是否要更新本地节点信息？
			// 1. 判断父节点是否为原父节点：
			//     1) 是， 判断父节点下边计算节点是否有变化
			//        a. 有变化， 更新本地缓存信息
			//        b. 没有变化， 跳过处理
			//     2) 否， 当父节点发生变化时将原父节点下的边计算节点权重进行下降设置（设置为10），
			//             将新的父节点下的边计算节点权重就行上调设置（设置为50）
			// 平台协同关系不影响缓存版本
			cache.AddNodeCache(&node)
		} else if topic == "app/instances" {
			var service []model.Service
			if err := json.Unmarshal([]byte(data), &service); err != nil {
				log.Printf("Failed to unmarshal service: %v \n", err)
			}
			// 根据上报信息，判断是否要更新本地服务状态
			// 服务状态分为 活跃/不活跃 两个状态
			// 通过服务名称和服务状态判断服务是否有变化
			// 1. 如果服务有变化，则更新服务状态的同时修改节点 version字段（版本变化）
			cache.AddServiceCache(service)
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
