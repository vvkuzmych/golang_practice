package main

import (
	"fmt"
	"sync"
	"time"
)

// ❌ BROKEN: Race condition on counter
type BrokenCounter struct {
	value int
}

func (c *BrokenCounter) Increment() {
	c.value++ // Race condition!
}

func (c *BrokenCounter) Get() int {
	return c.value // Race condition!
}

// Test broken counter
func testBrokenCounter() {
	counter := &BrokenCounter{}

	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counter.Increment()
		}()
	}

	wg.Wait()
	fmt.Printf("Broken Counter (expected 1000): %d\n", counter.Get())
}

func main() {
	fmt.Println("=== BROKEN CODE (has race conditions) ===")
	fmt.Println("Run with: go run -race broken_counter.go")
	fmt.Println()

	testBrokenCounter()

	// Small delay to see all goroutines finish
	time.Sleep(100 * time.Millisecond)
}
