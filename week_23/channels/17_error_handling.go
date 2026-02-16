package main

import (
	"errors"
	"fmt"
	"time"
)

// 17. Error Handling - Обробка помилок через канали

type Result struct {
	Value int
	Err   error
}

func process(id int) Result {
	time.Sleep(100 * time.Millisecond)

	if id%3 == 0 {
		return Result{Err: errors.New("divisible by 3")}
	}

	return Result{Value: id * 2}
}

func worker(jobs <-chan int, results chan<- Result) {
	for job := range jobs {
		results <- process(job)
	}
}

func main() {
	jobs := make(chan int, 10)
	results := make(chan Result, 10)

	// Запускаємо workers
	numWorkers := 3
	for w := 0; w < numWorkers; w++ {
		go worker(jobs, results)
	}

	// Відправляємо jobs
	for j := 1; j <= 10; j++ {
		jobs <- j
	}
	close(jobs)

	// Збираємо результати
	for i := 0; i < 10; i++ {
		result := <-results
		if result.Err != nil {
			fmt.Printf("Error: %v\n", result.Err)
		} else {
			fmt.Printf("Success: %d\n", result.Value)
		}
	}
}
