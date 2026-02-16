package main

import "fmt"

// 11. Generator Pattern - Канал як генератор

func fibonacci() <-chan int {
	ch := make(chan int)
	go func() {
		a, b := 0, 1
		for {
			ch <- a
			a, b = b, a+b
		}
	}()
	return ch
}

func take(n int, ch <-chan int) []int {
	result := make([]int, n)
	for i := 0; i < n; i++ {
		result[i] = <-ch
	}
	return result
}

func main() {
	// Генеруємо числа Фібоначчі
	fib := fibonacci()

	// Беремо перші 10
	numbers := take(10, fib)
	fmt.Println("First 10 Fibonacci numbers:", numbers)
}
