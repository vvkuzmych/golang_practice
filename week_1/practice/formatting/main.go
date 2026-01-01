package main

import "fmt"

type Person struct {
	FirstName string
	LastName  string
	Age       int
	City      string
	Salary    float64
	IsStudent bool
}

func main() {
	fmt.Println("=== Практика fmt.Printf ===\n")

	// Базові типи
	demonstrateBasicTypes()

	// Структури
	demonstrateStructs()

	// Форматування чисел
	demonstrateNumbers()

	// Спеціальні формати
	demonstrateSpecial()
}

func demonstrateBasicTypes() {
	fmt.Println("--- 1. Базові типи ---\n")

	name := "Іван"
	age := 25
	height := 180.5
	isStudent := true

	// %s - string
	fmt.Printf("%%s: Ім'я: %s\n", name)

	// %d - decimal integer
	fmt.Printf("%%d: Вік: %d років\n", age)

	// %f - float
	fmt.Printf("%%f: Зріст: %f см\n", height)
	fmt.Printf("%%.1f: Зріст: %.1f см\n", height)
	fmt.Printf("%%.2f: Зріст: %.2f см\n", height)

	// %t - boolean
	fmt.Printf("%%t: Студент: %t\n", isStudent)

	// %v - default format (будь-який тип)
	fmt.Printf("%%v: Ім'я=%v, Вік=%v, Зріст=%v, Студент=%v\n", name, age, height, isStudent)

	// %T - type
	fmt.Printf("%%T: Типи: %T, %T, %T, %T\n", name, age, height, isStudent)

	fmt.Println()
}

func demonstrateStructs() {
	fmt.Println("--- 2. Структури ---\n")

	person := Person{
		FirstName: "Марія",
		LastName:  "Коваленко",
		Age:       28,
		City:      "Київ",
		Salary:    50000.50,
		IsStudent: false,
	}

	// %v - default format
	fmt.Printf("%%v:  %v\n", person)

	// %+v - format з іменами полів
	fmt.Printf("%%+v: %+v\n", person)

	// %#v - Go-синтаксис
	fmt.Printf("%%#v: %#v\n", person)

	// %T - тип
	fmt.Printf("%%T:  %T\n", person)

	fmt.Println()
}

func demonstrateNumbers() {
	fmt.Println("--- 3. Форматування чисел ---\n")

	number := 42
	pi := 3.14159265359
	price := 1234.56

	// Цілі числа
	fmt.Printf("%%d:     %d\n", number)     // decimal
	fmt.Printf("%%b:     %b\n", number)     // binary
	fmt.Printf("%%o:     %o\n", number)     // octal
	fmt.Printf("%%x:     %x\n", number)     // hex lowercase
	fmt.Printf("%%X:     %X\n", number)     // hex uppercase
	fmt.Printf("%%5d:    |%5d|\n", number)  // ширина 5
	fmt.Printf("%%05d:   |%05d|\n", number) // padding нулями

	fmt.Println()

	// Дійсні числа
	fmt.Printf("%%f:       %f\n", pi)     // default precision
	fmt.Printf("%%.2f:     %.2f\n", pi)   // 2 знаки після коми
	fmt.Printf("%%.4f:     %.4f\n", pi)   // 4 знаки після коми
	fmt.Printf("%%8.2f:   |%8.2f|\n", pi) // ширина 8, precision 2
	fmt.Printf("%%e:       %e\n", pi)     // scientific notation
	fmt.Printf("%%E:       %E\n", pi)     // scientific notation uppercase

	fmt.Println()

	// Гроші
	fmt.Printf("Ціна: %.2f грн\n", price)
	fmt.Printf("Ціна: %8.2f грн\n", price)
	fmt.Printf("Ціна: %08.2f грн\n", price)

	fmt.Println()
}

func demonstrateSpecial() {
	fmt.Println("--- 4. Спеціальні формати ---\n")

	text := "Привіт"
	pointer := &text
	var nilPointer *string

	// Рядки
	fmt.Printf("%%s:    %s\n", text)     // string
	fmt.Printf("%%q:    %q\n", text)     // quoted string
	fmt.Printf("%%x:    %x\n", text)     // hex dump
	fmt.Printf("%%X:    %X\n", text)     // hex dump uppercase
	fmt.Printf("%%10s: |%10s|\n", text)  // ширина 10, right-aligned
	fmt.Printf("%%-10s:|%-10s|\n", text) // ширина 10, left-aligned

	fmt.Println()

	// Pointer
	fmt.Printf("%%p:    %p\n", pointer)    // pointer address
	fmt.Printf("%%v:    %v\n", pointer)    // pointer value
	fmt.Printf("%%v:    %v\n", nilPointer) // nil pointer

	fmt.Println()

	// Boolean
	fmt.Printf("%%t:    %t\n", true) // boolean
	fmt.Printf("%%t:    %t\n", false)

	fmt.Println()

	// Спеціальні символи
	fmt.Println("Новий рядок: перший\\nдругий")
	fmt.Printf("Табуляція: |%s\\t%s|\n", "один", "два")
	fmt.Println("Лапки: \"текст у лапках\"")

	fmt.Println()
}

/*
=== Довідка по форматуванню ===

Загальні:
  %v    - default format
  %+v   - з іменами полів (struct)
  %#v   - Go-синтаксис
  %T    - тип значення
  %%    - літерал %

Boolean:
  %t    - true або false

Integer:
  %d    - десяткове
  %b    - двійкове
  %o    - вісімкове
  %x    - hex (lowercase)
  %X    - hex (uppercase)
  %c    - character (rune)

Float:
  %f    - decimal point
  %e    - scientific notation
  %E    - scientific notation uppercase
  %.2f  - 2 знаки після коми

String:
  %s    - string
  %q    - quoted string
  %x    - hex dump
  %X    - hex dump uppercase

Pointer:
  %p    - адреса pointer

Ширина та precision:
  %5d   - ширина 5
  %05d  - padding нулями
  %.2f  - 2 знаки після коми
  %8.2f - ширина 8, precision 2
  %-10s - left-aligned
*/
