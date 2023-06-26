package main

import (
	"PhilippePitzClairoux/go-memory-reader/memory"
	"PhilippePitzClairoux/go-memory-reader/processes"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

var (
	pid            string
	filter         string
	hideEmptyLines bool
	readLoop       bool
	wg             sync.WaitGroup
)

func init() {
	flag.StringVar(&pid, "pid", "0", "pid of process to dump")
	flag.StringVar(&filter, "filter", "", "dump all addresses for pids that match filter")
	flag.BoolVar(&hideEmptyLines, "hide-empty", false, "hide lines that are completely empty")
	flag.BoolVar(&readLoop, "read-loop", false, "read continuously the data stored in memory")
}

func main() {
	flag.Parse()
	pids := processes.GetPIDList(filter, pid)
	pidChannel := make(chan []processes.Process)

	if len(pids) == 0 && filter == "" {
		fmt.Println("No pid found - could not start program")
		flag.PrintDefaults()
		os.Exit(1)
	}

	fmt.Println("Starting memory dump on following processes : ")
	for _, process := range pids {
		fmt.Printf("%+v\n", process)
	}

	go processes.RefreshPIDList(pidChannel, readLoop, filter, pid, &wg)
	fmt.Println("Memory scan starting...")
	pidChannelConsummer(pidChannel)

	fmt.Println("Waiting for memory reading to finish...")
	wg.Wait()
	fmt.Println("Done!")
}

func pidChannelConsummer(pidChannel chan []processes.Process) {
	for pidList := range pidChannel {
		for _, currentPid := range pidList {
			// Get the path to the process's memory map file
			memMapFilePath := filepath.Join("/proc", currentPid.PID, "maps")

			memMapData := memory.ReadMemoryMapFile(memMapFilePath)

			memoryRanges := memory.ParseMemoryMap(memMapData)
			go memory.DumpMemoryRanges(memoryRanges, &currentPid, hideEmptyLines, &wg)
		}
	}
}
