package main

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"sort"
)

// TODO(xendak): limit test this? maybe add in config? idk
const maxSuggestions = 10

// NOTE(xendak):maybe json is overkill and we can just use plainfile with dir:cmd:freq:cmd:freq, etc ? easy to parse
type Entry struct {
	Cmd  string `json:"cmd"`
	Freq int    `json:"freq"`
}

type history map[string][]Entry

func historyPath() string {
	base := os.Getenv("XDG_CACHE_HOME")
	if base == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return ""
		}
		base = filepath.Join(home, ".cache")
	}
	return filepath.Join(base, "gobuild", "gobuild.json")
}

// NOTE(xendak): we sort on save, since it teoretically only needs a single pass everytime ?
func saveHistory(cmd string) bool {
	histPath := historyPath()
	if histPath == "" {
		return false
	}

	if err := os.MkdirAll(filepath.Dir(histPath), 0755); err != nil {
		log.Printf("History mkdir error: %v", err)
		return false
	}

	dir, _ := os.Getwd()
	hist := loadHistory()

	entries := hist[dir]
	found := false

	// NOTE(xendak): while using hashmap for this gives o(1), go doesnt have ordered maps ? so sort on every suggestion seems a bad trade?
	for i := range entries {
		if entries[i].Cmd == cmd {
			entries[i].Freq++
			found = true
			break
		}
	}

	if !found {
		entries = append(entries, Entry{Cmd: cmd, Freq: 1})
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Freq > entries[j].Freq
	})

	if len(entries) > maxSuggestions {
		entries = entries[:maxSuggestions]
	}

	hist[dir] = entries

	data, err := json.MarshalIndent(hist, "", "  ")
	if err != nil {
		log.Printf("History marshal error: %v", err)
		return false
	}

	if err := os.WriteFile(histPath, data, 0644); err != nil {
		log.Printf("History write error: %v", err)
		return false
	}

	return true
}

func loadHistory() history {
	hist := make(history)
	path := historyPath()
	if path == "" {
		return hist
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if !os.IsNotExist(err) {
			log.Printf("History read error: %v", err)
		}
		return hist
	}

	if err := json.Unmarshal(data, &hist); err != nil {
		log.Printf("History parse error: %v", err)
	}
	return hist
}

func getSuggestions() []string {
	dir, _ := os.Getwd()
	hist := loadHistory()

	entries := hist[dir]
	if len(entries) == 0 {
		return nil
	}

	out := make([]string, len(entries))
	for i, e := range entries {
		out[i] = e.Cmd
	}

	return out
}
