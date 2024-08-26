package model

type Service struct {
	Name   string `json:"name"`
	Path   string `json:"path"`
	Pid    int    `json:"pid"`
	Status string `json:"status"`
	Meta   Meta   `json:"meta"`
}

type Meta struct {
	Namespace   string `json:"namespace"`
	Appname     string `json:"appname"`
	Version     string `json:"version"`
	Service     string `json:"service"`
	Instance    string `json:"instance"`
	Mode        string `json:"mode"`
	Type        string `json:"type"`
	ServiceType string `json:"serviceType"`
}
