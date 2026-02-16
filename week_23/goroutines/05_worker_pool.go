package main

import (
	"fmt"
	"sync"
	"time"
)

// 05. Worker Pool - Пул воркерів для обробки завдань

func worker(id int, jobs <-chan int, results chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()

	for job := range jobs {
		fmt.Printf("Worker %d processing job %d\n", id, job)
		time.Sleep(100 * time.Millisecond)
		results <- job * 2
	}
}

func main() {
	numWorkers := 3
	numJobs := 10

	jobs := make(chan int, numJobs)
	results := make(chan int, numJobs)
	var wg sync.WaitGroup

	// Запускаємо workers
	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go worker(w, jobs, results, &wg)
	}

	// Надсилаємо завдання
	for j := 1; j <= numJobs; j++ {
		jobs <- j
	}
	close(jobs)

	// Чекаємо завершення
	go func() {
		wg.Wait()
		close(results)
	}()

	// Збираємо результати
	for result := range results {
		fmt.Printf("Result: %d\n", result)
	}
}
