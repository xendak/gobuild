package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	tea "charm.land/bubbletea/v2"
)

func main() {
	var scanner *bufio.Scanner
	var files []string

	// TODO(xendak): interactive mode later
	if len(os.Args) > 1 {
		args := os.Args[1:]
		for i := 0; i < len(args); i++ {
			switch args[i] {
			// Note(xendak): pipe concept? but enables updating on demand
			case "--cmd", "-c":
				// NOTE(xendak): consume all args for cmd
				i = len(args)

			default:
				if !strings.HasPrefix(args[i], "-") {
					files = append(files, args[i])
					fmt.Println("Files?: %d", len(files))
				}

			}
		}

		if len(files) > 0 {
			var readFiles strings.Builder

			// NOTE(xendak): if file too big this would cause issues
			// maybe try different approach
			for i := 0; i < len(files); i++ {
				filename := files[i]
				file, err := os.ReadFile(filename)

				if err != nil {
					fmt.Fprintf(os.Stderr, "Error openin file: %v\n", err)
					os.Exit(1)
				}
				readFiles.Write(file)
			}

			fmt.Printf("Reading from file(s): %s\n", files)
			scanner = bufio.NewScanner(strings.NewReader(readFiles.String()))
		}
	} else {
		stat, _ := os.Stdin.Stat()

		if (stat.Mode() & os.ModeCharDevice) != 0 {
			fmt.Fprintln(os.Stderr, "Error: No input provided.")
			os.Exit(1)
		}

		fmt.Println("Reading from piped stdin.")
		scanner = bufio.NewScanner(os.Stdin)
	}

	var msg Message
	msg.count = 0
	first := -1
	for scanner.Scan() {
		text := scanner.Text()
		curr := parseLine(text)
		if first == -1 && curr.Match {
			first = msg.count
		}
		msg.Lines = append(msg.Lines, curr)
		msg.count++
	}

	if msg.count > 0 && first >= 0 {
		m := model{
			msg:    msg,
			cursor: first,
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
