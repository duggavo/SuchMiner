package main

import (
	"os"
	"strings"
)

type LogReader struct {
	File *os.File
}

func (l LogReader) Write(b []byte) (int, error) {
	n, err := os.Stdout.Write(b)
	if err != nil {
		return n, err
	}
	if l.File != nil {
		// remove all the colors, log file should not contain colors
		bstr := string(b)
		for _, col := range allColors {
			bstr = strings.ReplaceAll(bstr, col, "")
		}

		_, err := l.File.WriteString(bstr)
		return n, err
	}
	return n, err
}
