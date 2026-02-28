package main

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Task struct {
	ID int
}

type TaskResult struct {
	TaskID int
	Value  string
	Err    error
}

func processTask(ctx context.Context, task Task) TaskResult {
	// Simulate random processing time
	select {
	case <-time.After(time.Duration(rand.Intn(100)) * time.Millisecond):
	case <-ctx.Done():
		return TaskResult{TaskID: task.ID, Err: ctx.Err()}
	}

	// Simulate random errors
	if rand.Float32() < 0.2 {
		return TaskResult{
			TaskID: task.ID,
			Err:    errors.New("random processing error"),
		}
	}

	return TaskResult{
		TaskID: task.ID,
		Value:  fmt.Sprintf("Task %d completed", task.ID),
	}
}

func worker(ctx context.Context, id int, tasks <-chan Task, results chan<- TaskResult, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case task, ok := <-tasks:
			if !ok {
				return
			}
			result := processTask(ctx, task)
			select {
			case results <- result:
			case <-ctx.Done():
				return
			}
		}
	}
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	numWorkers := 5
	numTasks := 20

	tasks := make(chan Task, numTasks)
	results := make(chan TaskResult, numTasks)

	var wg sync.WaitGroup

	// Start workers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(ctx, i, tasks, results, &wg)
	}

	// Submit tasks
	go func() {
		for i := 0; i < numTasks; i++ {
			select {
			case tasks <- Task{ID: i}:
			case <-ctx.Done():
				return
			}
		}
		close(tasks)
	}()

	// Close results when workers done
	go func() {
		wg.Wait()
		close(results)
	}()

	// Process results
	var succeeded, failed int
	for result := range results {
		if result.Err != nil {
			fmt.Printf("Task %d failed: %v\n", result.TaskID, result.Err)
			failed++
		} else {
			fmt.Printf("Task %d: %s\n", result.TaskID, result.Value)
			succeeded++
		}
	}

	fmt.Printf("\nSummary: %d succeeded, %d failed\n", succeeded, failed)
}
