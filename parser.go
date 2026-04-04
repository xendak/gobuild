package main

import (
	"bufio"
	"regexp"
	"strconv"
	"strings"
)

type Severity int

const (
	None Severity = iota
	Note
	Info
	Hint
	Warning
	Error
)

var sevString = [...]string{
	"None",
	"Note",
	"Info",
	"Hint",
	"Warning",
	"Error",
}

func (s Severity) String() string {
	if s < 0 {
		panic("Assert Failed: Severity is less than 0")
	}
	if int(s) > len(sevString) {
		panic("Assert Failed: Severity exceed the sevString[] length")
	}

	return sevString[int(s)] + ": "
}

type Line struct {
	Raw   string
	File  string
	Lin   int
	Col   int
	Sev   Severity
	Msg   string
	Match bool
}

// TODO(xendak): add more functionaly ? maybe

type Message struct {
	Lines []Line
	match int
}

type pattern struct{ regex *regexp.Regexp }

// TODO(xendak): maybe we check the match in a different way ?
var patterns = []pattern{
	// Odin:  file.odin(10:8) Syntax Error: message
	{regexp.MustCompile(`^(?P<file>[^(\n]+)\((?P<line>\d+):(?P<col>\d+)\)\s*(?P<sev>Syntax Error|Error|Warning|Note)?:?\s*(?P<msg>.*)$`)},
	// Generic with col:  file:line:col: [sev:] message
	{regexp.MustCompile(`^(?P<file>[^:\n]+):(?P<line>\d+):(?P<col>\d+):\s*(?P<sev>error|warning|note|info|hint)?:?\s*(?P<msg>.*)$`)},
	// Generic no col:    file:line: message
	{regexp.MustCompile(`^(?P<file>[^:\n]+):(?P<line>\d+):\s*(?P<sev>error|warning|note|info|hint)?:?\s*(?P<msg>.*)$`)},
	// Python:  File "file.py", line 42
	{regexp.MustCompile(`^\s*File "(?P<file>[^"]+)", line (?P<line>\d+)`)},
}

func matchLine(regex *regexp.Regexp, line string) map[string]string {
	match := regex.FindStringSubmatch(line)
	if match == nil {
		return nil
	}

	// NOTE(xendak): empty hash table basically
	result := make(map[string]string)
	for i, name := range regex.SubexpNames() {
		if name != "" && i < len(match) {
			result[name] = match[i]
		}
	}
	return result
}

func getSeverity(msg string) Severity {
	lower := strings.ToLower(msg)
	switch {
	case strings.Contains(lower, "error"):
		return Error
	case strings.Contains(lower, "warning"):
		return Warning
	// TODO(xendak): expand to the other classifications
	case lower == "note" || lower == "hint" || lower == "info":
		return Note
	default:
		return None
	}
}

// TODO(xendak): create a ViewLine and a expandedLine, so we can wrap text
func parseMsg(raw string) Message {
	var msg Message
	scanner := bufio.NewScanner(strings.NewReader(raw))

	msg.match = -1
	for scanner.Scan() {
		curr := scanner.Text()
		// FIXME(xendak): make this a proper function or something
		// curr = regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`).ReplaceAllString(curr, "")
		parsed := parseLine(curr)

		if parsed.Match && msg.match < 0 {
			msg.match = len(msg.Lines)
		}

		msg.Lines = append(msg.Lines, parsed)

	}

	return msg
}

func parseLine(raw string) Line {
	result := Line{Raw: raw, Match: false}

	for _, pattern := range patterns {
		match := matchLine(pattern.regex, raw)

		if match == nil {
			continue
		}

		result.File = match["file"]
		result.Lin, _ = strconv.Atoi(match["line"])

		// failsafe so we can do {file}:{line}:1 instead of panic
		// Helix cant invoke open %s:line:0 for some reason :')
		result.Col, _ = strconv.Atoi(match["col"])
		result.Col = max(1, result.Col)

		result.Sev = getSeverity(match["sev"])

		result.Msg = match["msg"]
		result.Match = true

		// avoid stdin
		if strings.HasPrefix(result.File, "-") || result.Lin == 0 {
			result.Match = false
			continue
		}

		break
	}
	return result
}
