package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	var scanner *bufio.Scanner
	if len(os.Args) > 1 {
		filename := os.Args[1]
		file, err := os.Open(filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error openin file: %v\n", err)
			os.Exit(1)
		}
		defer file.Close()

		fmt.Printf("Reading from file: %s\n", filename)
		scanner = bufio.NewScanner(file)
	} else {
		stat, _ := os.Stdin.Stat()

		if (stat.Mode() & os.ModeCharDevice) != 0 {
			fmt.Fprintln(os.Stderr, "Error: No input provided. Pipe data or pass a file!")
			os.Exit(1)
		}

		fmt.Println("Reading from piped stdin:")
		scanner = bufio.NewScanner(os.Stdin)
	}

	lc := 0
	for scanner.Scan() {
		lc++
		text := scanner.Text()
		parsed := parseLine(text)
		if parsed.Match {
			fmt.Println("File = %s", parsed.File)
			fmt.Println("Line:Col = %d:%d", parsed.Lin, parsed.Col)
			fmt.Println("Msg = %s", parsed.Msg)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
	}
}
