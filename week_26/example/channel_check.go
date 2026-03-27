package main

//
//import (
//	"fmt"
//	"time"
//)
//
//// Типи для прикладу
//type DefaultSortingStruct struct{}
//
//type SortInterface interface {
//	Swap([]string) []string
//}
//
//func (s *DefaultSortingStruct) Swap(urls []string) []string {
//	return urls
//}
//
//type SorterService struct {
//	sorter SortInterface
//}
//
//func NewSorterService(s SortInterface) *SorterService {
//	return &SorterService{sorter: s}
//}
//
//func (newSorter *SorterService) SortMethod(urls []string) <-chan []string {
//	results := make(chan []string)
//
//	go func() {
//		defer close(results)
//		res := newSorter.sorter.Swap(urls)
//		results <- res
//	}()
//	return results
//}
//
//// ═══════════════════════════════════════════════════════════════════
//// ВАРІАНТИ ПЕРЕВІРКИ ЧИ ПРАЦЮЄ CHANNEL (БЕЗ TIMEOUT)
//// ═══════════════════════════════════════════════════════════════════
//
//// ═══════════════════════════════════════════════════════════════════
//// СПОСІБ 1: select з default (NON-BLOCKING CHECK)
//// ═══════════════════════════════════════════════════════════════════
//
//func example1_SelectWithDefault() {
//	fmt.Println("╔════════════════════════════════════════════════╗")
//	fmt.Println("║  Спосіб 1: select з default (non-blocking)    ║")
//	fmt.Println("╚════════════════════════════════════════════════╝")
//
//	ch := make(chan string, 1)
//
//	// Надсилаємо дані
//	ch <- "дані готові!"
//
//	// Перевіряємо чи є дані БЕЗ блокування
//	select {
//	case data := <-ch:
//		fmt.Println("✅ Отримав дані:", data)
//	default:
//		fmt.Println("⏳ Канал порожній, працюю далі...")
//	}
//
//	// Перевіряємо знову (channel тепер порожній)
//	select {
//	case data := <-ch:
//		fmt.Println("✅ Отримав дані:", data)
//	default:
//		fmt.Println("⏳ Канал порожній, працюю далі...")
//	}
//
//	fmt.Println()
//}
//
//// ═══════════════════════════════════════════════════════════════════
//// СПОСІБ 2: Перевірка чи channel ЗАКРИТИЙ
//// ═══════════════════════════════════════════════════════════════════
//
//func example2_CheckIfClosed() {
//	fmt.Println("╔════════════════════════════════════════════════╗")
//	fmt.Println("║  Спосіб 2: Перевірка чи закритий channel      ║")
//	fmt.Println("╚════════════════════════════════════════════════╝")
//
//	ch := make(chan string, 1)
//
//	// Надсилаємо дані і закриваємо
//	ch <- "останні дані"
//	close(ch)
//
//	// Перевіряємо з ok pattern
//	if data, ok := <-ch; ok {
//		fmt.Println("✅ Channel відкритий, дані:", data)
//	} else {
//		fmt.Println("❌ Channel закритий, дані:", data) // буде zero value
//	}
//
//	// Перевіряємо знову (після закриття завжди ok=false)
//	if data, ok := <-ch; ok {
//		fmt.Println("✅ Channel відкритий, дані:", data)
//	} else {
//		fmt.Println("❌ Channel закритий, zero value:", data)
//	}
//
//	fmt.Println()
//}
//
//// ═══════════════════════════════════════════════════════════════════
//// СПОСІБ 3: range - автоматична перевірка до закриття
//// ═══════════════════════════════════════════════════════════════════
//
//func example3_RangeOverChannel() {
//	fmt.Println("╔════════════════════════════════════════════════╗")
//	fmt.Println("║  Спосіб 3: range (читає до закриття)          ║")
//	fmt.Println("╚════════════════════════════════════════════════╝")
//
//	ch := make(chan int, 5)
//
//	// Надсилаємо кілька значень
//	go func() {
//		for i := 1; i <= 5; i++ {
//			fmt.Printf("  → Надсилаю %d\n", i)
//			ch <- i
//			time.Sleep(100 * time.Millisecond)
//		}
//		fmt.Println("  → Закриваю channel")
//		close(ch)
//	}()
//
//	// range автоматично перевіряє чи channel закритий
//	fmt.Println("\n  Отримую дані:")
//	for value := range ch {
//		fmt.Printf("  ✅ Отримав: %d\n", value)
//	}
//	fmt.Println("  ✅ Channel закритий, range завершився")
//
//	fmt.Println()
//}
//
//// ═══════════════════════════════════════════════════════════════════
//// СПОСІБ 4: Polling loop з select default
//// ═══════════════════════════════════════════════════════════════════
//
//func example4_PollingLoop() {
//	fmt.Println("╔════════════════════════════════════════════════╗")
//	fmt.Println("║  Спосіб 4: Polling loop (перевіряємо цикл)    ║")
//	fmt.Println("╚════════════════════════════════════════════════╝")
//
//	ch := make(chan string, 1)
//
//	// Запускаємо goroutine яка надішле дані через 500ms
//	go func() {
//		time.Sleep(500 * time.Millisecond)
//		ch <- "дані прийшли!"
//	}()
//
//	// Перевіряємо кілька разів
//	for i := 1; i <= 10; i++ {
//		select {
//		case data := <-ch:
//			fmt.Printf("  ✅ Спроба %d: Отримав дані: %s\n", i, data)
//			return // Отримали дані, виходимо
//		default:
//			fmt.Printf("  ⏳ Спроба %d: Ще немає даних, працюю далі...\n", i)
//			// Можемо робити іншу роботу тут
//			time.Sleep(100 * time.Millisecond)
//		}
//	}
//
//	fmt.Println()
//}
//
//// ═══════════════════════════════════════════════════════════════════
//// СПОСІБ 5: Комбінований - select з ok pattern
//// ═══════════════════════════════════════════════════════════════════
//
//func example5_SelectWithOkPattern() {
//	fmt.Println("╔════════════════════════════════════════════════╗")
//	fmt.Println("║  Спосіб 5: select + ok pattern                 ║")
//	fmt.Println("╚════════════════════════════════════════════════╝")
//
//	ch := make(chan int, 3)
//
//	// Надсилаємо дані і закриваємо
//	go func() {
//		ch <- 1
//		ch <- 2
//		ch <- 3
//		close(ch)
//	}()
//
//	time.Sleep(100 * time.Millisecond) // Даємо час goroutine
//
//	// Читаємо до закриття
//	for {
//		select {
//		case data, ok := <-ch:
//			if !ok {
//				fmt.Println("  ❌ Channel закритий, виходимо")
//				goto Done
//			}
//			fmt.Printf("  ✅ Отримав: %d\n", data)
//		default:
//			fmt.Println("  ⏳ Немає даних, працюю далі...")
//			time.Sleep(50 * time.Millisecond)
//		}
//	}
//
//Done:
//	fmt.Println()
//}
//
//// ═══════════════════════════════════════════════════════════════════
//// СПОСІБ 6: Перевірка len/cap (ОБЕРЕЖНО!)
//// ═══════════════════════════════════════════════════════════════════
//
//func example6_LenCapCheck() {
//	fmt.Println("╔════════════════════════════════════════════════╗")
//	fmt.Println("║  Спосіб 6: len/cap (тільки для buffered!)     ║")
//	fmt.Println("╚════════════════════════════════════════════════╝")
//
//	ch := make(chan string, 5) // buffered channel
//
//	ch <- "item1"
//	ch <- "item2"
//	ch <- "item3"
//
//	fmt.Printf("  📊 Capacity: %d\n", cap(ch))
//	fmt.Printf("  📊 Length (елементів в черзі): %d\n", len(ch))
//	fmt.Printf("  📊 Вільно місця: %d\n", cap(ch)-len(ch))
//
//	// ⚠️ ВАЖЛИВО: len() не каже чи channel закритий!
//	if len(ch) > 0 {
//		fmt.Println("  ✅ Є дані в channel")
//		data := <-ch
//		fmt.Println("  ✅ Отримав:", data)
//		fmt.Printf("  📊 Після читання length: %d\n", len(ch))
//	}
//
//	fmt.Println("\n  ⚠️  УВАГА: len() працює тільки для buffered channels!")
//	fmt.Println("  ⚠️  len() НЕ каже чи channel закритий!")
//
//	fmt.Println()
//}
//
//// ═══════════════════════════════════════════════════════════════════
//// ПРАКТИЧНИЙ ПРИКЛАД: Worker з перевіркою
//// ═══════════════════════════════════════════════════════════════════
//
//func example7_PracticalWorker() {
//	fmt.Println("╔════════════════════════════════════════════════╗")
//	fmt.Println("║  Практичний приклад: Worker                    ║")
//	fmt.Println("╚════════════════════════════════════════════════╝")
//
//	jobs := make(chan int, 10)
//	done := make(chan bool)
//
//	// Worker
//	go func() {
//		for {
//			select {
//			case job, ok := <-jobs:
//				if !ok {
//					// Channel закритий, всі роботи зроблено
//					fmt.Println("  ✅ Всі роботи виконано, worker завершився")
//					done <- true
//					return
//				}
//				// Обробляємо роботу
//				fmt.Printf("  🔧 Обробляю job %d\n", job)
//				time.Sleep(100 * time.Millisecond)
//			default:
//				// Немає робіт зараз, але channel не закритий
//				fmt.Println("  ⏳ Немає робіт, чекаю...")
//				time.Sleep(50 * time.Millisecond)
//			}
//		}
//	}()
//
//	// Надсилаємо роботи
//	for i := 1; i <= 3; i++ {
//		fmt.Printf("  → Надсилаю job %d\n", i)
//		jobs <- i
//		time.Sleep(150 * time.Millisecond)
//	}
//
//	// Закриваємо channel коли всі роботи надіслані
//	fmt.Println("  → Закриваю channel jobs")
//	close(jobs)
//
//	// Чекаємо завершення worker
//	<-done
//	fmt.Println()
//}
//
//// ═══════════════════════════════════════════════════════════════════
//// ПОРІВНЯЛЬНА ТАБЛИЦЯ
//// ═══════════════════════════════════════════════════════════════════
//
//func printComparison() {
//	fmt.Println("╔═══════════════════════════════════════════════════════════════╗")
//	fmt.Println("║              ПОРІВНЯННЯ СПОСОБІВ                              ║")
//	fmt.Println("╚═══════════════════════════════════════════════════════════════╝")
//	fmt.Println()
//	fmt.Println("┌───────────────────────────────────────────────────────────────┐")
//	fmt.Println("│ Спосіб              │ Блокує │ Перевіряє │ Коли використовувати│")
//	fmt.Println("│                     │        │ закриття  │                     │")
//	fmt.Println("├───────────────────────────────────────────────────────────────┤")
//	fmt.Println("│ 1. select + default │   ❌   │     ❌    │ Non-blocking read   │")
//	fmt.Println("│ 2. value, ok := <-ch│   ✅   │     ✅    │ Перевірка закриття  │")
//	fmt.Println("│ 3. range            │   ✅   │     ✅    │ Читати всі дані     │")
//	fmt.Println("│ 4. Polling loop     │   ❌   │     ❌    │ Періодична перевірка│")
//	fmt.Println("│ 5. select + ok      │   ❌   │     ✅    │ Non-block + закриття│")
//	fmt.Println("│ 6. len(ch)          │   ❌   │     ❌    │ Тільки buffered ch  │")
//	fmt.Println("└───────────────────────────────────────────────────────────────┘")
//	fmt.Println()
//	fmt.Println("РЕКОМЕНДАЦІЇ:")
//	fmt.Println("  ✅ select + default      - для non-blocking перевірки")
//	fmt.Println("  ✅ value, ok := <-ch     - щоб знати чи закритий")
//	fmt.Println("  ✅ range                 - найпростіше для читання всіх даних")
//	fmt.Println("  ✅ select + ok pattern   - комбінація переваг")
//	fmt.Println()
//}
//
//// ═══════════════════════════════════════════════════════════════════
//// ДЛЯ ТВОГО ВИПАДКУ (week_26/example/main.go)
//// ═══════════════════════════════════════════════════════════════════
//
//func yourExample() {
//	fmt.Println("╔═══════════════════════════════════════════════════════════════╗")
//	fmt.Println("║         ДЛЯ ТВОГО КОДУ (week_26/example/main.go)              ║")
//	fmt.Println("╚═══════════════════════════════════════════════════════════════╝")
//	fmt.Println()
//
//	urls := []string{"url1", "url2", "url3"}
//	s := &DefaultSortingStruct{}
//	newSorter := NewSorterService(s)
//
//	results := newSorter.SortMethod(urls)
//
//	fmt.Println("Варіант 1: Non-blocking check")
//	fmt.Println("──────────────────────────────")
//	select {
//	case url := <-results:
//		fmt.Println("✅ Є дані:", url)
//	default:
//		fmt.Println("⏳ Немає даних ще, працюю далі...")
//		// Можна робити іншу роботу
//		fmt.Println("   Роблю щось інше...")
//		time.Sleep(100 * time.Millisecond)
//		// Потім спробувати знову
//		url := <-results // Тепер блокуючий read
//		fmt.Println("✅ Отримав дані:", url)
//	}
//	fmt.Println()
//
//	fmt.Println("Варіант 2: З перевіркою закриття")
//	fmt.Println("─────────────────────────────────")
//	results2 := newSorter.SortMethod(urls)
//	if url, ok := <-results2; ok {
//		fmt.Println("✅ Channel відкритий, дані:", url)
//	} else {
//		fmt.Println("❌ Channel закритий")
//	}
//	fmt.Println()
//
//	fmt.Println("Варіант 3: Найпростіший (просто чекай)")
//	fmt.Println("──────────────────────────────────────")
//	results3 := newSorter.SortMethod(urls)
//	url := <-results3 // Блокує поки не прийдуть дані
//	fmt.Println("✅ Отримав:", url)
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
//	fmt.Println("║     ПЕРЕВІРКА ЧИ ПРАЦЮЄ CHANNEL (БЕЗ TIMEOUT)                 ║")
//	fmt.Println("╚═══════════════════════════════════════════════════════════════╝")
//	fmt.Println()
//
//	example1_SelectWithDefault()
//	example2_CheckIfClosed()
//	example3_RangeOverChannel()
//	example4_PollingLoop()
//	example5_SelectWithOkPattern()
//	example6_LenCapCheck()
//	example7_PracticalWorker()
//
//	printComparison()
//	yourExample()
//
//	fmt.Println("╔═══════════════════════════════════════════════════════════════╗")
//	fmt.Println("║                         ВИСНОВОК                              ║")
//	fmt.Println("╚═══════════════════════════════════════════════════════════════╝")
//	fmt.Println()
//	fmt.Println("💡 Для перевірки \"чи є дані БЕЗ блокування\":")
//	fmt.Println("   select { case data := <-ch: ...; default: ... }")
//	fmt.Println()
//	fmt.Println("💡 Для перевірки \"чи channel закритий\":")
//	fmt.Println("   if data, ok := <-ch; ok { ... }")
//	fmt.Println()
//	fmt.Println("💡 Для читання всіх даних до закриття:")
//	fmt.Println("   for data := range ch { ... }")
//	fmt.Println()
//	fmt.Println("💡 Для твого коду - найпростіше:")
//	fmt.Println("   url := <-results  // Просто чекай!")
//	fmt.Println()
//}
