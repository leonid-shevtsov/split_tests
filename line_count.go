package main

import (
	"bytes"
	"io"
	"os"
)

func estimateFileTimesByLineCount(currentFileSet map[string]bool, fileTimes map[string]float64) {
	for fileName := range currentFileSet {
		file, err := os.Open(fileName)
		if err != nil {
			printMsg("failed to count lines in file %s: %v\n", fileName, err)
			continue
		}
		defer file.Close()
		lineCount, err := lineCounter(file)
		if err != nil {
			printMsg("failed to count lines in file %s: %v\n", fileName, err)
			continue
		}
		fileTimes[fileName] = float64(lineCount)
	}
}

// Credit to http://stackoverflow.com/a/24563853/6678
func lineCounter(r io.Reader) (int, error) {
	buf := make([]byte, 32*1024)
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := r.Read(buf)
		count += bytes.Count(buf[:c], lineSep)

		switch {
		case err == io.EOF:
			return count, nil

		case err != nil:
			return count, err
		}
	}
}
