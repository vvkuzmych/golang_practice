package main

import "fmt"

// 13. Reverse - Reversing slice in place and creating reversed copy

func reverseInPlace[S ~[]E, E any](s S) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

func reverseCopy[S ~[]E, E any](s S) S {
	result := make(S, len(s))
	for i, v := range s {
		result[len(s)-1-i] = v
	}
	return result
}

func main() {
	nums := []int{1, 2, 3, 4, 5}
	fmt.Println("Original:", nums)

	reverseInPlace(nums)
	fmt.Println("Reversed in place:", nums)

	// Restore and create copy
	nums = []int{1, 2, 3, 4, 5}
	reversed := reverseCopy(nums)
	fmt.Println("Original:", nums)
	fmt.Println("Reversed copy:", reversed)

	words := []string{"a", "b", "c"}
	reverseInPlace(words)
	fmt.Println("Reversed strings:", words)
}
