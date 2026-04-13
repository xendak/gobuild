package main

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

const path = "/home/xendak/.local/cache/gobuild/gobuild.json"

// TODO(xendak):
// prune suggestions (over 10 maybe?) FIFO (?)
// implement frequency of use? maybe?

func save(key string) bool {
	dir, _ := os.Getwd()
	var hist map[string][]string

	dirPath := filepath.Dir(path)
	err := os.MkdirAll(dirPath, 0755)
	if err != nil {
		log.Printf("Mkdir error: %v\n", err)
	}

	_, err = os.Stat(path)
	if err == nil {
		hist = load()
		hist[dir] = append(hist[dir], key)
	} else {
		hist = make(map[string][]string)
		hist[dir] = append(hist[dir], key)
	}

	file, err := os.Create(path)
	if err != nil {
		log.Printf("Create error: %v\n", err)
		return false
	}

	defer file.Close()

	jsonBytes, err := json.Marshal(hist)
	if err != nil {
		log.Printf("Json error: %v\n", err)
	}

	_, err = file.Write(jsonBytes)
	if err != nil {
		log.Printf("Write error: %v\n", err)
	}

	return true
}

func load() map[string][]string {
	hist := make(map[string][]string)

	readBytes, err := os.ReadFile(path)
	if err != nil {
		log.Printf("ReadFile error: %v\n", err)
	}

	err = json.Unmarshal(readBytes, &hist)
	if err != nil {
		log.Printf("Json error: %v\n", err)
	}

	return hist
}

func getSuggestions() []string {
	dir, _ := os.Getwd()
	hist := load()
	suggestions := hist[dir]

	log.Printf("Completing %s with: ", dir)
	for _, str := range suggestions {
		log.Printf("---- %s", str)
	}

	return suggestions
}
