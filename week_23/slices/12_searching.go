package main

import (
	"fmt"
	"sort"
)

// 12. Searching - Binary search with sort.Search

func main() {
	nums := []int{1, 3, 5, 7, 9, 11}
	sort.Ints(nums)

	// sort.Search finds first index where nums[i] >= 5
	idx := sort.SearchInts(nums, 5)
	fmt.Printf("Index of 5: %d (value=%d)\n", idx, nums[idx])

	// Search for non-existent - returns insert position
	idx = sort.SearchInts(nums, 6)
	fmt.Printf("Index for 6 (insert pos): %d\n", idx)

	// sort.Search with custom predicate
	idx = sort.Search(len(nums), func(i int) bool { return nums[i] >= 7 })
	fmt.Printf("First >= 7 at index %d\n", idx)

	// Check if element exists
	target := 7
	idx = sort.SearchInts(nums, target)
	exists := idx < len(nums) && nums[idx] == target
	fmt.Printf("7 exists: %v\n", exists)
}
