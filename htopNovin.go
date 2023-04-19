package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
)

type ServerData struct {
	HostName           string  `json:"host_name"`
	Linux              string  `json:"linux"`
	OsName             string  `json:"osname"`
	Uptime             string  `json:"uptime"`
	Kernel             string  `json:"kernel"`
	BashVersion        string  `json:"bash_version"`
	CpuAvg1min         float64 `json:"cpu_avg_1min"`
	CpuAvg5min         float64 `json:"cpu_avg_5min"`
	CpuAvg15min        float64 `json:"cpu_avg_15min"`
	CpuInformation     string  `json:"cpu_information"`
	RamSize            float32 `json:"ram_size"`
	RamStatusTotal     uint64  `json:"ram_status_total"`
	RamStatusAvailable float32 `json:"ram_status_available"`
	RamStatusPercent   float32 `json:"ram_status_percent"`
	RamStatusUsed      float32 `json:"ram_status_used"`
	RamStatusFree      uint64  `json:"ram_status_free"`
	RamStatusActive    float32 `json:"ram_status_active"`
	RamStatusInactive  float32 `json:"ram_status_inactive"`
	RamStatusBuffers   float32 `json:"ram_status_buffers"`
	RamStatusCached    float32 `json:"ram_status_cached"`
	RamStatusShared    float32 `json:"ram_status_shared"`
	RamStatusSlab      float32 `json:"ram_status_slab"`
	Memory             string  `json:"memory"`
}

func int8ToStr(arr []int8) string {
	b := make([]byte, 0, len(arr))
	for _, v := range arr {
		if v == 0x00 {
			break
		}
		b = append(b, byte(v))
	}
	return string(b)
}

func getCpuAvg() []float64 {
	var cpuAvg []float64
	cpuAvg1min, _ := strconv.ParseFloat(loadAvg()[0], 64)
	cpuAvg5min, _ := strconv.ParseFloat(loadAvg()[1], 64)
	cpuAvg15min, _ := strconv.ParseFloat(loadAvg()[2], 64)
	cpuAvg = append(cpuAvg, cpuAvg1min, cpuAvg5min, cpuAvg15min)
	return cpuAvg
}

func uptime() string {
	// create a new *Cmd instance
	cmd := exec.Command("uptime")
	out, err := cmd.Output()
	if err != nil {
		fmt.Println("could not run command: ", err)
	}
	return string(out)
}

func bashVersion() string {
	// create a new *Cmd instance
	cmd := exec.Command("bash", "--version")
	out, err := cmd.Output()
	if err != nil {
		fmt.Println("could not run command: ", err)
	}
	return strings.Split(string(out), "\nCopyright")[0]
}

func cpuInformation() string {
	// create a new *Cmd instance
	cmd := exec.Command("lscpu")
	out, err := cmd.Output()
	if err != nil {
		fmt.Println("could not run command: ", err)
	}
	return string(out)
}

func memory() string {
	// create a new *Cmd instance
	cmd := exec.Command("vmstat", "-s")
	out, err := cmd.Output()
	if err != nil {
		fmt.Println("could not run command: ", err)
	}
	return string(out)
}

func sysTotalMemory() uint64 {
	in := &syscall.Sysinfo_t{}
	err := syscall.Sysinfo(in)
	if err != nil {
		return 0
	}
	// If this is a 32-bit system, then these fields are
	// uint32 instead of uint64.
	// So we always convert to uint64 to match signature.
	return uint64(in.Totalram) * uint64(in.Unit)
}

func sysFreeMemory() uint64 {
	in := &syscall.Sysinfo_t{}
	err := syscall.Sysinfo(in)
	if err != nil {
		return 0
	}
	// If this is a 32-bit system, then these fields are
	// uint32 instead of uint64.
	// So we always convert to uint64 to match signature.
	return uint64(in.Freeram) * uint64(in.Unit)
}

func ServerRequest() {
	var uname syscall.Utsname
	client := &http.Client{}
	hostname, _ := os.Hostname()
	cpuAvg := getCpuAvg()
	datas := ServerData{
		HostName:           hostname,
		Linux:              int8ToStr(uname.Version[:]),
		OsName:             int8ToStr(uname.Sysname[:]),
		Uptime:             uptime(),
		Kernel:             int8ToStr(uname.Release[:]),
		BashVersion:        bashVersion(),
		CpuAvg1min:         cpuAvg[0],
		CpuAvg5min:         cpuAvg[1],
		CpuAvg15min:        cpuAvg[2],
		CpuInformation:     cpuInformation(),
		Memory:             memory(),
		RamSize:            15.524925231933594,
		RamStatusTotal:     sysTotalMemory(),
		RamStatusAvailable: 3210944512,
		RamStatusPercent:   80.7,
		RamStatusUsed:      12685611008,
		RamStatusFree:      sysFreeMemory(),
		RamStatusActive:    2476175360,
		RamStatusInactive:  12848930816,
		RamStatusBuffers:   192327680,
		RamStatusCached:    3426385920,
		RamStatusShared:    431964160,
		RamStatusSlab:      601997312,
	}
	out, _ := json.Marshal(&datas)
	var data = strings.NewReader(string(out))
	req, err := http.NewRequest("POST", "https://stage.htop.ir/", data)
	if err != nil {
		log.Fatal(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
}

func main() {
	ServerRequest()
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
