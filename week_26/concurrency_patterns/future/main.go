package main

import (
	"fmt"
	"time"
)

type Future[T any] struct {
	resultCh chan T
}

func NewFuture[T any](action func() T) Future[T] {
	future := Future[T]{
		resultCh: make(chan T),
	}

	go func() {
		defer close(future.resultCh)
		future.resultCh <- action()
	}()

	return future
}

func (f *Future[T]) Get() T {
	return <-f.resultCh
}

func main() {
	asyncJob := func() interface{} {
		time.Sleep(time.Second)
		return "success"
	}

	future := NewFuture(asyncJob)
	result := future.Get()
	fmt.Println(result)
}
