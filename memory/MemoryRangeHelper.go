package memory

import (
	"PhilippePitzClairoux/go-memory-reader/processes"
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

var extractLine = regexp.MustCompile(`([0-9a-fA-F]{12,16})-([0-9a-fA-F]{12,16})\s+([rwxp-]{4})\s+([0-9-A-Fa-f]{8})(.*)`)

type MemoryRange struct {
	start       uint64
	end         uint64
	permissions string
	offset      uint64
}

// DumpMemoryRanges loops on all memoryRanges and dumps them in stdout
func DumpMemoryRanges(memoryRanges []MemoryRange, currentPid *processes.Process, hideEmptyLines bool, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()

	for _, memRange := range memoryRanges {
		localMemRange := &memRange

		if !strings.Contains(localMemRange.permissions, "r") {
			//fmt.Printf("skipping memory range (%s) - cannot read addresses\n", memRange)
			continue
		}

		go DumpMemoryRange(*localMemRange, currentPid, hideEmptyLines, wg)
	}
}

// ReadMemoryFile open's memory file
func ReadMemoryMapFile(memMapFilePath string) *os.File {
	memMapData, errr := os.Open(memMapFilePath) //os.ReadFile(memMapFilePath)
	if errr != nil {
		log.Fatalf("Failed to read memory map file: %v", errr)
	}
	return memMapData
}

// DumpMemoryRange dumps the content of a specific memory range
func DumpMemoryRange(memRange MemoryRange, currentPid *processes.Process, hideEmptyLines bool, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()

	memSize, errr := getMemSize(memRange)
	fmt.Printf("size : %d\nmemrange : %+v\n", memSize, memRange)
	if errr != nil {
		log.Println("cannot get memsize : ", errr)
		return
	}

	memFile, errr := processes.OpenPIDMemoryFile(currentPid.PID)
	defer processes.HandleFileClosure(memFile)
	if errr != nil {
		log.Fatalf("Failed to open memory file: %v", errr)
	}

	// Read memory from the specified range
	readMemoryAndPrintContent(memSize, memFile, memRange, hideEmptyLines)
}

func getMemSize(memRange MemoryRange) (uint64, error) {
	if memRange.start > memRange.end {
		return 0, errors.New(fmt.Sprintf("invalid memory range %+v", memRange))
	}

	memSize := memRange.end - memRange.start

	return memSize, nil
}

// ParseMemoryMap returns an array of MemoryRanges from read from memory file
func ParseMemoryMap(file *os.File) []MemoryRange {
	//lines := strings.Split(string(data), "\n")
	lines := make([][]string, 0)
	ranges := make([]MemoryRange, 0, len(lines))
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		var line = scanner.Text()

		if !extractLine.MatchString(line) {
			continue
		}

		matches := extractLine.FindStringSubmatch(line)
		fmt.Println(matches)
		memRange := MemoryRange{
			start:       parseStringUint64(matches[1]),
			end:         parseStringUint64(matches[2]),
			permissions: matches[3],
			offset:      parseStringUint64(matches[4]),
		}

		ranges = append(ranges, memRange)
	}

	return ranges
}

func parseStringUint64(input string) uint64 {
	out, err := strconv.ParseUint(input, 16, 64)
	if err != nil {
		log.Fatalf("Failed to parse string : %v", err)
	}

	return out
}

func readMemoryAndPrintContent(memSize uint64, memFile *os.File, memRange MemoryRange, hideEmptyLines bool) {
	memBuffer := make([]byte, memSize)
	_, errr := memFile.ReadAt(memBuffer, int64(memRange.start))
	if errr != nil {
		fmt.Printf("failed to read memory: %v\n", errr)
		memBuffer = make([]byte, 0)
		return
	}

	// Print the content of the read memory
	fmt.Printf("%s\n", FormatMemoryToReadableText(memBuffer, memRange, hideEmptyLines))
	memBuffer = make([]byte, 0)
	//time.Sleep(10 * time.Second)
}
