package main

import (
	"fmt"
	"sort"
)

// 11. Sorting - sort.Ints, sort.Slice, sort.SliceStable

func main() {
	// Built-in sort for ints
	nums := []int{3, 1, 4, 1, 5, 9, 2, 6}
	sort.Ints(nums)
	fmt.Println("Sorted ints:", nums)

	// sort.Slice for any type with custom comparator
	strs := []string{"banana", "apple", "cherry"}
	sort.Slice(strs, func(i, j int) bool { return strs[i] < strs[j] })
	fmt.Println("Sorted strings:", strs)

	// sort.SliceStable - keeps equal elements in original order
	people := []struct {
		Name string
		Age  int
	}{
		{"Alice", 30},
		{"Bob", 25},
		{"Charlie", 25},
	}
	sort.SliceStable(people, func(i, j int) bool { return people[i].Age < people[j].Age })
	fmt.Println("Stable sort by age:", people)

	// Reverse sort
	sort.Slice(nums, func(i, j int) bool { return nums[i] > nums[j] })
	fmt.Println("Reverse sorted:", nums)
}
