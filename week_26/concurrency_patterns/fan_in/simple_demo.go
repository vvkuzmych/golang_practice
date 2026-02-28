package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("=== Проста демонстрація проблеми ===\n")

	// Створюємо слайс значень
	values := []string{"channel1", "channel2", "channel3"}

	fmt.Println("❌ НЕПРАВИЛЬНО (без параметра):")
	fmt.Println("Запускаємо 3 горутини")
	fmt.Println("Очікуємо: channel1, channel2, channel3")
	fmt.Print("Отримуємо: ")

	// ❌ БАГ: змінна val захоплюється неправильно
	for _, val := range values {
		go func() {
			time.Sleep(10 * time.Millisecond) // Даємо циклу час закінчитися
			fmt.Print(val, " ")               // Всі горутини бачать останнє значення!
		}()
	}

	time.Sleep(100 * time.Millisecond)

	fmt.Println("\n\n✅ ПРАВИЛЬНО (з параметром):")
	fmt.Println("Запускаємо 3 горутини")
	fmt.Println("Очікуємо: channel1, channel2, channel3")
	fmt.Print("Отримуємо: ")

	// ✅ Передаємо як параметр
	for _, val := range values {
		go func(v string) {
			time.Sleep(10 * time.Millisecond)
			fmt.Print(v, " ")
		}(val) // Явна передача значення
	}

	time.Sleep(100 * time.Millisecond)

	fmt.Println("\n\n=== Пояснення ===")
	fmt.Println(`
У першому випадку:
- Змінна 'val' - це ОДНА змінна, яка перевикористовується
- Цикл закінчується швидше, ніж запускаються горутини
- Коли горутини нарешті читають 'val', вона вже = "channel3"
- Результат: усі горутини бачать "channel3"

У другому випадку:
- Кожна горутина отримує КОПІЮ значення через параметр
- Функція викликається негайно з поточним значенням
- Результат: кожна горутина має своє правильне значення
`)

	fmt.Println("=== Те саме з індексами ===\n")

	fmt.Println("❌ НЕПРАВИЛЬНО:")
	for i := range values {
		go func() {
			time.Sleep(10 * time.Millisecond)
			fmt.Printf("index=%d ", i) // Всі бачать i=2
		}()
	}
	time.Sleep(100 * time.Millisecond)

	fmt.Println("\n✅ ПРАВИЛЬНО:")
	for i := range values {
		go func(index int) {
			time.Sleep(10 * time.Millisecond)
			fmt.Printf("index=%d ", index)
		}(i) // Передаємо i як параметр
	}
	time.Sleep(100 * time.Millisecond)

	fmt.Println("\n\nЦе класична Go пастка! Завжди передавайте змінні циклу як параметри!")
}
