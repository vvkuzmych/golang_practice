package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	fmt.Println("=== Інтерактивний ввід даних ===\n")

	// Створюємо reader для зчитування з консолі
	reader := bufio.NewReader(os.Stdin)

	// 1. Зчитування рядка (string)
	fmt.Print("Введіть ваше ім'я: ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name) // Видаляємо \n в кінці

	// 2. Зчитування числа (int)
	fmt.Print("Введіть ваш вік: ")
	ageStr, _ := reader.ReadString('\n')
	ageStr = strings.TrimSpace(ageStr)
	age, err := strconv.Atoi(ageStr)

	if err != nil {
		fmt.Println("❌ Помилка: вік має бути числом")
		return
	}

	// 3. Зчитування yes/no (bool)
	fmt.Print("Ви студент? (так/ні): ")
	studentStr, _ := reader.ReadString('\n')
	studentStr = strings.TrimSpace(strings.ToLower(studentStr))
	isStudent := studentStr == "так" || studentStr == "yes" || studentStr == "y"

	// Вивід результату
	fmt.Println("\n--- Ваші дані ---")
	fmt.Printf("Ім'я: %s\n", name)
	fmt.Printf("Вік: %d років\n", age)
	fmt.Printf("Студент: %t\n", isStudent)

	// Привітання
	fmt.Printf("\nПривіт, %s! 👋\n", name)
	if isStudent {
		fmt.Println("Успіхів у навчанні! 📚")
	} else {
		fmt.Println("Гарного дня! ☀️")
	}
}
