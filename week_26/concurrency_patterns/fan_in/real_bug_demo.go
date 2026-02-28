package main

import (
	"fmt"
	"sync"
	"time"
)

// ❌ НЕПРАВИЛЬНО - всі горутини читають з останнього каналу
func MergeWrong(channels ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	wg.Add(len(channels))
	outputCh := make(chan int)

	// ❌ БАГ: for channel := range
	for channel := range channels {
		go func() {
			defer wg.Done()
			for value := range channels[channel] { // Всі горутини бачать останнє значення channel!
				outputCh <- value
			}
		}()
	}

	go func() {
		wg.Wait()
		close(outputCh)
	}()

	return outputCh
}

// ✅ ПРАВИЛЬНО - з параметром
func MergeCorrect(channels ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	wg.Add(len(channels))
	outputCh := make(chan int)

	// ✅ Передаємо як параметр
	for _, channel := range channels {
		go func(ch <-chan int) {
			defer wg.Done()
			for value := range ch {
				outputCh <- value
			}
		}(channel) // Явна передача
	}

	go func() {
		wg.Wait()
		close(outputCh)
	}()

	return outputCh
}

func main() {
	fmt.Println("=== Реальна демонстрація бага ===\n")

	// Створюємо 3 канали з різними даними
	ch1 := make(chan int, 5)
	ch2 := make(chan int, 5)
	ch3 := make(chan int, 5)

	// Заповнюємо канали
	for i := 0; i < 5; i++ {
		ch1 <- i * 100    // 0, 100, 200, 300, 400
		ch2 <- i*100 + 10 // 10, 110, 210, 310, 410
		ch3 <- i*100 + 20 // 20, 120, 220, 320, 420
	}
	close(ch1)
	close(ch2)
	close(ch3)

	// Тест 1: Неправильний варіант
	fmt.Println("❌ НЕПРАВИЛЬНО:")
	fmt.Println("Маємо: ch1=[0,100,200,300,400], ch2=[10,110,210,310,410], ch3=[20,120,220,320,420]")
	fmt.Println("Очікуємо побачити всі числа від усіх каналів (15 чисел)")
	fmt.Println("Реально отримуємо:")

	wrongResult := make([]int, 0)
	for value := range MergeWrong(ch1, ch2, ch3) {
		wrongResult = append(wrongResult, value)
	}
	fmt.Printf("Отримано %d чисел: %v\n", len(wrongResult), wrongResult)
	fmt.Println("❌ ПРОБЛЕМА: Ми втратили дані з перших каналів!")

	// Пересоздаємо канали для правильного тесту
	ch1 = make(chan int, 5)
	ch2 = make(chan int, 5)
	ch3 = make(chan int, 5)

	for i := 0; i < 5; i++ {
		ch1 <- i * 100
		ch2 <- i*100 + 10
		ch3 <- i*100 + 20
	}
	close(ch1)
	close(ch2)
	close(ch3)

	// Тест 2: Правильний варіант
	fmt.Println("\n✅ ПРАВИЛЬНО:")
	fmt.Println("Маємо: ch1=[0,100,200,300,400], ch2=[10,110,210,310,410], ch3=[20,120,220,320,420]")
	fmt.Println("Очікуємо побачити всі числа від усіх каналів (15 чисел)")
	fmt.Println("Реально отримуємо:")

	correctResult := make([]int, 0)
	for value := range MergeCorrect(ch1, ch2, ch3) {
		correctResult = append(correctResult, value)
	}
	fmt.Printf("Отримано %d чисел: %v\n", len(correctResult), correctResult)
	fmt.Println("✅ УСПІХ: Усі дані отримано!")

	time.Sleep(100 * time.Millisecond)
}
