package main

import (
	"fmt"
	"sync"
	"time"
)

// 01. Fan-In Pattern - Об'єднання кількох каналів в один

// Паттерн: Кілька producers → Один consumer
//
// Producer1 ──┐
// Producer2 ──┼─→ Fan-In ──→ Output
// Producer3 ──┘

func producer1(name string, count int) <-chan string {
	ch := make(chan string)
	go func() {
		defer close(ch)
		for i := 1; i <= count; i++ {
			time.Sleep(100 * time.Millisecond)
			ch <- fmt.Sprintf("%s: message %d", name, i)
		}
	}()
	return ch
}

// Fan-In: об'єднує кілька каналів в один
func fanIn(channels ...<-chan string) <-chan string {
	out := make(chan string)
	var wg sync.WaitGroup

	// Для кожного вхідного каналу
	wg.Add(len(channels))
	for _, ch := range channels {
		go func(c <-chan string) {
			defer wg.Done()
			for msg := range c {
				out <- msg
			}
		}(ch)
	}

	// Закриваємо вихідний канал після завершення всіх
	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func main() {
	//fmt.Println("=== Fan-In Pattern ===")
	//fmt.Println()
	//
	//// Створюємо 3 producers
	//p1 := producer1("Producer-1", 3)
	//p2 := producer1("Producer-2", 3)
	//p3 := producer1("Producer-3", 3)
	//
	//// Fan-In: об'єднуємо в один канал
	//merged := fanIn(p1, p2, p3)
	//
	//// Читаємо всі повідомлення
	//for msg := range merged {
	//	fmt.Println("Received:", msg)
	//}
	//
	//fmt.Println()
	//fmt.Println("✅ All messages received from 3 producers")
	p1 := make(chan string)
	p2 := make(chan string)
	p3 := make(chan string)

	go func() {
		defer func() {
			close(p1)
			close(p2)
			close(p3)
		}()

		for i := 0; i < 100; i++ {
			//p1 <- fmt.Sprintf("%s: message1 %d", fmt.Sprintf("%d", i), i)
			//p2 <- fmt.Sprintf("%s: message2 %d", fmt.Sprintf("%d", i+1), i+1)
			//p3 <- fmt.Sprintf("%s: message3 %d", fmt.Sprintf("%d", i+2), i+2)
			p1 <- fmt.Sprintf("%d", i)
			p2 <- fmt.Sprintf("%d", i+1)
			p3 <- fmt.Sprintf("%d", i+2)
		}
	}()

	for msg := range MergeChannels(p1, p2, p3) {
		fmt.Println(msg)
	}
}

// Use cases:
// - Aggregating logs from multiple services
// - Collecting metrics from multiple sources
// - Merging data streams
// - Load balancing responses

func MergeChannels[T any](channels ...<-chan T) <-chan T {
	//створюємо канал для результату
	//-----------------------------------------------
	var wg sync.WaitGroup
	wg.Add(len(channels))
	merged := make(chan T)
	//-----------------------------------------------

	//-----------------------------------------------
	//зчитуємо з вхідних каналів і записуємо в один
	for _, channel := range channels {
		go func() {
			defer wg.Done()

			for msg := range channel {
				merged <- msg
			}
		}()
	}
	//-----------------------------------------------

	//-----------------------------------------------
	go func() {
		wg.Wait()
		close(merged)
	}()
	//-----------------------------------------------
	return merged
}
