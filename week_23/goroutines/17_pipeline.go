package main

import "fmt"

// 17. Pipeline Pattern - Ланцюг обробки даних

// Stage 1: Generate numbers
func generate(nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		for _, n := range nums {
			out <- n
		}
		close(out)
	}()
	return out
}

// Stage 2: Square numbers
func square(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			out <- n * n
		}
		close(out)
	}()
	return out
}

// Stage 3: Filter even numbers
func filterEven(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			if n%2 == 0 {
				out <- n
			}
		}
		close(out)
	}()
	return out
}

func main() {
	// Pipeline: generate → square → filterEven
	numbers := generate(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
	squared := square(numbers)
	filtered := filterEven(squared)

	// Виводимо результати
	for n := range filtered {
		fmt.Println(n)
	}
}
