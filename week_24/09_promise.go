// package main
//
// import (
//
//	"fmt"
//	"time"
//
// )
//
// // 09. Promise Pattern - Відкладений результат
//
// // Паттерн: Контейнер для майбутнього результату
// //
// // Create Promise → Start Work → Resolve/Reject → Get Result
//
//	type Promise struct {
//		result chan interface{}
//		err    chan error
//	}
//
//	func NewPromise() *Promise {
//		return &Promise{
//			result: make(chan interface{}, 1),
//			err:    make(chan error, 1),
//		}
//	}
//
//	func (p *Promise) Resolve(value interface{}) {
//		p.result <- value
//	}
//
//	func (p *Promise) Reject(err error) {
//		p.err <- err
//	}
//
//	func (p *Promise) Await() (interface{}, error) {
//		select {
//		case val := <-p.result:
//			return val, nil
//		case err := <-p.err:
//			return nil, err
//		}
//	}
//
// // Асинхронна операція
//
//	func fetchUser(id int) *Promise {
//		promise := NewPromise()
//
//		go func() {
//			// Симулюємо запит
//			time.Sleep(500 * time.Millisecond)
//
//			if id <= 0 {
//				promise.Reject(fmt.Errorf("invalid user ID: %d", id))
//				return
//			}
//
//			user := map[string]interface{}{
//				"id":   id,
//				"name": fmt.Sprintf("User%d", id),
//			}
//			promise.Resolve(user)
//		}()
//
//		return promise
//	}
//
//	func main() {
//		fmt.Println("=== Promise Pattern ===")
//		fmt.Println()
//
//		// Success case
//		fmt.Println("Fetching user 1...")
//		p1 := fetchUser(1)
//		result1, err1 := p1.Await()
//		if err1 != nil {
//			fmt.Println("Error:", err1)
//		} else {
//			fmt.Printf("Success: %v\n", result1)
//		}
//
//		// Error case
//		fmt.Println()
//		fmt.Println("Fetching user -1...")
//		p2 := fetchUser(-1)
//		result2, err2 := p2.Await()
//		if err2 != nil {
//			fmt.Println("Error:", err2)
//		} else {
//			fmt.Printf("Success: %v\n", result2)
//		}
//
//		// Multiple promises
//		fmt.Println()
//		fmt.Println("Fetching multiple users...")
//		promises := []*Promise{
//			fetchUser(2),
//			fetchUser(3),
//			fetchUser(4),
//		}
//
//		for i, p := range promises {
//			result, err := p.Await()
//			if err != nil {
//				fmt.Printf("Promise %d error: %v\n", i+1, err)
//			} else {
//				fmt.Printf("Promise %d result: %v\n", i+1, result)
//			}
//		}
//
//		fmt.Println()
//		fmt.Println("✅ All promises resolved")
//	}
//
// // Use cases:
// // - Async API calls
// // - Database queries
// // - File operations
// // - External service calls
// // - JavaScript-like async pattern
package main

import (
	"fmt"
	"time"
)

type result[T any] struct {
	val T
	err error
}

type Promise[T any] struct {
	resultCh chan result[T]
}

func NewPromise[T any](asyncFn func() (T, error)) Promise[T] {
	promise := Promise[T]{
		resultCh: make(chan result[T]),
	}

	go func() {
		defer close(promise.resultCh)

		val, err := asyncFn()
		promise.resultCh <- result[T]{val: val, err: err}
		// can be in single goroutine
	}()

	return promise
}

func (p *Promise[T]) Then(successFn func(T), errorFn func(error)) {
	go func() {
		result := <-p.resultCh
		if result.err == nil {
			successFn(result.val)
		} else {
			errorFn(result.err)
		}
	}()
}

func main() {
	asyncJob := func() (string, error) {
		time.Sleep(time.Second)
		return "ok", nil
	}

	promise := NewPromise(asyncJob)
	promise.Then(
		func(value string) {
			fmt.Println("success", value)
		},
		func(err error) {
			fmt.Println("error", err.Error())
		},
	)

	time.Sleep(2 * time.Second)
}
