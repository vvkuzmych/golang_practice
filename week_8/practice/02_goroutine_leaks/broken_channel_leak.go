package main

import (
	"fmt"
	"runtime"
	"time"
)

// ❌ BROKEN: Goroutine leaks - blocked on channel receive
func leakyReceiver() {
	ch := make(chan int)

	go func() {
		val := <-ch // ❌ Blocks forever! No one sends
		fmt.Println(val)
	}()

	// Forgot to send or close channel!
}

// ❌ BROKEN: Goroutine leaks - blocked on channel send
func leakySender() {
	ch := make(chan int) // Unbuffered

	go func() {
		ch <- 42 // ❌ Blocks forever! No one receives
	}()

	// Forgot to receive!
}

// ❌ BROKEN: Multiple goroutines leak
func multipleLeaks(n int) {
	ch := make(chan int)

	for i := 0; i < n; i++ {
		go func(id int) {
			val := <-ch // ❌ All block forever!
			fmt.Printf("Goroutine %d got: %d\n", id, val)
		}(i)
	}

	// Never sends to channel!
}

func main() {
	fmt.Println("=== BROKEN CODE (has goroutine leaks) ===")
	fmt.Println()

	before := runtime.NumGoroutine()
	fmt.Printf("Goroutines before: %d\n", before)

	// Leak 1: Receiver blocked
	fmt.Println("\n1. Creating goroutine that blocks on receive...")
	leakyReceiver()
	time.Sleep(100 * time.Millisecond)

	after1 := runtime.NumGoroutine()
	fmt.Printf("Goroutines after leakyReceiver: %d (leaked: %d) ❌\n", after1, after1-before)

	// Leak 2: Sender blocked
	fmt.Println("\n2. Creating goroutine that blocks on send...")
	leakySender()
	time.Sleep(100 * time.Millisecond)

	after2 := runtime.NumGoroutine()
	fmt.Printf("Goroutines after leakySender: %d (leaked: %d) ❌\n", after2, after2-before)

	// Leak 3: Multiple goroutines
	fmt.Println("\n3. Creating 10 goroutines that all leak...")
	multipleLeaks(10)
	time.Sleep(100 * time.Millisecond)

	after3 := runtime.NumGoroutine()
	fmt.Printf("Goroutines after multipleLeaks: %d (leaked: %d) ❌\n", after3, after3-before)

	fmt.Println("\n⚠️ Memory leak! Goroutines will never exit!")
	fmt.Println("Press Ctrl+C to exit...")

	// Keep main running to see leaks
	select {}
}
