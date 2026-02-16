package main

import (
	"fmt"
	"sort"
)

// 14. Unique - Remove duplicate elements from slice

func uniqueInts(s []int) []int {
	seen := make(map[int]bool)
	result := make([]int, 0)
	for _, v := range s {
		if !seen[v] {
			seen[v] = true
			result = append(result, v)
		}
	}
	return result
}

func uniqueStrings(s []string) []string {
	seen := make(map[string]bool)
	result := make([]string, 0)
	for _, v := range s {
		if !seen[v] {
			seen[v] = true
			result = append(result, v)
		}
	}
	return result
}

// uniqueSorted - preserves order, works for any comparable type
func uniqueSorted[T comparable](s []T) []T {
	seen := make(map[T]bool)
	result := make([]T, 0)
	for _, v := range s {
		if !seen[v] {
			seen[v] = true
			result = append(result, v)
		}
	}
	return result
}

func main() {
	nums := []int{1, 2, 2, 3, 3, 3, 4, 4, 4, 4}
	fmt.Println("Original:", nums)
	fmt.Println("Unique:", uniqueInts(nums))

	words := []string{"go", "rust", "go", "python", "rust"}
	fmt.Println("Original:", words)
	fmt.Println("Unique:", uniqueStrings(words))

	// With sort for consistent output (optional)
	nums = []int{3, 1, 2, 3, 1, 2}
	unique := uniqueSorted(nums)
	sort.Ints(unique)
	fmt.Println("Unique sorted:", unique)
}
