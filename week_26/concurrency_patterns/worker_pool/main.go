package main

/*
WORKER POOL PATTERN - Step-by-Step Flow

============================================
START: Program Initialization
============================================

Step 1: CREATE WORKER POOL
   - Create WorkerPool with 5 workers
   - Initialize buffered channels (jobs: 100, results: 100)
   - WaitGroup counter starts at 0

Step 2: START WORKERS
   - Launch 5 goroutines (workers)
   - Each worker waits on jobs channel
   - WaitGroup counter = 5

============================================
PROCESSING: Concurrent Execution
============================================

Step 3: SUBMIT JOBS (Producer Goroutine)
   - Separate goroutine submits 20 jobs
   - Jobs sent to buffered jobs channel
   - If channel full (100), producer blocks (backpressure)

Step 4: WORKERS PROCESS JOBS
   - Each worker pulls job from jobs channel
   - First available worker gets the job
   - Worker processes job (simulated 100ms delay)
   - Worker sends result to results channel
   - Worker loops back for next job

Step 5: COLLECT RESULTS (Main Goroutine)
   - Main goroutine reads from results channel
   - Counts success/failure
   - Blocks until result available

============================================
END: Graceful Shutdown
============================================

Step 6: STOP SIGNAL
   - Producer goroutine finishes submitting all jobs
   - Calls pool.Stop()

Step 7: CLOSE JOBS CHANNEL
   - close(jobs) signals "no more jobs"
   - Workers finish current jobs
   - Workers exit their range loop

Step 8: WAIT FOR WORKERS
   - WaitGroup.Wait() blocks until all workers done
   - Each worker calls Done(), counter decrements
   - When counter = 0, Wait() unblocks

Step 9: CLOSE RESULTS CHANNEL
   - close(results) signals "no more results"
   - Main goroutine exits its range loop

Step 10: PRINT SUMMARY
   - Display total jobs, success count, fail count
   - Program exits

============================================
Key Concepts:
============================================
- Fixed Workers: 5 goroutines (not 20!)
- Buffered Channels: Prevent blocking, provide backpressure
- WaitGroup: Tracks worker completion
- Graceful Shutdown: Workers finish current work before stopping
- Concurrent Processing: 20 jobs processed by 5 workers (~4x faster)
*/

import (
	"fmt"
	"sync"
	"time"
)

// Job represents a unit of work to be processed
type Job struct {
	ID      int
	Message string
}

// Result represents the outcome of processing a job
type Result struct {
	JobID   int
	Success bool
	Error   error
}

// WorkerPool manages a fixed number of workers and job distribution
type WorkerPool struct {
	numWorkers int            // Number of concurrent workers
	jobs       chan Job       // Buffered channel for incoming jobs
	results    chan Result    // Buffered channel for job results
	wg         sync.WaitGroup // Tracks when all workers are done
}

// NewWorkerPool creates a new worker pool
// STEP 1: Initialize pool with buffered channels
func NewWorkerPool(numWorkers int) *WorkerPool {
	return &WorkerPool{
		numWorkers: numWorkers,
		jobs:       make(chan Job, 100),    // Buffer = 100 jobs
		results:    make(chan Result, 100), // Buffer = 100 results
	}
}

// Start launches all workers
// STEP 2: Start fixed number of worker goroutines
func (p *WorkerPool) Start() {
	for i := 0; i < p.numWorkers; i++ {
		p.wg.Add(1)    // Increment WaitGroup counter
		go p.worker(i) // Launch worker goroutine
	}
}

// worker processes jobs from the jobs channel
// STEP 4: Each worker continuously pulls jobs and processes them
func (p *WorkerPool) worker(id int) {
	defer p.wg.Done() // STEP 8: Decrement counter when worker exits

	// STEP 4a: Loop until jobs channel is closed
	for job := range p.jobs {
		fmt.Printf("Worker %d: processing job %d - %s\n", id, job.ID, job.Message)

		// STEP 4b: Process the job
		err := processJob(job)

		// STEP 4c: Send result back
		p.results <- Result{
			JobID:   job.ID,
			Success: err == nil,
			Error:   err,
		}
	}

	// STEP 7: Worker exits after jobs channel closed
	fmt.Printf("Worker %d: shutting down\n", id)
}

// Submit sends a job to the worker pool
// STEP 3: Producer sends jobs to workers via channel
func (p *WorkerPool) Submit(job Job) {
	p.jobs <- job // Blocks if channel buffer is full (backpressure)
}

// Stop gracefully shuts down the worker pool
// STEP 6-9: Shutdown sequence
func (p *WorkerPool) Stop() {
	close(p.jobs)    // STEP 7: Signal "no more jobs" to workers
	p.wg.Wait()      // STEP 8: Wait for all workers to finish
	close(p.results) // STEP 9: Signal "no more results"
}

// processJob simulates work (e.g., API call, DB query, file processing)
// STEP 4b: Actual job processing happens here
func processJob(job Job) error {
	time.Sleep(100 * time.Millisecond) // Simulate slow operation
	return nil                         // Return nil = success, error = failure
}

func main() {
	fmt.Println("Starting Worker Pool Example")
	fmt.Println("=============================")

	// ========================================
	// STEP 1: CREATE WORKER POOL
	// ========================================
	pool := NewWorkerPool(5) // 5 workers, not 20!

	// ========================================
	// STEP 2: START WORKERS
	// ========================================
	pool.Start() // Launch 5 goroutines, each waiting for jobs

	// ========================================
	// STEP 3: SUBMIT JOBS (Producer)
	// ========================================
	// Separate goroutine submits jobs concurrently
	go func() {
		for i := 1; i <= 20; i++ {
			pool.Submit(Job{
				ID:      i,
				Message: fmt.Sprintf("Message %d", i),
			})
		}
		// STEP 6: All jobs submitted, initiate shutdown
		pool.Stop()
	}()

	// ========================================
	// STEP 5: COLLECT RESULTS (Consumer)
	// ========================================
	successCount := 0
	failCount := 0

	// Main goroutine blocks here, reading results
	// Loop exits when results channel is closed
	for result := range pool.results {
		if result.Success {
			successCount++
			fmt.Printf("✓ Job %d completed successfully\n", result.JobID)
		} else {
			failCount++
			fmt.Printf("✗ Job %d failed: %v\n", result.JobID, result.Error)
		}
	}

	// ========================================
	// STEP 10: PRINT SUMMARY
	// ========================================
	fmt.Println("\n=============================")
	fmt.Printf("Total jobs: 20\n")
	fmt.Printf("Successful: %d\n", successCount)
	fmt.Printf("Failed: %d\n", failCount)
	fmt.Println("Worker Pool completed!")
}
