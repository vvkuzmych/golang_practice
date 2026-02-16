package main

import "fmt"

// 17. Remove - Remove element by index (preserving or not preserving order)

func remove[T any](s []T, i int) []T {
	return append(s[:i], s[i+1:]...)
}

func removeOrdered[T any](s []T, i int) []T {
	result := make([]T, len(s)-1)
	copy(result[:i], s[:i])
	copy(result[i:], s[i+1:])
	return result
}

func removeValue[T comparable](s []T, val T) []T {
	result := make([]T, 0)
	for _, v := range s {
		if v != val {
			result = append(result, v)
		}
	}
	return result
}

func main() {
	nums := []int{1, 2, 3, 4, 5}
	fmt.Println("Original:", nums)

	// Remove index 2
	removed := remove(nums, 2)
	fmt.Println("After remove(2):", removed)

	// Remove first element
	removed = remove(nums, 0)
	fmt.Println("After remove(0):", removed)

	// Remove by value
	nums = []int{1, 2, 2, 3, 2, 4}
	fmt.Println("Before:", nums)
	fmt.Println("Remove all 2s:", removeValue(nums, 2))
}
