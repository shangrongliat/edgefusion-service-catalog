package model

type Catalog struct {
	IP string                    `json:"ip"` // 本机IP
	ID string                    `json:"id"` // 本机ID（节点ID）
	SC map[string]ServiceCatalog `json:"sc"` // 服务目录（只存放“service”类型的服务）
}

type ServiceCatalog struct {
	Name   string `json:"name"`   // 服务名称
	Status string `json:"status"` //服务状态
}
