package main

//
//import (
//	"context"
//	"fmt"
//	"time"
//)
//
//// ═══════════════════════════════════════════════════════════════════
//// Типи з main.go (копія для незалежного запуску)
//// ═══════════════════════════════════════════════════════════════════
//
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
//func (newSorter *SorterService) SortMethod(ctx context.Context, urls []string) <-chan []string {
//	results := make(chan []string)
//
//	go func() {
//		defer close(results)
//		select {
//		case <-ctx.Done():
//			return
//		default:
//		}
//
//		res := newSorter.sorter.Swap(urls)
//
//		select {
//		case results <- res:
//		case <-ctx.Done():
//		}
//	}()
//	return results
//}
//
//// ═══════════════════════════════════════════════════════════════════
//// ПРАВИЛЬНА РЕАЛІЗАЦІЯ MOCK для SortInterface
//// ═══════════════════════════════════════════════════════════════════
//
//type MockSorter struct {
//	SwapFunc func(urls []string) []string
//	CallCount int
//	LastURLs  []string
//}
//
//// Swap реалізує SortInterface
//func (m *MockSorter) Swap(urls []string) []string {
//	m.CallCount++
//	m.LastURLs = append([]string{}, urls...)
//
//	if m.SwapFunc == nil {
//		return urls
//	}
//
//	return m.SwapFunc(urls)
//}
//
//// ═══════════════════════════════════════════════════════════════════
//// ПРИКЛАДИ ВИКОРИСТАННЯ
//// ═══════════════════════════════════════════════════════════════════
//
//func example1_ReverseOrder() {
//	fmt.Println("╔════════════════════════════════════════════════╗")
//	fmt.Println("║  Приклад 1: Mock з reverse logic              ║")
//	fmt.Println("╚════════════════════════════════════════════════╝")
//
//	ctx := context.Background()
//
//	// Створюємо mock з кастомною логікою
//	mock := &MockSorter{
//		SwapFunc: func(urls []string) []string {
//			// Реверсуємо порядок
//			reversed := make([]string, len(urls))
//			for i, url := range urls {
//				reversed[len(urls)-1-i] = url
//			}
//			return reversed
//		},
//	}
//
//	service := NewSorterService(mock)
//	results := service.SortMethod(ctx, []string{"url1", "url2", "url3"})
//
//	select {
//	case urls := <-results:
//		fmt.Printf("Input:  [url1, url2, url3]\n")
//		fmt.Printf("Output: %v\n", urls)
//		fmt.Printf("CallCount: %d\n", mock.CallCount)
//	case <-ctx.Done():
//		fmt.Println("Timeout")
//	}
//
//	fmt.Println()
//}
//
//func example2_AddPrefix() {
//	fmt.Println("╔════════════════════════════════════════════════╗")
//	fmt.Println("║  Приклад 2: Mock додає префікс                 ║")
//	fmt.Println("╚════════════════════════════════════════════════╝")
//
//	ctx := context.Background()
//
//	// Mock додає https:// до всіх URLs
//	mock := &MockSorter{
//		SwapFunc: func(urls []string) []string {
//			modified := make([]string, len(urls))
//			for i, url := range urls {
//				modified[i] = "https://" + url
//			}
//			return modified
//		},
//	}
//
//	service := NewSorterService(mock)
//	results := service.SortMethod(ctx, []string{"example.com", "test.com"})
//
//	select {
//	case urls := <-results:
//		fmt.Printf("Input:  [example.com, test.com]\n")
//		fmt.Printf("Output: %v\n", urls)
//	case <-ctx.Done():
//		fmt.Println("Timeout")
//	}
//
//	fmt.Println()
//}
//
//func example3_FilterURLs() {
//	fmt.Println("╔════════════════════════════════════════════════╗")
//	fmt.Println("║  Приклад 3: Mock фільтрує URLs                 ║")
//	fmt.Println("╚════════════════════════════════════════════════╝")
//
//	ctx := context.Background()
//
//	// Mock залишає тільки .com домени
//	mock := &MockSorter{
//		SwapFunc: func(urls []string) []string {
//			var filtered []string
//			for _, url := range urls {
//				if len(url) > 4 && url[len(url)-4:] == ".com" {
//					filtered = append(filtered, url)
//				}
//			}
//			return filtered
//		},
//	}
//
//	service := NewSorterService(mock)
//	results := service.SortMethod(ctx, []string{
//		"example.com",
//		"test.org",
//		"google.com",
//		"site.net",
//	})
//
//	select {
//	case urls := <-results:
//		fmt.Printf("Input:  [example.com, test.org, google.com, site.net]\n")
//		fmt.Printf("Output: %v (тільки .com)\n", urls)
//	case <-ctx.Done():
//		fmt.Println("Timeout")
//	}
//
//	fmt.Println()
//}
//
//func example4_DefaultBehavior() {
//	fmt.Println("╔════════════════════════════════════════════════╗")
//	fmt.Println("║  Приклад 4: Mock без функції (default)        ║")
//	fmt.Println("╚════════════════════════════════════════════════╝")
//
//	ctx := context.Background()
//
//	// Mock без SwapFunc - повертає URLs як є
//	mock := &MockSorter{
//		// Не задаємо SwapFunc
//	}
//
//	service := NewSorterService(mock)
//	results := service.SortMethod(ctx, []string{"a.com", "b.com", "c.com"})
//
//	select {
//	case urls := <-results:
//		fmt.Printf("Input:  [a.com, b.com, c.com]\n")
//		fmt.Printf("Output: %v (без змін)\n", urls)
//	case <-ctx.Done():
//		fmt.Println("Timeout")
//	}
//
//	fmt.Println()
//}
//
//func example5_Tracking() {
//	fmt.Println("╔════════════════════════════════════════════════╗")
//	fmt.Println("║  Приклад 5: Mock з tracking викликів           ║")
//	fmt.Println("╚════════════════════════════════════════════════╝")
//
//	ctx := context.Background()
//
//	mock := &MockSorter{
//		SwapFunc: func(urls []string) []string {
//			return urls
//		},
//	}
//
//	service := NewSorterService(mock)
//
//	// Викликаємо 3 рази
//	for i := 1; i <= 3; i++ {
//		results := service.SortMethod(ctx, []string{fmt.Sprintf("url%d.com", i)})
//		<-results // чекаємо результат
//	}
//
//	fmt.Printf("Кількість викликів: %d\n", mock.CallCount)
//	fmt.Printf("Останні URLs: %v\n", mock.LastURLs)
//
//	fmt.Println()
//}
//
//func example6_SlowOperation() {
//	fmt.Println("╔════════════════════════════════════════════════╗")
//	fmt.Println("║  Приклад 6: Mock з затримкою (slow)            ║")
//	fmt.Println("╚════════════════════════════════════════════════╝")
//
//	ctx := context.Background()
//
//	// Mock з затримкою
//	mock := &MockSorter{
//		SwapFunc: func(urls []string) []string {
//			fmt.Println("  → Починаю повільну обробку...")
//			time.Sleep(500 * time.Millisecond)
//			fmt.Println("  → Завершив!")
//			return urls
//		},
//	}
//
//	service := NewSorterService(mock)
//	start := time.Now()
//	results := service.SortMethod(ctx, []string{"slow.com"})
//
//	select {
//	case urls := <-results:
//		fmt.Printf("Output: %v\n", urls)
//		fmt.Printf("Час виконання: %v\n", time.Since(start))
//	case <-ctx.Done():
//		fmt.Println("Timeout")
//	}
//
//	fmt.Println()
//}
//
//func example7_ContextCancellation() {
//	fmt.Println("╔════════════════════════════════════════════════╗")
//	fmt.Println("║  Приклад 7: Тестування timeout                 ║")
//	fmt.Println("╚════════════════════════════════════════════════╝")
//
//	// Context з коротким timeout
//	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
//	defer cancel()
//
//	// Mock з довгою затримкою
//	mock := &MockSorter{
//		SwapFunc: func(urls []string) []string {
//			time.Sleep(500 * time.Millisecond) // Довше за timeout
//			return urls
//		},
//	}
//
//	service := NewSorterService(mock)
//	results := service.SortMethod(ctx, []string{"timeout-test.com"})
//
//	select {
//	case urls := <-results:
//		fmt.Printf("✅ Отримано: %v\n", urls)
//	case <-ctx.Done():
//		fmt.Printf("⏱️  Timeout! Context cancelled: %v\n", ctx.Err())
//	}
//
//	fmt.Println()
//}
//
//func example8_EmptyURLs() {
//	fmt.Println("╔════════════════════════════════════════════════╗")
//	fmt.Println("║  Приклад 8: Mock з порожнім списком            ║")
//	fmt.Println("╚════════════════════════════════════════════════╝")
//
//	ctx := context.Background()
//
//	mock := &MockSorter{
//		SwapFunc: func(urls []string) []string {
//			if len(urls) == 0 {
//				return []string{"default.com"}
//			}
//			return urls
//		},
//	}
//
//	service := NewSorterService(mock)
//
//	// Тест 1: порожній список
//	results1 := service.SortMethod(ctx, []string{})
//	select {
//	case urls := <-results1:
//		fmt.Printf("Порожній input → Output: %v\n", urls)
//	case <-ctx.Done():
//		fmt.Println("Timeout")
//	}
//
//	// Тест 2: непорожній список
//	results2 := service.SortMethod(ctx, []string{"test.com"})
//	select {
//	case urls := <-results2:
//		fmt.Printf("Непорожній input → Output: %v\n", urls)
//	case <-ctx.Done():
//		fmt.Println("Timeout")
//	}
//
//	fmt.Println()
//}
//
//// ═══════════════════════════════════════════════════════════════════
//// ПОРІВНЯННЯ: DefaultSortingStruct vs MockSorter
//// ═══════════════════════════════════════════════════════════════════
//
//func comparison() {
//	fmt.Println("╔═══════════════════════════════════════════════════════════════╗")
//	fmt.Println("║           ПОРІВНЯННЯ: Default vs Mock                        ║")
//	fmt.Println("╚═══════════════════════════════════════════════════════════════╝")
//	fmt.Println()
//
//	ctx := context.Background()
//	urls := []string{"url1", "url2", "url3"}
//
//	// Default sorter
//	fmt.Println("1. DefaultSortingStruct (без змін):")
//	defaultSorter := &DefaultSortingStruct{}
//	svc1 := NewSorterService(defaultSorter)
//	results1 := svc1.SortMethod(ctx, urls)
//	url1 := <-results1
//	fmt.Printf("   Input:  %v\n", urls)
//	fmt.Printf("   Output: %v\n", url1)
//	fmt.Println()
//
//	// Mock sorter - reverse
//	fmt.Println("2. MockSorter (reverse):")
//	mockReverse := &MockSorter{
//		SwapFunc: func(urls []string) []string {
//			reversed := make([]string, len(urls))
//			for i, url := range urls {
//				reversed[len(urls)-1-i] = url
//			}
//			return reversed
//		},
//	}
//	svc2 := NewSorterService(mockReverse)
//	results2 := svc2.SortMethod(ctx, urls)
//	url2 := <-results2
//	fmt.Printf("   Input:  %v\n", urls)
//	fmt.Printf("   Output: %v\n", url2)
//	fmt.Println()
//
//	// Mock sorter - add prefix
//	fmt.Println("3. MockSorter (add prefix):")
//	mockPrefix := &MockSorter{
//		SwapFunc: func(urls []string) []string {
//			modified := make([]string, len(urls))
//			for i, url := range urls {
//				modified[i] = "https://" + url
//			}
//			return modified
//		},
//	}
//	svc3 := NewSorterService(mockPrefix)
//	results3 := svc3.SortMethod(ctx, urls)
//	url3 := <-results3
//	fmt.Printf("   Input:  %v\n", urls)
//	fmt.Printf("   Output: %v\n", url3)
//	fmt.Println()
//
//	fmt.Println("┌──────────────────────────────────────────────────────────┐")
//	fmt.Println("│  Властивість           │ Default │ Mock                  │")
//	fmt.Println("├──────────────────────────────────────────────────────────┤")
//	fmt.Println("│ Фіксована логіка       │   ✅    │   ❌                  │")
//	fmt.Println("│ Кастомна логіка        │   ❌    │   ✅                  │")
//	fmt.Println("│ Для production         │   ✅    │   ❌                  │")
//	fmt.Println("│ Для тестів             │   ❌    │   ✅                  │")
//	fmt.Println("│ Можна inject логіку    │   ❌    │   ✅                  │")
//	fmt.Println("│ Tracking викликів      │   ❌    │   ✅                  │")
//	fmt.Println("└──────────────────────────────────────────────────────────┘")
//	fmt.Println()
//}
//
//// ═══════════════════════════════════════════════════════════════════
//// ІНСТРУКЦІЇ
//// ═══════════════════════════════════════════════════════════════════
//
//func printInstructions() {
//	fmt.Println("╔═══════════════════════════════════════════════════════════════╗")
//	fmt.Println("║              ЯК ВИКОРИСТОВУВАТИ MOCK                          ║")
//	fmt.Println("╚═══════════════════════════════════════════════════════════════╝")
//	fmt.Println()
//	fmt.Println("1️⃣  Створи mock з кастомною логікою:")
//	fmt.Println("   mock := &MockSorter{")
//	fmt.Println("       SwapFunc: func(urls []string) []string {")
//	fmt.Println("           // Твоя логіка тут")
//	fmt.Println("           return urls")
//	fmt.Println("       },")
//	fmt.Println("   }")
//	fmt.Println()
//	fmt.Println("2️⃣  Передай mock в service:")
//	fmt.Println("   service := NewSorterService(mock)")
//	fmt.Println()
//	fmt.Println("3️⃣  Використовуй як звичайний service:")
//	fmt.Println("   results := service.SortMethod(ctx, urls)")
//	fmt.Println("   urls := <-results")
//	fmt.Println()
//	fmt.Println("💡 Mock реалізує той самий SortInterface інтерфейс!")
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
//	fmt.Println("║              MOCK ПРИКЛАДИ ДЛЯ SortInterface                  ║")
//	fmt.Println("╚═══════════════════════════════════════════════════════════════╝")
//	fmt.Println()
//
//	printInstructions()
//
//	example1_ReverseOrder()
//	example2_AddPrefix()
//	example3_FilterURLs()
//	example4_DefaultBehavior()
//	example5_Tracking()
//	example6_SlowOperation()
//	example7_ContextCancellation()
//	example8_EmptyURLs()
//
//	comparison()
//
//	fmt.Println("╔═══════════════════════════════════════════════════════════════╗")
//	fmt.Println("║                         ВИСНОВОК                              ║")
//	fmt.Println("╚═══════════════════════════════════════════════════════════════╝")
//	fmt.Println()
//	fmt.Println("✅ MockSorter реалізує SortInterface")
//	fmt.Println("✅ Можна inject будь-яку логіку через SwapFunc")
//	fmt.Println("✅ Має tracking викликів (CallCount, LastURLs)")
//	fmt.Println("✅ Ідеально для тестування!")
//	fmt.Println()
//}
