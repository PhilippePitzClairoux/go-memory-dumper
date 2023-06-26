package processes

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type Process struct {
	PID     string
	Command string
}

// GetRunningProcesses main call that returns an array of processes
func GetRunningProcesses() []Process {
	cmd := exec.Command("ps", "aux")
	output, err := cmd.Output()
	strOutput := string(output)
	if err != nil {
		return nil
	}

	lines := strings.Split(strOutput, "\n")
	for i := range lines {
		lines[i] = strings.Join(strings.Fields(lines[i]), " ")
	}

	var processes []Process

	for _, line := range lines {
		fields := strings.Fields(line)

		if len(fields) > 0 && fields[0] != "USER" {
			process := Process{
				PID:     fields[1],
				Command: strings.Join(fields[10:], " "),
			}
			processes = append(processes, process)
		}
	}

	return processes
}

// HandleFileClosure handle erros for `defer File.Close()`
func HandleFileClosure(memFile *os.File) {
	errr := memFile.Close()
	if errr != nil {
		log.Fatalf("Failed to handleFileClosure memory file: %v", errr)
	}
}

// OpenPIDMemoryFile open the memory file of a specific memory file
func OpenPIDMemoryFile(currentPid string) (*os.File, error) {
	memFilePath := filepath.Join("/proc", currentPid, "mem")
	memFile, errr := os.Open(memFilePath)
	if errr != nil {
		return nil, errr
	}

	return memFile, nil
}

// GetPIDList get current list of PIDs
func GetPIDList(filter string, pid string) []Process {
	processes := GetRunningProcesses()
	pids := make([]Process, 0)

	if filter != "" {
		for _, process := range processes {
			if strings.Contains(process.Command, filter) && !strings.Contains(process.Command, os.Args[0]) {
				pids = append(pids, process)
			}
		}
	} else if pid != "0" {
		pids = append(pids, Process{pid, ""})
	}

	return pids
}

// RefreshPIDList get a new list of pid every 1 min and send it to channel
func RefreshPIDList(pidChannel chan []Process, loop bool, filter string, pid string, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()
	defer close(pidChannel)

	for {
		pidChannel <- GetPIDList(filter, pid)
		time.Sleep(time.Second * 60)

		if !loop {

			break
		}
	}
}
