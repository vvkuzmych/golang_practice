package main

//
//import (
//	"context"
//	"fmt"
//	"time"
//)
//
//// Типи з оригінального main.go
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
//
//		res := newSorter.sorter.Swap(urls)
//		results <- res
//	}()
//	return results
//}
//
//// ═══════════════════════════════════════════════════════════════════
//// ВАРІАНТ 1: НАЙПРОСТІШИЙ - Просто чекати (БЕЗ timeout)
//// ═══════════════════════════════════════════════════════════════════
//
//func example1_SimpleBlocking() {
//	fmt.Println("╔════════════════════════════════════════════════╗")
//	fmt.Println("║  Варіант 1: Простий блокуючий receive         ║")
//	fmt.Println("╚════════════════════════════════════════════════╝")
//
//	urls := []string{"url1", "url2", "url3"}
//	s := &DefaultSortingStruct{}
//	newSorter := NewSorterService(s)
//
//	results := newSorter.SortMethod(urls)
//
//	// ✅ Найпростіше - просто чекаємо
//	url := <-results
//	fmt.Println("✅ Отримано:", url)
//	fmt.Println()
//}
//
//// ═══════════════════════════════════════════════════════════════════
//// ВАРІАНТ 2: Context з timeout (ПРОФЕСІЙНИЙ ПІДХІД)
//// ═══════════════════════════════════════════════════════════════════
//
//// Модифікований метод з context
//func (newSorter *SorterService) SortMethodWithContext(ctx context.Context, urls []string) <-chan []string {
//	results := make(chan []string, 1)
//
//	go func() {
//		defer close(results)
//
//		// Перевіряємо чи не скасовано
//		select {
//		case <-ctx.Done():
//			return
//		default:
//		}
//
//		res := newSorter.sorter.Swap(urls)
//
//		// Надсилаємо з можливістю cancellation
//		select {
//		case results <- res:
//		case <-ctx.Done():
//		}
//	}()
//	return results
//}
//
//func example2_WithContext() {
//	fmt.Println("╔════════════════════════════════════════════════╗")
//	fmt.Println("║  Варіант 2: Context з timeout                  ║")
//	fmt.Println("╚════════════════════════════════════════════════╝")
//
//	// ✅ Timeout через context (більш гнучко)
//	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//	defer cancel()
//
//	urls := []string{"url1", "url2", "url3"}
//	s := &DefaultSortingStruct{}
//	newSorter := NewSorterService(s)
//
//	results := newSorter.SortMethodWithContext(ctx, urls)
//
//	// Чекаємо результат АБО context timeout
//	select {
//	case url := <-results:
//		fmt.Println("✅ Отримано:", url)
//	case <-ctx.Done():
//		fmt.Println("❌ Timeout або cancellation:", ctx.Err())
//	}
//	fmt.Println()
//}
//
//// ═══════════════════════════════════════════════════════════════════
//// ВАРІАНТ 3: Context з deadline (конкретний час)
//// ═══════════════════════════════════════════════════════════════════
//
//func example3_WithDeadline() {
//	fmt.Println("╔════════════════════════════════════════════════╗")
//	fmt.Println("║  Варіант 3: Context з deadline                 ║")
//	fmt.Println("╚════════════════════════════════════════════════╝")
//
//	// ✅ Deadline - конкретний час завершення
//	deadline := time.Now().Add(5 * time.Second)
//	ctx, cancel := context.WithDeadline(context.Background(), deadline)
//	defer cancel()
//
//	urls := []string{"url1", "url2", "url3"}
//	s := &DefaultSortingStruct{}
//	newSorter := NewSorterService(s)
//
//	results := newSorter.SortMethodWithContext(ctx, urls)
//
//	select {
//	case url := <-results:
//		fmt.Println("✅ Отримано:", url)
//	case <-ctx.Done():
//		fmt.Println("❌ Досягнуто deadline:", ctx.Err())
//	}
//	fmt.Println()
//}
//
//// ═══════════════════════════════════════════════════════════════════
//// ВАРІАНТ 4: Done channel (ручний контроль)
//// ═══════════════════════════════════════════════════════════════════
//
//func (newSorter *SorterService) SortMethodWithDone(urls []string, done <-chan struct{}) <-chan []string {
//	results := make(chan []string, 1)
//
//	go func() {
//		defer close(results)
//
//		// Перевіряємо done channel
//		select {
//		case <-done:
//			return
//		default:
//		}
//
//		res := newSorter.sorter.Swap(urls)
//
//		// Надсилаємо з можливістю cancellation
//		select {
//		case results <- res:
//		case <-done:
//		}
//	}()
//	return results
//}
//
//func example4_WithDoneChannel() {
//	fmt.Println("╔════════════════════════════════════════════════╗")
//	fmt.Println("║  Варіант 4: Done channel                       ║")
//	fmt.Println("╚════════════════════════════════════════════════╝")
//
//	// ✅ Власний done channel
//	done := make(chan struct{})
//
//	// Можна вручну контролювати cancellation
//	go func() {
//		time.Sleep(5 * time.Second)
//		close(done) // Скасовуємо через 5 секунд
//	}()
//
//	urls := []string{"url1", "url2", "url3"}
//	s := &DefaultSortingStruct{}
//	newSorter := NewSorterService(s)
//
//	results := newSorter.SortMethodWithDone(urls, done)
//
//	select {
//	case url := <-results:
//		fmt.Println("✅ Отримано:", url)
//	case <-done:
//		fmt.Println("❌ Операція скасована")
//	}
//	fmt.Println()
//}
//
//// ═══════════════════════════════════════════════════════════════════
//// ВАРІАНТ 5: Context з cancellation (повний контроль)
//// ═══════════════════════════════════════════════════════════════════
//
//func example5_WithCancellation() {
//	fmt.Println("╔════════════════════════════════════════════════╗")
//	fmt.Println("║  Варіант 5: Context з cancel function          ║")
//	fmt.Println("╚════════════════════════════════════════════════╝")
//
//	// ✅ Context з можливістю скасування
//	ctx, cancel := context.WithCancel(context.Background())
//
//	// Симулюємо користувача який скасовує операцію
//	go func() {
//		time.Sleep(100 * time.Millisecond)
//		fmt.Println("→ Користувач натиснув Cancel!")
//		cancel()
//	}()
//
//	urls := []string{"url1", "url2", "url3"}
//	s := &DefaultSortingStruct{}
//	newSorter := NewSorterService(s)
//
//	results := newSorter.SortMethodWithContext(ctx, urls)
//
//	select {
//	case url := <-results:
//		fmt.Println("✅ Отримано:", url)
//	case <-ctx.Done():
//		fmt.Println("❌ Користувач скасував:", ctx.Err())
//	}
//	fmt.Println()
//}
//
//// ═══════════════════════════════════════════════════════════════════
//// ВАРІАНТ 6: Retry логіка (якщо перша спроба не вдалась)
//// ═══════════════════════════════════════════════════════════════════
//
//func example6_WithRetry() {
//	fmt.Println("╔════════════════════════════════════════════════╗")
//	fmt.Println("║  Варіант 6: З retry логікою                    ║")
//	fmt.Println("╚════════════════════════════════════════════════╝")
//
//	urls := []string{"url1", "url2", "url3"}
//	s := &DefaultSortingStruct{}
//	newSorter := NewSorterService(s)
//
//	maxRetries := 3
//	for attempt := 1; attempt <= maxRetries; attempt++ {
//		fmt.Printf("→ Спроба %d/%d\n", attempt, maxRetries)
//
//		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
//		results := newSorter.SortMethodWithContext(ctx, urls)
//
//		select {
//		case url := <-results:
//			fmt.Println("✅ Успіх! Отримано:", url)
//			cancel()
//			return
//		case <-ctx.Done():
//			fmt.Printf("❌ Спроба %d не вдалась\n", attempt)
//			cancel()
//			if attempt < maxRetries {
//				time.Sleep(500 * time.Millisecond) // Пауза перед retry
//			}
//		}
//	}
//	fmt.Println("❌ Всі спроби вичерпано")
//	fmt.Println()
//}
//
//// ═══════════════════════════════════════════════════════════════════
//// ВАРІАНТ 7: БЕЗ select взагалі (найпростіший)
//// ═══════════════════════════════════════════════════════════════════
//
//func example7_NoSelect() {
//	fmt.Println("╔════════════════════════════════════════════════╗")
//	fmt.Println("║  Варіант 7: Без select (найпростіше)          ║")
//	fmt.Println("╚════════════════════════════════════════════════╝")
//
//	urls := []string{"url1", "url2", "url3"}
//	s := &DefaultSortingStruct{}
//	newSorter := NewSorterService(s)
//
//	results := newSorter.SortMethod(urls)
//
//	// ✅ Просто отримуємо з channel - НАЙПРОСТІШЕ!
//	url := <-results
//	fmt.Println("✅ Отримано:", url)
//	fmt.Println()
//}
//
//// ═══════════════════════════════════════════════════════════════════
//// ПОРІВНЯЛЬНА ТАБЛИЦЯ
//// ═══════════════════════════════════════════════════════════════════
//
//func printComparison() {
//	fmt.Println("╔════════════════════════════════════════════════════════════════╗")
//	fmt.Println("║              ПОРІВНЯННЯ ВСІХ ВАРІАНТІВ                         ║")
//	fmt.Println("╚════════════════════════════════════════════════════════════════╝")
//	fmt.Println()
//	fmt.Println("┌──────────────────────────────────────────────────────────────┐")
//	fmt.Println("│  Варіант                 │ Складність │ Timeout │ Cancel   │")
//	fmt.Println("├──────────────────────────────────────────────────────────────┤")
//	fmt.Println("│ 1. Simple blocking       │     ⭐     │    ❌   │    ❌    │")
//	fmt.Println("│ 2. Context timeout       │    ⭐⭐    │    ✅   │    ✅    │")
//	fmt.Println("│ 3. Context deadline      │    ⭐⭐    │    ✅   │    ✅    │")
//	fmt.Println("│ 4. Done channel          │   ⭐⭐⭐   │    ✅   │    ✅    │")
//	fmt.Println("│ 5. Context cancel        │    ⭐⭐    │    ❌   │    ✅    │")
//	fmt.Println("│ 6. Retry logic           │   ⭐⭐⭐   │    ✅   │    ✅    │")
//	fmt.Println("│ 7. No select             │     ⭐     │    ❌   │    ❌    │")
//	fmt.Println("└──────────────────────────────────────────────────────────────┘")
//	fmt.Println()
//	fmt.Println("РЕКОМЕНДАЦІЇ:")
//	fmt.Println("  • Варіант 1 або 7 - найпростіше для навчання ✅")
//	fmt.Println("  • Варіант 2 - найкраще для production коду ✅✅")
//	fmt.Println("  • Варіант 5 - для UI з можливістю Cancel ✅")
//	fmt.Println("  • Варіант 6 - для ненадійних операцій ✅")
//	fmt.Println()
//}
//
//// ═══════════════════════════════════════════════════════════════════
//// MAIN
//// ═══════════════════════════════════════════════════════════════════
//
//func main() {
//	fmt.Println()
//	fmt.Println("╔════════════════════════════════════════════════════════════════╗")
//	fmt.Println("║     АЛЬТЕРНАТИВИ ДО time.After В SELECT                        ║")
//	fmt.Println("╚════════════════════════════════════════════════════════════════╝")
//	fmt.Println()
//
//	example1_SimpleBlocking()
//	example2_WithContext()
//	example3_WithDeadline()
//	example4_WithDoneChannel()
//	example5_WithCancellation()
//	example6_WithRetry()
//	example7_NoSelect()
//
//	printComparison()
//
//	fmt.Println("╔════════════════════════════════════════════════════════════════╗")
//	fmt.Println("║                    ВИСНОВОК                                    ║")
//	fmt.Println("╚════════════════════════════════════════════════════════════════╝")
//	fmt.Println()
//	fmt.Println("💡 Для простого коду (як у твоєму прикладі):")
//	fmt.Println("   → Використовуй Варіант 1 або 7 (без select)")
//	fmt.Println()
//	fmt.Println("💡 Для production коду:")
//	fmt.Println("   → Використовуй Варіант 2 (context.WithTimeout)")
//	fmt.Println()
//	fmt.Println("💡 Для UI з кнопкою Cancel:")
//	fmt.Println("   → Використовуй Варіант 5 (context.WithCancel)")
//	fmt.Println()
//}
