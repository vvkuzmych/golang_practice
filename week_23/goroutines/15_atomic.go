package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

// 15. Atomic Operations - Атомарні операції без mutex

type AtomicCounter struct {
	value int64
}

func (c *AtomicCounter) Increment() {
	atomic.AddInt64(&c.value, 1)
}

func (c *AtomicCounter) Value() int64 {
	return atomic.LoadInt64(&c.value)
}

func main() {
	var wg sync.WaitGroup
	counter := &AtomicCounter{}

	// 10000 горутин інкрементують
	for i := 0; i < 10000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counter.Increment()
		}()
	}

	wg.Wait()
	fmt.Printf("Final count: %d\n", counter.Value())
}
