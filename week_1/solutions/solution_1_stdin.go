package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	// Читаємо ім'я
	fmt.Fprintln(os.Stderr, "=== Привітання користувача (STDIN) ===")
	fmt.Fprintln(os.Stderr, "\nВведіть ваше ім'я:")

	if !scanner.Scan() {
		fmt.Fprintln(os.Stderr, "❌ Помилка: не вдалося прочитати ім'я")
		os.Exit(1)
	}

	name := strings.TrimSpace(scanner.Text())
	if name == "" {
		fmt.Fprintln(os.Stderr, "❌ Помилка: ім'я не може бути порожнім")
		os.Exit(1)
	}

	// Читаємо вік (опційно)
	fmt.Fprintln(os.Stderr, "Введіть вік (або Enter для пропуску):")

	var age int
	var hasAge bool

	if scanner.Scan() {
		ageStr := strings.TrimSpace(scanner.Text())

		if ageStr != "" {
			parsedAge, err := strconv.Atoi(ageStr)

			if err != nil {
				fmt.Fprintf(os.Stderr, "⚠️  Попередження: '%s' не є числом, пропускаємо вік\n", ageStr)
			} else if parsedAge < 0 || parsedAge > 120 {
				fmt.Fprintf(os.Stderr, "⚠️  Попередження: вік %d поза діапазоном (0-120), пропускаємо\n", parsedAge)
			} else {
				age = parsedAge
				hasAge = true
			}
		}
	}

	// Перевірка помилок сканера
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Помилка читання: %v\n", err)
		os.Exit(1)
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

	// Вивід результату (в STDOUT)
	fmt.Println("\n═══════════════════════════════════════")
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

/*
ВИКОРИСТАННЯ:

1. Інтерактивно (вводите з клавіатури):
   go run solution_1_stdin.go

2. Через echo:
   echo -e "Іван\n25" | go run solution_1_stdin.go

3. Через printf:
   printf "Марія\n22\n" | go run solution_1_stdin.go

4. Heredoc:
   go run solution_1_stdin.go << EOF
   Петро
   30
   EOF

5. З файлу:
   cat > user_data.txt << EOF
   Олена
   28
   EOF

   go run solution_1_stdin.go < user_data.txt

6. Без віку:
   echo -e "Андрій\n" | go run solution_1_stdin.go

7. Pipe з іншою командою:
   cat user_data.txt | go run solution_1_stdin.go
*/
