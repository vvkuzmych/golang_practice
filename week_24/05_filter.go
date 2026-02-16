// package main
//
// import "fmt"
//
// // 05. Filter Pattern - Фільтрація даних з каналу
//
// // Паттерн: Input → Filter(predicate) → Output
// //
// // [1,2,3,4,5] → (x % 2 == 0) → [2,4]
//
//	func filter(in <-chan int, predicate func(int) bool) <-chan int {
//		out := make(chan int)
//		go func() {
//			defer close(out)
//			for val := range in {
//				if predicate(val) {
//					out <- val
//				}
//			}
//		}()
//		return out
//	}
//
//	func main() {
//		fmt.Println("=== Filter Pattern ===")
//		fmt.Println()
//
//		// Producer: числа від 1 до 20
//		numbers := make(chan int)
//		go func() {
//			defer close(numbers)
//			for i := 1; i <= 20; i++ {
//				numbers <- i
//			}
//		}()
//
//		// Filter 1: Тільки парні
//		even := filter(numbers, func(n int) bool {
//			return n%2 == 0
//		})
//
//		// Filter 2: Більше 10
//		large := filter(even, func(n int) bool {
//			return n > 10
//		})
//
//		// Виводимо результати
//		fmt.Println("Numbers (even and > 10):")
//		for n := range large {
//			fmt.Println(n)
//		}
//
//		fmt.Println()
//		fmt.Println("✅ Filters applied: even && > 10")
//
//		// Додатковий приклад: фільтр strings
//		fmt.Println()
//		fmt.Println("String filter example:")
//
//		words := make(chan string)
//		go func() {
//			defer close(words)
//			for _, w := range []string{"apple", "banana", "avocado", "cherry", "apricot"} {
//				words <- w
//			}
//		}()
//
//		// Filter: слова що починаються з 'a'
//		startsWithA := make(chan string)
//		go func() {
//			defer close(startsWithA)
//			for word := range words {
//				if len(word) > 0 && word[0] == 'a' {
//					startsWithA <- word
//				}
//			}
//		}()
//
//		fmt.Println("Words starting with 'a':")
//		for word := range startsWithA {
//			fmt.Println(" -", word)
//		}
//	}
//
// // Use cases:
// // - Data validation
// // - Search results filtering
// // - Event filtering
// // - Log filtering
// // - Stream processing
package main

import (
	"fmt"
)

func Filter[T any](inputCh <-chan T, predicate func(T) bool) <-chan T {
	outputCh := make(chan T)

	go func() {
		defer close(outputCh)
		for value := range inputCh {
			if predicate(value) {
				outputCh <- value
			}
		}
	}()

	return outputCh
}

func main() {
	channel := make(chan int)

	go func() {
		defer close(channel)
		for i := 0; i < 10; i++ {
			channel <- i
		}
	}()

	isOdd := func(value int) bool {
		return value%2 != 0
	}

	for value := range Filter(channel, isOdd) {
		fmt.Println(value)
	}
}
