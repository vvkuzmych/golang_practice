// package main
//
// import "fmt"
//
// // 04. Transform Pattern - Перетворення даних в каналі
//
// // Паттерн: Input → Transform → Output
// //
// // [1,2,3] → (x * 2) → [2,4,6]
//
//	func transform(in <-chan int, fn func(int) int) <-chan int {
//		out := make(chan int)
//		go func() {
//			defer close(out)
//			for val := range in {
//				out <- fn(val)
//			}
//		}()
//		return out
//	}
//
// // Трансформація в string
//
//	func transformToString(in <-chan int, fn func(int) string) <-chan string {
//		out := make(chan string)
//		go func() {
//			defer close(out)
//			for val := range in {
//				out <- fn(val)
//			}
//		}()
//		return out
//	}
//
//	func main() {
//		fmt.Println("=== Transform Pattern ===")
//		fmt.Println()
//
//		// Producer
//		numbers := make(chan int)
//		go func() {
//			defer close(numbers)
//			for i := 1; i <= 5; i++ {
//				numbers <- i
//			}
//		}()
//
//		// Transform 1: Помножити на 2
//		doubled := transform(numbers, func(n int) int {
//			return n * 2
//		})
//
//		// Transform 2: Додати 10
//		added := transform(doubled, func(n int) int {
//			return n + 10
//		})
//
//		// Transform 3: Конвертувати в string
//		strings := transformToString(added, func(n int) string {
//			return fmt.Sprintf("Number: %d", n)
//		})
//
//		// Виводимо результати
//		fmt.Println("Original → *2 → +10 → toString:")
//		for s := range strings {
//			fmt.Println(s)
//		}
//
//		fmt.Println()
//		fmt.Println("✅ Pipeline transformations applied")
//	}
//
// // Use cases:
// // - Data normalization
// // - Format conversion
// // - Encoding/Decoding
// // - ETL pipelines
// // - Stream processing
package main

import "fmt"

func Transform[T any](inputCh <-chan T, action func(T) T) <-chan T {
	outputCh := make(chan T)

	go func() {
		defer close(outputCh)
		for number := range inputCh {
			outputCh <- action(number)
		}
	}()

	return outputCh
}

func main() {
	channel := make(chan int)

	go func() {
		defer close(channel)
		for i := 0; i <= 5; i++ {
			channel <- i
		}
	}()

	mul := func(value int) int {
		return value * value
	}

	for number := range Transform(channel, mul) {
		fmt.Println(number)
	}
}
