package main

import "fmt"

// 18. Insert - Insert element at index

func insert[T any](s []T, idx int, val T) []T {
	if idx < 0 || idx > len(s) {
		return s
	}
	s = append(s, val)
	copy(s[idx+1:], s[idx:len(s)-1])
	s[idx] = val
	return s
}

func insertAtEnd[T any](s []T, val T) []T {
	return append(s, val)
}

func insertAtStart[T any](s []T, val T) []T {
	return append([]T{val}, s...)
}

func main() {
	nums := []int{1, 2, 3, 4, 5}
	fmt.Println("Original:", nums)

	nums = insert(nums, 2, 99)
	fmt.Println("Insert 99 at index 2:", nums)

	nums = insertAtStart(nums, 0)
	fmt.Println("Insert 0 at start:", nums)

	nums = insertAtEnd(nums, 100)
	fmt.Println("Insert 100 at end:", nums)
}
