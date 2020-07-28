package main

import (
	"os"
	"strings"
	"time"
)

func date(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

func fileExist(fullPath string) bool {
	_, err := os.Open(fullPath)
	if err != nil {
		return false
	}

	return true
}

func fileNotExists(fullPath string) bool {
	return !fileExist(fullPath)
}

func splitTexts(longText string, size int) []string {
	parts := []string{}
	cells := strings.Split(longText, "\n")
	tmpPart := ""
	for _, cell := range cells {
		if len(tmpPart)+len(cell) > size {
			parts = append(parts, tmpPart)
			tmpPart = cell
		} else {
			tmpPart += "\n" + cell
		}
	}
	parts = append(parts, tmpPart)

	return parts
}
