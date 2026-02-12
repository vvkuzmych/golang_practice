package main

import (
	"fmt"
	"sync"
	"time"
)

// ❌ BUGGY Counter with race condition
type UnsafeCounter struct {
	value int
}

func (c *UnsafeCounter) Increment() {
	c.value++ // Race condition!
}

func (c *UnsafeCounter) GetValue() int {
	return c.value // Race condition!
}

// ✅ SAFE Counter with Mutex
type SafeCounter struct {
	value int
	mu    sync.Mutex
}

func NewSafeCounter() *SafeCounter {
	return &SafeCounter{}
}

func (c *SafeCounter) Increment() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value++
}

func (c *SafeCounter) Decrement() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value--
}

func (c *SafeCounter) GetValue() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.value
}

func (c *SafeCounter) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value = 0
}

// ✅ BETTER: SafeCounter with RWMutex (optimized for many reads)
type SafeCounterRW struct {
	value int
	mu    sync.RWMutex
}

func NewSafeCounterRW() *SafeCounterRW {
	return &SafeCounterRW{}
}

func (c *SafeCounterRW) Increment() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value++
}

func (c *SafeCounterRW) Decrement() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value--
}

func (c *SafeCounterRW) GetValue() int {
	c.mu.RLock() // Read lock (multiple readers OK)
	defer c.mu.RUnlock()
	return c.value
}

func (c *SafeCounterRW) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value = 0
}

func main() {
	fmt.Println("=== Demonstrating Race Condition ===\n")

	// Test 1: Unsafe counter (with race condition)
	fmt.Println("Test 1: Unsafe Counter (with race condition)")
	unsafeCounter := &UnsafeCounter{}

	for i := 0; i < 1000; i++ {
		go unsafeCounter.Increment()
	}

	time.Sleep(1 * time.Second)
	unsafeValue := unsafeCounter.GetValue()
	fmt.Printf("Unsafe Counter: %d (expected 1000) ❌\n", unsafeValue)
	if unsafeValue != 1000 {
		fmt.Println("⚠️  Race condition detected! Run with 'go run -race' to see details")
	}

	// Test 2: Safe counter with Mutex
	fmt.Println("\nTest 2: Safe Counter with Mutex")
	safeCounter := NewSafeCounter()

	var wg sync.WaitGroup
	wg.Add(1000)

	for i := 0; i < 1000; i++ {
		go func() {
			defer wg.Done()
			safeCounter.Increment()
		}()
	}

	wg.Wait()
	safeValue := safeCounter.GetValue()
	fmt.Printf("Safe Counter: %d (expected 1000) ✅\n", safeValue)

	// Test 3: Concurrent increments and decrements
	fmt.Println("\nTest 3: Concurrent Increments and Decrements")
	safeCounter.Reset()

	wg.Add(2000)

	// 1000 increments
	for i := 0; i < 1000; i++ {
		go func() {
			defer wg.Done()
			safeCounter.Increment()
		}()
	}

	// 1000 decrements
	for i := 0; i < 1000; i++ {
		go func() {
			defer wg.Done()
			safeCounter.Decrement()
		}()
	}

	wg.Wait()
	finalValue := safeCounter.GetValue()
	fmt.Printf("Final value: %d (expected 0) ✅\n", finalValue)

	// Test 4: RWMutex with many readers
	fmt.Println("\nTest 4: RWMutex with Concurrent Reads and Writes")
	rwCounter := NewSafeCounterRW()

	wg.Add(10000)

	// 1000 writes
	for i := 0; i < 1000; i++ {
		go func() {
			defer wg.Done()
			rwCounter.Increment()
		}()
	}

	// 9000 reads (many more reads than writes)
	for i := 0; i < 9000; i++ {
		go func() {
			defer wg.Done()
			_ = rwCounter.GetValue()
		}()
	}

	wg.Wait()
	rwValue := rwCounter.GetValue()
	fmt.Printf("RWMutex Counter: %d (expected 1000) ✅\n", rwValue)

	// Performance comparison
	fmt.Println("\n=== Performance Comparison ===")

	// Mutex
	fmt.Println("\nWith sync.Mutex:")
	safeCounter.Reset()
	start := time.Now()

	wg.Add(10000)
	for i := 0; i < 5000; i++ {
		go func() {
			defer wg.Done()
			safeCounter.Increment()
		}()
	}
	for i := 0; i < 5000; i++ {
		go func() {
			defer wg.Done()
			_ = safeCounter.GetValue()
		}()
	}
	wg.Wait()

	mutexDuration := time.Since(start)
	fmt.Printf("Time: %v\n", mutexDuration)

	// RWMutex
	fmt.Println("\nWith sync.RWMutex:")
	rwCounter.Reset()
	start = time.Now()

	wg.Add(10000)
	for i := 0; i < 5000; i++ {
		go func() {
			defer wg.Done()
			rwCounter.Increment()
		}()
	}
	for i := 0; i < 5000; i++ {
		go func() {
			defer wg.Done()
			_ = rwCounter.GetValue()
		}()
	}
	wg.Wait()

	rwMutexDuration := time.Since(start)
	fmt.Printf("Time: %v\n", rwMutexDuration)

	fmt.Println("\n📊 Summary:")
	fmt.Printf("- Mutex:   %v\n", mutexDuration)
	fmt.Printf("- RWMutex: %v\n", rwMutexDuration)
	if rwMutexDuration < mutexDuration {
		improvement := float64(mutexDuration-rwMutexDuration) / float64(mutexDuration) * 100
		fmt.Printf("- RWMutex is %.1f%% faster (for read-heavy workloads)\n", improvement)
	}

	fmt.Println("\n💡 Tip: Run with 'go run -race solution_05_race_condition.go' to detect race conditions")
}
