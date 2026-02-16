package main

import "fmt"

// 10. Reduce - Aggregate slice to single value

func reduceInts(s []int, fn func(acc, val int) int, initial int) int {
	acc := initial
	for _, v := range s {
		acc = fn(acc, v)
	}
	return acc
}

func main() {
	nums := []int{1, 2, 3, 4, 5}

	sum := reduceInts(nums, func(acc, val int) int { return acc + val }, 0)
	fmt.Println("Sum:", sum)

	product := reduceInts(nums, func(acc, val int) int { return acc * val }, 1)
	fmt.Println("Product:", product)

	max := reduceInts(nums, func(acc, val int) int {
		if val > acc {
			return val
		}
		return acc
	}, nums[0])
	fmt.Println("Max:", max)

	// Count with reduce
	count := reduceInts(nums, func(acc, _ int) int { return acc + 1 }, 0)
	fmt.Println("Count:", count)
}
