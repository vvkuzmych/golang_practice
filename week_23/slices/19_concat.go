package main

import "fmt"

// 19. Concat - Concatenate multiple slices

func concat[T any](slices ...[]T) []T {
	size := 0
	for _, s := range slices {
		size += len(s)
	}
	result := make([]T, 0, size)
	for _, s := range slices {
		result = append(result, s...)
	}
	return result
}

func main() {
	a := []int{1, 2, 3}
	b := []int{4, 5, 6}
	c := []int{7, 8, 9}

	// Using append
	combined := append(append([]int{}, a...), b...)
	combined = append(combined, c...)
	fmt.Println("Append:", combined)

	// Using concat helper
	combined = concat(a, b, c)
	fmt.Println("Concat:", combined)

	// Concat strings
	words1 := []string{"hello"}
	words2 := []string{"world", "!"}
	fmt.Println("Concat strings:", concat(words1, words2))
}
