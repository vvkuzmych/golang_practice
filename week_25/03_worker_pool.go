package main

import (
	"fmt"
	"sync"
	"time"
)

// 03. Worker Pool in Go

func main() {
	fmt.Println("=== Go Worker Pool ===")
	fmt.Println()

	const workerCount = 3
	tasks := make(chan int, 10)
	var wg sync.WaitGroup

	// Start workers
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go worker(i, tasks, &wg)
	}

	// Send tasks
	fmt.Println("Adding 10 tasks...")
	for i := 0; i < 10; i++ {
		tasks <- i
	}
	close(tasks) // Signal workers to stop
	fmt.Println()

	// Wait for all workers
	wg.Wait()

	fmt.Println()
	fmt.Println("✅ All workers finished")
	fmt.Println()

	// Alternative: Worker pool with results
	fmt.Println("=== Worker Pool with Results ===")
	fmt.Println()

	resultCh := make(chan int, 5)
	taskCh := make(chan int, 5)

	// Start workers
	wg.Add(3)
	for i := 0; i < 3; i++ {
		go func(id int) {
			defer wg.Done()
			for task := range taskCh {
				time.Sleep(100 * time.Millisecond)
				result := task * 2
				resultCh <- result
				fmt.Printf("Worker %d: task %d completed: %d\n", id, task, result)
			}
		}(i)
	}

	// Send tasks
	go func() {
		for i := 0; i < 5; i++ {
			taskCh <- i
		}
		close(taskCh)
	}()

	// Collect results
	go func() {
		wg.Wait()
		close(resultCh)
	}()

	// Gather results
	results := []int{}
	for result := range resultCh {
		results = append(results, result)
	}

	fmt.Println()
	fmt.Printf("Results: %v\n", results)
	fmt.Println()
	fmt.Println("✅ Worker pool with results complete")
}

func worker(id int, tasks <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()

	for task := range tasks {
		fmt.Printf("Worker %d: processing task %d\n", id, task)
		time.Sleep(200 * time.Millisecond)
		fmt.Printf("Worker %d: task %d done\n", id, task)
	}
	fmt.Printf("Worker %d: shutting down\n", id)
}

// Key points:
// - Channels for task distribution
// - close(ch) to signal completion
// - WaitGroup to wait for all workers
// - Buffered channels for backpressure
// - No manual shutdown needed (close channel)
// - Simpler than Ruby Queue
