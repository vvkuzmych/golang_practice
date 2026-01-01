package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Helper функції для читання даних
func readString(reader *bufio.Reader, prompt string) string {
	fmt.Print(prompt)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func readInt(reader *bufio.Reader, prompt string, min, max int) (int, error) {
	for {
		input := readString(reader, prompt)
		num, err := strconv.Atoi(input)

		if err != nil {
			fmt.Printf("❌ Помилка: введіть число\n")
			continue
		}

		if num < min || num > max {
			fmt.Printf("❌ Помилка: число має бути між %d і %d\n", min, max)
			continue
		}

		return num, nil
	}
}

func readOptionalInt(reader *bufio.Reader, prompt string) (int, bool) {
	input := readString(reader, prompt)

	if input == "" {
		return 0, false
	}

	num, err := strconv.Atoi(input)
	if err != nil {
		fmt.Println("⚠️  Некоректне число, пропускаємо")
		return 0, false
	}

	if num < 0 || num > 120 {
		fmt.Println("⚠️  Вік має бути 0-120, пропускаємо")
		return 0, false
	}

	return num, true
}

func main() {
	fmt.Println("╔════════════════════════════════════════╗")
	fmt.Println("║  Інтерактивне привітання користувача   ║")
	fmt.Println("╚════════════════════════════════════════╝\n")

	reader := bufio.NewReader(os.Stdin)

	// Введення імені (обов'язково)
	var name string
	for {
		name = readString(reader, "Введіть ваше ім'я: ")
		if name != "" {
			break
		}
		fmt.Println("❌ Ім'я не може бути порожнім!")
	}

	// Введення віку (опційно)
	fmt.Print("Введіть ваш вік (Enter щоб пропустити): ")
	age, hasAge := readOptionalInt(reader, "")

	fmt.Println() // Порожній рядок для красоти

	// Привітання залежно від часу доби
	hour := time.Now().Hour()
	var greeting string

	switch {
	case hour >= 5 && hour < 12:
		greeting = "Доброго ранку"
	case hour >= 12 && hour < 17:
		greeting = "Доброго дня"
	case hour >= 17 && hour < 23:
		greeting = "Доброго вечора"
	default:
		greeting = "Доброї ночі"
	}

	// Вивід привітання
	fmt.Println("═══════════════════════════════════════")
	fmt.Printf("%s, %s! 👋\n", greeting, name)

	if hasAge {
		fmt.Printf("Тобі %d років.\n", age)

		// Додаткова інформація залежно від віку
		switch {
		case age < 18:
			fmt.Println("Ти ще молодий, багато всього попереду! 🌟")
		case age >= 18 && age < 65:
			fmt.Println("Продуктивного дня! 💼")
		default:
			fmt.Println("Бажаю здоров'я та гарного настрою! 🌺")
		}
	}

	fmt.Println("Радий тебе бачити!")
	fmt.Println("═══════════════════════════════════════")

	// Інформація про змінні
	fmt.Println("\n--- Інформація про змінні ---")
	fmt.Printf("  name: %q (тип: %T)\n", name, name)

	if hasAge {
		fmt.Printf("  age: %d (тип: %T)\n", age, age)
		fmt.Printf("  hasAge: %t (тип: %T)\n", hasAge, hasAge)
	} else {
		fmt.Println("  age: не вказано")
	}

	fmt.Printf("  greeting: %q (тип: %T)\n", greeting, greeting)
	fmt.Printf("  hour: %d (тип: %T)\n", hour, hour)

	// Демонстрація різних форматів
	fmt.Println("\n--- Різні формати виводу ---")
	fmt.Printf("%%s:  %s\n", name)
	fmt.Printf("%%q:  %q\n", name)
	fmt.Printf("%%v:  %v\n", name)
	fmt.Printf("%%#v: %#v\n", name)
	fmt.Printf("%%T:  %T\n", name)

	if hasAge {
		fmt.Println("\nВік:")
		fmt.Printf("%%d:  %d\n", age)
		fmt.Printf("%%v:  %v\n", age)
		fmt.Printf("%%T:  %T\n", age)
	}

	fmt.Println("\n✅ Програма завершена")
}
