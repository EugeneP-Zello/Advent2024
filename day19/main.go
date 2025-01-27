package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func isValidWord(word string, alphabet []string) bool {
	for _, letter := range alphabet {
		if subWord, found := strings.CutPrefix(word, letter); found {
			if subWord == "" || isValidWord(subWord, alphabet) {
				return true
			}
		}
	}
	return false
}

func getValidCollectionsCount(word string, alphabet []string, known map[string]int) int {
	if cached, ok := known[word]; ok {
		return cached
	}
	count := 0
	for _, letter := range alphabet {
		if subWord, found := strings.CutPrefix(word, letter); found {
			if subWord == "" {
				count += 1
			} else {
				count += getValidCollectionsCount(subWord, alphabet, known)
			}
		}
	}
	known[word] = count
	return count
}


func getValidOutputs(filename string) (int, int) {
	file, _ := os.Open(filename)
	content, _ := io.ReadAll(file)
	file.Close()
	parts := strings.Split(string(content), "\r\n\r\n")
	alphabet := strings.Split(parts[0], ", ")
	words := strings.Split(parts[1], "\r\n")
	count := 0
	for _, word := range words {
		if isValidWord(word, alphabet) {
			count++
		}
	}
	totalCount := 0
	known := make(map[string]int)
	for idx, word := range words {
		totalCount += getValidCollectionsCount(word, alphabet, known)
		fmt.Printf("Word#%d: %d", idx+1, totalCount)
	}
	return count, totalCount
}

func runForFile(filename string) {
	validCount, totalCount := getValidOutputs(filename)
	fmt.Printf("File: %s, Valid designs: %d, total %d\n", filename, validCount, totalCount)
}

func main() {
	runForFile("test.txt")
	runForFile("input.txt")
}
