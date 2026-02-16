package main

import (
	"fmt"
	"time"
)

// 02. Multiple Goroutines - Запуск кількох горутин

func printNumber(n int) {
	fmt.Printf("Number: %d\n", n)
}

func main() {
	// Запускаємо 5 горутин
	for i := 1; i <= 5; i++ {
		go printNumber(i)
	}

	// Чекаємо завершення
	time.Sleep(100 * time.Millisecond)
	fmt.Println("All goroutines done")
}
