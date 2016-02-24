package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os/exec"
	"strconv"
	"strings"
)

type CPUInfo struct {
	Sockets        int32   `json:"sockets"`
	CoresPerSocket int32   `json:"cores_per_socket"`
	ThreadsPerCore int32   `json:"threads_per_core"`
	TotalCores     int32   `json:"total_cores"`
	ClockSpeed     float32 `json:"clockspeed"`
}

type NetworkInterface struct {
	Name       string   `json:"name"`
	MacAddress string   `json:"macaddress"`
	Addrs      []string `json:"ipaddress"`
}

type Master struct {
	CPUInfo
	NetworkInterface
}

func convertStringToInteger(s string) int32 {
	i, err := strconv.Atoi(s)
	checkErr(err)
	return int32(i)
}

func convertStringToFloat(s string) float32 {
	i, err := strconv.ParseFloat(s, 32)
	checkErr(err)
	return float32(i)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func NetworkInterfacesToMap() (networkInterfaceMap map[string]NetworkInterface) {
	networkInterfaceMap = make(map[string]NetworkInterface)

	networkInterfaces, err := net.Interfaces()
	if err != nil {
		log.Fatal(err.Error())
	}

	for _, networkInterface := range networkInterfaces {
		var n NetworkInterface
		n.Name = networkInterface.Name
		n.MacAddress = networkInterface.HardwareAddr.String()

		addrs, _ := networkInterface.Addrs()
		for _, addr := range addrs {
			n.Addrs = append(n.Addrs, addr.String())
		}

		networkInterfaceMap[networkInterface.Name] = n
	}
	return
}

func basicCPUInfo() (basicCPUInfoMap CPUInfo) {
	out, err := exec.Command("lscpu").Output()
	outstring := strings.TrimSpace(string(out))
	lines := strings.Split(outstring, "\n")
	c := CPUInfo{}

	for _, line := range lines {
		fields := strings.Split(line, ":")
		if len(fields) < 2 {
			continue
		}
		key := strings.TrimSpace(fields[0])
		value := strings.TrimSpace(fields[1])

		switch key {
		case "Socket(s)":
			c.Sockets = convertStringToInteger(value)
		case "Core(s) per socket":
			c.CoresPerSocket = convertStringToInteger(value)
		case "Thread(s) per core":
			c.ThreadsPerCore = convertStringToInteger(value)
		case "CPU(s)":
			c.TotalCores = convertStringToInteger(value)
		case "CPU MHz":
			c.ClockSpeed = convertStringToFloat(value)
		}
	}
	checkErr(err)
	return c
}

func main() {
	basicCPUInfoMap := basicCPUInfo()
	networkInterfaceMap := NetworkInterfacesToMap()

	CPUInfoJSON, _ := json.MarshalIndent(basicCPUInfoMap, "", "  ")
	networkInterfaceJSON, _ := json.MarshalIndent(networkInterfaceMap, "", "  ")
	fmt.Println(string(CPUInfoJSON))
	fmt.Println(string(networkInterfaceJSON))
}
