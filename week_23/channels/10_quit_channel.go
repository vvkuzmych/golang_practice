package main

import (
	"fmt"
	"time"
)

// 10. Quit/Done Channel - Сигнал зупинки

func worker1(quit <-chan bool) {
	for {
		select {
		case <-quit:
			fmt.Println("Worker: received quit signal")
			return
		default:
			fmt.Println("Worker: working...")
			time.Sleep(500 * time.Millisecond)
		}
	}
}

func main() {
	quit := make(chan bool)

	go worker1(quit)

	// Даємо працювати 2 секунди
	time.Sleep(2 * time.Second)

	// Відправляємо сигнал зупинки
	fmt.Println("Main: sending quit signal")
	quit <- true

	// Даємо час на завершення
	time.Sleep(100 * time.Millisecond)
	fmt.Println("Main: done")
}
