package main

import "fmt"

// Student представляє студента
type Student struct {
	FirstName  string
	LastName   string
	Age        int
	GPA        float64
	IsActive   bool
	University *string // pointer для демонстрації nil
}

// FullName повертає повне ім'я студента
func (s Student) FullName() string {
	return s.FirstName + " " + s.LastName
}

// IsExcellent перевіряє чи студент відмінник
func (s Student) IsExcellent() bool {
	return s.GPA >= 4.0
}

// Print виводить інформацію про студента
func (s Student) Print() {
	fmt.Printf("  Ім'я: %s %s\n", s.FirstName, s.LastName)
	fmt.Printf("  Вік: %d років\n", s.Age)
	fmt.Printf("  Середній бал: %.2f\n", s.GPA)
	fmt.Printf("  Активний: %t\n", s.IsActive)

	if s.University != nil {
		fmt.Printf("  Університет: %s\n", *s.University)
	} else {
		fmt.Println("  Університет: не вказано")
	}

	if s.IsExcellent() {
		fmt.Println("  🌟 Відмінник!")
	}
}

// NewStudent створює нового студента з дефолтними значеннями
func NewStudent(firstName, lastName string, age int) Student {
	return Student{
		FirstName: firstName,
		LastName:  lastName,
		Age:       age,
		IsActive:  true,
		GPA:       0.0,
	}
}

func main() {
	fmt.Println("=== Інформація про студентів ===\n")
	var university string = "DPI"
	// 1. Повна ініціалізація
	student1 := Student{
		FirstName:  "Іван",
		LastName:   "Петренко",
		Age:        20,
		GPA:        3.8,
		IsActive:   true,
		University: &university,
	}

	// 2. Часткова ініціалізація (інші поля отримають zero values)
	student2 := Student{
		FirstName: "Марія",
		LastName:  "Коваленко",
	}

	// 3. Zero value (всі поля мають zero values)
	var student3 Student

	// 4. З університетом (pointer)
	uni := "КНУ імені Тараса Шевченка"
	student4 := Student{
		FirstName:  "Петро",
		LastName:   "Сидоренко",
		Age:        22,
		GPA:        4.2,
		IsActive:   true,
		University: &uni,
	}

	// 5. Використання конструктора
	student5 := NewStudent("Олена", "Мельник", 19)
	student5.GPA = 3.5

	// Вивід інформації про кожного студента
	printStudentInfo("Студент 1 (повна ініціалізація)", student1)
	printStudentInfo("Студент 2 (часткова ініціалізація)", student2)
	printStudentInfo("Студент 3 (zero value)", student3)
	printStudentInfo("Студент 4 (з університетом)", student4)
	printStudentInfo("Студент 5 (через конструктор)", student5)

	// Демонстрація методів
	fmt.Println("\n=== Методи структури ===\n")
	students := []Student{student1, student2, student3, student4, student5}

	for i, s := range students {
		fmt.Printf("Студент %d:\n", i+1)
		fmt.Printf("  Повне ім'я: %s\n", s.FullName())
		fmt.Printf("  Відмінник: %t\n", s.IsExcellent())
		fmt.Println()
	}

	// Таблиця Zero Values
	printZeroValuesTable()

	// Порівняння студентів
	fmt.Println("\n=== Порівняння студентів ===\n")
	compareStudents(student1, student4)
}

func printStudentInfo(title string, s Student) {
	fmt.Printf("--- %s ---\n\n", title)

	// %v - default format
	fmt.Printf("%%v:  %v\n", s)

	// %+v - з іменами полів
	fmt.Printf("%%+v: %+v\n", s)

	// %#v - Go syntax
	fmt.Printf("%%#v: %#v\n", s)

	// %T - тип
	fmt.Printf("%%T:  %T\n", s)

	fmt.Println("\nДетальна інформація:")
	s.Print()

	fmt.Println()
}

func printZeroValuesTable() {
	fmt.Println("\n=== Таблиця Zero Values ===\n")

	var zeroStudent Student

	fmt.Println("Поле         | Тип       | Zero Value | Значення")
	fmt.Println("-------------|-----------|------------|------------------")
	fmt.Printf("FirstName    | %-9s | %-10s | %q\n", "string", `""`, zeroStudent.FirstName)
	fmt.Printf("LastName     | %-9s | %-10s | %q\n", "string", `""`, zeroStudent.LastName)
	fmt.Printf("Age          | %-9s | %-10d | %d\n", "int", 0, zeroStudent.Age)
	fmt.Printf("GPA          | %-9s | %-10.1f | %.1f\n", "float64", 0.0, zeroStudent.GPA)
	fmt.Printf("IsActive     | %-9s | %-10t | %t\n", "bool", false, zeroStudent.IsActive)
	fmt.Printf("University   | %-9s | %-10s | %v\n", "*string", "nil", zeroStudent.University)
}

func compareStudents(s1, s2 Student) {
	fmt.Printf("Порівняння: %s vs %s\n\n", s1.FullName(), s2.FullName())

	// Порівняння GPA
	if s1.GPA > s2.GPA {
		fmt.Printf("🏆 %s має вищий середній бал (%.2f vs %.2f)\n", s1.FullName(), s1.GPA, s2.GPA)
	} else if s1.GPA < s2.GPA {
		fmt.Printf("🏆 %s має вищий середній бал (%.2f vs %.2f)\n", s2.FullName(), s2.GPA, s1.GPA)
	} else {
		fmt.Printf("Однаковий середній бал: %.2f\n", s1.GPA)
	}

	// Порівняння віку
	if s1.Age > s2.Age {
		fmt.Printf("📅 %s старший на %d років\n", s1.FullName(), s1.Age-s2.Age)
	} else if s1.Age < s2.Age {
		fmt.Printf("📅 %s старший на %d років\n", s2.FullName(), s2.Age-s1.Age)
	} else {
		fmt.Println("📅 Однаковий вік")
	}

	// Статус відмінників
	if s1.IsExcellent() && s2.IsExcellent() {
		fmt.Println("⭐ Обидва відмінники!")
	} else if s1.IsExcellent() {
		fmt.Printf("⭐ %s є відмінником\n", s1.FullName())
	} else if s2.IsExcellent() {
		fmt.Printf("⭐ %s є відмінником\n", s2.FullName())
	}
}
