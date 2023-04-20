package main

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

func ramUsed(size string, free string) int {
	Size, _ := strconv.Atoi(size)
	Free, _ := strconv.Atoi(free)
	Used := Size - Free
	return Used
}

func ramPer(size string, used string) float32 {
	const bitSize = 64
	Size, _ := strconv.ParseFloat(size, bitSize)
	Used, _ := strconv.ParseFloat(used, bitSize)
	Per := Used / Size * 100
	return float32(Per)
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
	// RamFree
	RamFree := strings.Split(ramStatusStr[1], ": ")
	RamFreeMain := strings.ReplaceAll(RamFree[1], " ", "")
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
	// RamUsed
	RamUsed := strconv.Itoa(ramUsed(RamSizeMain, RamFreeMain))
	// Shared Memory
	MemoryShared := strings.Split(ramSharedStr[5], "\nm")
	MemorySharedMain := strings.ReplaceAll(MemoryShared[0], " ", "")
	// RamPer
	RamPer := ramPer(RamSizeMain, RamUsed)
	RamPerMain := fmt.Sprintf("%f", RamPer)

	ramStatusList = append(ramStatusList, RamSizeMain, AvailableMain, RamPerMain, ActiveMain, InactiveMain, BufferMain, CachedMain, SlabMain, RamUsed, MemorySharedMain)

	return ramStatusList
}
