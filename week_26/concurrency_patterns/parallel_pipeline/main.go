package main

import (
	"fmt"
	"sync"
)

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

	outputCh := make(chan string)

	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			for data := range inputCh {
				outputCh <- fmt.Sprintf("sent - %s", data)
			}
		}()
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
