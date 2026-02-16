//package main
//
//import "fmt"
//
//// 06. Pipeline Pattern - Послідовна обробка даних
//
//// Паттерн: Stage1 → Stage2 → Stage3 → ... → Result
////
//// Generate → Square → Filter → Format → Print
//
//// Stage 1: Generate numbers
//func generate(nums ...int) <-chan int {
//	out := make(chan int)
//	go func() {
//		defer close(out)
//		for _, n := range nums {
//			out <- n
//		}
//	}()
//	return out
//}
//
//// Stage 2: Square numbers
//func square(in <-chan int) <-chan int {
//	out := make(chan int)
//	go func() {
//		defer close(out)
//		for n := range in {
//			out <- n * n
//		}
//	}()
//	return out
//}
//
//// Stage 3: Filter (тільки парні)
//func filterEven(in <-chan int) <-chan int {
//	out := make(chan int)
//	go func() {
//		defer close(out)
//		for n := range in {
//			if n%2 == 0 {
//				out <- n
//			}
//		}
//	}()
//	return out
//}
//
//// Stage 4: Format to string
//func format(in <-chan int) <-chan string {
//	out := make(chan string)
//	go func() {
//		defer close(out)
//		for n := range in {
//			out <- fmt.Sprintf("Result: %d", n)
//		}
//	}()
//	return out
//}
//
//// Stage 5: Add prefix
//func addPrefix(in <-chan string, prefix string) <-chan string {
//	out := make(chan string)
//	go func() {
//		defer close(out)
//		for s := range in {
//			out <- prefix + s
//		}
//	}()
//	return out
//}
//
//func main() {
//	fmt.Println("=== Pipeline Pattern ===")
//	fmt.Println()
//
//	// Pipeline: Generate → Square → FilterEven → Format → AddPrefix
//	numbers := generate(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
//	squared := square(numbers)
//	filtered := filterEven(squared)
//	formatted := format(filtered)
//	final := addPrefix(formatted, "✅ ")
//
//	// Виводимо результати
//	fmt.Println("Pipeline: numbers → square → filter(even) → format → prefix")
//	fmt.Println()
//	for result := range final {
//		fmt.Println(result)
//	}
//
//	fmt.Println()
//	fmt.Println("Pipeline completed!")
//
//	// Compact pipeline example
//	fmt.Println()
//	fmt.Println("Compact version:")
//
//	compact := addPrefix(
//		format(
//			filterEven(
//				square(
//					generate(11, 12, 13, 14, 15),
//				),
//			),
//		),
//		"→ ",
//	)
//
//	for result := range compact {
//		fmt.Println(result)
//	}
//}
//
//// Use cases:
//// - ETL (Extract, Transform, Load)
//// - Image processing pipeline
//// - Data processing streams
//// - Request/Response transformation
//// - Video encoding stages

package main

import (
	"fmt"
	"sync"
)

func parse(inputCh <-chan string) <-chan string {
	outputCh := make(chan string)

	go func() {
		defer close(outputCh)
		for data := range inputCh {
			outputCh <- fmt.Sprintf("parsed - %s", data)
		}
	}()

	return outputCh
}

func send(inputCh <-chan string, n int) <-chan string {
	var wg sync.WaitGroup
	wg.Add(n)

	outputCh := make(chan string)

	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			for data := range inputCh {
				outputCh <- fmt.Sprintf("sent - %s", data)
			}
		}()
	}

	go func() {
		wg.Wait()
		close(outputCh)
	}()

	return outputCh
}

func main() {
	channel := make(chan string)

	go func() {
		defer close(channel)
		for i := 0; i < 5; i++ {
			channel <- "value"
		}
	}()

	for value := range send(parse(channel), 2) {
		fmt.Println(value)
	}
}
