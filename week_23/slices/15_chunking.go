package main

import "fmt"

// 15. Chunking - Split slice into chunks of specified size

func chunk[T any](s []T, size int) [][]T {
	if size <= 0 {
		return nil
	}
	result := make([][]T, 0, (len(s)+size-1)/size)
	for i := 0; i < len(s); i += size {
		end := i + size
		if end > len(s) {
			end = len(s)
		}
		result = append(result, s[i:end])
	}
	return result
}

func main() {
	nums := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	chunks := chunk(nums, 3)
	fmt.Println("Chunks of 3:", chunks)

	chunks = chunk(nums, 4)
	fmt.Println("Chunks of 4:", chunks)

	chunks = chunk(nums, 7)
	fmt.Println("Chunks of 7:", chunks)

	words := []string{"a", "b", "c", "d"}
	fmt.Println("Chunks of 2:", chunk(words, 2))
}
