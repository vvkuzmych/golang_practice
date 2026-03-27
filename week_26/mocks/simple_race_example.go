// Простий приклад для демонстрації race condition
//
// Запустити БЕЗ race detector:
//   go run simple_race_example.go
//
// Запустити З race detector:
//   go run -race simple_race_example.go

package main

import (
	"fmt"
	"sync"
)

func main() {
	fmt.Println("╔════════════════════════════════════════╗")
	fmt.Println("║  ПРОСТИЙ ПРИКЛАД RACE CONDITION        ║")
	fmt.Println("╚════════════════════════════════════════╝")
	fmt.Println()

	// ❌ НЕБЕЗПЕЧНИЙ КОД
	counter := 0
	var wg sync.WaitGroup

	// Запускаємо 1000 goroutines
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counter++ // ❌ RACE CONDITION!
		}()
	}

	wg.Wait()

	fmt.Printf("Очікували: 1000\n")
	fmt.Printf("Отримали:  %d\n", counter)
	fmt.Println()

	if counter != 1000 {
		fmt.Printf("❌ Втрачено %d інкрементів!\n", 1000-counter)
		fmt.Println()
		fmt.Println("💡 Запусти з race detector:")
		fmt.Println("   go run -race simple_race_example.go")
	} else {
		fmt.Println("✓ Результат правильний (але тільки випадково!)")
		fmt.Println()
		fmt.Println("💡 Все одно запусти з race detector:")
		fmt.Println("   go run -race simple_race_example.go")
	}
	fmt.Println()
}
