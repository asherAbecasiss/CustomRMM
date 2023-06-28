package main

import (
	"github.com/docker/docker/api/types/swarm"
)

func CheckForNode(nodes []swarm.Node) int {

	for _, v := range nodes {
		if v.Status.State != "ready" {
			return 1
		}

	}
	return 0
}

func CheckForServices(services []swarm.Service) int {
	return 0
}

func CheckForNumberOfFilesInDir(files []DirFiles) int {
	for _, v := range files {
		if v.Count == -100 {
			//: no such file or directory
			return 3
		} else {
			if (v.Count - 3) > 10 {
				return 3
			}

		}

	}
	return 0
}
func CheckForDisk(disk uint64) int {

	if disk < 50 {
		return 4
	}
	return 0
}

func CheckForMemPercent(memorypercent int) int {

	if memorypercent > 90 {
		return 5
	}
	return 0
}
