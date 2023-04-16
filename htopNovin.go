package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
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
	Cpu_avg_1min         float32 `json:"cpu_avg_1min"`
	Cpu_avg_5min         float32 `json:"cpu_avg_5min"`
	Cpu_avg_15min        float32 `json:"cpu_avg_15min"`
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
	datas := ServerData{
		Host_name:            hostname,
		Linux:                int8ToStr(uname.Version[:]),
		OsName:               int8ToStr(uname.Sysname[:]),
		Uptime:               "18:54:09 up  6:18,  2 users,  load average: 0.91, 0.95, 0.98",
		Kernel:               int8ToStr(uname.Release[:]),
		Bash_version:         "GNU bash, version 5.1.16(1)-release (x86_64-pc-linux-gnu)",
		Cpu_avg_1min:         0,
		Cpu_avg_5min:         0,
		Cpu_avg_15min:        0,
		Cpu_information:      "Architecture:                    x86_64\nCPU op-mode(s):                  32-bit, 64-bit\nAddress sizes:                   39 bits physical, 48 bits virtual\nByte Order:                      Little Endian\nCPU(s):                          4\nOn-line CPU(s) list:             0-3\nVendor ID:                       GenuineIntel\nModel name:                      Intel(R) Core(TM) i7-4600M CPU @ 2.90GHz\nCPU family:                      6\nModel:                           60\nThread(s) per core:              2\nCore(s) per socket:              2\nSocket(s):                       1\nStepping:                        3\nCPU max MHz:                     3600.0000\nCPU min MHz:                     800.0000\nBogoMIPS:                        5787.09\nFlags:                           fpu vme de pse tsc msr pae mce cx8 apic sep mtrr pge mca cmov pat pse36 clflush dts acpi mmx fxsr sse sse2 ss ht tm pbe syscall nx pdpe1gb rdtscp lm constant_tsc arch_perfmon pebs bts rep_good nopl xtopology nonstop_tsc cpuid aperfmperf pni pclmulqdq dtes64 monitor ds_cpl smx est tm2 ssse3 sdbg fma cx16 xtpr pdcm pcid sse4_1 sse4_2 x2apic movbe popcnt tsc_deadline_timer aes xsave avx f16c rdrand lahf_lm abm cpuid_fault epb invpcid_single pti ssbd ibrs ibpb stibp fsgsbase tsc_adjust bmi1 avx2 smep bmi2 erms invpcid xsaveopt dtherm ida arat pln pts md_clear flush_l1d\nL1d cache:                       64 KiB (2 instances)\nL1i cache:                       64 KiB (2 instances)\nL2 cache:                        512 KiB (2 instances)\nL3 cache:                        4 MiB (1 instance)\nNUMA node(s):                    1\nNUMA node0 CPU(s):               0-3\nVulnerability Itlb multihit:     KVM: Mitigation: VMX unsupported\nVulnerability L1tf:              Mitigation; PTE Inversion\nVulnerability Mds:               Mitigation; Clear CPU buffers; SMT vulnerable\nVulnerability Meltdown:          Mitigation; PTI\nVulnerability Mmio stale data:   Unknown: No mitigations\nVulnerability Retbleed:          Not affected\nVulnerability Spec store bypass: Mitigation; Speculative Store Bypass disabled via prctl\nVulnerability Spectre v1:        Mitigation; usercopy/swapgs barriers and __user pointer sanitization\nVulnerability Spectre v2:        Mitigation; Retpolines, IBPB conditional, IBRS_FW, STIBP conditional, RSB filling, PBRSB-eIBRS Not affected\nVulnerability Srbds:             Mitigation; Microcode\nVulnerability Tsx async abort:   Not affected', 'memory': '16279064 K total memory\n     12388292 K used memory\n      2418140 K active memory\n     12547784 K inactive memory\n       356872 K free memory\n       187820 K buffer memory\n      3346080 K swap cache\n      2097148 K total swap\n       589568 K used swap\n      1507580 K free swap\n       649484 non-nice user cpu ticks\n         3240 nice user cpu ticks\n       137750 system cpu ticks\n      3444062 idle cpu ticks\n         8645 IO-wait cpu ticks\n            0 IRQ cpu ticks\n         5124 softirq cpu ticks\n            0 stolen cpu ticks\n      7015721 pages paged in\n      5727013 pages paged out\n         3242 pages swapped in\n       149299 pages swapped out\n     34425951 interrupts\n     96476236 CPU context switches\n   1681635942 boot time\n        53105 forks",
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
