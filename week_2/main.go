package main

import (
	"fmt"
	"strings"
)

// ============= Simple Example of Struct, Methods, and Interfaces =============

// Person struct
type Person struct {
	Name string
	Age  int
}

// Method on Person (value receiver)
func (p Person) Greet() string {
	return fmt.Sprintf("Привіт! Мене звати %s, мені %d років", p.Name, p.Age)
}

// Method with pointer receiver
func (p *Person) HaveBirthday() {
	p.Age++
}

// Greeter interface
type Greeter interface {
	Greet() string
}

// Function that accepts interface
func SayHello(g Greeter) {
	fmt.Println(g.Greet())
}

func main() {
	fmt.Println("╔════════════════════════════════════════════╗")
	fmt.Println("║    Week 2: Struct & Interface - Demo      ║")
	fmt.Println("╚════════════════════════════════════════════╝")

	// Create a person
	person := Person{Name: "Іван", Age: 25}

	// Call method
	fmt.Println("\n🔹 Method call:")
	fmt.Println(person.Greet())

	// Pointer receiver method
	fmt.Println("\n🔹 Pointer receiver (змінює дані):")
	fmt.Printf("До: %d років\n", person.Age)
	person.HaveBirthday()
	fmt.Printf("Після HaveBirthday(): %d років\n", person.Age)

	// Interface usage
	fmt.Println("\n🔹 Interface (Person реалізує Greeter):")
	SayHello(person)

	// Show that Person implements Greeter implicitly
	fmt.Println("\n💡 Person автоматично реалізує Greeter!")
	fmt.Println("   (неявна реалізація інтерфейсів)")

	// Instructions
	fmt.Println("\n" + strings.Repeat("═", 44))
	fmt.Println("📚 Навчальні матеріали:")
	fmt.Println(strings.Repeat("═", 44))
	fmt.Println("\n1️⃣  Теорія:")
	fmt.Println("   cd theory")
	fmt.Println("   cat 01_methods_vs_functions.md")
	fmt.Println("   cat 02_pointer_receivers.md")
	fmt.Println("   cat 03_implicit_interfaces.md")

	fmt.Println("\n2️⃣  Практика:")
	fmt.Println("   cd practice/methods_demo")
	fmt.Println("   go run main.go")
	fmt.Println()
	fmt.Println("   cd ../interface_demo")
	fmt.Println("   go run main.go")
	fmt.Println()
	fmt.Println("   cd ../user_service")
	fmt.Println("   go run main.go")

	fmt.Println("\n3️⃣  Вправи:")
	fmt.Println("   cd exercises")
	fmt.Println("   cat exercise_1.md  # Calculator")
	fmt.Println("   cat exercise_2.md  # Shape Interface")
	fmt.Println("   cat exercise_3.md  # Storage Interface")

	fmt.Println("\n4️⃣  Рішення:")
	fmt.Println("   cd solutions")
	fmt.Println("   go run solution_1.go")
	fmt.Println("   go run solution_2.go")
	fmt.Println("   go run solution_3.go")

	fmt.Println("\n" + strings.Repeat("═", 44))
	fmt.Println("🎯 Ключові концепції:")
	fmt.Println(strings.Repeat("═", 44))
	fmt.Println("✅ Methods - функції з receiver")
	fmt.Println("✅ Value receiver vs Pointer receiver")
	fmt.Println("✅ Interfaces - набір методів")
	fmt.Println("✅ Неявна реалізація (implicit)")
	fmt.Println("✅ Polymorphism через інтерфейси")
	fmt.Println("✅ Dependency Injection")

	fmt.Println("\n🚀 Почніть з:")
	fmt.Println("   cat README.md")
	fmt.Println("   cat QUICK_START.md")

	fmt.Println("\n" + strings.Repeat("═", 44))
	fmt.Println("Удачі у навчанні! 🎉")
	fmt.Println(strings.Repeat("═", 44))
}
