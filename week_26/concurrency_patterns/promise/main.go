package main

import (
	"fmt"
	"time"
)

type result[T any] struct {
	val T
	err error
}

type Promise[T any] struct {
	resultCh chan result[T]
}

func NewPromise[T any](asyncFn func() (T, error)) Promise[T] {
	promise := Promise[T]{
		resultCh: make(chan result[T]),
	}

	go func() {
		defer close(promise.resultCh)

		val, err := asyncFn()
		promise.resultCh <- result[T]{val: val, err: err}
		// can be in single goroutine
	}()

	return promise
}

func (p *Promise[T]) Then(successFn func(T), errorFn func(error)) {
	go func() {
		result := <-p.resultCh
		if result.err == nil {
			successFn(result.val)
		} else {
			errorFn(result.err)
		}
	}()
}

func main() {
	asyncJob := func() (string, error) {
		time.Sleep(time.Second)
		return "ok", nil
	}

	promise := NewPromise(asyncJob)
	promise.Then(
		func(value string) {
			fmt.Println("success", value)
		},
		func(err error) {
			fmt.Println("error", err.Error())
		},
	)

	time.Sleep(2 * time.Second)
}
