// package main
//
// import (
//
//	"fmt"
//	"sync"
//	"time"
//
// )
//
// // 14. SingleFlight Pattern - Дедуплікація одночасних запитів
//
// // Паттерн: Багато запитів → Один реальний виклик → Спільний результат
// //
// // Request1 ──┐
// // Request2 ──┼→ SingleFlight → Single Call → Shared Result
// // Request3 ──┘
//
//	type call struct {
//		wg  sync.WaitGroup
//		val interface{}
//		err error
//	}
//
//	type SingleFlight struct {
//		mu    sync.Mutex
//		calls map[string]*call
//	}
//
//	func NewSingleFlight() *SingleFlight {
//		return &SingleFlight{
//			calls: make(map[string]*call),
//		}
//	}
//
//	func (sf *SingleFlight) Do(key string, fn func() (interface{}, error)) (interface{}, error) {
//		sf.mu.Lock()
//
//		// Перевіряємо чи є вже виклик для цього ключа
//		if c, ok := sf.calls[key]; ok {
//			sf.mu.Unlock()
//			// Чекаємо на результат існуючого виклику
//			c.wg.Wait()
//			return c.val, c.err
//		}
//
//		// Створюємо новий виклик
//		c := &call{}
//		c.wg.Add(1)
//		sf.calls[key] = c
//		sf.mu.Unlock()
//
//		// Виконуємо функцію
//		c.val, c.err = fn()
//		c.wg.Done()
//
//		// Видаляємо з map
//		sf.mu.Lock()
//		delete(sf.calls, key)
//		sf.mu.Unlock()
//
//		return c.val, c.err
//	}
//
// // Симулюємо повільну операцію (наприклад, DB запит)
//
//	func expensiveQuery(id int) (string, error) {
//		fmt.Printf("  → Executing REAL query for ID: %d\n", id)
//		time.Sleep(1 * time.Second)
//		return fmt.Sprintf("Data for ID %d", id), nil
//	}
//
//	func main() {
//		fmt.Println("=== SingleFlight Pattern ===")
//		fmt.Println()
//
//		sf := NewSingleFlight()
//		var wg sync.WaitGroup
//
//		// 10 горутин запитують ОДИН і ТОЙ ЖЕ ресурс
//		fmt.Println("10 goroutines requesting same resource (ID: 1)...")
//		fmt.Println()
//
//		for i := 1; i <= 10; i++ {
//			wg.Add(1)
//			go func(n int) {
//				defer wg.Done()
//
//				fmt.Printf("Goroutine %d: requesting...\n", n)
//
//				// SingleFlight гарантує що буде ОДИН реальний виклик
//				result, err := sf.Do("user:1", func() (interface{}, error) {
//					return expensiveQuery(1)
//				})
//
//				if err != nil {
//					fmt.Printf("Goroutine %d: error: %v\n", n, err)
//				} else {
//					fmt.Printf("Goroutine %d: got result: %v\n", n, result)
//				}
//			}(i)
//		}
//
//		wg.Wait()
//
//		fmt.Println()
//		fmt.Println("✅ Only 1 real query executed, 10 goroutines got result!")
//
//		// Різні ключі
//		fmt.Println()
//		fmt.Println("Different keys (parallel execution):")
//		fmt.Println()
//
//		wg.Add(2)
//		go func() {
//			defer wg.Done()
//			result, _ := sf.Do("user:2", func() (interface{}, error) {
//				return expensiveQuery(2)
//			})
//			fmt.Println("Got:", result)
//		}()
//
//		go func() {
//			defer wg.Done()
//			result, _ := sf.Do("user:3", func() (interface{}, error) {
//				return expensiveQuery(3)
//			})
//			fmt.Println("Got:", result)
//		}()
//
//		wg.Wait()
//
//		fmt.Println()
//		fmt.Println("✅ Different keys execute in parallel")
//	}
//
// // Use cases:
// // - Cache warming (prevent thundering herd)
// // - Database queries deduplication
// // - API call deduplication
// // - Resource loading
// // - Preventing duplicate work
package main

import (
	"fmt"
	"sync"
	"time"
)

type call struct {
	err error
	val interface{}

	done chan struct{}
}

type SingleFlight struct {
	mutex sync.Mutex
	calls map[string]*call
}

func NewSingleFlight() *SingleFlight {
	return &SingleFlight{
		calls: make(map[string]*call),
	}
}

func (s *SingleFlight) Do(key string, action func() (interface{}, error)) (interface{}, error) {
	s.mutex.Lock()
	if call, found := s.calls[key]; found {
		s.mutex.Unlock()
		return s.wait(call)
	}

	call := &call{
		done: make(chan struct{}),
	}

	s.calls[key] = call
	s.mutex.Unlock()

	go func() {
		defer func() {
			s.mutex.Lock()
			close(call.done) // can be outside or omitted
			delete(s.calls, key)
			s.mutex.Unlock()
		}()

		call.val, call.err = action()
	}()

	return s.wait(call)
}

func (s *SingleFlight) wait(call *call) (interface{}, error) {
	<-call.done
	return call.val, call.err
}

func main() {
	const inFlightRequests = 5
	var wg sync.WaitGroup
	wg.Add(inFlightRequests)

	singleFlight := NewSingleFlight()

	const key = "same_key"
	for i := 0; i < inFlightRequests; i++ {
		go func() {
			defer wg.Done()
			value, err := singleFlight.Do(key, func() (interface{}, error) {
				fmt.Println("single flight")
				time.Sleep(5 * time.Second)
				return "result", nil
			})

			fmt.Println(i, "=", value, err)
		}()
	}

	wg.Wait()
}
