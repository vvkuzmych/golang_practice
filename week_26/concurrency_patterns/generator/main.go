package main

import "fmt"

func GenerateWithChannel(start, end int) <-chan int {
	outputCh := make(chan int)

	go func() {
		defer close(outputCh)
		for number := start; number <= end; number++ {
			outputCh <- number
		}
	}()

	return outputCh
}

func main() {
	for number := range GenerateWithChannel(100, 200) {
		fmt.Println(number)
	}
}
