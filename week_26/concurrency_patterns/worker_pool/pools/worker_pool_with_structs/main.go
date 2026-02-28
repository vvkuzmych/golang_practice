package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Job struct {
	ID      int
	Payload string
}

type Result struct {
	JobID  int
	Output string
	Error  error
}

type WorkerPool struct {
	workers int
	jobs    chan Job
	results chan Result
	wg      sync.WaitGroup
}

func NewWorkerPool(workers int, jobBuffer int) *WorkerPool {
	return &WorkerPool{
		workers: workers,
		jobs:    make(chan Job, jobBuffer),
		results: make(chan Result, jobBuffer),
	}
}

func (wp *WorkerPool) Start(ctx context.Context) {
	for i := 0; i < wp.workers; i++ {
		wp.wg.Add(1)
		go wp.worker(ctx, i)
	}
}

func (wp *WorkerPool) worker(ctx context.Context, id int) {
	defer wp.wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case job, ok := <-wp.jobs:
			if !ok {
				return
			}

			// Process job
			result := wp.processJob(job)

			select {
			case wp.results <- result:
			case <-ctx.Done():
				return
			}
		}
	}
}

func (wp *WorkerPool) processJob(job Job) Result {
	time.Sleep(100 * time.Millisecond) // Simulate work
	return Result{
		JobID:  job.ID,
		Output: fmt.Sprintf("Processed: %s", job.Payload),
	}
}

func (wp *WorkerPool) Submit(job Job) {
	wp.jobs <- job
}

func (wp *WorkerPool) Results() <-chan Result {
	return wp.results
}

func (wp *WorkerPool) Close() {
	close(wp.jobs)
	wp.wg.Wait()
	close(wp.results)
}

func main() {
	ctx := context.Background()
	pool := NewWorkerPool(3, 100)
	pool.Start(ctx)

	// Submit jobs
	go func() {
		for i := 0; i < 10; i++ {
			pool.Submit(Job{ID: i, Payload: fmt.Sprintf("data-%d", i)})
		}
		pool.Close()
	}()

	// Collect results
	for result := range pool.Results() {
		fmt.Printf("Job %d: %s\n", result.JobID, result.Output)
	}
}
