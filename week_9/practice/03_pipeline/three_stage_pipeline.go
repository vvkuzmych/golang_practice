package main

import (
	"context"
	"fmt"
	"time"
)

// Stage 1: Generate numbers
func generate(ctx context.Context, nums []int) <-chan int {
	out := make(chan int)

	go func() {
		defer close(out)
		defer fmt.Println("Generator: done")

		for _, n := range nums {
			select {
			case <-ctx.Done():
				fmt.Println("Generator: cancelled")
				return
			case out <- n:
				fmt.Printf("Generator: sent %d\n", n)
			}
		}
	}()

	return out
}

// Stage 2: Square numbers (with some processing time)
func square(ctx context.Context, in <-chan int) <-chan int {
	out := make(chan int)

	go func() {
		defer close(out)
		defer fmt.Println("Square: done")

		for {
			select {
			case <-ctx.Done():
				fmt.Println("Square: cancelled")
				return
			case n, ok := <-in:
				if !ok {
					return
				}

				fmt.Printf("Square: processing %d\n", n)

				// Simulate work
				select {
				case <-ctx.Done():
					fmt.Println("Square: cancelled during work")
					return
				case <-time.After(200 * time.Millisecond):
					result := n * n
					out <- result
					fmt.Printf("Square: sent %d\n", result)
				}
			}
		}
	}()

	return out
}

// Stage 3: Filter (only even numbers)
func filterEven(ctx context.Context, in <-chan int) <-chan int {
	out := make(chan int)

	go func() {
		defer close(out)
		defer fmt.Println("Filter: done")

		for {
			select {
			case <-ctx.Done():
				fmt.Println("Filter: cancelled")
				return
			case n, ok := <-in:
				if !ok {
					return
				}

				fmt.Printf("Filter: checking %d\n", n)

				if n%2 == 0 {
					out <- n
					fmt.Printf("Filter: sent %d (even)\n", n)
				} else {
					fmt.Printf("Filter: dropped %d (odd)\n", n)
				}
			}
		}
	}()

	return out
}

func main() {
	fmt.Println("=== Three-Stage Pipeline with Context ===\n")

	// Example 1: Normal completion
	fmt.Println("--- Example 1: Normal Completion ---")
	{
		ctx := context.Background()
		nums := []int{1, 2, 3, 4, 5}

		// Build pipeline
		stage1 := generate(ctx, nums)
		stage2 := square(ctx, stage1)
		stage3 := filterEven(ctx, stage2)

		// Collect results
		fmt.Println("\nResults:")
		for result := range stage3 {
			fmt.Printf("  Final result: %d\n", result)
		}
		fmt.Println()
	}

	time.Sleep(500 * time.Millisecond)

	// Example 2: Timeout cancellation
	fmt.Println("\n--- Example 2: Timeout Cancellation ---")
	{
		ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
		defer cancel()

		nums := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

		// Build pipeline
		stage1 := generate(ctx, nums)
		stage2 := square(ctx, stage1)
		stage3 := filterEven(ctx, stage2)

		// Collect results until timeout
		fmt.Println("\nResults (until timeout):")
		for result := range stage3 {
			fmt.Printf("  Final result: %d\n", result)
		}

		fmt.Println("\n⏱️ Pipeline cancelled by timeout")
	}

	time.Sleep(500 * time.Millisecond)

	// Example 3: Manual cancellation
	fmt.Println("\n--- Example 3: Manual Cancellation ---")
	{
		ctx, cancel := context.WithCancel(context.Background())

		nums := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

		// Build pipeline
		stage1 := generate(ctx, nums)
		stage2 := square(ctx, stage1)
		stage3 := filterEven(ctx, stage2)

		// Collect 2 results then cancel
		fmt.Println("\nResults (first 2 only):")
		count := 0
		for result := range stage3 {
			fmt.Printf("  Final result: %d\n", result)
			count++
			if count == 2 {
				fmt.Println("\n🛑 Manually cancelling pipeline...")
				cancel()
				// Drain remaining
				for range stage3 {
					// Discard
				}
				break
			}
		}
	}

	fmt.Println("\n✅ All examples completed!")
}
