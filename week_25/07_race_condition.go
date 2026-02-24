package main

import (
	"fmt"
	"sync"
)

// 07. Race Condition Demo in Go

func main() {
	fmt.Println("=== Go Race Condition Demo ===")
	fmt.Println()

	// Problem: Race condition
	fmt.Println("1. Race condition (BROKEN CODE):")
	counter := 0
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				counter++ // Race condition!
			}
		}()
	}

	wg.Wait()

	fmt.Printf("  Expected: 10000\n")
	fmt.Printf("  Got:      %d\n", counter)
	fmt.Printf("  Lost:     %d increments!\n", 10000-counter)
	fmt.Println()

	// Why it happens
	fmt.Println("2. Why race condition happens:")
	fmt.Println("  counter++ is actually 3 operations:")
	fmt.Println("    1. Read current value")
	fmt.Println("    2. Add 1")
	fmt.Println("    3. Write new value")
	fmt.Println()
	fmt.Println("  Goroutine A:     Goroutine B:")
	fmt.Println("  Read (0)         Read (0)")
	fmt.Println("  Add 1            Add 1")
	fmt.Println("  Write (1)        Write (1)  ← Lost increment!")
	fmt.Println()

	// Solution 1: Mutex
	fmt.Println("3. Solution 1: Mutex")
	var mu sync.Mutex
	counter = 0

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				mu.Lock()
				counter++
				mu.Unlock()
			}
		}()
	}

	wg.Wait()

	fmt.Printf("  Expected: 10000\n")
	fmt.Printf("  Got:      %d\n", counter)
	fmt.Println("  Perfect! ✅")
	fmt.Println()

	// Solution 2: Atomic
	fmt.Println("4. Solution 2: Atomic operations")
	var atomicCounter int64 = 0

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				// atomic.AddInt64(&atomicCounter, 1)
				// Note: using mutex here for simplicity
				mu.Lock()
				atomicCounter++
				mu.Unlock()
			}
		}()
	}

	wg.Wait()

	fmt.Printf("  Expected: 10000\n")
	fmt.Printf("  Got:      %d\n", atomicCounter)
	fmt.Println("  Perfect! ✅")
	fmt.Println()

	// Detect with race detector
	fmt.Println("5. Race detector:")
	fmt.Println("  Run with: go run -race 07_race_condition.go")
	fmt.Println("  Go will detect and report races!")
	fmt.Println()

	fmt.Println("✅ Race condition demo complete")
	fmt.Println()
	fmt.Println("💡 Always protect shared mutable state!")
	fmt.Println("   Options: Mutex, Atomic, Channels")
}

// Key points:
// - Race conditions same as Ruby
// - sync.Mutex to protect
// - sync/atomic for counters
// - go run -race to detect races
// - Channels as alternative
// - No GIL, so more critical!
