package main

import (
	"fmt"
	"strings"
)

// 09. Map Operation - Transform slice values (like map in functional programming)

func mapInts(s []int, fn func(int) int) []int {
	result := make([]int, len(s))
	for i, v := range s {
		result[i] = fn(v)
	}
	return result
}

func mapStrings(s []string, fn func(string) string) []string {
	result := make([]string, len(s))
	for i, v := range s {
		result[i] = fn(v)
	}
	return result
}

func main() {
	nums := []int{1, 2, 3, 4, 5}

	doubled := mapInts(nums, func(n int) int { return n * 2 })
	fmt.Println("Doubled:", doubled)

	squared := mapInts(nums, func(n int) int { return n * n })
	fmt.Println("Squared:", squared)

	words := []string{"hello", "world", "go"}
	upper := mapStrings(words, strings.ToUpper)
	fmt.Println("Uppercase:", upper)

	withPrefix := mapStrings(words, func(s string) string { return "prefix_" + s })
	fmt.Println("With prefix:", withPrefix)
}
