package main

import (
	"bufio"
	"bytes"
	"context"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/shirou/gopsutil/disk"
)

const ShellToUse = "bash"

func Shellout(command string) (string, string, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command(ShellToUse, "-c", command)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return stdout.String(), stderr.String(), err
}
func toInt(raw string) int {
	if raw == "" {
		return 0
	}
	res, err := strconv.Atoi(raw)
	if err != nil {
		panic(err)
	}
	return res
}

func parseLine(raw string) (key string, value int) {
	// fmt.Println(raw)
	text := strings.ReplaceAll(raw[:len(raw)-2], " ", "")
	keyValue := strings.Split(text, ":")
	return keyValue[0], toInt(keyValue[1])
}

func ReadMemoryStats() Memory {
	file, err := os.Open("/proc/meminfo")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	bufio.NewScanner(file)
	scanner := bufio.NewScanner(file)
	res := Memory{}
	for scanner.Scan() {
		key, value := parseLine(scanner.Text())
		switch key {
		case "MemTotal":
			res.MemTotal = float64(value)
		case "MemFree":
			res.MemFree = float64(value)
		case "MemAvailable":
			res.MemAvailable = float64(value)
		}
	}
	return res
}
func GetDiskServices(path string) disk.UsageStat {
	diskInfo, _ := disk.Usage(path)
	return *diskInfo
}
func GetLocalIP() Ip {
	addrs, err := net.InterfaceAddrs()
	var ip Ip
	if err != nil {
		ip.LocalIp = "error"
		return ip
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {

				ip.LocalIp = ipnet.IP.String()
				return ip
			}
		}
	}
	ip.LocalIp = "error"
	return ip
}

func countFileInDir(path string) int {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Print(err)
		return -100
	}

	return len(files)
}

//-----------------------------------

func GetInfoForPc() ServerInfo {

	se := ServerInfo{}

	res := GetDiskServices("/")
	res.Free = res.Free / 1000000000

	memModel := ReadMemoryStats()
	memModel.MemAvailable = memModel.MemAvailable / 1000000
	memModel.MemFree = memModel.MemFree / 1000000
	memModel.MemTotal = memModel.MemTotal / 1000000
	memModel.MemPercent = int((100 - (memModel.MemAvailable/memModel.MemTotal)*100))
	t := time.Now()

	se.Date = string(t.Format("2006-01-02 15:04:05"))
	se.ServerIp = GetLocalIP().LocalIp
	se.MemoryPercent = memModel.MemPercent
	se.FreeDiskSpace = res.Free

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		log.Print(err)
	}

	se.SwarmNode, err = cli.NodeList(context.Background(), types.NodeListOptions{})
	if err != nil {
		log.Print(err)
	}
	se.DockerServices, err = cli.ServiceList(context.Background(), types.ServiceListOptions{})
	if err != nil {
		log.Print(err)
	}

	file, err := os.Open("path.txt")
	if err != nil {
		log.Print(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var filesCount []DirFiles
	for scanner.Scan() {
		fs := DirFiles{}
		fs.Path = scanner.Text()
		fs.Count = countFileInDir(scanner.Text())
		filesCount = append(filesCount, fs)

	}

	if err := scanner.Err(); err != nil {
		log.Print(err)
	}

	se.CountFils = filesCount
	return se
}
