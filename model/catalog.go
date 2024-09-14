package model

// 目录结构
type Catalog struct {
	IP      string                     `json:"ip"`      // 本机IP
	ID      string                     `json:"id"`      // 本机ID（节点ID）
	Version string                     `json:"version"` // 数据版本(当该结构为本地缓存时，字段为空)
	SC      map[string]*ServiceCatalog `json:"sc"`      // 服务目录（只存放“service”类型的服务） key 服务ID value 服务信息
}

// 目录属性
type ServiceCatalog struct {
	ID             string `json:"id"`              // 服务ID
	Name           string `json:"name"`            // 服务名称
	Status         string `json:"status"`          // 服务状态
	CheckInterface string `json:"check_interface"` // 检测接口
	CheckInterval  string `json:"check_interval"`  // 检测间隔
	Port           string `json:"port"`            // 服务端口
}

// 本地持久化对象
type CacheData struct {
	LocalCache *Catalog
	Ecache     map[string]*Catalog
	ParentId   string
	Version    string
}

// 广播信息
type Broadcast struct {
	ID      string `json:"id"`      // 节点ID
	Version string `json:"version"` // 信息版本
}
