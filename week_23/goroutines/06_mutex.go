package main

import (
	"fmt"
	"sync"
)

// 06. Mutex - Синхронізація доступу до спільних даних

type Counter struct {
	mu    sync.Mutex
	value int
}

func (c *Counter) Increment() {
	c.mu.Lock()
	c.value++
	c.mu.Unlock()
}

func (c *Counter) Value() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.value
}

func main() {
	var wg sync.WaitGroup
	counter := &Counter{}

	// 1000 горутин інкрементують лічильник
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counter.Increment()
		}()
	}

	wg.Wait()
	fmt.Printf("Final count: %d\n", counter.Value())
}
