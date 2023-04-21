package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

var version = "0.0.1"

type ServerData struct {
	HostName           string  `json:"host_name"`
	Ip                 string  `json:"ip"`
	Linux              string  `json:"linux"`
	OsName             string  `json:"osname"`
	Uptime             string  `json:"uptime"`
	Kernel             string  `json:"kernel"`
	BashVersion        string  `json:"bash_version"`
	CpuAvg1min         float64 `json:"cpu_avg_1min"`
	CpuAvg5min         float64 `json:"cpu_avg_5min"`
	CpuAvg15min        float64 `json:"cpu_avg_15min"`
	CpuInformation     string  `json:"cpu_information"`
	RamSize            string  `json:"ram_size"`
	RamStatusTotal     uint64  `json:"ram_status_total"`
	RamStatusAvailable string  `json:"ram_status_available"`
	RamStatusPercent   string  `json:"ram_status_percent"`
	RamStatusUsed      string  `json:"ram_status_used"`
	RamStatusFree      uint64  `json:"ram_status_free"`
	RamStatusActive    string  `json:"ram_status_active"`
	RamStatusInactive  string  `json:"ram_status_inactive"`
	RamStatusBuffers   string  `json:"ram_status_buffers"`
	RamStatusCached    string  `json:"ram_status_cached"`
	RamStatusShared    string  `json:"ram_status_shared"`
	RamStatusSlab      string  `json:"ram_status_slab"`
	Memory             string  `json:"memory"`
}

func main() {
	args := os.Args[1]
	if args == "version" {
		fmt.Println(version)
	} else {
		status_runner()
	}
}
func status_runner() {

	client := &http.Client{}
	hostname, _ := os.Hostname()
	cpuAvg := getCpuAvg()
	ramStatuses := ramStatus()
	osNameInfo := strings.Split(osInformation(), "\n")[0]
	osName := strings.Split(osNameInfo, "\"")
	osVersionInfo := strings.Split(osInformation(), "\n")[1]
	osVersion := strings.Split(osVersionInfo, "\"")
	fmt.Println()
	datas := ServerData{
		HostName:           hostname,
		Ip:                 getIp(),
		Linux:              osVersion[1],
		OsName:             osName[1],
		Uptime:             uptime(),
		Kernel:             kernelVersion(),
		BashVersion:        bashVersion(),
		CpuAvg1min:         cpuAvg[0],
		CpuAvg5min:         cpuAvg[1],
		CpuAvg15min:        cpuAvg[2],
		CpuInformation:     cpuInformation(),
		Memory:             memory(),
		RamSize:            ramStatuses[0],
		RamStatusTotal:     sysTotalMemory(),
		RamStatusAvailable: ramStatuses[1],
		RamStatusPercent:   ramStatuses[2],
		RamStatusUsed:      ramStatuses[8],
		RamStatusFree:      sysFreeMemory(),
		RamStatusActive:    ramStatuses[3],
		RamStatusInactive:  ramStatuses[4],
		RamStatusBuffers:   ramStatuses[5],
		RamStatusCached:    ramStatuses[6],
		RamStatusShared:    ramStatuses[9],
		RamStatusSlab:      ramStatuses[7],
	}
	out, _ := json.Marshal(&datas)
	var data = strings.NewReader(string(out))
	req, err := http.NewRequest("POST", "https://stage.htop.ir/", data)
	req.Header.Set("HTOP-AGENT-VERSION", version)
	if err != nil {
		log.Fatal(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
}

func loadAvg() []string {
	// create a new *Cmd instance
	cmd := exec.Command("cat", "/proc/loadavg")
	out, err := cmd.Output()
	if err != nil {
		fmt.Println("could not run command: ", err)
	}
	return strings.Split(string(out), " ")
}
