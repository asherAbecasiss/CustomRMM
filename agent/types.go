package main


type Memory struct {
	MemTotal     float64 `json:"total"`
	MemFree      float64 `json:"free"`
	MemAvailable float64 `json:"avilable"`
	MemPercent   int     `json:"mempercent"`
}


type Ip struct {
	LocalIp string `json:"hostip"`
}