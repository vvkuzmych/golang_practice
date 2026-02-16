package main

import (
	"fmt"
	"time"
)

func OrDone[T any](inputCh chan T, doneCh chan struct{}) <-chan T {
	outputCh := make(chan T)

	go func() {
		defer close(outputCh)

		for {
			select {
			case <-doneCh:
				return
			default:
			}

			select {
			case value, opened := <-inputCh:
				if !opened {
					return
				}

				outputCh <- value
			case <-doneCh:
				return
			}
		}
	}()

	return outputCh
}

func main() {
	channel := make(chan string)

	go func() {
		for {
			channel <- "test"
			time.Sleep(200 * time.Millisecond)
		}
	}()

	done := make(chan struct{})

	go func() {
		time.Sleep(time.Second)
		close(done)
	}()

	for value := range OrDone(channel, done) {
		fmt.Println(value)
	}

	/*
		for {
			select {
			case value, opened := <-inputCh:
				if !opened {
					return
				}

				// processing...
			case <-doneCh:
				return
			}
		}
	*/
}
