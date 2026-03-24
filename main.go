package main

import (
	"bufio"
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
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
	var parsed []Data
	for scanner.Scan() {
		text := scanner.Text()
		curr := parseLine(text)
		if curr.Match {
			parsed = append(parsed, curr)
			lc++
		}
	}

	if lc > 0 {
		m := model{
			lines: parsed,
		}

		p := tea.NewProgram(m)
		if _, err := p.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
			os.Exit(1)
		}
	} else {
		fmt.Fprintf(os.Stderr, "Nothing matched.\n")
	}
}
