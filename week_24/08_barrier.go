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
// // 08. Barrier Pattern - Синхронізація точки зустрічі
//
// // Паттерн: Всі горутини чекають одна одну
// //
// // Worker1 ──┐
// // Worker2 ──┼─→ Barrier (wait all) ──→ Continue
// // Worker3 ──┘
//
//	type Barrier struct {
//		mu      sync.Mutex
//		cond    *sync.Cond
//		count   int
//		current int
//	}
//
//	func NewBarrier(count int) *Barrier {
//		b := &Barrier{count: count}
//		b.cond = sync.NewCond(&b.mu)
//		return b
//	}
//
//	func (b *Barrier) Wait() {
//		b.mu.Lock()
//		defer b.mu.Unlock()
//
//		b.current++
//
//		if b.current == b.count {
//			// Всі дійшли до бар'єру
//			b.current = 0
//			b.cond.Broadcast() // Розбудити всіх
//		} else {
//			// Чекаємо інших
//			b.cond.Wait()
//		}
//	}
//
//	func worker(id int, barrier *Barrier, wg *sync.WaitGroup) {
//		defer wg.Done()
//
//		// Phase 1
//		fmt.Printf("Worker %d: Phase 1 started\n", id)
//		time.Sleep(time.Duration(id*100) * time.Millisecond)
//		fmt.Printf("Worker %d: Phase 1 done, waiting at barrier...\n", id)
//
//		barrier.Wait() // Чекаємо всіх
//
//		// Phase 2 (починається тільки після того як всі завершили Phase 1)
//		fmt.Printf("Worker %d: Phase 2 started\n", id)
//		time.Sleep(time.Duration(id*50) * time.Millisecond)
//		fmt.Printf("Worker %d: Phase 2 done\n", id)
//	}
//
//	func main() {
//		fmt.Println("=== Barrier Pattern ===")
//		fmt.Println()
//
//		numWorkers := 5
//		barrier := NewBarrier(numWorkers)
//		var wg sync.WaitGroup
//
//		// Запускаємо workers
//		for i := 1; i <= numWorkers; i++ {
//			wg.Add(1)
//			go worker(i, barrier, &wg)
//		}
//
//		wg.Wait()
//
//		fmt.Println()
//		fmt.Println("✅ All workers completed both phases")
//		fmt.Println("   Phase 2 started only after ALL finished Phase 1")
//	}
//
// // Use cases:
// // - Multi-phase computations
// // - Distributed algorithms
// // - Parallel testing stages
// // - Game loop synchronization
// // - Map-Reduce operations
package main

import (
	"log"
	"sync"
)

type Barrier struct {
	mutex sync.Mutex
	count int
	size  int

	beforeCh chan struct{}
	afterCh  chan struct{}
}

func NewBarrier(size int) *Barrier {
	return &Barrier{
		size:     size,
		beforeCh: make(chan struct{}, size),
		afterCh:  make(chan struct{}, size),
	}
}

func (b *Barrier) Before() {
	b.mutex.Lock()

	b.count++
	if b.count == b.size {
		for i := 0; i < b.size; i++ {
			b.beforeCh <- struct{}{}
		}
	}

	b.mutex.Unlock()
	<-b.beforeCh
}

func (b *Barrier) After() {
	b.mutex.Lock()

	b.count--
	if b.count == 0 {
		for i := 0; i < b.size; i++ {
			b.afterCh <- struct{}{}
		}
	}

	b.mutex.Unlock()
	<-b.afterCh
}

func main() {
	var wg sync.WaitGroup
	wg.Add(3)

	bootstrap := func() {
		log.Println("bootstrap")
	}

	work := func() {
		log.Println("work")
	}

	count := 3
	barrier := NewBarrier(count)
	for i := 0; i < count; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < count; j++ {
				// wait for all workers to finish previous loop
				barrier.Before()
				bootstrap()
				// wait for other workers to bootstrap
				barrier.After()
				work()
			}
		}()
	}

	wg.Wait()
}
