package main

import (
	"fmt"
	"time"
)

// 08. Multiple Channels - Робота з кількома каналами

func main() {
	ch1 := make(chan int)
	ch2 := make(chan int)
	ch3 := make(chan int)

	// Producer 1
	go func() {
		for i := 0; i < 3; i++ {
			ch1 <- i
			time.Sleep(100 * time.Millisecond)
		}
		close(ch1)
	}()

	// Producer 2
	go func() {
		for i := 10; i < 13; i++ {
			ch2 <- i
			time.Sleep(150 * time.Millisecond)
		}
		close(ch2)
	}()

	// Producer 3
	go func() {
		for i := 20; i < 23; i++ {
			ch3 <- i
			time.Sleep(200 * time.Millisecond)
		}
		close(ch3)
	}()

	// Consumer
	for {
		select {
		case val, ok := <-ch1:
			if ok {
				fmt.Println("From ch1:", val)
			} else {
				ch1 = nil // Вимикаємо канал
			}
		case val, ok := <-ch2:
			if ok {
				fmt.Println("From ch2:", val)
			} else {
				ch2 = nil
			}
		case val, ok := <-ch3:
			if ok {
				fmt.Println("From ch3:", val)
			} else {
				ch3 = nil
			}
		}

		// Всі канали закриті
		if ch1 == nil && ch2 == nil && ch3 == nil {
			break
		}
	}
}
