package main

import (
	"context"
	"sync"
	"testing"
	"time"
)

// ═══════════════════════════════════════════════════════════════════
// ТЕСТИ ДЛЯ ПЕРЕВІРКИ RACE CONDITIONS В URLService
// ═══════════════════════════════════════════════════════════════════

// Запустити тести з race detector:
//   go test -race -v

// TestProcessAsyncRaceCondition перевіряє чи немає race condition
// при одночасному виклику ProcessAsync з багатьох goroutines
func TestProcessAsyncRaceCondition(t *testing.T) {
	sorter := &DefaultURLSorter{}
	service := NewURLService(sorter)
	ctx := context.Background()

	urls := []string{"example.com", "test.com", "google.com"}
	numGoroutines := 100

	var wg sync.WaitGroup

	// Запускаємо багато goroutines одночасно
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			// Викликаємо ProcessAsync
			resultCh := service.ProcessAsync(ctx, urls)
			result := <-resultCh

			if result.Err != nil {
				t.Errorf("Unexpected error: %v", result.Err)
				return
			}

			if len(result.URLs) != len(urls) {
				t.Errorf("Expected %d URLs, got %d", len(urls), len(result.URLs))
			}
		}()
	}

	wg.Wait()
}

// TestProcessWithCallbackRaceCondition перевіряє race condition
// при використанні callback pattern
func TestProcessWithCallbackRaceCondition(t *testing.T) {
	sorter := &DefaultURLSorter{}
	service := NewURLService(sorter)
	ctx := context.Background()

	urls := []string{"example.com", "test.com", "google.com"}
	numGoroutines := 100

	var wg sync.WaitGroup
	done := make(chan bool, numGoroutines)

	// Запускаємо багато goroutines одночасно
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			service.ProcessWithCallback(ctx, urls, func(result []string, err error) {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
					return
				}

				if len(result) != len(urls) {
					t.Errorf("Expected %d URLs, got %d", len(urls), len(result))
				}

				done <- true
			})
		}()
	}

	// Чекаємо всі callbacks
	go func() {
		wg.Wait()
		close(done)
	}()

	// Збираємо результати
	count := 0
	for range done {
		count++
	}

	if count != numGoroutines {
		t.Errorf("Expected %d callbacks, got %d", numGoroutines, count)
	}
}

// TestMockSorterRaceCondition перевіряє race condition з mock
func TestMockSorterRaceCondition(t *testing.T) {
	// Створюємо mock з shared state
	callCount := 0
	var mu sync.Mutex

	mockSorter := &MockURLSorter{
		SortFunc: func(ctx context.Context, urls []string) ([]string, error) {
			// ❌ БЕЗ mutex тут був би race condition!
			mu.Lock()
			callCount++
			mu.Unlock()
			return urls, nil
		},
	}

	service := NewURLService(mockSorter)
	ctx := context.Background()
	urls := []string{"test.com"}

	numGoroutines := 100
	var wg sync.WaitGroup

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			resultCh := service.ProcessAsync(ctx, urls)
			<-resultCh
		}()
	}

	wg.Wait()

	mu.Lock()
	finalCount := callCount
	mu.Unlock()

	if finalCount != numGoroutines {
		t.Errorf("Expected %d calls, got %d", numGoroutines, finalCount)
	}
}

// TestContextCancellationRaceCondition перевіряє race condition
// при одночасній відміні контексту
func TestContextCancellationRaceCondition(t *testing.T) {
	sorter := &SlowSorterForTest{delay: 100 * time.Millisecond}
	service := NewURLService(sorter)

	urls := []string{"example.com", "test.com"}
	numGoroutines := 50

	var wg sync.WaitGroup

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
			defer cancel()

			resultCh := service.ProcessAsync(ctx, urls)
			result := <-resultCh

			// Очікуємо помилку через timeout
			if result.Err == nil {
				t.Error("Expected timeout error, got nil")
			}
		}()
	}

	wg.Wait()
}

// TestSharedStateRaceCondition - приклад НЕБЕЗПЕЧНОГО коду
// (цей тест покаже race condition якщо запустити з -race)
func TestSharedStateRaceCondition(t *testing.T) {
	t.Skip("Цей тест НАВМИСНО має race condition для демонстрації")

	// ❌ НЕБЕЗПЕЧНО! Shared state без синхронізації
	counter := 0

	numGoroutines := 100
	var wg sync.WaitGroup

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// ❌ RACE CONDITION!
			counter++
		}()
	}

	wg.Wait()

	// Результат буде непередбачуваним!
	t.Logf("Counter: %d (очікували %d)", counter, numGoroutines)
}

// ═══════════════════════════════════════════════════════════════════
// HELPER TYPES
// ═══════════════════════════════════════════════════════════════════

type SlowSorterForTest struct {
	delay time.Duration
}

func (s *SlowSorterForTest) Sort(ctx context.Context, urls []string) ([]string, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-time.After(s.delay):
		return urls, nil
	}
}

// ═══════════════════════════════════════════════════════════════════
// BENCHMARK З RACE DETECTOR
// ═══════════════════════════════════════════════════════════════════

// go test -race -bench=. -benchtime=1s

func BenchmarkProcessAsyncConcurrent(b *testing.B) {
	sorter := &DefaultURLSorter{}
	service := NewURLService(sorter)
	ctx := context.Background()
	urls := []string{"example.com", "test.com"}

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			resultCh := service.ProcessAsync(ctx, urls)
			<-resultCh
		}
	})
}

func BenchmarkProcessCallbackConcurrent(b *testing.B) {
	sorter := &DefaultURLSorter{}
	service := NewURLService(sorter)
	ctx := context.Background()
	urls := []string{"example.com", "test.com"}

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			done := make(chan bool, 1)
			service.ProcessWithCallback(ctx, urls, func(result []string, err error) {
				done <- true
			})
			<-done
		}
	})
}
