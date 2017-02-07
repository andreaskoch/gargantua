package main

import (
	"fmt"
	"strings"
)

const (
	BYTE     = 1.0
	KILOBYTE = 1024 * BYTE
	MEGABYTE = 1024 * KILOBYTE
	GIGABYTE = 1024 * MEGABYTE
	TERABYTE = 1024 * GIGABYTE
)

func formatBytes(numberOfBytes int) string {
	unit := ""
	value := float32(numberOfBytes)

	switch {
	case numberOfBytes >= TERABYTE:
		unit = "T"
		value = value / TERABYTE
	case numberOfBytes >= GIGABYTE:
		unit = "G"
		value = value / GIGABYTE
	case numberOfBytes >= MEGABYTE:
		unit = "M"
		value = value / MEGABYTE
	case numberOfBytes >= KILOBYTE:
		unit = "K"
		value = value / KILOBYTE
	case numberOfBytes >= BYTE:
		unit = "B"
	case numberOfBytes == 0:
		return "0"
	}

	stringValue := fmt.Sprintf("%.1f", value)
	stringValue = strings.TrimSuffix(stringValue, ".0")
	return fmt.Sprintf("%s%s", stringValue, unit)
}
