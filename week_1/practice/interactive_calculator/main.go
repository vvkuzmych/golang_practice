package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("╔════════════════════════════════════╗")
	fmt.Println("║   Інтерактивний Калькулятор        ║")
	fmt.Println("╚════════════════════════════════════╝\n")

	for {
		// Меню
		fmt.Println("\nВиберіть операцію:")
		fmt.Println("  1 - Додавання (+)")
		fmt.Println("  2 - Віднімання (-)")
		fmt.Println("  3 - Множення (×)")
		fmt.Println("  4 - Ділення (÷)")
		fmt.Println("  0 - Вихід")
		fmt.Print("\nВаш вибір: ")

		choiceStr, _ := reader.ReadString('\n')
		choiceStr = strings.TrimSpace(choiceStr)

		// Вихід з програми
		if choiceStr == "0" {
			fmt.Println("\n👋 До побачення!")
			break
		}

		// Перевірка вибору
		choice, err := strconv.Atoi(choiceStr)
		if err != nil || choice < 1 || choice > 4 {
			fmt.Println("❌ Невірний вибір! Спробуйте ще раз.")
			continue
		}

		// Введення першого числа
		fmt.Print("\nВведіть перше число: ")
		num1Str, _ := reader.ReadString('\n')
		num1Str = strings.TrimSpace(num1Str)
		num1, err := strconv.ParseFloat(num1Str, 64)

		if err != nil {
			fmt.Println("❌ Помилка: перше число некоректне")
			continue
		}

		// Введення другого числа
		fmt.Print("Введіть друге число: ")
		num2Str, _ := reader.ReadString('\n')
		num2Str = strings.TrimSpace(num2Str)
		num2, err := strconv.ParseFloat(num2Str, 64)

		if err != nil {
			fmt.Println("❌ Помилка: друге число некоректне")
			continue
		}

		// Виконання операції
		var result float64

		switch choice {
		case 1:
			result = num1 + num2
			fmt.Printf("\n%.2f + %.2f = %.2f\n", num1, num2, result)
		case 2:
			result = num1 - num2
			fmt.Printf("\n%.2f - %.2f = %.2f\n", num1, num2, result)
		case 3:
			result = num1 * num2
			fmt.Printf("\n%.2f × %.2f = %.2f\n", num1, num2, result)
		case 4:
			if num2 == 0 {
				fmt.Println("\n❌ Помилка: Ділення на нуль!")
				continue
			}
			result = num1 / num2
			fmt.Printf("\n%.2f ÷ %.2f = %.2f\n", num1, num2, result)
		}

		fmt.Printf("✅ Результат: %.2f\n", result)

		// Продовжити?
		fmt.Print("\nВиконати ще одну операцію? (так/ні): ")
		continueStr, _ := reader.ReadString('\n')
		continueStr = strings.TrimSpace(strings.ToLower(continueStr))

		if continueStr != "так" && continueStr != "yes" && continueStr != "y" {
			fmt.Println("\n👋 До побачення!")
			break
		}
	}
}
