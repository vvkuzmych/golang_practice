package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func rateLimitedWorker(
	ctx context.Context,
	id int,
	jobs <-chan int,
	results chan<- int,
	limiter <-chan time.Time,
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case job, ok := <-jobs:
			if !ok {
				return
			}

			// Wait for rate limiter
			select {
			case <-limiter:
			case <-ctx.Done():
				return
			}

			// Process job
			fmt.Printf("Worker %d processing job %d at %v\n",
				id, job, time.Now().Format("15:04:05.000"))
			results <- job * 2
		}
	}
}

func main() {
	ctx := context.Background()
	numWorkers := 5
	numJobs := 10
	rateLimit := 2 * time.Second // 1 job per 2 seconds per worker

	jobs := make(chan int, numJobs)
	results := make(chan int, numJobs)
	limiter := time.Tick(rateLimit / time.Duration(numWorkers))

	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go rateLimitedWorker(ctx, i, jobs, results, limiter, &wg)
	}

	// Submit jobs
	go func() {
		for i := 0; i < numJobs; i++ {
			jobs <- i
		}
		close(jobs)
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	for result := range results {
		fmt.Println("Result:", result)
	}
}
