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

func ramStatus() []string {
	ram, errCat := exec.Command("cat", "/proc/meminfo").Output()
	ramShared, errShared := exec.Command("ipcs", "-l").Output()
	var ramStatusList []string
	if errCat != nil {
		fmt.Println("could not run command: ", errCat)
	}
	if errShared != nil {
		fmt.Println("could not run command: ", errShared)
	}
	ramStatusStr := strings.Split(string(ram), " kB")
	ramSharedStr := strings.Split(string(ramShared), " =")
	// RamSize
	RamSize := strings.Split(ramStatusStr[0], ": ")
	RamSizeMain := strings.ReplaceAll(RamSize[1], " ", "")
	// RamStatusAvailable
	Available := strings.Split(ramStatusStr[2], ": ")
	AvailableMain := strings.ReplaceAll(Available[1], " ", "")
	// per cpu
	Per := strings.Split(ramStatusStr[35], ": ")
	PerMain := strings.ReplaceAll(Per[1], " ", "")
	// Active
	Active := strings.Split(ramStatusStr[6], ": ")
	ActiveMain := strings.ReplaceAll(Active[1], " ", "")
	// Inactive
	Inactive := strings.Split(ramStatusStr[7], ": ")
	InactiveMain := strings.ReplaceAll(Inactive[1], " ", "")
	// Buffers
	Buffers := strings.Split(ramStatusStr[3], ": ")
	BufferMain := strings.ReplaceAll(Buffers[1], " ", "")
	// Cached
	Cached := strings.Split(ramStatusStr[4], ": ")
	CachedMain := strings.ReplaceAll(Cached[1], " ", "")
	// Slab
	Slab := strings.Split(ramStatusStr[22], ": ")
	SlabMain := strings.ReplaceAll(Slab[1], " ", "")
	// Swpd
	Swpd := strings.Split(ramStatusStr[5], ": ")
	SwpdMain := strings.ReplaceAll(Swpd[1], " ", "")
	// Shared Memory
	memoryShared := strings.Split(ramSharedStr[5], "\nm")
	memorySharedMain := strings.ReplaceAll(memoryShared[0], " ", "")

	ramStatusList = append(ramStatusList, RamSizeMain, AvailableMain, PerMain, ActiveMain, InactiveMain, BufferMain, CachedMain, SlabMain, SwpdMain, memorySharedMain)

	return ramStatusList
}

func kernelVersion() string {
	kernel, err := exec.Command("uname", "-r").Output()
	if err != nil {
		fmt.Println("could not run command: ", err)
	}
	return string(kernel)
}

func osVersion() string {
	osV, err := exec.Command("cat", "/etc/os-release").Output()
	osVString := strings.Split(string(osV), "ID=")
	if err != nil {
		fmt.Println("could not run command: ", err)
	}
	return osVString[0]
}

func ServerRequest() {
	var uname syscall.Utsname
	client := &http.Client{}
	hostname, _ := os.Hostname()
	cpuAvg := getCpuAvg()
	ramStatuses := ramStatus()
	datas := ServerData{
		HostName:           hostname,
		Linux:              osVersion(),
		OsName:             int8ToStr(uname.Sysname[:]),
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
