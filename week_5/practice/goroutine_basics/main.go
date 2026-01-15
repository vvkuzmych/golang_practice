package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// ============= Examples =============

func example1_BasicGoroutine() {
	fmt.Println("1️⃣ Базовий Goroutine")
	fmt.Println("─────────────────────────────────────────")

	// Синхронний виклик
	fmt.Println("Main: Start")

	// Асинхронний виклик (goroutine)
	go func() {
		fmt.Println("   Goroutine: Running...")
		time.Sleep(50 * time.Millisecond)
		fmt.Println("   Goroutine: Done!")
	}()

	// Даємо час goroutine виконатись
	time.Sleep(100 * time.Millisecond)
	fmt.Println("Main: End\n")
}

func example2_MultipleGoroutines() {
	fmt.Println("2️⃣ Кілька Goroutines")
	fmt.Println("─────────────────────────────────────────")

	for i := 1; i <= 5; i++ {
		go func(id int) {
			fmt.Printf("   Goroutine %d: started\n", id)
			time.Sleep(50 * time.Millisecond)
			fmt.Printf("   Goroutine %d: finished\n", id)
		}(i) // ВАЖЛИВО: передаємо i як параметр!
	}

	time.Sleep(100 * time.Millisecond)
	fmt.Println()
}

func example3_WaitGroupBasic() {
	fmt.Println("3️⃣ WaitGroup для очікування")
	fmt.Println("─────────────────────────────────────────")

	var wg sync.WaitGroup

	// Запускаємо 3 goroutines
	for i := 1; i <= 3; i++ {
		wg.Add(1) // Збільшуємо counter ПЕРЕД запуском goroutine
		go func(id int) {
			defer wg.Done() // Зменшуємо counter при завершенні
			fmt.Printf("   Worker %d: working...\n", id)
			time.Sleep(50 * time.Millisecond)
			fmt.Printf("   Worker %d: done\n", id)
		}(i)
	}

	// Чекаємо поки ВСІ goroutines завершаться
	fmt.Println("Main: waiting for workers...")
	wg.Wait()
	fmt.Println("✓ All workers finished\n")
}

func example4_RaceCondition() {
	fmt.Println("4️⃣ Race Condition (WITHOUT sync)")
	fmt.Println("─────────────────────────────────────────")

	counter := 0

	// Запускаємо 100 goroutines, кожна збільшує counter
	for i := 0; i < 100; i++ {
		go func() {
			counter++ // ❌ Race condition!
		}()
	}

	time.Sleep(100 * time.Millisecond)
	fmt.Printf("Counter: %d (очікувалось 100)\n", counter)
	fmt.Println("⚠️  Результат може бути меншим через race condition!")
	fmt.Println()
}

func example5_FixRaceWithMutex() {
	fmt.Println("5️⃣ Виправлення Race Condition (Mutex)")
	fmt.Println("─────────────────────────────────────────")

	counter := 0
	var mu sync.Mutex
	var wg sync.WaitGroup

	// Запускаємо 100 goroutines
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mu.Lock() // ✅ Захищаємо critical section
			counter++
			mu.Unlock() // ✅ Звільняємо lock
		}()
	}

	wg.Wait()
	fmt.Printf("Counter: %d (точно 100!)\n", counter)
	fmt.Println("✓ Mutex забезпечив безпечний доступ\n")
}

func example6_FixRaceWithAtomic() {
	fmt.Println("6️⃣ Виправлення Race Condition (Atomic)")
	fmt.Println("─────────────────────────────────────────")

	var counter int64 // atomic працює з int32/int64
	var wg sync.WaitGroup

	// Запускаємо 100 goroutines
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			atomic.AddInt64(&counter, 1) // ✅ Атомарна операція
		}()
	}

	wg.Wait()
	fmt.Printf("Counter: %d (точно 100!)\n", counter)
	fmt.Println("✓ Atomic забезпечив lock-free синхронізацію\n")
}

