package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// 01. Basic Goroutines in Go

func main() {
	fmt.Println("=== Go Goroutines - Basic ===")
	fmt.Println()

	// 1. Simple goroutine
	fmt.Println("1. Simple goroutine:")
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("  Hello from goroutine!")
		time.Sleep(500 * time.Millisecond)
		fmt.Println("  Goroutine finished")
	}()

	fmt.Println("  Main goroutine continues...")
	wg.Wait()
	fmt.Println()

	// 2. Goroutine with parameters
	fmt.Println("2. Goroutine with parameters:")
	wg.Add(1)
	go func(num int, text string) {
		defer wg.Done()
		fmt.Printf("  Received: num=%d, text=%s\n", num, text)
	}(5, "test")
	wg.Wait()
	fmt.Println()

	// 3. Multiple goroutines
	fmt.Println("3. Multiple goroutines:")
	wg.Add(5)
	for i := 0; i < 5; i++ {
		go func(n int) {
			defer wg.Done()
			time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
			fmt.Printf("  Goroutine %d finished\n", n)
		}(i)
	}
	wg.Wait()
	fmt.Println()

	// 4. Goroutine return value (via channel)
	fmt.Println("4. Goroutine return value:")
	resultCh := make(chan int)
	go func() {
		time.Sleep(200 * time.Millisecond)
		resultCh <- 42
	}()
	result := <-resultCh
	fmt.Printf("  Result: %d\n", result)
	fmt.Println()

	// 5. No direct status check, but can use channels
	fmt.Println("5. Communication via channels:")
	done := make(chan bool)
	go func() {
		time.Sleep(100 * time.Millisecond)
		done <- true
	}()
	<-done
	fmt.Println("  Goroutine completed (signaled via channel)")
	fmt.Println()

	fmt.Println("✅ Go goroutines complete")
}

// Key points:
// - Lightweight (millions possible)
// - No GIL - true parallelism
// - Channels for communication
// - WaitGroup to wait for completion
// - go func() { }() syntax
// - Better for CPU-bound tasks
