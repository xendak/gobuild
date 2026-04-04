package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	tea "charm.land/bubbletea/v2"
)

func main() {
	m := newModel()
	var text string
	var files []string

	appState := Interactive

	stat, _ := os.Stdin.Stat()

	// :debug
	logPath := "/tmp/gobuild.log"
	os.Remove(logPath)
	f, _ := tea.LogToFile(logPath, "debug")
	defer f.Close()

	if (stat.Mode() & os.ModeCharDevice) == 0 {
		appState = Results
		bytes, err := io.ReadAll(os.Stdin)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		}
		text = string(bytes)

		// TODO(xendak): can we get this on linux(?)
		// cmdText, err := getSourceCommand()
	}

	if len(os.Args) > 1 {
		appState = Results

		args := os.Args[1:]
		for i := 0; i < len(args); i++ {
			switch args[i] {

			// NOTE(xendak): pipe concept? but enables updating on demand
			case "--cmd", "-c":
				cmdStr := strings.Join(args[i+1:], " ")
				log.Printf("Cmd: %s\n", cmdStr)

				m.cmd = runCommand(cmdStr)
				// NOTE(xendak): consume all args for cmd
				i = len(args)

			default:
				if !strings.HasPrefix(args[i], "-") {
					files = append(files, args[i])
					log.Printf("Files?: %d\n", len(files))
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

			log.Printf("Reading from file(s): %s\n", files)
			text = readFiles.String()
		}
	}

	if appState == Results {
		m.msg = parseMsg(text)
		m.cursor = max(0, m.msg.match)
		if m.msg.match < 0 {
			appState = Passthrough
		}
	}

	if (len(m.msg.Lines) > 0) || appState == Interactive || m.cmd != nil {
		m.state = AppState(appState)

		p := tea.NewProgram(m)

		_, err := p.Run()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
			os.Exit(1)
		}

	} else {
		fmt.Fprintf(os.Stderr, "Nothing matched.\n")
	}
}
