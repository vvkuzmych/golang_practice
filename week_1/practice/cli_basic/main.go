package main

import (
	"fmt"
	"os"
)

// Назва програми (для допомоги)
const programName = "cli-demo"

func main() {
	// Приклад роботи з аргументами командного рядка

	// os.Args - slice з усіма аргументами
	// os.Args[0] - назва програми
	// os.Args[1:] - аргументи користувача

	fmt.Println("=== Базовий CLI Приклад ===\n")

	// Показати всі аргументи
	fmt.Printf("Загальна кількість аргументів: %d\n", len(os.Args))
	fmt.Printf("Назва програми: %s\n\n", os.Args[0])

	// Якщо є аргументи від користувача
	if len(os.Args) > 1 {
		fmt.Println("Отримані аргументи:")
		for i, arg := range os.Args[1:] {
			fmt.Printf("  [%d] %s\n", i+1, arg)
		}
	} else {
		fmt.Println("Аргументи не передано.")
		fmt.Println("\nВикористання:")
		fmt.Printf("  %s <аргумент1> <аргумент2> ...\n", programName)
		fmt.Println("\nПриклад:")
		fmt.Printf("  %s привіт світ\n", programName)
	}

	fmt.Println("\n=== Кінець програми ===")
}
