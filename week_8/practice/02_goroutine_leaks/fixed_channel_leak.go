package main

import (
	"context"
	"fmt"
	"runtime"
	"time"
)

// ✅ FIXED Solution 1: Close channel
func fixedReceiver1() {
	ch := make(chan int)

	go func() {
		val, ok := <-ch
		if !ok {
			fmt.Println("Channel closed, exiting")
			return
		}
		fmt.Println(val)
	}()

	close(ch) // ✅ Close channel to unblock goroutine
}

// ✅ FIXED Solution 2: Send value
func fixedReceiver2() {
	ch := make(chan int)

	go func() {
		val := <-ch
		fmt.Printf("Received: %d\n", val)
	}()

	ch <- 42 // ✅ Send value
	time.Sleep(10 * time.Millisecond)
}

// ✅ FIXED Solution 3: Buffered channel
func fixedSender() {
	ch := make(chan int, 1) // ✅ Buffered channel

	go func() {
		ch <- 42 // Won't block
		fmt.Println("Sent to buffered channel")
	}()

	time.Sleep(10 * time.Millisecond)
	// Can ignore channel if needed
}

// ✅ FIXED Solution 4: Context for cancellation
func fixedWithContext() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch := make(chan int)

	go func() {
		select {
		case val := <-ch:
			fmt.Println(val)
		case <-ctx.Done():
			fmt.Println("Context canceled, exiting")
			return // ✅ Clean exit
		}
	}()

	time.Sleep(10 * time.Millisecond)
	cancel() // ✅ Signal goroutine to exit
	time.Sleep(10 * time.Millisecond)
}

// ✅ FIXED Solution 5: Done channel pattern
func fixedWithDone() {
	ch := make(chan int)
	done := make(chan struct{})

	go func() {
		for {
			select {
			case val := <-ch:
				fmt.Println(val)
			case <-done:
				fmt.Println("Done signal received, exiting")
				return // ✅ Clean exit
			}
		}
	}()

	ch <- 42
	time.Sleep(10 * time.Millisecond)
	close(done) // ✅ Signal done
	time.Sleep(10 * time.Millisecond)
}

// ✅ FIXED Solution 6: Timeout
func fixedWithTimeout() {
	ch := make(chan int)

	go func() {
		select {
		case val := <-ch:
			fmt.Println(val)
		case <-time.After(100 * time.Millisecond):
			fmt.Println("Timeout, exiting")
			return // ✅ Exit after timeout
		}
	}()

	// Don't send anything - goroutine will exit after timeout
	time.Sleep(150 * time.Millisecond)
}

func main() {
	fmt.Println("=== FIXED CODE (no goroutine leaks) ===\n")

	before := runtime.NumGoroutine()
	fmt.Printf("Goroutines before: %d\n", before)

	// Solution 1: Close channel
	fmt.Println("\n1. Close channel to unblock receiver")
	fixedReceiver1()
	time.Sleep(50 * time.Millisecond)
	fmt.Printf("Goroutines: %d ✅\n", runtime.NumGoroutine())

	// Solution 2: Send value
	fmt.Println("\n2. Send value to channel")
	fixedReceiver2()
	time.Sleep(50 * time.Millisecond)
	fmt.Printf("Goroutines: %d ✅\n", runtime.NumGoroutine())

	// Solution 3: Buffered channel
	fmt.Println("\n3. Buffered channel for sender")
	fixedSender()
	time.Sleep(50 * time.Millisecond)
	fmt.Printf("Goroutines: %d ✅\n", runtime.NumGoroutine())

	// Solution 4: Context
	fmt.Println("\n4. Context for cancellation")
	fixedWithContext()
	time.Sleep(50 * time.Millisecond)
	fmt.Printf("Goroutines: %d ✅\n", runtime.NumGoroutine())

	// Solution 5: Done channel
	fmt.Println("\n5. Done channel pattern")
	fixedWithDone()
	time.Sleep(50 * time.Millisecond)
	fmt.Printf("Goroutines: %d ✅\n", runtime.NumGoroutine())

	// Solution 6: Timeout
	fmt.Println("\n6. Timeout to prevent forever blocking")
	fixedWithTimeout()
	time.Sleep(50 * time.Millisecond)
	fmt.Printf("Goroutines: %d ✅\n", runtime.NumGoroutine())

	after := runtime.NumGoroutine()
	fmt.Printf("\n✅ No leaks! Goroutines before: %d, after: %d\n", before, after)
}
