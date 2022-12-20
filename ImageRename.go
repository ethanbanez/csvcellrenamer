package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"regexp"
)

func check(e error) {
	if e != nil {
		if e == io.EOF {
			os.Exit(0)
		}
		fmt.Printf("Error %e", e)
		os.Exit(1)
	}
}

// reads a line from a file
// should read up to cells which are comma delimited (csv)
func readCell(reader *bytes.Reader) []byte {
	line := make([]byte, 24)
	b, err := reader.ReadByte()

	// end of line = 0b1010
	// comma = 0b101100
	for b != byte(0b101100) && b != byte(0b1010) {
		if err != nil {
			if err == io.EOF {
				return line
			}
			fmt.Printf("Error %e", err)
			os.Exit(1)
		}
		line = append(line, b)
		b, err = reader.ReadByte()
	}
	line = applyRegex(line)
	// reinserts the comma or the return line to maintain the ordering of the cells
	if b == byte(0b101100) || b == byte(0b1010) {
		line = append(line, b)
	}
	return line
}

func applyRegex(line []byte) []byte {
	reg := regexp.MustCompile(os.Args[1])
	return reg.ReplaceAll(line, []byte(os.Args[2]))
}

func main() {
	// take in two arguments
	//	1 regex to apply
	//	2 file to apply it to
	filePath := os.Args[3] // interestingly, os.Args[0] takes in the executable itself

	fileContents, err := os.ReadFile(filePath)
	check(err)
	file := bytes.NewReader(fileContents)

	// reads the whole file
	lines := make([]byte, 30)
	for file.Len() != 0 {
		line := readCell(file)
		for b := range line {
			lines = append(lines, line[b])
		}
	}
	info, err := os.Stat(filePath)
	check(err)
	err = os.WriteFile("renamed_images.csv", lines, info.Mode()) // this just creates a new file. In order to do it in line we could use pointers...?
	check(err)
}
