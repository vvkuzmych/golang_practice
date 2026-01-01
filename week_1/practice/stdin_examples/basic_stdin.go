package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Базовий приклад читання з STDIN
func main() {
	fmt.Println("=== Читання з STDIN ===\n")

	// Створюємо scanner для читання з STDIN
	scanner := bufio.NewScanner(os.Stdin)

	// Читаємо перший рядок (ім'я)
	fmt.Fprintln(os.Stderr, "Введіть ім'я:")
	if scanner.Scan() {
		name := strings.TrimSpace(scanner.Text())
		fmt.Printf("Привіт, %s! 👋\n", name)
	}

	// Читаємо другий рядок (вік)
	fmt.Fprintln(os.Stderr, "Введіть вік:")
	if scanner.Scan() {
		age := strings.TrimSpace(scanner.Text())
		fmt.Printf("Вік: %s років\n", age)
	}

	// Перевірка на помилки
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Помилка читання: %v\n", err)
		os.Exit(1)
	}
}

/*
ВИКОРИСТАННЯ:

1. Інтерактивно:
   go run basic_stdin.go
   # Вводите дані вручну

2. Через echo (pipe):
   echo -e "Іван\n25" | go run basic_stdin.go

3. Через printf:
   printf "Марія\n22\n" | go run basic_stdin.go

4. Через heredoc:
   go run basic_stdin.go << EOF
   Петро
   30
   EOF

5. З файлу:
   echo -e "Олена\n28" > input.txt
   go run basic_stdin.go < input.txt

6. Через cat:
   cat input.txt | go run basic_stdin.go
*/
