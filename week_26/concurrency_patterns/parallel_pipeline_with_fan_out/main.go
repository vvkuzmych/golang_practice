package main

import (
	"fmt"
	"sync"
)

func split[T any](inputCh <-chan T, n int) []<-chan T {
	outputChs := make([]chan T, n)
	for i := 0; i < n; i++ {
		outputChs[i] = make(chan T)
	}

	go func() {
		idx := 0
		for value := range inputCh {
			outputChs[idx] <- value // can be non-blocking
			idx = (idx + 1) % n
		}

		for _, ch := range outputChs {
			close(ch)
		}
	}()

	// cannot cast []chan T to []<-chan T
	resultChs := make([]<-chan T, n)
	for i := 0; i < n; i++ {
		resultChs[i] = outputChs[i]
	}

	return resultChs
}

func parse(inputCh <-chan string) <-chan string {
	outputCh := make(chan string)

	go func() {
		defer close(outputCh)
		for data := range inputCh {
			outputCh <- fmt.Sprintf("parsed - %s", data)
		}
	}()

	return outputCh
}

func send(inputCh <-chan string, n int) <-chan string {
	var wg sync.WaitGroup
	wg.Add(n)

	splittedChs := split(inputCh, n)
	outputCh := make(chan string)

	for i := 0; i < n; i++ {
		go func(idx int) {
			defer wg.Done()
			for data := range splittedChs[idx] {
				outputCh <- fmt.Sprintf("sent - %s", data)
			}
		}(i)
	}

	go func() {
		wg.Wait()
		close(outputCh)
	}()

	return outputCh
}

func main() {
	channel := make(chan string)

	go func() {
		defer close(channel)
		for i := 0; i < 5; i++ {
			channel <- "value"
		}
	}()

	for value := range send(parse(channel), 2) {
		fmt.Println(value)
	}
}
