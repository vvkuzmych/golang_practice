package main

import (
	"fmt"
	"sync"
	"time"
)

// ═══════════════════════════════════════════════════════════════════
// ПРИКЛАД 1: НЕБЕЗПЕЧНИЙ КОД З RACE CONDITION ❌
// ═══════════════════════════════════════════════════════════════════

type UnsafeCounter struct {
	count int // ❌ Небезпечний доступ з різних goroutines!
}

func (c *UnsafeCounter) Increment() {
	c.count++ // ❌ RACE CONDITION! Не атомарна операція!
}

func (c *UnsafeCounter) Value() int {
	return c.count // ❌ RACE CONDITION! Читання без синхронізації!
}

func demoRaceCondition() {
	fmt.Println("╔══════════════════════════════════════════════════════════╗")
	fmt.Println("║   ПРИКЛАД 1: НЕБЕЗПЕЧНИЙ КОД (Race Condition) ❌         ║")
	fmt.Println("╚══════════════════════════════════════════════════════════╝")
	fmt.Println()

	counter := &UnsafeCounter{}

	// Запускаємо 100 goroutines, кожна інкрементує 1000 разів
	numGoroutines := 100
	incrementsPerGoroutine := 1000

	for i := 0; i < numGoroutines; i++ {
		go func() {
			for j := 0; j < incrementsPerGoroutine; j++ {
				counter.Increment() // ❌ RACE CONDITION!
			}
		}()
	}

	time.Sleep(1 * time.Second) // Чекаємо завершення

	expected := numGoroutines * incrementsPerGoroutine
	actual := counter.Value()

	fmt.Printf("Очікуваний результат: %d\n", expected)
	fmt.Printf("Фактичний результат:  %d\n", actual)

	if actual != expected {
		fmt.Printf("❌ ПОМИЛКА! Втрачено %d інкрементів через race condition!\n", expected-actual)
	} else {
		fmt.Println("✓ Результат правильний (але тільки випадково!)")
	}
	fmt.Println()
}

// ═══════════════════════════════════════════════════════════════════
// ПРИКЛАД 2: БЕЗПЕЧНИЙ КОД З MUTEX ✅
// ═══════════════════════════════════════════════════════════════════

type SafeCounterMutex struct {
	mu    sync.Mutex
	count int
}

func (c *SafeCounterMutex) Increment() {
	c.mu.Lock()         // ✅ Блокуємо доступ
	defer c.mu.Unlock() // ✅ Розблокуємо після виконання
	c.count++
}

func (c *SafeCounterMutex) Value() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.count
}

func demoSafeWithMutex() {
	fmt.Println("╔══════════════════════════════════════════════════════════╗")
	fmt.Println("║   ПРИКЛАД 2: БЕЗПЕЧНИЙ КОД (Mutex) ✅                    ║")
	fmt.Println("╚══════════════════════════════════════════════════════════╝")
	fmt.Println()

	counter := &SafeCounterMutex{}
	var wg sync.WaitGroup

	numGoroutines := 100
	incrementsPerGoroutine := 1000

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < incrementsPerGoroutine; j++ {
				counter.Increment() // ✅ Безпечно!
			}
		}()
	}

	wg.Wait() // Чекаємо всі goroutines

	expected := numGoroutines * incrementsPerGoroutine
	actual := counter.Value()

	fmt.Printf("Очікуваний результат: %d\n", expected)
	fmt.Printf("Фактичний результат:  %d\n", actual)

	if actual == expected {
		fmt.Println("✅ Результат ЗАВЖДИ правильний! Mutex працює!")
	}
	fmt.Println()
}

// ═══════════════════════════════════════════════════════════════════
// ПРИКЛАД 3: БЕЗПЕЧНИЙ КОД З CHANNELS ✅
// ═══════════════════════════════════════════════════════════════════

type SafeCounterChannel struct {
	ch chan int
}

func NewSafeCounterChannel() *SafeCounterChannel {
	c := &SafeCounterChannel{
		ch: make(chan int),
	}
	// Запускаємо goroutine для обробки запитів
	go c.run()
	return c
}

func (c *SafeCounterChannel) run() {
	count := 0
	for increment := range c.ch {
		count += increment
	}
}

