package main    

import (
"os"
"io"
"fmt"
"bufio"
"bytes"
"path/filepath"
)

const version = "0.1.0"
const lines = 1_000_000

func main() {
	err := run()
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

func run() error { 

	var err error
	err = nil

    if len(os.Args) > 1 {
        switch os.Args[1] {
        case "-v", "--version":
            Version()
        case "-h", "--help":
            Usage()
        default:
			err = CutFile()
        } 
    } else {
		Usage()
    }

	return err
}

func Version() {
    fmt.Println(filepath.Base(os.Args[0]), "version", version)
}

func Usage() {
    fmt.Println("Usage:", filepath.Base(os.Args[0]), "[options] <file>")
    fmt.Println("  <file>           The path to the text file to be cut")
    fmt.Println("Options:")
    fmt.Println("  -v, --version    Print version information and exit")
    fmt.Println("  -h, --help       Print this message and exit")
}

func CutFile() error {

	filepath := os.Args[1]

	// check if the file exists
	if !fileExists(filepath) {
		return fmt.Errorf("file does not exist: %s", filepath)
	}

	// check if the file is a text file
	if !isTextFile(filepath) {
		return fmt.Errorf("file is not a text file: %s", filepath)
	}

	// check if the file contains more than 1_000_000 lines
	largeEnough, err := isLargeEnough(filepath)
	if err != nil {
		return err
	}
	if !largeEnough {
		return fmt.Errorf("file does not contain more than %d lines: %s", lines, filepath)
	}

	// get the first line from the file and return it as a string
	header, err := getTableHeader(filepath)
	if err != nil {
		return err
	}

	// open the file for reading
	file, err := os.Open(filepath)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()
	
	// read the file 1_000_000 lines at a time and save it to a new file with the
	// same name as the original file but with a suffix of _part1, _part2, etc.
	// inject header into each new file after the first
	scanner := bufio.NewScanner(file)
	part := 1
	lineCount := 0
	var outputFile *os.File
	var writer *bufio.Writer
	
	for scanner.Scan() {
		// if this is the first file
		if lineCount == 0 {
			outPath := createOutputFilePath(filepath, part)
			outputFile, err = os.Create(outPath)
			if err != nil {
				return fmt.Errorf("failed to create output file: %v", err)
			}
			writer = bufio.NewWriter(outputFile)
			_, err = writer.WriteString(header + "\n")
			if err != nil {
				return fmt.Errorf("failed to write header to output file: %v", err)
			}
		}

		_, err = writer.WriteString(scanner.Text() + "\n")
		if err != nil {
			return fmt.Errorf("failed to write line to output file: %v", err)
		}

		lineCount++

		// if we have reached the maximum number of lines, flush the writer 
		// close the file and reset the counter
		if lineCount >= lines {
			writer.Flush()
			outputFile.Close()
			part++
			lineCount = 0
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("failed to read input file: %v", err)
	}

	if lineCount > 0 {
		writer.Flush()
		outputFile.Close()
	}

	return nil

}

// check if a file exists
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}


// check if the file is a text file
func isTextFile(filename string) bool {
	file, err := os.Open(filename)
	if err != nil {
		return false
	}
	defer file.Close()

	// read the first 512 bytes of the file
	buf := make([]byte, 512)
	n, err := file.Read(buf)
	if err != nil && err != io.EOF {
		return false
	}

	// check if the buffer contains any null bytes
	return !bytes.Contains(buf[:n], []byte{0})
}	


// get the first line from the file and return it as a string
func getTableHeader(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	reader := bufio.NewScanner(file)
	if reader.Scan() {
		return reader.Text(), nil
	}
	if err := reader.Err(); err != nil {
		return "", fmt.Errorf("failed to read file: %v", err)
	}
	return "", fmt.Errorf("file is empty")
}

// check if the file contains more than 1_000_000 lines
func isLargeEnough(filename string) (bool, error) {
	file, err := os.Open(filename)
	if err != nil {
		return false, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	reader := bufio.NewScanner(file)
	count := 0
	for reader.Scan() {
		count++
		if count > lines {
			return true, nil
		}
	}
	if err := reader.Err(); err != nil {
		return false, fmt.Errorf("failed to read file: %v", err)
	}
	return false, nil
}

// create a file path with the same name as the original file but with a suffix
// of _part1, _part2, right before the file extension
func createOutputFilePath(filename string, part int) string {
	ext := filepath.Ext(filename)
	name := filename[:len(filename)-len(ext)]
	return fmt.Sprintf("%s_part%d%s", name, part, ext)
}
