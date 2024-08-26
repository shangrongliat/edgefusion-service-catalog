package model

type Node struct {
	Node      Info   `json:"node"`
	Parent    Info   `json:"parent"`
	ENodeList []Info `json:"eNodeList"`
}

type Info struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	ParentID string `json:"parentId"`
	CoopID   string `json:"coopId"`
	IP       string `json:"ip"`
	State    string `json:"state"`
	Online   int    `json:"online"`
}

type NodeCache struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	State   string `json:"state"`
	Passing int    `json:"passing"`
	Warning int    `json:"warning"`
}
