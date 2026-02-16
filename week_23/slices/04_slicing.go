package main

import "fmt"

// 04. Slicing - slice operations [1:3], [:2], [2:]

func main() {
	nums := []int{10, 20, 30, 40, 50}

	// [1:3] - elements from index 1 up to (not including) 3
	slice1 := nums[1:3]
	fmt.Println("nums[1:3]:", slice1) // [20 30]

	// [:2] - from start to index 2 (exclusive)
	slice2 := nums[:2]
	fmt.Println("nums[:2]:", slice2) // [10 20]

	// [2:] - from index 2 to end
	slice3 := nums[2:]
	fmt.Println("nums[2:]:", slice3) // [30 40 50]

	// Full slice shorthand
	full := nums[:]
	fmt.Println("nums[:]:", full)

	// Slices share underlying array - modifications visible
	slice1[0] = 99
	fmt.Println("After slice1[0]=99, nums:", nums)
}
