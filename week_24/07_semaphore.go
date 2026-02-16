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
// // 07. Semaphore Pattern - Обмеження кількості одночасних операцій
//
// // Паттерн: Обмеження доступу до ресурсу
// //
// // Task1 ─┐
// // Task2 ─┼→ Semaphore (max 3) ─→ Execute
// // Task3 ─┤     ┌─ Waiting...
// // Task4 ─┘     │
// // Task5 ───────┘
//
//	type Semaphore struct {
//		sem chan struct{}
//	}
//
//	func NewSemaphore(max int) *Semaphore {
//		return &Semaphore{
//			sem: make(chan struct{}, max),
//		}
//	}
//
//	func (s *Semaphore) Acquire() {
//		s.sem <- struct{}{}
//	}
//
//	func (s *Semaphore) Release() {
//		<-s.sem
//	}
//
//	func worker(id int, sem *Semaphore, wg *sync.WaitGroup) {
//		defer wg.Done()
//
//		fmt.Printf("Task %d waiting for semaphore...\n", id)
//		sem.Acquire()
//		defer sem.Release()
//
//		fmt.Printf("Task %d acquired semaphore, working...\n", id)
//		time.Sleep(time.Second)
//		fmt.Printf("Task %d done\n", id)
//	}
//
//	func main() {
//		fmt.Println("=== Semaphore Pattern ===")
//		fmt.Println()
//		fmt.Println("Max concurrent tasks: 3")
//		fmt.Println()
//
//		// Семафор: максимум 3 одночасні операції
//		sem := NewSemaphore(3)
//		var wg sync.WaitGroup
//
//		// Запускаємо 10 tasks
//		numTasks := 10
//		for i := 1; i <= numTasks; i++ {
//			wg.Add(1)
//			go worker(i, sem, &wg)
//		}
//
//		wg.Wait()
//
//		fmt.Println()
//		fmt.Println("✅ All 10 tasks completed with max 3 concurrent")
//	}
//
// // Альтернативна реалізація через channel
//
//	func simpleSemaphore() {
//		fmt.Println()
//		fmt.Println("=== Simple Semaphore (channel-based) ===")
//		fmt.Println()
//
//		maxConcurrent := 2
//		sem := make(chan struct{}, maxConcurrent)
//		var wg sync.WaitGroup
//
//		for i := 1; i <= 5; i++ {
//			wg.Add(1)
//			go func(id int) {
//				defer wg.Done()
//
//				// Acquire
//				sem <- struct{}{}
//				defer func() { <-sem }() // Release
//
//				fmt.Printf("Task %d running\n", id)
//				time.Sleep(500 * time.Millisecond)
//			}(i)
//		}
//
//		wg.Wait()
//		fmt.Println("✅ Simple semaphore done")
//	}
//
// // Use cases:
// // - Database connection pool
// // - API rate limiting
// // - Resource management
// // - Concurrent downloads
// // - Thread pool simulation
//
//	func init() {
//		go simpleSemaphore()
//		time.Sleep(3 * time.Second)
//	}
package main

import (
	"fmt"
	"sync"
	"time"
)

type Semaphore struct {
	tickets chan struct{}
}

func NewSemaphore(ticketsNumber int) Semaphore {
	return Semaphore{
		tickets: make(chan struct{}, ticketsNumber),
	}
}

func (s *Semaphore) Acquire() {
	s.tickets <- struct{}{}
}

func (s *Semaphore) Release() {
	<-s.tickets
}

func main() {
	wg := sync.WaitGroup{}
	wg.Add(16)

	semaphore := NewSemaphore(5)
	for i := 0; i < 16; i++ {
		semaphore.Acquire()
		go func() {
			defer func() {
				wg.Done()
				semaphore.Release()
			}()

			fmt.Println("working...")
			time.Sleep(time.Second * 2)
			fmt.Println("exiting...")
		}()
	}

	wg.Wait()
}
