package main

import (
	"context"
	"fmt"
	"sync"
)

type Pool[T any, R any] struct {
	workers int
	jobs    chan T
	results chan R
	process func(T) R
	wg      sync.WaitGroup
}

func NewPool[T any, R any](workers int, process func(T) R) *Pool[T, R] {
	return &Pool[T, R]{
		workers: workers,
		jobs:    make(chan T, workers*2),
		results: make(chan R, workers*2),
		process: process,
	}
}

func (p *Pool[T, R]) Start(ctx context.Context) {
	for i := 0; i < p.workers; i++ {
		p.wg.Add(1)
		go func() {
			defer p.wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case job, ok := <-p.jobs:
					if !ok {
						return
					}
					result := p.process(job)
					select {
					case p.results <- result:
					case <-ctx.Done():
						return
					}
				}
			}
		}()
	}
}

func (p *Pool[T, R]) Submit(job T) {
	p.jobs <- job
}

func (p *Pool[T, R]) Results() <-chan R {
	return p.results
}

func (p *Pool[T, R]) Close() {
	close(p.jobs)
	p.wg.Wait()
	close(p.results)
}

func main() {
	// String processing pool
	pool := NewPool[string, int](3, func(s string) int {
		return len(s)
	})

	pool.Start(context.Background())

	go func() {
		words := []string{"hello", "world", "foo", "bar", "baz"}
		for _, w := range words {
			pool.Submit(w)
		}
		pool.Close()
	}()

	for length := range pool.Results() {
		fmt.Println("Length:", length)
	}
}
