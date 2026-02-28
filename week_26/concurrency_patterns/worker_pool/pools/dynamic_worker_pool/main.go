package main

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type DynamicPool struct {
	jobs        chan int
	results     chan int
	workerCount int32
	maxWorkers  int32
	wg          sync.WaitGroup
	ctx         context.Context
	cancel      context.CancelFunc
}

func NewDynamicPool(maxWorkers int) *DynamicPool {
	ctx, cancel := context.WithCancel(context.Background())
	return &DynamicPool{
		jobs:       make(chan int, 100),
		results:    make(chan int, 100),
		maxWorkers: int32(maxWorkers),
		ctx:        ctx,
		cancel:     cancel,
	}
}

func (p *DynamicPool) AddWorker() bool {
	if atomic.LoadInt32(&p.workerCount) >= p.maxWorkers {
		return false
	}

	p.wg.Add(1)
	atomic.AddInt32(&p.workerCount, 1)

	go func() {
		defer p.wg.Done()
		defer atomic.AddInt32(&p.workerCount, -1)

		for {
			select {
			case <-p.ctx.Done():
				return
			case job, ok := <-p.jobs:
				if !ok {
					return
				}
				time.Sleep(100 * time.Millisecond)
				p.results <- job * 2
			}
		}
	}()

	return true
}

func (p *DynamicPool) WorkerCount() int {
	return int(atomic.LoadInt32(&p.workerCount))
}

func (p *DynamicPool) Submit(job int) {
	p.jobs <- job
}

func (p *DynamicPool) Stop() {
	close(p.jobs)
	p.cancel()
	p.wg.Wait()
	close(p.results)
}

func main() {
	pool := NewDynamicPool(10)

	// Start with 2 workers
	pool.AddWorker()
	pool.AddWorker()

	fmt.Println("Workers:", pool.WorkerCount())

	// Add more workers based on load
	go func() {
		for i := 0; i < 50; i++ {
			pool.Submit(i)

			// Scale up if needed
			if i%10 == 0 && pool.AddWorker() {
				fmt.Println("Added worker, total:", pool.WorkerCount())
			}
		}
		pool.Stop()
	}()

	for result := range pool.results {
		fmt.Println("Result:", result)
	}
}
