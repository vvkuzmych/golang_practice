package main

import (
	"fmt"
	"sync"
	"time"
)

// Fan-Out: Split work across multiple workers
func fanOut(input <-chan int, numWorkers int) []<-chan int {
	outputs := make([]<-chan int, numWorkers)

	for i := 0; i < numWorkers; i++ {
		ch := make(chan int)
		outputs[i] = ch

		go func(workerID int, out chan<- int) {
			defer close(out)
			defer fmt.Printf("Worker %d: done\n", workerID)

			for val := range input {
				fmt.Printf("Worker %d: processing %d\n", workerID, val)

				// Simulate work
				time.Sleep(100 * time.Millisecond)
				result := val * val

				out <- result
				fmt.Printf("Worker %d: sent %d\n", workerID, result)
			}
		}(i+1, ch)
	}

	return outputs
}

// Fan-In: Merge results from multiple workers
func fanIn(channels ...<-chan int) <-chan int {
	out := make(chan int)
	var wg sync.WaitGroup

	for i, ch := range channels {
		wg.Add(1)
		go func(id int, c <-chan int) {
			defer wg.Done()
			for val := range c {
				fmt.Printf("Fan-In (from worker %d): received %d\n", id+1, val)
				out <- val
			}
		}(i, ch)
	}

	go func() {
		wg.Wait()
		close(out)
		fmt.Println("Fan-In: all workers done")
	}()

	return out
}

// Generator: Create input stream
func generator(nums []int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for _, num := range nums {
			out <- num
			fmt.Printf("Generator: sent %d\n", num)
		}
		fmt.Println("Generator: done")
	}()
	return out
}

func main() {
	fmt.Println("=== Fan-Out / Fan-In Pattern ===\n")

	// Example 1: Basic Fan-Out/Fan-In
	fmt.Println("--- Example 1: Basic Fan-Out/Fan-In ---")
	{
		input := generator([]int{1, 2, 3, 4, 5, 6})

		// Fan-Out to 3 workers
		fmt.Println("\nFanning out to 3 workers...")
		workers := fanOut(input, 3)

		// Fan-In results
		fmt.Println("\nFanning in results...")
		results := fanIn(workers...)

		// Collect all results
		fmt.Println("\nFinal Results:")
		collected := make([]int, 0)
		for result := range results {
			fmt.Printf("  Result: %d\n", result)
			collected = append(collected, result)
		}

		fmt.Printf("\nTotal results: %d\n", len(collected))
	}

	time.Sleep(500 * time.Millisecond)

	// Example 2: Large workload
	fmt.Println("\n--- Example 2: Large Workload ---")
	{
		// Generate 20 numbers
		nums := make([]int, 20)
		for i := range nums {
			nums[i] = i + 1
		}

		input := generator(nums)

		// Fan-Out to 5 workers
		fmt.Println("\nFanning out to 5 workers...")
		workers := fanOut(input, 5)

		// Fan-In results
		fmt.Println("\nFanning in results...")
		results := fanIn(workers...)

		// Count results
		count := 0
		sum := 0
		for result := range results {
			count++
			sum += result
		}

		fmt.Printf("\nProcessed %d results\n", count)
		fmt.Printf("Sum of squares: %d\n", sum)
	}

	fmt.Println("\n✅ All examples completed!")
}
