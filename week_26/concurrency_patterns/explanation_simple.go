//go:build ignore
// +build ignore

package main

import "fmt"

// ============================================
// ПРИКЛАД 1: БЕЗ КАНАЛУ (НЕ ПРАЦЮЄ)
// ============================================

func tryToReturnValue(x int) {
	x = x * 2
	// Нічого не повертаємо
}

func exampleWithoutChannel() {
	fmt.Println("=== ПРИКЛАД 1: БЕЗ КАНАЛУ ===")

	num := 5
	tryToReturnValue(num)
	fmt.Printf("Результат: %d (не змінився!)\n", num)
	// Виведе: 5, бо функція працювала з копією
}

// ============================================
// ПРИКЛАД 2: З КАНАЛОМ (ПРАЦЮЄ!)
// ============================================

// Функція НЕ повертає нічого, але ПИШЕ В КАНАЛ!
func writeToChannel(resultCh chan int) {
	// resultCh - це як поштова скринька
	// Ми просто кладемо туди лист (число)
	resultCh <- 42 // Кладемо число 42 в скриньку

	fmt.Println("✅ Worker: Я поклав 42 у канал!")
	// Функція закінчилась, але 42 залишилось в каналі!
}

func exampleWithChannel() {
	fmt.Println("\n=== ПРИКЛАД 2: З КАНАЛОМ ===")

	// Створюємо поштову скриньку
	resultCh := make(chan int, 1)

	// Запускаємо функцію в окремій горутині
	go writeToChannel(resultCh)

	// Читаємо з поштової скриньки
	result := <-resultCh
	fmt.Printf("✅ Main: Я прочитав з каналу: %d\n", result)
}

// ============================================
// ПРИКЛАД 3: БАГАТО ПИСАЛЬНИКІВ (як у твоєму коді)
// ============================================

func worker(id int, resultCh chan int) {
	value := id * 10
	resultCh <- value // Кожен worker кидає число в ТУ САМУ скриньку
	fmt.Printf("Worker %d: поклав %d у канал\n", id, value)
}

func exampleMultipleWorkers() {
	fmt.Println("\n=== ПРИКЛАД 3: БАГАТО ПИСАЛЬНИКІВ ===")

	// Одна спільна поштова скринька для всіх
	resultCh := make(chan int, 3)

	// Запускаємо 3 воркерів - всі пишуть в ТУ САМУ скриньку
	go worker(1, resultCh)
	go worker(2, resultCh)
	go worker(3, resultCh)

	// Читаємо 3 результати з тієї самої скриньки
	fmt.Println("\nMain читає результати:")
	for i := 0; i < 3; i++ {
		result := <-resultCh
		fmt.Printf("  Прочитав: %d\n", result)
	}
}

// ============================================
// ПРИКЛАД 4: ВІЗУАЛІЗАЦІЯ (покроково)
// ============================================

func exampleStepByStep() {
	fmt.Println("\n=== ПРИКЛАД 4: ПОКРОКОВО ===")

	// Крок 1: Main створює канал
	resultCh := make(chan int, 1)
	fmt.Println("Крок 1: Main створив канал (поштову скриньку)")
	fmt.Printf("        Адреса каналу в пам'яті: %p\n", resultCh)

	// Крок 2: Запускаємо worker
	go func(ch chan int) {
		fmt.Println("Крок 2: Worker запустився")
		fmt.Printf("        Worker отримав канал з адресою: %p (та сама!)\n", ch)

		fmt.Println("Крок 3: Worker пише 999 в канал")
		ch <- 999

		fmt.Println("Крок 4: Worker закінчився (але 999 залишилось в каналі!)")
	}(resultCh) // Передаємо канал як параметр

	// Крок 5: Main читає з каналу
	fmt.Println("Крок 5: Main чекає на дані з каналу...")
	result := <-resultCh
	fmt.Printf("Крок 6: Main прочитав: %d\n", result)
	fmt.Println("\n💡 Worker нічого не повернув, але Main отримав дані!")
}

// ============================================
// ПРИКЛАД 5: АНАЛОГІЯ З РЕАЛЬНИМ СВІТОМ
// ============================================

func exampleRealWorldAnalogy() {
	fmt.Println("\n=== ПРИКЛАД 5: АНАЛОГІЯ ===")
	fmt.Println(`
Уяви:
1. Main (головний офіс) створює ПОШТОВУ СКРИНЬКУ на вулиці
2. Worker (працівник) отримує АДРЕСУ цієї скриньки
3. Worker кладе туди лист (число)
4. Worker йде додому (функція закінчується)
5. Main приходить до скриньки і забирає лист

Worker НЕ повертає лист в руки Main!
Worker просто КЛАДЕ лист у СПІЛЬНУ скриньку!
Main сам ЗАБИРАЄ лист з тієї скриньки!
	`)
}

func main() {
	exampleWithoutChannel()
	exampleWithChannel()
	exampleMultipleWorkers()
	exampleStepByStep()
	exampleRealWorldAnalogy()

	fmt.Println("\n==================================================")
	fmt.Println("ГОЛОВНЕ:")
	fmt.Println("- Канал створюється в main")
	fmt.Println("- Worker отримує ТОЙ САМИЙ канал (не копію!)")
	fmt.Println("- Worker пише в канал: ch <- значення")
	fmt.Println("- Main читає з каналу: значення <- ch")
	fmt.Println("- Worker НЕ повертає канал - він вже існує!")
	fmt.Println("==================================================")
}