func (c *SafeCounterChannel) Increment() {
	c.ch <- 1 // ✅ Безпечно! Channel thread-safe
}

func demoSafeWithChannels() {
	fmt.Println("╔══════════════════════════════════════════════════════════╗")
	fmt.Println("║   ПРИКЛАД 3: БЕЗПЕЧНИЙ КОД (Channels) ✅                 ║")
	fmt.Println("╚══════════════════════════════════════════════════════════╝")
	fmt.Println()

	counter := NewSafeCounterChannel()
	var wg sync.WaitGroup

	numGoroutines := 100
	incrementsPerGoroutine := 1000

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < incrementsPerGoroutine; j++ {
				counter.Increment() // ✅ Безпечно!
			}
		}()
	}

	wg.Wait()
	close(counter.ch)
	time.Sleep(100 * time.Millisecond)

	fmt.Println("✅ Channels теж працюють безпечно!")
	fmt.Println()
}

// ═══════════════════════════════════════════════════════════════════
// ПРИКЛАД 4: RACE CONDITION В МАПАХ ❌
// ═══════════════════════════════════════════════════════════════════

func demoMapRaceCondition() {
	fmt.Println("╔══════════════════════════════════════════════════════════╗")
	fmt.Println("║   ПРИКЛАД 4: Race Condition в Maps ❌                    ║")
	fmt.Println("╚══════════════════════════════════════════════════════════╝")
	fmt.Println()

	// ❌ НЕБЕЗПЕЧНО!
	unsafeMap := make(map[string]int)

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			key := fmt.Sprintf("key-%d", id)
			unsafeMap[key] = id // ❌ RACE CONDITION!
		}(i)
	}
	wg.Wait()

	fmt.Println("❌ Цей код може паніку викликати при race detector!")
	fmt.Println("   Запустіть з -race щоб побачити")
	fmt.Println()
}

// ═══════════════════════════════════════════════════════════════════
// ПРИКЛАД 5: БЕЗПЕЧНА МАПА З SYNC.MAP ✅
// ═══════════════════════════════════════════════════════════════════

func demoSyncMap() {
	fmt.Println("╔══════════════════════════════════════════════════════════╗")
	fmt.Println("║   ПРИКЛАД 5: Безпечна Map (sync.Map) ✅                  ║")
	fmt.Println("╚══════════════════════════════════════════════════════════╝")
	fmt.Println()

	// ✅ БЕЗПЕЧНО!
	var safeMap sync.Map

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			key := fmt.Sprintf("key-%d", id)
			safeMap.Store(key, id) // ✅ Thread-safe!
		}(i)
	}
	wg.Wait()

	// Читання
	safeMap.Range(func(key, value interface{}) bool {
		fmt.Printf("  %s = %d\n", key, value)
		return true
	})

	fmt.Println()
	fmt.Println("✅ sync.Map безпечна для concurrent доступу!")
	fmt.Println()
}

// ═══════════════════════════════════════════════════════════════════
// ІНСТРУКЦІЇ
// ═══════════════════════════════════════════════════════════════════

