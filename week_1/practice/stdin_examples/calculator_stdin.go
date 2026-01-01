package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	// Читаємо перше число
	fmt.Fprintln(os.Stderr, "Введіть перше число:")
	if !scanner.Scan() {
		fmt.Fprintln(os.Stderr, "Помилка: не вдалося прочитати перше число")
		os.Exit(1)
	}
	num1Str := strings.TrimSpace(scanner.Text())
	num1, err := strconv.ParseFloat(num1Str, 64)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Помилка: '%s' не є числом\n", num1Str)
		os.Exit(1)
	}

	// Читаємо операцію
	fmt.Fprintln(os.Stderr, "Введіть операцію (+, -, *, /):")
	if !scanner.Scan() {
		fmt.Fprintln(os.Stderr, "Помилка: не вдалося прочитати операцію")
		os.Exit(1)
	}
	operation := strings.TrimSpace(scanner.Text())

	// Читаємо друге число
	fmt.Fprintln(os.Stderr, "Введіть друге число:")
	if !scanner.Scan() {
		fmt.Fprintln(os.Stderr, "Помилка: не вдалося прочитати друге число")
		os.Exit(1)
	}
	num2Str := strings.TrimSpace(scanner.Text())
	num2, err := strconv.ParseFloat(num2Str, 64)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Помилка: '%s' не є числом\n", num2Str)
		os.Exit(1)
	}

	// Виконання операції
	var result float64
	var valid bool = true

	switch operation {
	case "+":
		result = num1 + num2
	case "-":
		result = num1 - num2
	case "*", "×":
		result = num1 * num2
	case "/", "÷":
		if num2 == 0 {
			fmt.Fprintln(os.Stderr, "Помилка: Ділення на нуль!")
			os.Exit(1)
		}
		result = num1 / num2
	default:
		fmt.Fprintf(os.Stderr, "Помилка: невідома операція '%s'\n", operation)
		valid = false
		os.Exit(1)
	}

	if valid {
		// Вивід результату (в STDOUT)
		fmt.Printf("%.2f %s %.2f = %.2f\n", num1, operation, num2, result)
	}

	// Перевірка на помилки сканера
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Помилка читання: %v\n", err)
		os.Exit(1)
	}
}

/*
ВИКОРИСТАННЯ:

1. Інтерактивно:
   go run calculator_stdin.go
   10
   +
   5

2. Через echo:
   echo -e "10\n+\n5" | go run calculator_stdin.go

3. Через printf:
   printf "20\n*\n3\n" | go run calculator_stdin.go

4. Heredoc:
   go run calculator_stdin.go << EOF
   100
   /
   4
   EOF

5. З файлу:
   cat > calc_input.txt << EOF
   15
   -
   7
   EOF

   go run calculator_stdin.go < calc_input.txt
   cat calc_input.txt | go run calculator_stdin.go

6. В одну лінію:
   echo "10" | go run calculator_stdin.go
   # (потім введете +, потім 5)
*/
