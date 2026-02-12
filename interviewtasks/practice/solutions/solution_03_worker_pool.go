package main

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

type Job struct {
	ID   int
	Data string
}

type Result struct {
	JobID  int
	Output string
	Error  error
}

type ProcessFunc func(job Job) (string, error)

// WorkerPool обробляє jobs паралельно з фіксованою кількістю workers
func WorkerPool(jobs []Job, numWorkers int, process ProcessFunc) []Result {
	if len(jobs) == 0 {
		return []Result{}
	}

	// Channels
	jobChan := make(chan Job, len(jobs))
	resultChan := make(chan Result, len(jobs))

	var wg sync.WaitGroup

	// Запустити workers
	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()

			// Worker обробляє jobs з channel до закриття
			for job := range jobChan {
				fmt.Printf("Worker %d processing job %d\n", workerID, job.ID)

				output, err := process(job)

				resultChan <- Result{
					JobID:  job.ID,
					Output: output,
					Error:  err,
				}
			}

			fmt.Printf("Worker %d finished\n", workerID)
		}(w)
	}

	// Надіслати всі jobs в channel
	for _, job := range jobs {
		jobChan <- job
	}
	close(jobChan) // Закрити channel (workers завершаться після обробки)

	// Goroutine для закриття resultChan після всіх workers
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Зібрати results
	resultMap := make(map[int]Result)
	for result := range resultChan {
		resultMap[result.JobID] = result
	}

	// Відновити порядок
	results := make([]Result, len(jobs))
	for i, job := range jobs {
		results[i] = resultMap[job.ID]
	}

	return results
}

func main() {
	// Processing function: convert to uppercase
	processFunc := func(job Job) (string, error) {
		time.Sleep(100 * time.Millisecond) // Simulate work

		if job.Data == "error" {
			return "", fmt.Errorf("processing failed for job %d", job.ID)
		}

		return strings.ToUpper(job.Data), nil
	}

	// Test 1: Normal processing
	fmt.Println("=== Test 1: Normal Processing ===")
	jobs := []Job{
		{ID: 1, Data: "hello"},
		{ID: 2, Data: "world"},
		{ID: 3, Data: "golang"},
		{ID: 4, Data: "rocks"},
	}

	start := time.Now()
	results := WorkerPool(jobs, 2, processFunc)
	elapsed := time.Since(start)

	for _, r := range results {
		if r.Error != nil {
			fmt.Printf("Job %d: Error - %v\n", r.JobID, r.Error)
		} else {
			fmt.Printf("Job %d: %s\n", r.JobID, r.Output)
		}
	}
	fmt.Printf("Time: %v (with 2 workers)\n\n", elapsed)

	// Test 2: With errors
	fmt.Println("=== Test 2: With Errors ===")
	jobs = []Job{
		{ID: 1, Data: "success"},
		{ID: 2, Data: "error"},
		{ID: 3, Data: "another"},
	}

	results = WorkerPool(jobs, 2, processFunc)

	for _, r := range results {
		if r.Error != nil {
			fmt.Printf("Job %d: ❌ Error - %v\n", r.JobID, r.Error)
		} else {
			fmt.Printf("Job %d: ✅ %s\n", r.JobID, r.Output)
		}
	}

	// Test 3: Many jobs, few workers
	fmt.Println("\n=== Test 3: 10 Jobs, 3 Workers ===")
	jobs = make([]Job, 10)
	for i := range jobs {
		jobs[i] = Job{ID: i + 1, Data: fmt.Sprintf("task%d", i+1)}
	}

	start = time.Now()
	results = WorkerPool(jobs, 3, processFunc)
	elapsed = time.Since(start)

	fmt.Printf("Processed %d jobs with 3 workers in %v\n", len(results), elapsed)
}
