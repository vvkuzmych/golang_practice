package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// ============= Job & Result =============

type Job struct {
	ID    int
	Value int
}

type Result struct {
	JobID  int
	Result int
	Worker int
}

// ============= Worker Pool Implementations =============

// Example 1: Basic Worker Pool
func example1_BasicWorkerPool() {
	fmt.Println("1️⃣ Basic Worker Pool")
	fmt.Println("─────────────────────────────────────────")

	const numJobs = 10
	const numWorkers = 3

	jobs := make(chan int, numJobs)
	results := make(chan int, numJobs)

	// Worker function
	worker := func(id int, jobs <-chan int, results chan<- int) {
		for job := range jobs {
			fmt.Printf("   Worker %d: processing job %d\n", id, job)
			time.Sleep(100 * time.Millisecond) // Simulate work
			results <- job * 2
		}
		fmt.Printf("   Worker %d: finished\n", id)
	}

	// Start workers
	for w := 1; w <= numWorkers; w++ {
		go worker(w, jobs, results)
	}

	// Send jobs
	fmt.Println("Sending jobs...")
	for j := 1; j <= numJobs; j++ {
		jobs <- j
	}
	close(jobs)

	// Collect results
	fmt.Println("\nCollecting results...")
	for r := 1; r <= numJobs; r++ {
		result := <-results
		fmt.Printf("   Result: %d\n", result)
	}
	fmt.Println()
}

// Example 2: Worker Pool with WaitGroup
func example2_WorkerPoolWithWaitGroup() {
	fmt.Println("2️⃣ Worker Pool with WaitGroup")
	fmt.Println("─────────────────────────────────────────")

	const numJobs = 8
	const numWorkers = 3

	jobs := make(chan Job, numJobs)
	results := make(chan Result, numJobs)
	var wg sync.WaitGroup

	// Worker function
	worker := func(id int, jobs <-chan Job, results chan<- Result, wg *sync.WaitGroup) {
		defer wg.Done()
		for job := range jobs {
			fmt.Printf("   Worker %d: job %d started\n", id, job.ID)
			time.Sleep(50 * time.Millisecond)
			results <- Result{
				JobID:  job.ID,
				Result: job.Value * job.Value,
				Worker: id,
			}
		}
	}

	// Start workers
	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go worker(w, jobs, results, &wg)
	}

	// Send jobs
	go func() {
		for j := 1; j <= numJobs; j++ {
			jobs <- Job{ID: j, Value: j}
		}
		close(jobs)
	}()

	// Close results after all workers finish
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results
	fmt.Println("Results:")
	for result := range results {
		fmt.Printf("   Job %d → %d (by worker %d)\n",
			result.JobID, result.Result, result.Worker)
	}
	fmt.Println()
}

// Example 3: Worker Pool with Error Handling
func example3_WorkerPoolWithErrors() {
	fmt.Println("3️⃣ Worker Pool with Error Handling")
	fmt.Println("─────────────────────────────────────────")

	type JobResult struct {
		JobID  int
		Value  int
		Error  error
		Worker int
	}

	const numJobs = 6
	const numWorkers = 2

	jobs := make(chan Job, numJobs)
	results := make(chan JobResult, numJobs)
	var wg sync.WaitGroup

	// Worker with error handling
	worker := func(id int, jobs <-chan Job, results chan<- JobResult, wg *sync.WaitGroup) {
		defer wg.Done()
		for job := range jobs {
			// Simulate random errors
			if rand.Intn(3) == 0 {
				results <- JobResult{
					JobID:  job.ID,
					Error:  fmt.Errorf("processing failed"),
					Worker: id,
				}
				fmt.Printf("   Worker %d: job %d ❌ FAILED\n", id, job.ID)
			} else {
				time.Sleep(50 * time.Millisecond)
				results <- JobResult{
					JobID:  job.ID,
					Value:  job.Value * 10,
					Worker: id,
				}
				fmt.Printf("   Worker %d: job %d ✅ SUCCESS\n", id, job.ID)
			}
		}
	}

	// Start workers
	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go worker(w, jobs, results, &wg)
	}

	// Send jobs
	go func() {
		for j := 1; j <= numJobs; j++ {
			jobs <- Job{ID: j, Value: j}
		}
		close(jobs)
	}()

	// Close results after all workers finish
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results
	successCount := 0
	errorCount := 0
	for result := range results {
		if result.Error != nil {
			errorCount++
		} else {
			successCount++
		}
	}

	fmt.Printf("\n✓ Success: %d, ❌ Errors: %d\n\n", successCount, errorCount)
}

