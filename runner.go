package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"syscall"

	tea "charm.land/bubbletea/v2"
)

type cmdOut struct {
	out []byte
	err error
}

func runCommand(rawCommand string) tea.Cmd {
	return func() tea.Msg {
		shell := os.Getenv("SHELL")
		if shell == "" {
			shell = "sh"
		}

		cmd := exec.Command(shell, "-ic", rawCommand)
		out, err := cmd.CombinedOutput()

		log.Printf("Executing: %s -ic %s", shell, rawCommand)

		return cmdOut{
			out: out,
			err: err,
		}
	}
}

func openEditor(editor string, arg string) tea.Cmd {
	editorPath, err := exec.LookPath(editor)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	log.Printf("Opening: %s %s", editor, arg)

	return func() tea.Msg {
		syscall.Exec(editorPath, []string{editor, arg}, os.Environ())
		return nil
	}
}

func openEditorAsync (editor string, arg string) tea.Cmd {
	editorPath, err := exec.LookPath(editor)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	log.Printf("Opening Async: %s %s", editor, arg)

	return func() tea.Msg {
		exec.Command(editorPath, arg).Start()
		return nil
	}
}
