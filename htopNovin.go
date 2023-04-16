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

type ServerData struct {
	Host_name            string  `json:"host_name"`
	Linux                string  `json:"linux"`
	OsName               string  `json:"osname"`
	Uptime               string  `json:"uptime"`
	Kernel               string  `json:"kernel"`
	Bash_version         string  `json:"bash_version"`
	Cpu_avg_1min         float64 `json:"cpu_avg_1min"`
	Cpu_avg_5min         float64 `json:"cpu_avg_5min"`
	Cpu_avg_15min        float64 `json:"cpu_avg_15min"`
	Cpu_information      string  `json:"cpu_information"`
	Ram_size             float32 `json:"ram_size"`
	Ram_status_total     uint64  `json:"ram_status_total"`
	Ram_status_available float32 `json:"ram_status_available"`
	Ram_status_percent   float32 `json:"ram_status_percent"`
	Ram_status_used      float32 `json:"ram_status_used"`
	Ram_status_free      uint64  `json:"ram_status_free"`
	Ram_status_active    float32 `json:"ram_status_active"`
	Ram_status_inactive  float32 `json:"ram_status_inactive"`
	Ram_status_buffers   float32 `json:"ram_status_buffers"`
	Ram_status_cached    float32 `json:"ram_status_cached"`
	Ram_status_shared    float32 `json:"ram_status_shared"`
	Ram_status_slab      float32 `json:"ram_status_slab"`
	Memory               string  `json:"memory"`
}

func main() {

	var uname syscall.Utsname
	if err := syscall.Uname(&uname); err == nil {
		// extract members:
		// type Utsname struct {
		//  Sysname    [65]int8
		//  Nodename   [65]int8
		//  Release    [65]int8
		//  Version    [65]int8
		//  Machine    [65]int8
		//  Domainname [65]int8
		// }
	}

	client := &http.Client{}
	hostname, _ := os.Hostname()
	cpu_avg_1min, _ := strconv.ParseFloat(loadAvg()[0], 32)
	cpu_avg_5min, _ := strconv.ParseFloat(loadAvg()[1], 64)
	cpu_avg_15min, _ := strconv.ParseFloat(loadAvg()[2], 64)
	fmt.Println(cpu_avg_1min)
	datas := ServerData{
		Host_name:            hostname,
		Linux:                int8ToStr(uname.Version[:]),
		OsName:               int8ToStr(uname.Sysname[:]),
		Uptime:               uptime(),
		Kernel:               int8ToStr(uname.Release[:]),
		Bash_version:         bash_version(),
		Cpu_avg_1min:         cpu_avg_1min,
		Cpu_avg_5min:         cpu_avg_5min,
		Cpu_avg_15min:        cpu_avg_15min,
		Cpu_information:      cpu_information(),
		Memory:               memory(),
		Ram_size:             15.524925231933594,
		Ram_status_total:     sysTotalMemory(),
		Ram_status_available: 3210944512,
		Ram_status_percent:   80.7,
		Ram_status_used:      12685611008,
		Ram_status_free:      sysFreeMemory(),
		Ram_status_active:    2476175360,
		Ram_status_inactive:  12848930816,
		Ram_status_buffers:   192327680,
		Ram_status_cached:    3426385920,
		Ram_status_shared:    431964160,
		Ram_status_slab:      601997312,
	}
	// Marshal back to json (as original)
	out, _ := json.Marshal(&datas)
	var data = strings.NewReader(string(out))
	fmt.Println(string(out))
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

func memory() string {
	// create a new *Cmd instance
	// here we pass the command as the first argument and the arguments to pass to the command as the
	// remaining arguments in the function
	cmd := exec.Command("vmstat", "-s")

	// The `Output` method executes the command and
	// collects the output, returning its value
	out, err := cmd.Output()
	if err != nil {
		// if there was any error, print it here
		fmt.Println("could not run command: ", err)
	}
	// otherwise, print the output from running the command
	return string(out)
}

func bash_version() string {
	// create a new *Cmd instance
	// here we pass the command as the first argument and the arguments to pass to the command as the
	// remaining arguments in the function
	cmd := exec.Command("bash", "--version")

	// The `Output` method executes the command and
	// collects the output, returning its value
	out, err := cmd.Output()
	if err != nil {
		// if there was any error, print it here
		fmt.Println("could not run command: ", err)
	}
	// otherwise, print the output from running the command
	return strings.Split(string(out), "\nCopyright")[0]
}
func cpu_information() string {
	// create a new *Cmd instance
	// here we pass the command as the first argument and the arguments to pass to the command as the
	// remaining arguments in the function
	cmd := exec.Command("lscpu")

	// The `Output` method executes the command and
	// collects the output, returning its value
	out, err := cmd.Output()
	if err != nil {
		// if there was any error, print it here
		fmt.Println("could not run command: ", err)
	}
	// otherwise, print the output from running the command
	return string(out)
}
func uptime() string {
	// create a new *Cmd instance
	// here we pass the command as the first argument and the arguments to pass to the command as the
	// remaining arguments in the function
	cmd := exec.Command("uptime")

	// The `Output` method executes the command and
	// collects the output, returning its value
	out, err := cmd.Output()
	if err != nil {
		// if there was any error, print it here
		fmt.Println("could not run command: ", err)
	}
	// otherwise, print the output from running the command
	return string(out)
}
func loadAvg() []string {
	// create a new *Cmd instance
	// here we pass the command as the first argument and the arguments to pass to the command as the
	// remaining arguments in the function
	cmd := exec.Command("cat", "/proc/loadavg")

	// The `Output` method executes the command and
	// collects the output, returning its value
	out, err := cmd.Output()
	if err != nil {
		// if there was any error, print it here
		fmt.Println("could not run command: ", err)
	}
	// otherwise, print the output from running the command
	return strings.Split(string(out), " ")
}