// Example 4: Dynamic Worker Pool (add/remove workers)
func example4_DynamicWorkerPool() {
	fmt.Println("4️⃣ Dynamic Worker Pool")
	fmt.Println("─────────────────────────────────────────")

	jobs := make(chan int, 20)
	results := make(chan int, 20)
	var wg sync.WaitGroup

	worker := func(id int) {
		defer wg.Done()
		fmt.Printf("   Worker %d: started\n", id)
		for job := range jobs {
			time.Sleep(100 * time.Millisecond)
			results <- job * 2
		}
		fmt.Printf("   Worker %d: stopped\n", id)
	}

	// Start with 2 workers
	fmt.Println("Starting 2 workers...")
	for w := 1; w <= 2; w++ {
		wg.Add(1)
		go worker(w)
	}

	// Send first batch of jobs
	for j := 1; j <= 5; j++ {
		jobs <- j
	}

	// Add 2 more workers dynamically
	time.Sleep(150 * time.Millisecond)
	fmt.Println("\n+ Adding 2 more workers...")
	for w := 3; w <= 4; w++ {
		wg.Add(1)
		go worker(w)
	}

	// Send more jobs
	for j := 6; j <= 10; j++ {
		jobs <- j
	}

	close(jobs)

	// Collect results
	go func() {
		wg.Wait()
		close(results)
	}()

	count := 0
	for range results {
		count++
	}
	fmt.Printf("\n✓ Processed %d jobs with dynamic worker pool\n\n", count)
}

// Example 5: Worker Pool with Rate Limiting
func example5_WorkerPoolWithRateLimit() {
	fmt.Println("5️⃣ Worker Pool with Rate Limiting")
	fmt.Println("─────────────────────────────────────────")

	const numJobs = 10
	const numWorkers = 3
	const rateLimit = 2 // Max 2 concurrent operations

	jobs := make(chan int, numJobs)
	results := make(chan int, numJobs)
	semaphore := make(chan struct{}, rateLimit) // Rate limiter
	var wg sync.WaitGroup

	worker := func(id int) {
		defer wg.Done()
		for job := range jobs {
			// Acquire semaphore
			semaphore <- struct{}{}
			fmt.Printf("   Worker %d: processing job %d (rate limited)\n", id, job)
			time.Sleep(100 * time.Millisecond)
			results <- job * 2
			// Release semaphore
			<-semaphore
		}
	}

	// Start workers
	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go worker(w)
	}

	// Send jobs
	go func() {
		for j := 1; j <= numJobs; j++ {
			jobs <- j
		}
		close(jobs)
	}()

	// Close results after workers finish
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results
	count := 0
	for range results {
		count++
	}
	fmt.Printf("✓ Processed %d jobs with rate limiting (max %d concurrent)\n\n",
		count, rateLimit)
}

// Example 6: Worker Pool with Priority Queue
func example6_WorkerPoolWithPriority() {
	fmt.Println("6️⃣ Worker Pool with Priority")
	fmt.Println("─────────────────────────────────────────")

	type PriorityJob struct {
		ID       int
		Priority int // Lower = higher priority
	}

	highPriority := make(chan PriorityJob, 10)
	lowPriority := make(chan PriorityJob, 10)
	results := make(chan int, 20)
	var wg sync.WaitGroup

	worker := func(id int) {
		defer wg.Done()
		for {
			select {
			case job, ok := <-highPriority:
				if !ok {
					highPriority = nil
					continue
				}
				fmt.Printf("   Worker %d: HIGH priority job %d\n", id, job.ID)
				time.Sleep(50 * time.Millisecond)
				results <- job.ID

			case job, ok := <-lowPriority:
				if !ok {
					lowPriority = nil
					continue
				}
				fmt.Printf("   Worker %d: low priority job %d\n", id, job.ID)
				time.Sleep(50 * time.Millisecond)
				results <- job.ID
			}

			if highPriority == nil && lowPriority == nil {
				return
			}
		}
	}

	// Start workers
	for w := 1; w <= 2; w++ {
		wg.Add(1)
		go worker(w)
	}

	// Send jobs
	go func() {
		// Low priority jobs
		for j := 1; j <= 5; j++ {
			lowPriority <- PriorityJob{ID: j, Priority: 10}
		}
		close(lowPriority)

		// High priority jobs (arrive later but processed first)
		time.Sleep(100 * time.Millisecond)
		for j := 10; j <= 12; j++ {
			highPriority <- PriorityJob{ID: j, Priority: 1}
		}
		close(highPriority)
	}()

	// Close results after workers finish
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results
	count := 0
	for range results {
		count++
	}
	fmt.Printf("✓ Processed %d jobs with priority\n\n", count)
}

// ============= Main =============

func main() {
	fmt.Println("╔══════════════════════════════════════════╗")
	fmt.Println("║        Worker Pool Patterns              ║")
	fmt.Println("╚══════════════════════════════════════════╝")
	fmt.Println()

	rand.Seed(time.Now().UnixNano())

	example1_BasicWorkerPool()
	example2_WorkerPoolWithWaitGroup()
	example3_WorkerPoolWithErrors()
	example4_DynamicWorkerPool()
	example5_WorkerPoolWithRateLimit()
	example6_WorkerPoolWithPriority()

	fmt.Println("📝 Висновки:")
	fmt.Println("─────────────────────────────────────────")
	fmt.Println("✅ Worker pool обмежує concurrency")
	fmt.Println("✅ WaitGroup для синхронізації workers")
	fmt.Println("✅ Error handling через Result channels")
	fmt.Println("✅ Rate limiting через semaphore")
	fmt.Println("✅ Priority через multiple channels")
	fmt.Println("✅ Dynamic scaling можливий")
}
