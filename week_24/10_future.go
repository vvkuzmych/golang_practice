// package main
//
// import (
//
//	"fmt"
//	"time"
//
// )
//
// // 10. Future Pattern - Асинхронне обчислення
//
// // Паттерн: Запустити обчислення і отримати результат пізніше
// //
// // Start Future → Do Work → Get Result (blocks until ready)
//
//	type Future struct {
//		result chan interface{}
//	}
//
//	func NewFuture(fn func() interface{}) *Future {
//		f := &Future{
//			result: make(chan interface{}, 1),
//		}
//
//		go func() {
//			// Виконуємо функцію асинхронно
//			res := fn()
//			f.result <- res
//		}()
//
//		return f
//	}
//
//	func (f *Future) Get() interface{} {
//		return <-f.result
//	}
//
// // Future з timeout
//
//	func (f *Future) GetWithTimeout(timeout time.Duration) (interface{}, bool) {
//		select {
//		case res := <-f.result:
//			return res, true
//		case <-time.After(timeout):
//			return nil, false
//		}
//	}
//
//	func main() {
//		fmt.Println("=== Future Pattern ===")
//		fmt.Println()
//
//		// Future 1: Складне обчислення
//		fmt.Println("Starting computation...")
//		future1 := NewFuture(func() interface{} {
//			time.Sleep(1 * time.Second)
//			sum := 0
//			for i := 1; i <= 1000; i++ {
//				sum += i
//			}
//			return sum
//		})
//
//		// Робимо інші речі поки обчислюється
//		fmt.Println("Doing other work...")
//		time.Sleep(500 * time.Millisecond)
//		fmt.Println("Other work done")
//
//		// Отримуємо результат (чекаємо якщо ще не готовий)
//		fmt.Println("Getting future result...")
//		result := future1.Get()
//		fmt.Printf("Sum 1..1000 = %v\n", result)
//
//		// Future 2: З timeout
//		fmt.Println()
//		fmt.Println("Starting slow computation...")
//		future2 := NewFuture(func() interface{} {
//			time.Sleep(3 * time.Second)
//			return "slow result"
//		})
//
//		result2, ok := future2.GetWithTimeout(1 * time.Second)
//		if ok {
//			fmt.Println("Result:", result2)
//		} else {
//			fmt.Println("Timeout! Computation too slow")
//		}
//
//		// Multiple futures
//		fmt.Println()
//		fmt.Println("Multiple futures in parallel:")
//
//		futures := []*Future{
//			NewFuture(func() interface{} {
//				time.Sleep(300 * time.Millisecond)
//				return "Result A"
//			}),
//			NewFuture(func() interface{} {
//				time.Sleep(100 * time.Millisecond)
//				return "Result B"
//			}),
//			NewFuture(func() interface{} {
//				time.Sleep(200 * time.Millisecond)
//				return "Result C"
//			}),
//		}
//
//		// Отримуємо всі результати
//		for i, f := range futures {
//			fmt.Printf("Future %d: %v\n", i+1, f.Get())
//		}
//
//		fmt.Println()
//		fmt.Println("✅ All futures resolved")
//	}
//
// // Use cases:
// // - Lazy evaluation
// // - Parallel computations
// // - Async API calls
// // - Background tasks
// // - Deferred execution
package main

import (
	"fmt"
	"time"
)

type Future[T any] struct {
	resultCh chan T
}

func NewFuture[T any](action func() T) Future[T] {
	future := Future[T]{
		resultCh: make(chan T),
	}

	go func() {
		defer close(future.resultCh)
		future.resultCh <- action()
	}()

	return future
}

func (f *Future[T]) Get() T {
	return <-f.resultCh
}

func main() {
	asyncJob := func() interface{} {
		time.Sleep(time.Second)
		return "success"
	}

	future := NewFuture(asyncJob)
	result := future.Get()
	fmt.Println(result)
}
