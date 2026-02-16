package main

import "fmt"

// 02. Buffered Channel - Буферизований канал

func main() {
	// Канал з буфером на 3 елементи
	messages := make(chan string, 3)

	// Можемо відправити 3 значення без блокування
	messages <- "buffered"
	messages <- "channel"
	messages <- "example"

	// Читаємо всі значення
	fmt.Println(<-messages)
	fmt.Println(<-messages)
	fmt.Println(<-messages)
}
