package main

import (
	"fmt"
	"sync"
	"time"
)

// 16. Bounded Parallelism - Обмеження паралелізму

func worker(id int, jobs <-chan int, results chan<- int) {
	for j := range jobs {
		fmt.Printf("Worker %d started job %d\n", id, j)
		time.Sleep(time.Second)
		fmt.Printf("Worker %d finished job %d\n", id, j)
		results <- j * 2
	}
}

func main() {
	const numJobs = 10
	const numWorkers = 3

	jobs := make(chan int, numJobs)
	results := make(chan int, numJobs)

	// Запускаємо обмежену кількість workers
	var wg sync.WaitGroup
	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			worker(id, jobs, results)
		}(w)
	}

	// Відправляємо jobs
	for j := 1; j <= numJobs; j++ {
		jobs <- j
	}
	close(jobs)

	// Закриваємо results після завершення
	go func() {
		wg.Wait()
		close(results)
	}()

	// Збираємо результати
	for r := range results {
		fmt.Println("Result:", r)
	}
}
