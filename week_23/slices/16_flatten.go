package main

import "fmt"

// 16. Flatten - Flatten nested slices into single slice

func flattenInts(matrix [][]int) []int {
	result := make([]int, 0)
	for _, row := range matrix {
		result = append(result, row...)
	}
	return result
}

func flatten[T any](matrix [][]T) []T {
	size := 0
	for _, row := range matrix {
		size += len(row)
	}
	result := make([]T, 0, size)
	for _, row := range matrix {
		result = append(result, row...)
	}
	return result
}

func main() {
	matrix := [][]int{
		{1, 2, 3},
		{4, 5},
		{6, 7, 8, 9},
	}
	fmt.Println("Matrix:", matrix)
	fmt.Println("Flattened:", flattenInts(matrix))

	words := [][]string{
		{"hello"},
		{"world", "!"},
		{},
		{"go"},
	}
	fmt.Println("Words:", words)
	fmt.Println("Flattened:", flatten(words))
}
