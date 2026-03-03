package main

import (
	"fmt"
	"sync"
	"time"
)

// ProcessFunc is a function that processes a number
type ProcessFunc func(int) int

// ProcessSequential - processes numbers one by one (NO fan-out)
func ProcessSequential(numbers []int, process ProcessFunc) []int {
	results := make([]int, len(numbers))

	for i, num := range numbers {
		results[i] = process(num)
	}

	return results
}

// ProcessFanOut - processes numbers concurrently using FAN-OUT pattern
// Multiple workers read from the SAME input channel
func ProcessFanOut(numbers []int, process ProcessFunc, numWorkers int) []int {
	type job struct {
		index int
		value int
	}

	// Input channel for jobs
	jobsCh := make(chan job, len(numbers))

	// Pre-allocate results slice to maintain order
	results := make([]int, len(numbers))
	var mu sync.Mutex
	var wg sync.WaitGroup

	// FAN-OUT: Start multiple workers reading from SAME channel
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := range jobsCh {
				result := process(j.value)
				// Write result at correct index
				mu.Lock()
				results[j.index] = result
				mu.Unlock()
			}
		}()
	}

	// Feed jobs with their indices
	for i, num := range numbers {
		jobsCh <- job{index: i, value: num}
	}
	close(jobsCh)

	// Wait for all workers to finish
	wg.Wait()

	return results
}

func main() {
	fmt.Println("=== Fan-Out Pattern Demo ===")
	fmt.Println()

	numbers := []int{1, 2, 3, 4, 5, 6, 7, 8}

	// Simple processing function: double the number
	double := func(n int) int {
		time.Sleep(100 * time.Millisecond) // Simulate work
		return n * 2
	}

	// Example 1: Sequential Processing
	fmt.Println("1. Sequential Processing (NO Fan-Out):")
	start := time.Now()
	resultsSeq := ProcessSequential(numbers, double)
	durationSeq := time.Since(start)

	fmt.Printf("   Processed %d numbers in %v\n", len(resultsSeq), durationSeq)
	for i, result := range resultsSeq {
		fmt.Printf("   %d → %d\n", numbers[i], result)
	}
	fmt.Println()

	// Example 2: Fan-Out Processing (4 Workers)
	fmt.Println("2. Fan-Out Processing (4 Workers):")
	start = time.Now()
	resultsFanOut := ProcessFanOut(numbers, double, 4)
	durationFanOut := time.Since(start)

	fmt.Printf("   Processed %d numbers in %v\n", len(resultsFanOut), durationFanOut)
	for i, result := range resultsFanOut {
		fmt.Printf("   %d → %d\n", numbers[i], result)
	}
	fmt.Println()

	// Performance comparison
	fmt.Println("3. Performance Comparison:")
	speedup := float64(durationSeq) / float64(durationFanOut)
	fmt.Printf("   Sequential: %v\n", durationSeq)
	fmt.Printf("   Fan-Out:    %v\n", durationFanOut)
	fmt.Printf("   Speedup:    %.2fx faster\n", speedup)
	fmt.Println()

	// Example 3: Different processing function
	fmt.Println("4. Different Function (multiply by 10):")
	multiply10 := func(n int) int {
		return n * 10
	}
	resultsMult := ProcessFanOut([]int{1, 2, 3}, multiply10, 2)
	for i, result := range resultsMult {
		fmt.Printf("   %d → %d\n", i+1, result)
	}
	fmt.Println()

	fmt.Println("=== Demo Complete ===")
}
