package main

import (
	"context"
	"fmt"
	"time"
)

// 06. Timeout Pattern in Go

func main() {
	fmt.Println("=== Go Timeout with Context ===")
	fmt.Println()

	// 1. Basic timeout with context
	fmt.Println("1. Basic timeout:")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	done := make(chan string)
	go func() {
		fmt.Println("  Starting long operation...")
		time.Sleep(1 * time.Second)
		done <- "Completed"
	}()

	select {
	case result := <-done:
		fmt.Printf("  ✓ %s in time\n", result)
	case <-ctx.Done():
		fmt.Println("  ✗ Operation timed out!")
	}
	fmt.Println()

	// 2. Timeout that exceeds
	fmt.Println("2. Timeout that exceeds:")
	ctx, cancel = context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	done = make(chan string)
	go func() {
		fmt.Println("  Starting long operation...")
		time.Sleep(3 * time.Second)
		done <- "This won't be received"
	}()

	select {
	case result := <-done:
		fmt.Printf("  Result: %s\n", result)
	case <-ctx.Done():
		fmt.Println("  ✗ Timed out after 1 second!")
	}
	fmt.Println()

	// 3. Context propagation
	fmt.Println("3. Context propagation:")
	ctx, cancel = context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err := performTask(ctx)
	if err != nil {
		fmt.Printf("  ✗ Error: %v\n", err)
	} else {
		fmt.Println("  ✓ Task completed")
	}
	fmt.Println()

	// 4. Multiple operations with shared timeout
	fmt.Println("4. Multiple operations with shared timeout:")
	ctx, cancel = context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	results := make(chan string, 3)

	// Start 3 operations
	for i := 1; i <= 3; i++ {
		go func(id int) {
			select {
			case <-ctx.Done():
				results <- fmt.Sprintf("Op %d: cancelled", id)
				return
			case <-time.After(time.Duration(id) * time.Second):
				results <- fmt.Sprintf("Op %d: completed", id)
			}
		}(i)
	}

	// Collect results
	for i := 0; i < 3; i++ {
		fmt.Println(" ", <-results)
	}

	fmt.Println()
	fmt.Println("✅ Go timeout complete")
}

func performTask(ctx context.Context) error {
	// Simulate work with context awareness
	select {
	case <-time.After(500 * time.Millisecond):
		fmt.Println("  Task work completed")
		return nil
	case <-ctx.Done():
		return ctx.Err() // context.DeadlineExceeded
	}
}

// Key points:
// - context.WithTimeout for timeout
// - select with ctx.Done() for cancellation
// - Context propagates through function calls
// - cancel() to cleanup
// - More powerful than Ruby Timeout
// - Built-in cancellation support
