//package main
//
//import (
//	"fmt"
//	"time"
//)
//
//// 03. Tee Pattern - Дублювання потоку даних
//
//// Паттерн: Один вхід → Кілька виходів (дублювання)
////
////            ┌→ Output1
//// Input ─────┼→ Output2
////            └→ Output3
//
//func producer(count int) <-chan int {
//	ch := make(chan int)
//	go func() {
//		defer close(ch)
//		for i := 1; i <= count; i++ {
//			time.Sleep(100 * time.Millisecond)
//			ch <- i
//		}
//	}()
//	return ch
//}
//
//// Tee: дублює дані з вхідного каналу в кілька вихідних
//func tee(in <-chan int, n int) []<-chan int {
//	outputs := make([]chan int, n)
//	outChans := make([]<-chan int, n)
//
//	// Створюємо вихідні канали
//	for i := 0; i < n; i++ {
//		outputs[i] = make(chan int)
//		outChans[i] = outputs[i]
//	}
//
//	// Читаємо з вхідного і дублюємо в усі вихідні
//	go func() {
//		defer func() {
//			for _, ch := range outputs {
//				close(ch)
//			}
//		}()
//
//		for val := range in {
//			for _, ch := range outputs {
//				ch <- val
//			}
//		}
//	}()
//
//	return outChans
//}
//
//func consumer(id int, ch <-chan int) {
//	for val := range ch {
//		fmt.Printf("Consumer %d received: %d\n", id, val)
//	}
//	fmt.Printf("Consumer %d finished\n", id)
//}
//
//func main() {
//	fmt.Println("=== Tee Pattern ===")
//	fmt.Println()
//
//	// Producer генерує дані
//	input := producer(5)
//
//	// Tee: дублюємо в 3 канали
//	outputs := tee(input, 3)
//
//	// Запускаємо 3 consumers
//	for i, ch := range outputs {
//		go consumer(i+1, ch)
//	}
//
//	// Чекаємо завершення
//	time.Sleep(2 * time.Second)
//
//	fmt.Println()
//	fmt.Println("✅ Data duplicated to 3 consumers")
//}
//
//// Use cases:
//// - Logging (console + file + network)
//// - Metrics (multiple collectors)
//// - Broadcasting events
//// - Real-time monitoring
//// - Audit trails

package main

import (
	"fmt"
	"sync"
)

func Tee[T any](inputCh <-chan T, n int) []<-chan T {
	outputChs := make([]chan T, n)
	for i := 0; i < n; i++ {
		outputChs[i] = make(chan T)
	}

	go func() {
		for value := range inputCh {
			for i := 0; i < n; i++ {
				outputChs[i] <- value // can be non-blocking
			}
		}

		for _, channel := range outputChs {
			close(channel)
		}
	}()

	// cannot cast []chan T to []<-chan T
	resultChs := make([]<-chan T, n)
	for i := 0; i < n; i++ {
		resultChs[i] = outputChs[i]
	}

	return resultChs
}

func main() {
	channel := make(chan int)

	go func() {
		defer close(channel)
		for i := 0; i < 5; i++ {
			channel <- i
		}
	}()

	wg := sync.WaitGroup{}
	wg.Add(2)

	channels := Tee(channel, 2)

	go func() {
		defer wg.Done()
		for value := range channels[0] {
			fmt.Println("ch1: ", value)
		}
	}()

	go func() {
		defer wg.Done()
		for value := range channels[1] {
			fmt.Println("ch2: ", value)
		}
	}()

	wg.Wait()
}
