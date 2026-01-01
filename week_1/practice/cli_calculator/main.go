package main

import (
	"fmt"
	"os"
	"strconv"
)

// Назва програми
const programName = "calculator"

func main() {
	fmt.Println("=== CLI Калькулятор ===\n")

	// Перевірка кількості аргументів
	if len(os.Args) != 4 {
		printUsage()
		os.Exit(1)
	}

	// Отримання аргументів
	num1Str := os.Args[1]
	operation := os.Args[2]
	num2Str := os.Args[3]

	// Конвертація рядків у числа
	num1, err := strconv.ParseFloat(num1Str, 64)
	if err != nil {
		fmt.Printf("Помилка: '%s' не є числом\n", num1Str)
		os.Exit(1)
	}

	num2, err := strconv.ParseFloat(num2Str, 64)
	if err != nil {
		fmt.Printf("Помилка: '%s' не є числом\n", num2Str)
		os.Exit(1)
	}

	// Виконання операції
	var result float64
	var validOperation bool = true

	switch operation {
	case "+", "add", "додати":
		result = num1 + num2
		fmt.Printf("%.2f + %.2f = %.2f\n", num1, num2, result)
	case "-", "sub", "відняти":
		result = num1 - num2
		fmt.Printf("%.2f - %.2f = %.2f\n", num1, num2, result)
	case "*", "mul", "множити":
		result = num1 * num2
		fmt.Printf("%.2f × %.2f = %.2f\n", num1, num2, result)
	case "/", "div", "ділити":
		if num2 == 0 {
			fmt.Println("Помилка: Ділення на нуль!")
			os.Exit(1)
		}
		result = num1 / num2
		fmt.Printf("%.2f ÷ %.2f = %.2f\n", num1, num2, result)
	default:
		validOperation = false
		fmt.Printf("Помилка: Невідома операція '%s'\n", operation)
		fmt.Println("\nДоступні операції: +, -, *, /")
		os.Exit(1)
	}

	if validOperation {
		fmt.Printf("\n✅ Результат: %.2f\n", result)
	}
}

func printUsage() {
	fmt.Println("Використання:")
	fmt.Printf("  %s <число1> <операція> <число2>\n\n", programName)

	fmt.Println("Операції:")
	fmt.Println("  +  або  add     - Додавання")
	fmt.Println("  -  або  sub     - Віднімання")
	fmt.Println("  *  або  mul     - Множення")
	fmt.Println("  /  або  div     - Ділення")

	fmt.Println("\nПриклади:")
	fmt.Printf("  %s 10 + 5\n", programName)
	fmt.Printf("  %s 20 - 7\n", programName)
	fmt.Printf("  %s 3.14 mul 2\n", programName)
	fmt.Printf("  %s 100 / 4\n", programName)
}
