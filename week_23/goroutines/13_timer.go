package main

import (
	"fmt"
	"time"
)

// 13. Timer - Одноразова затримка

func main() {
	// Timer спрацює через 2 секунди
	timer1 := time.NewTimer(2 * time.Second)

	fmt.Println("Waiting for timer...")
	<-timer1.C
	fmt.Println("Timer 1 fired")

	// Timer з можливістю скасування
	timer2 := time.NewTimer(time.Second)
	go func() {
		<-timer2.C
		fmt.Println("Timer 2 fired")
	}()

	// Скасовуємо timer2
	stop2 := timer2.Stop()
	if stop2 {
		fmt.Println("Timer 2 stopped")
	}

	time.Sleep(2 * time.Second)
}
