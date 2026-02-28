package main

import (
	"fmt"
	"sync"
	"time"
)

// ❌ НЕПРАВИЛЬНО - демонстрація бага
func MergeChannelsBuggy[T any](channels ...<-chan T) <-chan T {
	var wg sync.WaitGroup
	wg.Add(len(channels))

	outputCh := make(chan T)

	// ❌ БАГ: змінна channel перезаписується
	for channel := range channels {
		go func() {
			defer wg.Done()
			// Всі горутини бачать останнє значення channel!
			fmt.Printf("Goroutine читає з індексу: %d\n", channel)
			time.Sleep(100 * time.Millisecond) // Даємо циклу час закінчитися
		}()
	}

	go func() {
		wg.Wait()
		close(outputCh)
	}()

	return outputCh
}

// ✅ ПРАВИЛЬНО - з параметром
func MergeChannelsCorrect[T any](channels ...<-chan T) <-chan T {
	var wg sync.WaitGroup
	wg.Add(len(channels))

	outputCh := make(chan T)

	// ✅ Передаємо як параметр
	for index, channel := range channels {
		go func(idx int, ch <-chan T) {
			defer wg.Done()
			fmt.Printf("Goroutine правильно читає з індексу: %d\n", idx)
			for value := range ch {
				outputCh <- value
			}
		}(index, channel) // Явна передача параметрів
	}

	go func() {
		wg.Wait()
		close(outputCh)
	}()

	return outputCh
}

func main() {
	fmt.Println("=== Демонстрація бага ===\n")

	// Тест 1: Показати проблему з індексами
	fmt.Println("❌ НЕПРАВИЛЬНО (без параметра):")
	fmt.Println("Запускаємо 3 горутини для індексів 0, 1, 2")
	fmt.Println("Очікуємо побачити: 0, 1, 2")
	fmt.Println("Насправді бачимо:")

	channels := make([]<-chan int, 3)
	MergeChannelsBuggy(channels...)
	time.Sleep(200 * time.Millisecond)

	fmt.Println("\n✅ ПРАВИЛЬНО (з параметром):")
	fmt.Println("Запускаємо 3 горутини для індексів 0, 1, 2")
	fmt.Println("Очікуємо побачити: 0, 1, 2")
	fmt.Println("Насправді бачимо:")

	ch1 := make(chan int)
	ch2 := make(chan int)
	ch3 := make(chan int)

	go func() {
		close(ch1)
		close(ch2)
		close(ch3)
	}()

	MergeChannelsCorrect(ch1, ch2, ch3)
	time.Sleep(200 * time.Millisecond)

	fmt.Println("\n=========================")
	fmt.Println("Бачите різницю? У неправильному всі горутини бачать індекс 2!")
}
