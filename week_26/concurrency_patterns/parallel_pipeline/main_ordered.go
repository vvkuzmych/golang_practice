package main

import (
	"fmt"
	"sort"
	"sync"
)

// ═══════════════════════════════════════════════════════════════════
// СПОСІБ 1: Додаємо індекс до кожного елементу
// ═══════════════════════════════════════════════════════════════════

type IndexedValue struct {
	Index int
	Value string
}

func parseWithIndex(inputCh <-chan IndexedValue) <-chan IndexedValue {
	outputCh := make(chan IndexedValue)

	go func() {
		defer close(outputCh)
		for item := range inputCh {
			outputCh <- IndexedValue{
				Index: item.Index,
				Value: fmt.Sprintf("parsed - %s", item.Value),
			}
		}
	}()

	return outputCh
}

func sendWithIndex(inputCh <-chan IndexedValue, n int) <-chan IndexedValue {
	var wg sync.WaitGroup
	wg.Add(n)

	outputCh := make(chan IndexedValue)

	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			for item := range inputCh {
				outputCh <- IndexedValue{
					Index: item.Index,
					Value: fmt.Sprintf("sent - %s", item.Value),
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(outputCh)
	}()

	return outputCh
}

func mainWithIndex() {
	fmt.Println("╔════════════════════════════════════════════════╗")
	fmt.Println("║  Спосіб 1: З індексами + сортування           ║")
	fmt.Println("╚════════════════════════════════════════════════╝")
	fmt.Println()

	channel := make(chan IndexedValue)

	// Генеруємо дані з індексами
	go func() {
		defer close(channel)
		for i := 0; i < 5; i++ {
			channel <- IndexedValue{
				Index: i,
				Value: fmt.Sprintf("value %d", i),
			}
		}
	}()

	// Збираємо всі результати
	var results []IndexedValue
	for value := range sendWithIndex(parseWithIndex(channel), 2) {
		results = append(results, value)
	}

	// Сортуємо по індексу
	sort.Slice(results, func(i, j int) bool {
		return results[i].Index < results[j].Index
	})

	// Виводимо в правильному порядку
	fmt.Println("Результат (відсортований):")
	for _, r := range results {
		fmt.Printf("[%d] %s\n", r.Index, r.Value)
	}

	fmt.Println()
}

//
//// ═══════════════════════════════════════════════════════════════════
//// СПОСІБ 2: Ordered Results з масивом
//// ═══════════════════════════════════════════════════════════════════
//
//func sendOrdered(inputCh <-chan IndexedValue, n int, totalItems int) []string {
//	// Створюємо масив для результатів
//	results := make([]string, totalItems)
//	var mu sync.Mutex
//	var wg sync.WaitGroup
//	wg.Add(n)
//
//	// Workers обробляють паралельно
//	for i := 0; i < n; i++ {
//		go func() {
//			defer wg.Done()
//			for item := range inputCh {
//				processed := fmt.Sprintf("sent - %s", item.Value)
//
//				// Зберігаємо на правильній позиції
//				mu.Lock()
//				results[item.Index] = processed
//				mu.Unlock()
//			}
//		}()
//	}
//
//	wg.Wait()
//	return results
//}
//
//func mainOrderedResults() {
//	fmt.Println("╔════════════════════════════════════════════════╗")
//	fmt.Println("║  Спосіб 2: Ordered Results з масивом          ║")
//	fmt.Println("╚════════════════════════════════════════════════╝")
//	fmt.Println()
//
//	channel := make(chan IndexedValue)
//
//	go func() {
//		defer close(channel)
//		for i := 0; i < 5; i++ {
//			channel <- IndexedValue{
//				Index: i,
//				Value: fmt.Sprintf("parsed - value %d", i),
//			}
//		}
//	}()
//
//	// Отримуємо відсортовані результати
//	results := sendOrdered(channel, 2, 5)
//
//	fmt.Println("Результат (гарантовано по порядку):")
//	for i, r := range results {
//		fmt.Printf("[%d] %s\n", i, r)
//	}
//
//	fmt.Println()
//}
//
//// ═══════════════════════════════════════════════════════════════════
//// СПОСІБ 3: Sequential ordering channel
//// ═══════════════════════════════════════════════════════════════════
//
//func sendWithOrdering(inputCh <-chan IndexedValue, n int) <-chan string {
//	resultsCh := make(chan IndexedValue, 100)
//	var wg sync.WaitGroup
//	wg.Add(n)
//
//	// Workers надсилають в resultsCh (не впорядковано)
//	for i := 0; i < n; i++ {
//		go func() {
//			defer wg.Done()
//			for item := range inputCh {
//				resultsCh <- IndexedValue{
//					Index: item.Index,
//					Value: fmt.Sprintf("sent - %s", item.Value),
//				}
//			}
//		}()
//	}
//
//	// Закриваємо після всіх workers
//	go func() {
//		wg.Wait()
//		close(resultsCh)
//	}()
//
//	// Orderer goroutine - впорядковує результати
//	orderedCh := make(chan string)
//	go func() {
//		defer close(orderedCh)
//
//		buffer := make(map[int]string)
//		nextIndex := 0
//
//		for item := range resultsCh {
//			buffer[item.Index] = item.Value
//
//			// Відправляємо всі послідовні елементи
//			for {
//				if val, exists := buffer[nextIndex]; exists {
//					orderedCh <- val
//					delete(buffer, nextIndex)
//					nextIndex++
//				} else {
//					break
//				}
//			}
//		}
//
//		// Відправляємо залишки (якщо є)
//		for i := nextIndex; ; i++ {
//			if val, exists := buffer[i]; exists {
//				orderedCh <- val
//				delete(buffer, i)
//			} else {
//				break
//			}
//		}
//	}()
//
//	return orderedCh
//}
//
//func mainSequentialOrdering() {
//	fmt.Println("╔════════════════════════════════════════════════╗")
//	fmt.Println("║  Спосіб 3: Sequential Ordering Channel        ║")
//	fmt.Println("╚════════════════════════════════════════════════╝")
//	fmt.Println()
//
//	channel := make(chan IndexedValue)
//
//	go func() {
//		defer close(channel)
//		for i := 0; i < 5; i++ {
//			channel <- IndexedValue{
//				Index: i,
//				Value: fmt.Sprintf("parsed - value %d", i),
//			}
//		}
//	}()
//
//	fmt.Println("Результат (впорядкований на льоту):")
//	for value := range sendWithOrdering(channel, 2) {
//		fmt.Println(value)
//	}
//
//	fmt.Println()
//}
//
//// ═══════════════════════════════════════════════════════════════════
//// ПОРІВНЯННЯ СПОСОБІВ
//// ═══════════════════════════════════════════════════════════════════
//
//func printComparison() {
//	fmt.Println("╔═══════════════════════════════════════════════════════════════╗")
//	fmt.Println("║              ПОРІВНЯННЯ СПОСОБІВ                              ║")
//	fmt.Println("╚═══════════════════════════════════════════════════════════════╝")
//	fmt.Println()
//	fmt.Println("┌───────────────────────────────────────────────────────────────┐")
//	fmt.Println("│ Спосіб              │ Пам'ять │ Латентність │ Складність    │")
//	fmt.Println("├───────────────────────────────────────────────────────────────┤")
//	fmt.Println("│ 1. Collect + Sort   │   ⭐⭐  │    ⭐⭐⭐   │    ⭐         │")
//	fmt.Println("│ 2. Ordered Array    │   ⭐⭐  │    ⭐⭐⭐   │    ⭐⭐       │")
//	fmt.Println("│ 3. Sequential Order │   ⭐⭐⭐ │    ⭐      │    ⭐⭐⭐     │")
//	fmt.Println("└───────────────────────────────────────────────────────────────┘")
//	fmt.Println()
//	fmt.Println("РЕКОМЕНДАЦІЇ:")
//	fmt.Println("  • Спосіб 1 - найпростіший для початку ✅")
//	fmt.Println("  • Спосіб 2 - коли знаєш кількість елементів ✅")
//	fmt.Println("  • Спосіб 3 - для streaming з low latency ✅")
//	fmt.Println()
//	fmt.Println("Пам'ять:      ⭐ = мало, ⭐⭐⭐ = багато")
//	fmt.Println("Латентність:  ⭐ = низька (швидко), ⭐⭐⭐ = висока (чекати)")
//	fmt.Println("Складність:   ⭐ = легко, ⭐⭐⭐ = складно")
//	fmt.Println()
//}
//
//// ═══════════════════════════════════════════════════════════════════
//// MAIN
//// ═══════════════════════════════════════════════════════════════════
//
//func main() {
//	fmt.Println()
//	fmt.Println("╔═══════════════════════════════════════════════════════════════╗")
//	fmt.Println("║         PARALLEL PIPELINE З ЗБЕРЕЖЕННЯМ ПОРЯДКУ               ║")
//	fmt.Println("╚═══════════════════════════════════════════════════════════════╝")
//	fmt.Println()
//
//	mainWithIndex()
//	mainOrderedResults()
//	mainSequentialOrdering()
//	printComparison()
//
//	fmt.Println("╔═══════════════════════════════════════════════════════════════╗")
//	fmt.Println("║                    ВИСНОВОК                                   ║")
//	fmt.Println("╚═══════════════════════════════════════════════════════════════╝")
//	fmt.Println()
//	fmt.Println("💡 Для збереження порядку в parallel pipeline:")
//	fmt.Println("   1. Додай індекс до кожного елементу")
//	fmt.Println("   2. Використай спосіб який підходить твоїм потребам")
//	fmt.Println("   3. Спосіб 1 - найпростіший для початку")
//	fmt.Println()
//}
