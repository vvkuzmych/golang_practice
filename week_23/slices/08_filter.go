package main

import "fmt"

// 08. Filter - Filtering slice elements by condition

func filterInts(s []int, predicate func(int) bool) []int {
	result := make([]int, 0)
	for _, v := range s {
		if predicate(v) {
			result = append(result, v)
		}
	}
	return result
}

func filterStrings(s []string, predicate func(string) bool) []string {
	result := make([]string, 0)
	for _, v := range s {
		if predicate(v) {
			result = append(result, v)
		}
	}
	return result
}

func main() {
	nums := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	evens := filterInts(nums, func(n int) bool { return n%2 == 0 })
	fmt.Println("Evens:", evens)

	odds := filterInts(nums, func(n int) bool { return n%2 != 0 })
	fmt.Println("Odds:", odds)

	gt5 := filterInts(nums, func(n int) bool { return n > 5 })
	fmt.Println("Greater than 5:", gt5)

	words := []string{"go", "rust", "python", "java", "c"}
	long := filterStrings(words, func(s string) bool { return len(s) > 3 })
	fmt.Println("Words len>3:", long)
}
