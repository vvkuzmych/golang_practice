// package main
//
// import (
//
//	"errors"
//	"fmt"
//	"time"
//
// )
//
// // 11. Future + Promise Together - Комбінація паттернів
//
// // Паттерн: Promise для контролю, Future для отримання
// //
// // Create Promise → Return Future → Async Work → Resolve/Reject
//
//	type Result struct {
//		Value interface{}
//		Err   error
//	}
//
//	type FuturePromise struct {
//		result chan Result
//	}
//
//	func NewFuturePromise() *FuturePromise {
//		return &FuturePromise{
//			result: make(chan Result, 1),
//		}
//	}
//
// // Promise methods (setter)
//
//	func (fp *FuturePromise) Resolve(value interface{}) {
//		fp.result <- Result{Value: value}
//	}
//
//	func (fp *FuturePromise) Reject(err error) {
//		fp.result <- Result{Err: err}
//	}
//
// // Future methods (getter)
//
//	func (fp *FuturePromise) Get() (interface{}, error) {
//		res := <-fp.result
//		return res.Value, res.Err
//	}
//
//	func (fp *FuturePromise) GetWithTimeout(timeout time.Duration) (interface{}, error, bool) {
//		select {
//		case res := <-fp.result:
//			return res.Value, res.Err, true
//		case <-time.After(timeout):
//			return nil, nil, false
//		}
//	}
//
// // Async операція яка повертає Future і використовує Promise
//
//	func fetchData(url string) *FuturePromise {
//		fp := NewFuturePromise()
//
//		go func() {
//			// Симулюємо HTTP запит
//			time.Sleep(500 * time.Millisecond)
//
//			if url == "" {
//				fp.Reject(errors.New("empty URL"))
//				return
//			}
//
//			data := fmt.Sprintf("Data from %s", url)
//			fp.Resolve(data)
//		}()
//
//		return fp
//	}
//
// // Обробка множини futures
//
//	func fetchAll(urls []string) []Result {
//		futures := make([]*FuturePromise, len(urls))
//
//		// Запускаємо всі запити паралельно
//		for i, url := range urls {
//			futures[i] = fetchData(url)
//		}
//
//		// Збираємо результати
//		results := make([]Result, len(urls))
//		for i, f := range futures {
//			value, err := f.Get()
//			results[i] = Result{Value: value, Err: err}
//		}
//
//		return results
//	}
//
//	func main() {
//		fmt.Println("=== Future + Promise Pattern ===")
//		fmt.Println()
//
//		// Single request
//		fmt.Println("Fetching data...")
//		future := fetchData("https://api.example.com/users")
//
//		fmt.Println("Doing other work...")
//		time.Sleep(200 * time.Millisecond)
//
//		value, err := future.Get()
//		if err != nil {
//			fmt.Println("Error:", err)
//		} else {
//			fmt.Println("Success:", value)
//		}
//
//		// Error case
//		fmt.Println()
//		fmt.Println("Fetching with empty URL...")
//		future2 := fetchData("")
//		value2, err2 := future2.Get()
//		if err2 != nil {
//			fmt.Println("Error:", err2)
//		} else {
//			fmt.Println("Success:", value2)
//		}
//
//		// Multiple parallel requests
//		fmt.Println()
//		fmt.Println("Fetching multiple URLs in parallel...")
//		urls := []string{
//			"https://api1.com/data",
//			"https://api2.com/data",
//			"https://api3.com/data",
//		}
//
//		results := fetchAll(urls)
//		for i, res := range results {
//			if res.Err != nil {
//				fmt.Printf("URL %d error: %v\n", i+1, res.Err)
//			} else {
//				fmt.Printf("URL %d success: %v\n", i+1, res.Value)
//			}
//		}
//
//		// With timeout
//		fmt.Println()
//		fmt.Println("Fetching with timeout...")
//		slowFuture := NewFuturePromise()
//		go func() {
//			time.Sleep(2 * time.Second)
//			slowFuture.Resolve("Slow data")
//		}()
//
//		value3, err3, ok := slowFuture.GetWithTimeout(500 * time.Millisecond)
//		if !ok {
//			fmt.Println("Timeout!")
//		} else if err3 != nil {
//			fmt.Println("Error:", err3)
//		} else {
//			fmt.Println("Success:", value3)
//		}
//
//		fmt.Println()
//		fmt.Println("✅ Future + Promise pattern complete")
//	}
//
// // Use cases:
// // - HTTP client libraries
// // - Database async queries
// // - RPC calls
// // - Background jobs
// // - Promise.all() equivalent
package main

import (
	"fmt"
	"time"
)

type Future1[T any] struct {
	resultCh <-chan T
}

func NewFuture1[T any](resultCh <-chan T) Future1[T] {
	return Future1[T]{
		resultCh: resultCh,
	}
}

func (f *Future1[T]) Get() T {
	return <-f.resultCh
}

type Promise1[T any] struct {
	resultCh chan T
}

func NewPromise1[T any]() Promise1[T] {
	return Promise1[T]{
		resultCh: make(chan T),
	}
}

func (p *Promise1[T]) Set(value T) {
	p.resultCh <- value
	close(p.resultCh)
}

func (p *Promise1[T]) GetFuture() Future1[T] {
	return NewFuture1(p.resultCh)
}

func main() {
	promise := NewPromise1[string]()

	go func() {
		time.Sleep(time.Second)
		promise.Set("agreement")
	}()

	future := promise.GetFuture()
	value := future.Get()
	fmt.Println(value)
}
