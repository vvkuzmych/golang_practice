package main

import "fmt"

// 15. Pipeline Pattern - Ланцюг обробки через канали

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

func add(in <-chan int, value int) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			out <- n + value
		}
		close(out)
	}()
	return out
}

func toString(in <-chan int) <-chan string {
	out := make(chan string)
	go func() {
		for n := range in {
			out <- fmt.Sprintf("Number: %d", n)
		}
		close(out)
	}()
	return out
}

func main() {
	// Pipeline: generate → square → add 10 → toString
	nums := generate(1, 2, 3, 4, 5)
	squared := square(nums)
	added := add(squared, 10)
	strings := toString(added)

	// Читаємо результати
	for s := range strings {
		fmt.Println(s)
	}
}
