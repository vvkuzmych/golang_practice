// package main
//
// import (
//
//	"fmt"
//	"time"
//
// )
//
// // 12. Generator Pattern - Нескінченний потік даних
//
// // Паттерн: Функція повертає канал який генерує дані
// //
// // Generator() → <-chan T → Infinite stream
//
// // Простий генератор
//
//	func counter(start int) <-chan int {
//		ch := make(chan int)
//		go func() {
//			n := start
//			for {
//				ch <- n
//				n++
//			}
//		}()
//		return ch
//	}
//
// // Генератор з інтервалом
//
//	func ticker(interval time.Duration) <-chan time.Time {
//		ch := make(chan time.Time)
//		go func() {
//			for {
//				time.Sleep(interval)
//				ch <- time.Now()
//			}
//		}()
//		return ch
//	}
//
// // Генератор Fibonacci
//
//	func fibonacci() <-chan int {
//		ch := make(chan int)
//		go func() {
//			a, b := 0, 1
//			for {
//				ch <- a
//				a, b = b, a+b
//				time.Sleep(100 * time.Millisecond)
//			}
//		}()
//		return ch
//	}
//
// // Генератор з зупинкою
//
//	func generatorWithStop(start int, stop <-chan struct{}) <-chan int {
//		ch := make(chan int)
//		go func() {
//			defer close(ch)
//			n := start
//			for {
//				select {
//				case <-stop:
//					return
//				case ch <- n:
//					n++
//				}
//			}
//		}()
//		return ch
//	}
//
// // Take n values from channel
//
//	func take(n int, ch <-chan int) []int {
//		result := make([]int, n)
//		for i := 0; i < n; i++ {
//			result[i] = <-ch
//		}
//		return result
//	}
//
//	func main() {
//		fmt.Println("=== Generator Pattern ===")
//		fmt.Println()
//
//		// 1. Counter generator
//		fmt.Println("Counter from 100:")
//		cnt := counter(100)
//		for i := 0; i < 5; i++ {
//			fmt.Println(<-cnt)
//		}
//
//		// 2. Fibonacci generator
//		fmt.Println()
//		fmt.Println("First 10 Fibonacci numbers:")
//		fib := fibonacci()
//		fibNumbers := take(10, fib)
//		fmt.Println(fibNumbers)
//
//		// 3. Generator with stop
//		fmt.Println()
//		fmt.Println("Generator with stop signal:")
//		stop := make(chan struct{})
//		gen := generatorWithStop(1, stop)
//
//		go func() {
//			time.Sleep(500 * time.Millisecond)
//			close(stop)
//			fmt.Println("Stop signal sent")
//		}()
//
//		for val := range gen {
//			fmt.Println("Value:", val)
//			time.Sleep(100 * time.Millisecond)
//		}
//
//		fmt.Println()
//		fmt.Println("✅ Generator pattern complete")
//	}
//
// // Use cases:
// // - Infinite sequences
// // - Event streams
// // - ID generators
// // - Random data generation
// // - Pagination with cursors
package main

import "fmt"

func GenerateWithChannel(start, end int) <-chan int {
	outputCh := make(chan int)

	go func() {
		defer close(outputCh)
		for number := start; number <= end; number++ {
			outputCh <- number
		}
	}()

	return outputCh
}

func main() {
	for number := range GenerateWithChannel(100, 200) {
		fmt.Println(number)
	}
}
