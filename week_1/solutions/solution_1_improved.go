package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Константа для назви програми
const programName = "greet"

func main() {
	fmt.Println("=== Привітання користувача ===\n")

	// Перевірка кількості аргументів
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	// Отримання імені (обов'язковий аргумент)
	name := os.Args[1]

	// Отримання віку (опційний аргумент)
	var age int
	var hasAge bool = false

	if len(os.Args) >= 3 {
		ageStr := os.Args[2]
		parsedAge, err := strconv.Atoi(ageStr)

		if err != nil {
			fmt.Printf("⚠️  Помилка: '%s' не є коректним віком\n", ageStr)
			fmt.Println("Вік має бути числом від 0 до 120")
			os.Exit(1)
		}

		// Валідація віку
		if parsedAge < 0 || parsedAge > 120 {
			fmt.Printf("⚠️  Помилка: Вік %d поза допустимим діапазоном (0-120)\n", parsedAge)
			os.Exit(1)
		}

		age = parsedAge
		hasAge = true
	}

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

	// Інформація про змінні
	fmt.Println("\n--- Інформація про змінні ---")
	fmt.Printf("  name: %q (тип: %T)\n", name, name)

	if hasAge {
		fmt.Printf("  age: %d (тип: %T)\n", age, age)
		fmt.Printf("  hasAge: %t (тип: %T)\n", hasAge, hasAge)
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
}

func printUsage() {
	fmt.Println("❌ Помилка: не вказано ім'я\n")
	fmt.Println("Використання:")
	fmt.Printf("  %s <ім'я> [вік]\n\n", programName)
	fmt.Println("Аргументи:")
	fmt.Println("  <ім'я>  - Ваше ім'я (обов'язково)")
	fmt.Println("  [вік]   - Ваш вік, число 0-120 (опційно)")
	fmt.Println("\nПриклади:")
	fmt.Printf("  %s Іван\n", programName)
	fmt.Printf("  %s Марія 25\n", programName)
}
