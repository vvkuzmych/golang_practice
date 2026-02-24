package main

import (
	"fmt"
)

func Filter[T any](inputCh <-chan T, predicate func(T) bool) <-chan T {
	outputCh := make(chan T)

	go func() {
		defer close(outputCh)
		for value := range inputCh {
			if predicate(value) {
				outputCh <- value
			}
		}
	}()

	return outputCh
}

func main() {
	channel := make(chan int)

	go func() {
		defer close(channel)
		for i := 0; i < 10; i++ {
			channel <- i
		}
	}()

	isOdd := func(value int) bool {
		return value%2 != 0
	}

	for value := range Filter(channel, isOdd) {
		fmt.Println(value)
	}
}
