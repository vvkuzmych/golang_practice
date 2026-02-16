package main

import (
	"fmt"
	"sync"
	"time"
)

// 13. Fan-In Pattern - Об'єднання кількох каналів

func producer(name string, delay time.Duration) <-chan string {
	ch := make(chan string)
	go func() {
		for i := 1; i <= 3; i++ {
			time.Sleep(delay)
			ch <- fmt.Sprintf("%s: message %d", name, i)
		}
		close(ch)
	}()
	return ch
}

func fanIn(channels ...<-chan string) <-chan string {
	out := make(chan string)
	var wg sync.WaitGroup

	// Для кожного вхідного каналу
	for _, ch := range channels {
		wg.Add(1)
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
	// 3 producers
	ch1 := producer("Producer1", 100*time.Millisecond)
	ch2 := producer("Producer2", 150*time.Millisecond)
	ch3 := producer("Producer3", 200*time.Millisecond)

	// Fan-In
	merged := fanIn(ch1, ch2, ch3)

	// Читаємо всі повідомлення
	for msg := range merged {
		fmt.Println(msg)
	}
}