func example7_GoroutineLeak() {
	fmt.Println("7️⃣ Goroutine Leak (blocked goroutine)")
	fmt.Println("─────────────────────────────────────────")

	ch := make(chan int)

	// ❌ Goroutine блокується назавжди (leak!)
	go func() {
		fmt.Println("   Goroutine: waiting to send...")
		ch <- 42 // Блокується - ніхто не читає
		fmt.Println("   Goroutine: sent (never happens)")
	}()

	time.Sleep(100 * time.Millisecond)
	fmt.Println("⚠️  Goroutine заблокована і ніколи не завершиться!")
	fmt.Println("⚠️  Це goroutine leak - витік ресурсів!")
	fmt.Println()
}

func example8_PreventGoroutineLeak() {
	fmt.Println("8️⃣ Запобігання Goroutine Leak")
	fmt.Println("─────────────────────────────────────────")

	ch := make(chan int)
	done := make(chan bool)

	// Goroutine з можливістю завершення
	go func() {
		select {
		case ch <- 42:
			fmt.Println("   Goroutine: sent")
		case <-done:
			fmt.Println("   Goroutine: cancelled (no leak!)")
			return
		}
	}()

	time.Sleep(50 * time.Millisecond)

	// Сигналізуємо про завершення
	close(done)

	time.Sleep(50 * time.Millisecond)
	fmt.Println("✓ Goroutine завершилась коректно\n")
}

func example9_AnonymousVsClosure() {
	fmt.Println("9️⃣ Anonymous goroutine vs Closure")
	fmt.Println("─────────────────────────────────────────")

	// ❌ ПОМИЛКА: closure захоплює змінну за посиланням
	fmt.Println("Closure (WRONG):")
	for i := 1; i <= 3; i++ {
		go func() {
			fmt.Printf("   Value: %d\n", i) // i змінюється!
		}()
	}
	time.Sleep(50 * time.Millisecond)

	// ✅ ПРАВИЛЬНО: передаємо як параметр
	fmt.Println("\nParameter (CORRECT):")
	for i := 1; i <= 3; i++ {
		go func(id int) {
			fmt.Printf("   Value: %d\n", id) // Фіксоване значення
		}(i)
	}
	time.Sleep(50 * time.Millisecond)
	fmt.Println()
}

func example10_GoroutineLifecycle() {
	fmt.Println("🔟 Goroutine Lifecycle")
	fmt.Println("─────────────────────────────────────────")

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		fmt.Println("   1. Created (go keyword)")
		time.Sleep(20 * time.Millisecond)

		fmt.Println("   2. Running (scheduled by Go runtime)")
		time.Sleep(20 * time.Millisecond)

		fmt.Println("   3. Potentially blocked (I/O, channel, etc.)")
		time.Sleep(20 * time.Millisecond)

		fmt.Println("   4. Resumed (unblocked)")
		time.Sleep(20 * time.Millisecond)

		fmt.Println("   5. Finished (function returns)")
	}()

	wg.Wait()
	fmt.Println("✓ Goroutine lifecycle completed\n")
}

// ============= Main =============

func main() {
	fmt.Println("╔══════════════════════════════════════════╗")
	fmt.Println("║       Goroutine Basics Examples         ║")
	fmt.Println("╚══════════════════════════════════════════╝")
	fmt.Println()

	example1_BasicGoroutine()
	example2_MultipleGoroutines()
	example3_WaitGroupBasic()
	example4_RaceCondition()
	example5_FixRaceWithMutex()
	example6_FixRaceWithAtomic()
	example7_GoroutineLeak()
	example8_PreventGoroutineLeak()
	example9_AnonymousVsClosure()
	example10_GoroutineLifecycle()

	fmt.Println("📝 Висновки:")
	fmt.Println("─────────────────────────────────────────")
	fmt.Println("✅ Goroutines - легкі та швидкі")
	fmt.Println("✅ WaitGroup для очікування завершення")
	fmt.Println("✅ Mutex або Atomic для race conditions")
	fmt.Println("✅ Завжди забезпечуйте завершення goroutines")
	fmt.Println("✅ Передавайте змінні як параметри в loops")
	fmt.Println("✅ Уникайте goroutine leaks (blocked forever)")
}
