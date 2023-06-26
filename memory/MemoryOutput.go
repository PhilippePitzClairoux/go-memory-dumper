package memory

import (
	"bytes"
	"fmt"
)

func FormatMemoryToReadableText(data []byte, memoryRange MemoryRange, hideEmptyLines bool) string {
	currentMemoryRange := memoryRange.start
	emptyData := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	text := ""
	buffer := ""

	for i, b := range data {
		buffer = changeByteToChar(b, buffer)
		if (i+1)%30 == 0 {
			if (hideEmptyLines && !bytes.Equal(data[i-29:i+1], emptyData)) || !hideEmptyLines {
				text += fmt.Sprintf("0x%09x - 0x%09x => %s | %s\n",
					currentMemoryRange,
					currentMemoryRange+30,
					formatArray(data[i-29:i+1]),
					buffer,
				)
			}

			buffer = ""
			currentMemoryRange += 30
		}
	}
	return text
}

func changeByteToChar(b byte, buffer string) string {
	if b >= 32 && b <= 126 {
		buffer += string(b)
	} else {
		buffer += "."
	}
	return buffer
}

func formatArray(data []byte) string {
	output := ""
	for _, num := range data {
		output += fmt.Sprintf("%03d ", num)
	}

	return output
}
