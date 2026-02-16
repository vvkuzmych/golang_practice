package main

import (
	"fmt"
	"time"
)

// 01. Basic Goroutine - Простий запуск горутини

func sayHello() {
	fmt.Println("Hello from goroutine!")
}

func main() {
	// Запускаємо функцію в окремій горутині
	go sayHello()

	// Чекаємо, щоб горутина встигла виконатись
	time.Sleep(100 * time.Millisecond)

	fmt.Println("Main function done")
}
