package httplink

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"edgefusion-service-catalog/model"
)

func GetNode() (data model.NodeResult) {
	// 发送HTTP GET请求
	resp, err := http.Get("http://127.0.0.1:19300/ef/engine/node/np")
	if err != nil {
		log.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()
	// 检查HTTP响应状态码
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Unexpected status code: %d", resp.StatusCode)
	}
	// 读取响应体
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response body: %v", err)
	}
	// 解析JSON数据
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Fatalf("Failed to parse JSON: %v", err)
	}
	return
}

func GetService() (data model.ServiceResult) {
	// 发送HTTP GET请求
	resp, err := http.Get("http://127.0.0.1:19300/ef/engine/app/instances")
	if err != nil {
		log.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()
	// 检查HTTP响应状态码
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Unexpected status code: %d", resp.StatusCode)
	}
	// 读取响应体
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response body: %v", err)
	}
	// 解析JSON数据
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Fatalf("Failed to parse JSON: %v", err)
	}
	return
}
