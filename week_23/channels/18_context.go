package main

import (
	"context"
	"fmt"
	"time"
)

// 18. Context with Channels - Інтеграція context і каналів

func worker(ctx context.Context, id int, jobs <-chan int, results chan<- int) {
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("Worker %d: cancelled\n", id)
			return
		case job, ok := <-jobs:
			if !ok {
				return
			}
			fmt.Printf("Worker %d processing %d\n", id, job)
			time.Sleep(200 * time.Millisecond)

			// Перевіряємо context перед відправкою
			select {
			case <-ctx.Done():
				fmt.Printf("Worker %d: cancelled before result\n", id)
				return
			case results <- job * 2:
			}
		}
	}
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	jobs := make(chan int, 10)
	results := make(chan int, 10)

	// Запускаємо workers
	for w := 1; w <= 3; w++ {
		go worker(ctx, w, jobs, results)
	}

	// Відправляємо jobs
	go func() {
		for j := 1; j <= 10; j++ {
			jobs <- j
			time.Sleep(100 * time.Millisecond)
		}
		close(jobs)
	}()

	// Збираємо результати поки є
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Main: context cancelled")
			return
		case result, ok := <-results:
			if !ok {
				return
			}
			fmt.Printf("Result: %d\n", result)
		}
	}
}
