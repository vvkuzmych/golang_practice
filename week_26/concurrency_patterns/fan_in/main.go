package main

import (
	"fmt"
	"sync"
)

func MergeChannels[T any](channels ...<-chan T) <-chan T {
	var wg sync.WaitGroup
	wg.Add(len(channels))

	outputCh := make(chan T)
	for _, channel := range channels {
		go func() {
			defer wg.Done()
			for value := range channel {
				outputCh <- value
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
	channel1 := make(chan int)
	channel2 := make(chan int)
	channel3 := make(chan int)

	go func() {
		defer func() {
			close(channel1)
			close(channel2)
			close(channel3)
		}()

		for i := 0; i < 100; i += 3 {
			channel1 <- i
			channel2 <- i + 1
			channel3 <- i + 2
		}
	}()

	for value := range MergeChannels(channel1, channel2, channel3) {
		fmt.Println(value)
	}
}
