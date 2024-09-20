package httplink

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"edgefusion-service-catalog/model"
	"edgefusion-service-catalog/model/consul"
)

func GetNode() (data model.NodeResult) {
	// 发送HTTP GET请求
	resp, err := http.Get("http://127.0.0.1:19300/ef/engine/node/np")
	if err != nil {
		log.Printf("Failed to send request: %v \n", err)
	}
	defer resp.Body.Close()
	// 检查HTTP响应状态码
	if resp.StatusCode != http.StatusOK {
		log.Printf("Unexpected status code: %d \n", resp.StatusCode)
	}
	// 读取响应体
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read response body: %v", err)
	}
	// 解析JSON数据
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Printf("Failed to parse JSON: %v \n", err)
	}
	return
}

func GetService() (data model.ServiceResult) {
	// 发送HTTP GET请求
	resp, err := http.Get("http://127.0.0.1:19300/ef/engine/app/instances")
	if err != nil {
		log.Printf("Failed to send request: %v \n", err)
	}
	defer resp.Body.Close()
	// 检查HTTP响应状态码
	if resp.StatusCode != http.StatusOK {
		log.Printf("Unexpected status code: %d \n", resp.StatusCode)
	}
	// 读取响应体
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read response body: %v \n", err)
	}
	// 解析JSON数据
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Printf("Failed to parse JSON: %v \n", err)
	}
	return
}

func Put(msg *consul.Register) {
	log.Println("服务注册", msg)
	body, err := json.Marshal(msg)
	if err != nil {
		log.Println("目录服务注册失败", msg)
	}
	// 设置请求的数据
	reader := bytes.NewReader(body)
	// 创建一个PUT请求
	req, err := http.NewRequest("PUT", "http://127.0.0.1:8500/v1/agent/service/register", reader)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}
	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	defer resp.Body.Close()
	// 读取响应体
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	} // 打印响应状态码和响应体
	fmt.Printf("Response Status: %s \n", resp.Status)
	fmt.Printf("Response Body: %s \n", responseBody)
}
