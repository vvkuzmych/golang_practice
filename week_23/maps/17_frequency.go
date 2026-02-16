package main

import "fmt"

// 17. Frequency - Frequency counter / word count

func main() {
	// Character frequency
	text := "hello world"
	freq := make(map[rune]int)
	for _, c := range text {
		freq[c]++
	}
	fmt.Println("Char frequency:", freq)

	// Word frequency
	words := []string{"go", "rust", "go", "python", "rust", "go"}
	wordFreq := make(map[string]int)
	for _, w := range words {
		wordFreq[w]++
	}
	fmt.Println("Word frequency:", wordFreq)

	// Find most frequent
	maxWord, maxCount := "", 0
	for w, c := range wordFreq {
		if c > maxCount {
			maxCount = c
			maxWord = w
		}
	}
	fmt.Printf("Most frequent: %s (%d times)\n", maxWord, maxCount)
}