func printInstructions() {
	fmt.Println("╔══════════════════════════════════════════════════════════╗")
	fmt.Println("║         ЯК ПЕРЕВІРИТИ RACE CONDITION В GO                ║")
	fmt.Println("╚══════════════════════════════════════════════════════════╝")
	fmt.Println()
	fmt.Println("🔍 Go має вбудований Race Detector!")
	fmt.Println()
	fmt.Println("┌─────────────────────────────────────────────────────────┐")
	fmt.Println("│  КОМАНДИ ДЛЯ ПЕРЕВІРКИ:                                 │")
	fmt.Println("└─────────────────────────────────────────────────────────┘")
	fmt.Println()
	fmt.Println("1️⃣  Запустити програму з race detector:")
	fmt.Println("   go run -race race_condition_demo.go")
	fmt.Println()
	fmt.Println("2️⃣  Запустити тести з race detector:")
	fmt.Println("   go test -race")
	fmt.Println()
	fmt.Println("3️⃣  Збілдити з race detector:")
	fmt.Println("   go build -race myprogram.go")
	fmt.Println()
	fmt.Println("4️⃣  Перевірити весь проект:")
	fmt.Println("   go test -race ./...")
	fmt.Println()
	fmt.Println("┌─────────────────────────────────────────────────────────┐")
	fmt.Println("│  ЩО ШУКАЄ RACE DETECTOR:                                │")
	fmt.Println("└─────────────────────────────────────────────────────────┘")
	fmt.Println()
	fmt.Println("  ❌ Одночасний запис в одну змінну з різних goroutines")
	fmt.Println("  ❌ Одночасне читання/запис без синхронізації")
	fmt.Println("  ❌ Небезпечний доступ до maps")
	fmt.Println("  ❌ Небезпечний доступ до slices")
	fmt.Println()
	fmt.Println("┌─────────────────────────────────────────────────────────┐")
	fmt.Println("│  ПРИКЛАД ВИВОДУ RACE DETECTOR:                          │")
	fmt.Println("└─────────────────────────────────────────────────────────┘")
	fmt.Println()
	fmt.Println("  ==================")
	fmt.Println("  WARNING: DATA RACE")
	fmt.Println("  Write at 0x00c000014098 by goroutine 7:")
	fmt.Println("    main.(*UnsafeCounter).Increment()")
	fmt.Println("      /path/to/file.go:15 +0x38")
	fmt.Println()
	fmt.Println("  Previous write at 0x00c000014098 by goroutine 6:")
	fmt.Println("    main.(*UnsafeCounter).Increment()")
	fmt.Println("      /path/to/file.go:15 +0x38")
	fmt.Println("  ==================")
	fmt.Println()
	fmt.Println("┌─────────────────────────────────────────────────────────┐")
	fmt.Println("│  ЯК ВИПРАВИТИ RACE CONDITIONS:                          │")
	fmt.Println("└─────────────────────────────────────────────────────────┘")
	fmt.Println()
	fmt.Println("  ✅ 1. Використовуй sync.Mutex / sync.RWMutex")
	fmt.Println("  ✅ 2. Використовуй channels для комунікації")
	fmt.Println("  ✅ 3. Використовуй sync.Map для concurrent maps")
	fmt.Println("  ✅ 4. Використовуй sync/atomic для простих операцій")
	fmt.Println("  ✅ 5. Уникай shared state (якщо можливо)")
	fmt.Println()
	fmt.Println("┌─────────────────────────────────────────────────────────┐")
	fmt.Println("│  ВАЖЛИВО:                                               │")
	fmt.Println("└─────────────────────────────────────────────────────────┘")
	fmt.Println()
	fmt.Println("  ⚠️  Race detector збільшує час виконання ~10x")
	fmt.Println("  ⚠️  Race detector збільшує використання пам'яті ~10x")
	fmt.Println("  ⚠️  Використовуй тільки для тестування/розробки")
	fmt.Println("  ⚠️  НЕ використовуй в production білдах")
	fmt.Println()
}

// ═══════════════════════════════════════════════════════════════════
// MAIN
// ═══════════════════════════════════════════════════════════════════

func main() {
	printInstructions()

	fmt.Println("╔══════════════════════════════════════════════════════════╗")
	fmt.Println("║              ДЕМОНСТРАЦІЯ ПРИКЛАДІВ                      ║")
	fmt.Println("╚══════════════════════════════════════════════════════════╝")
	fmt.Println()

	demoRaceCondition()
	time.Sleep(500 * time.Millisecond)

	demoSafeWithMutex()
	time.Sleep(500 * time.Millisecond)

	demoSafeWithChannels()
	time.Sleep(500 * time.Millisecond)

	demoMapRaceCondition()
	time.Sleep(500 * time.Millisecond)

	demoSyncMap()

	fmt.Println("╔══════════════════════════════════════════════════════════╗")
	fmt.Println("║                   СПРОБУЙ ЗАРАЗ:                         ║")
	fmt.Println("╚══════════════════════════════════════════════════════════╝")
	fmt.Println()
	fmt.Println("Запусти цю програму з race detector:")
	fmt.Println()
	fmt.Println("  go run -race race_condition_demo.go")
	fmt.Println()
	fmt.Println("Race detector покаже всі проблеми! 🔍")
	fmt.Println()
}
