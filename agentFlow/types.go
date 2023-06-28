package main

import "github.com/docker/docker/api/types/swarm"

type Memory struct {
	MemTotal     float64 `json:"total"`
	MemFree      float64 `json:"free"`
	MemAvailable float64 `json:"avilable"`
	MemPercent   int     `json:"mempercent"`
}

type Ip struct {
	LocalIp string `json:"hostip"`
}

type DirFiles struct {
	Path  string
	Count int
}
type ServerInfo struct {
	Date string `json:"date"`
	// Node          string       `json:"node"`
	// ServicesLs    string       `json:"servicesls"`
	ServerIp       string          `json:"serverIp"`
	FreeDiskSpace  uint64          `json:"freediskspace"`
	MemoryPercent  int             `json:"memorypercent"`
	SwarmNode      []swarm.Node    `json:"swarmnode"`
	DockerServices []swarm.Service `json:"dockerservices"`
	CountFils      []DirFiles      `json:"countfils"`
}
