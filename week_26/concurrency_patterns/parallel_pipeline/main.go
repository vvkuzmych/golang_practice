package main

import (
	"fmt"
	"sort"
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
	//channel := make(chan string)
	//
	//go func() {
	//	defer close(channel)
	//	for i := 0; i < 5; i++ {
	//		channel <- fmt.Sprintf("value %d", i)
	//	}
	//}()
	//
	//for value := range send(parse(channel), 2) {
	//	fmt.Println(value)
	//}
	fmt.Println("╔════════════════════════════════════════════════╗")
	fmt.Println("║  Спосіб 1: З індексами + сортування           ║")
	fmt.Println("╚════════════════════════════════════════════════╝")
	fmt.Println()

	channel := make(chan IndexedValue)

	// Генеруємо дані з індексами
	go func() {
		defer close(channel)
		for i := 0; i < 5; i++ {
			channel <- IndexedValue{
				Index: i,
				Value: fmt.Sprintf("value %d", i),
			}
		}
	}()

	// Збираємо всі результати
	var results []IndexedValue
	for value := range sendWithIndex(parseWithIndex(channel), 6) {
		results = append(results, value)
	}

	// Сортуємо по індексу
	sort.Slice(results, func(i, j int) bool {
		return results[i].Index < results[j].Index
	})

	// Виводимо в правильному порядку
	fmt.Println("Результат (відсортований):")
	for _, r := range results {
		fmt.Printf("[%d] %s\n", r.Index, r.Value)
	}

	fmt.Println()
}
