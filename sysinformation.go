package main

import (
	"fmt"
	"net"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
)

func getCpuAvg() []float64 {
	var cpuAvg []float64
	cpuAvg1min, _ := strconv.ParseFloat(loadAvg()[0], 64)
	cpuAvg5min, _ := strconv.ParseFloat(loadAvg()[1], 64)
	cpuAvg15min, _ := strconv.ParseFloat(loadAvg()[2], 64)
	cpuAvg = append(cpuAvg, cpuAvg1min, cpuAvg5min, cpuAvg15min)
	return cpuAvg
}

func getIp() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println("could not run command: ", err)
	}
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
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
