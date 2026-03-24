package main

import (
	"regexp"
	"strconv"
	"strings"
)

type Data struct {
	Raw   string
	File  string
	Lin   int
	Col   int
	Msg   string
	Match bool
}

type pattern struct{ regex *regexp.Regexp }

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

	// empty hash table basically
	result := make(map[string]string)
	for i, name := range regex.SubexpNames() {
		if name != "" && i < len(match) {
			result[name] = match[i]
		}
	}
	return result
}

func parseLine(raw string) Data {
	result := Data{Raw: raw, Match: false}
	for _, pattern := range patterns {
		match := matchLine(pattern.regex, raw)
		// if this line is not a match, we move on
		if match == nil {
			continue
		}

		result.File = match["file"]
		result.Lin, _ = strconv.Atoi(match["line"])

		// failsafe so we dont need to deal with this on cases where its only {file} and {line}
		// since we can just call open editor at col 1 without any issue
		result.Col, _ = strconv.Atoi(match["col"])
		if result.Col == 0 {
			result.Col = 1
		}

		result.Msg = match["msg"]
		result.Match = true

		// avoid stdin
		if strings.HasPrefix(result.File, "-") || result.Lin == 0 {
			result.Match = false
			continue
		}

	}
	return result
}
