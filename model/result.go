package model

type NodeResult struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data Node   `json:"data"`
}

type ServiceResult struct {
	Code int     `json:"code"`
	Msg  string  `json:"msg"`
	Data Service `json:"data"`
}
