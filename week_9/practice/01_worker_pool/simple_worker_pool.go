package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Job struct {
	ID   int
	Data string
}

type Result struct {
	JobID  int
	Result string
}

// Worker function
func worker(ctx context.Context, id int, jobs <-chan Job, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()

	fmt.Printf("Worker %d: started\n", id)

	for {
		select {
		case <-ctx.Done():
			fmt.Printf("Worker %d: cancelled\n", id)
			return
		case job, ok := <-jobs:
			if !ok {
				fmt.Printf("Worker %d: no more jobs\n", id)
				return
			}

			fmt.Printf("Worker %d: processing job %d\n", id, job.ID)

			// Simulate work
			select {
			case <-ctx.Done():
				fmt.Printf("Worker %d: job %d cancelled\n", id, job.ID)
				return
			case <-time.After(100 * time.Millisecond):
				result := fmt.Sprintf("Processed: %s", job.Data)
				results <- Result{JobID: job.ID, Result: result}
				fmt.Printf("Worker %d: finished job %d\n", id, job.ID)
			}
		}
	}
}

func main() {
	fmt.Println("=== Worker Pool Pattern ===\n")

	// Example 1: Simple worker pool
	fmt.Println("--- Example 1: Simple Worker Pool ---")
	{
		numWorkers := 3
		numJobs := 10

		jobs := make(chan Job, numJobs)
		results := make(chan Result, numJobs)

		var wg sync.WaitGroup
		ctx := context.Background()

		// Start workers
		for w := 1; w <= numWorkers; w++ {
			wg.Add(1)
			go worker(ctx, w, jobs, results, &wg)
		}

		// Send jobs
		go func() {
			for j := 1; j <= numJobs; j++ {
				jobs <- Job{
					ID:   j,
					Data: fmt.Sprintf("task-%d", j),
				}
			}
			close(jobs)
		}()

		// Close results when all workers done
		go func() {
			wg.Wait()
			close(results)
		}()

		// Collect results
		fmt.Println("\nResults:")
		for result := range results {
			fmt.Printf("  Job %d: %s\n", result.JobID, result.Result)
		}

		fmt.Println()
	}

	time.Sleep(500 * time.Millisecond)

	// Example 2: Worker pool with timeout
	fmt.Println("\n--- Example 2: Worker Pool with Timeout ---")
	{
		numWorkers := 2
		numJobs := 20

		jobs := make(chan Job, numJobs)
		results := make(chan Result, numJobs)

		var wg sync.WaitGroup
		ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
		defer cancel()

		// Start workers
		for w := 1; w <= numWorkers; w++ {
			wg.Add(1)
			go worker(ctx, w, jobs, results, &wg)
		}

		// Send jobs
		go func() {
			for j := 1; j <= numJobs; j++ {
				jobs <- Job{
					ID:   j,
					Data: fmt.Sprintf("task-%d", j),
				}
			}
			close(jobs)
		}()

		// Close results when all workers done
		go func() {
			wg.Wait()
			close(results)
		}()

		// Collect results
		fmt.Println("\nResults (until timeout):")
		count := 0
		for result := range results {
			fmt.Printf("  Job %d: %s\n", result.JobID, result.Result)
			count++
		}

		fmt.Printf("\n⏱️ Completed %d jobs before timeout\n", count)
	}

	time.Sleep(500 * time.Millisecond)

	// Example 3: Dynamic cancellation
	fmt.Println("\n--- Example 3: Dynamic Cancellation ---")
	{
		numWorkers := 3
		numJobs := 15

		jobs := make(chan Job, numJobs)
		results := make(chan Result, numJobs)

		var wg sync.WaitGroup
		ctx, cancel := context.WithCancel(context.Background())

		// Start workers
		for w := 1; w <= numWorkers; w++ {
			wg.Add(1)
			go worker(ctx, w, jobs, results, &wg)
		}

		// Send jobs
		go func() {
			for j := 1; j <= numJobs; j++ {
				jobs <- Job{
					ID:   j,
					Data: fmt.Sprintf("task-%d", j),
				}
			}
			close(jobs)
		}()

		// Close results when all workers done
		go func() {
			wg.Wait()
			close(results)
		}()

		// Collect results and cancel after 5
		fmt.Println("\nResults (first 5 only):")
		count := 0
		for result := range results {
			fmt.Printf("  Job %d: %s\n", result.JobID, result.Result)
			count++
			if count == 5 {
				fmt.Println("\n🛑 Cancelling after 5 results...")
				cancel()
			}
		}

		fmt.Printf("\nCompleted %d jobs before cancellation\n", count)
	}

	fmt.Println("\n✅ All examples completed!")
}
