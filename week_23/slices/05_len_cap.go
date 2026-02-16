package main

import "fmt"

// 05. Length vs Capacity - len() and cap() for slices

func main() {
	// Make slice with length 3, capacity 5
	s := make([]int, 3, 5)
	s[0], s[1], s[2] = 1, 2, 3
	fmt.Printf("len=%d cap=%d s=%v\n", len(s), cap(s), s)

	// Append within capacity - no reallocation
	s = append(s, 4)
	fmt.Printf("After append(4): len=%d cap=%d s=%v\n", len(s), cap(s), s)

	// Append beyond capacity - doubles capacity
	s = append(s, 5, 6)
	fmt.Printf("After append(5,6): len=%d cap=%d s=%v\n", len(s), cap(s), s)

	// Slice from array - capacity extends to end of array
	arr := [5]int{1, 2, 3, 4, 5}
	slice := arr[1:3]
	fmt.Printf("arr[1:3]: len=%d cap=%d slice=%v\n", len(slice), cap(slice), slice)
}
