package main

import (
	"fmt"
	"math"
	"strings"
)

// ============= Shape Interface =============

type Shape interface {
	Area() float64
	Perimeter() float64
	Name() string
}

// ============= Rectangle =============

type Rectangle struct {
	Width  float64
	Height float64
}

func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

func (r Rectangle) Perimeter() float64 {
	return 2 * (r.Width + r.Height)
}

func (r Rectangle) Name() string {
	return "Прямокутник"
}

// ============= Circle =============

type Circle struct {
	Radius float64
}

func (c Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}

func (c Circle) Perimeter() float64 {
	return 2 * math.Pi * c.Radius
}

func (c Circle) Name() string {
	return "Коло"
}

// ============= Triangle =============

type Triangle struct {
	A, B, C float64 // сторони
}

func (t Triangle) Area() float64 {
	// Формула Герона
	s := (t.A + t.B + t.C) / 2
	return math.Sqrt(s * (s - t.A) * (s - t.B) * (s - t.C))
}

func (t Triangle) Perimeter() float64 {
	return t.A + t.B + t.C
}

func (t Triangle) Name() string {
	return "Трикутник"
}

// ============= Helper Functions =============

func PrintShapeInfo(s Shape) {
	fmt.Printf("📊 %s\n", s.Name())
	fmt.Printf("   Площа: %.2f\n", s.Area())
	fmt.Printf("   Периметр: %.2f\n", s.Perimeter())
}

func TotalArea(shapes []Shape) float64 {
	total := 0.0
	for _, shape := range shapes {
		total += shape.Area()
	}
	return total
}

func LargestShape(shapes []Shape) Shape {
	if len(shapes) == 0 {
		return nil
	}

	largest := shapes[0]
	for _, shape := range shapes[1:] {
		if shape.Area() > largest.Area() {
			largest = shape
		}
	}
	return largest
}

// ============= Writer Interface =============

type Writer interface {
	Write(data string) error
}

// Console Writer
type ConsoleWriter struct{}

func (c ConsoleWriter) Write(data string) error {
	fmt.Println("📺 Console:", data)
	return nil
}

// File Writer (mock)
type FileWriter struct {
	filename string
}

func (f FileWriter) Write(data string) error {
	fmt.Printf("💾 File [%s]: %s\n", f.filename, data)
	return nil
}

// Uppercase Writer
type UppercaseWriter struct {
	writer Writer
}

func (u UppercaseWriter) Write(data string) error {
	return u.writer.Write(strings.ToUpper(data))
}

// ============= Logger =============

type Logger struct {
	writer Writer
}

func NewLogger(w Writer) *Logger {
	return &Logger{writer: w}
}

func (l *Logger) Log(message string) {
	l.writer.Write(message)
}

// ============= Greeter Interface =============

type Greeter interface {
	Greet() string
}

type Person struct {
	Name string
	Age  int
}

func (p Person) Greet() string {
	return fmt.Sprintf("Привіт! Мене звати %s, мені %d років", p.Name, p.Age)
}

type Dog struct {
	Name string
}

func (d Dog) Greet() string {
	return fmt.Sprintf("Гав-гав! Я собака %s!", d.Name)
}

type Robot struct {
	Model string
}

func (r Robot) Greet() string {
	return fmt.Sprintf("Beep-boop. Модель робота: %s", r.Model)
}

func SayHello(g Greeter) {
	fmt.Println("💬", g.Greet())
}

// ============= Main =============

func main() {
	fmt.Println("╔══════════════════════════════════════════╗")
	fmt.Println("║  Interface Demo: Polymorphism            ║")
	fmt.Println("╚══════════════════════════════════════════╝")

	// ===== Shape Interface Demo =====
	fmt.Println("\n🔷 SHAPE INTERFACE")
	fmt.Println("─────────────────────────────────────────")

	rect := Rectangle{Width: 10, Height: 5}
	circle := Circle{Radius: 7}
	triangle := Triangle{A: 3, B: 4, C: 5}

	fmt.Println("\nІнформація про фігури:")
	PrintShapeInfo(rect)
	fmt.Println()
	PrintShapeInfo(circle)
	fmt.Println()
	PrintShapeInfo(triangle)

	// Slice різних фігур через інтерфейс!
	shapes := []Shape{rect, circle, triangle}

	fmt.Println("\n📈 Загальна статистика:")
	fmt.Printf("   Кількість фігур: %d\n", len(shapes))
	fmt.Printf("   Загальна площа: %.2f\n", TotalArea(shapes))

	largest := LargestShape(shapes)
	fmt.Printf("   Найбільша фігура: %s (площа: %.2f)\n", largest.Name(), largest.Area())

	// ===== Writer Interface Demo =====
	fmt.Println("\n\n✍️  WRITER INTERFACE")
	fmt.Println("─────────────────────────────────────────")

	console := ConsoleWriter{}
	file := FileWriter{filename: "app.log"}

	console.Write("Hello from console")
	file.Write("Hello from file")

	// Wrapper
	fmt.Println("\n🔠 Uppercase Wrapper:")
	uppercase := UppercaseWriter{writer: console}
	uppercase.Write("this will be uppercase")

	// Logger with different writers
	fmt.Println("\n📝 Logger:")
	logger1 := NewLogger(console)
	logger2 := NewLogger(file)

	logger1.Log("Лог на консоль")
	logger2.Log("Лог у файл")

	// ===== Greeter Interface Demo =====
	fmt.Println("\n\n👋 GREETER INTERFACE")
	fmt.Println("─────────────────────────────────────────")

	person := Person{Name: "Іван", Age: 25}
	dog := Dog{Name: "Рекс"}
	robot := Robot{Model: "T-800"}

	// Різні типи через один інтерфейс!
	greeters := []Greeter{person, dog, robot}

	for _, greeter := range greeters {
		SayHello(greeter)
	}

	// ===== Type Assertion Demo =====
	fmt.Println("\n\n🔍 TYPE ASSERTION")
	fmt.Println("─────────────────────────────────────────")

	var s Shape = circle

	// Type assertion з перевіркою
	if c, ok := s.(Circle); ok {
		fmt.Printf("✅ Це коло з радіусом %.2f\n", c.Radius)
	}

	if _, ok := s.(Rectangle); ok {
		fmt.Println("Це прямокутник")
	} else {
		fmt.Println("❌ Це НЕ прямокутник")
	}

	// Type switch
	fmt.Println("\n🔀 Type Switch:")
	DescribeShape(rect)
	DescribeShape(circle)
	DescribeShape(triangle)

	// ===== Empty Interface Demo =====
	fmt.Println("\n\n📦 EMPTY INTERFACE (interface{})")
	fmt.Println("─────────────────────────────────────────")

	PrintAnything("Hello, Go!")
	PrintAnything(42)
	PrintAnything(3.14)
	PrintAnything(true)
	PrintAnything([]int{1, 2, 3})
	PrintAnything(Person{Name: "Марія", Age: 30})

	// ===== Summary =====
	fmt.Println("\n\n📝 ВИСНОВКИ")
	fmt.Println("─────────────────────────────────────────")
	fmt.Println("✅ Інтерфейси дозволяють:")
	fmt.Println("   • Polymorphism (різні типи через один інтерфейс)")
	fmt.Println("   • Dependency Injection")
	fmt.Println("   • Тестування через моки")
	fmt.Println("   • Гнучку архітектуру")
	fmt.Println()
	fmt.Println("💡 Неявна реалізація:")
	fmt.Println("   • Не потрібно явно вказувати 'implements'")
	fmt.Println("   • Якщо тип має всі методи → реалізує інтерфейс")
	fmt.Println("   • 'Duck typing': якщо виглядає як качка...")
}

func DescribeShape(s Shape) {
	switch v := s.(type) {
	case Rectangle:
		fmt.Printf("   Прямокутник %v x %v\n", v.Width, v.Height)
	case Circle:
		fmt.Printf("   Коло з радіусом %v\n", v.Radius)
	case Triangle:
		fmt.Printf("   Трикутник зі сторонами %v, %v, %v\n", v.A, v.B, v.C)
	default:
		fmt.Printf("   Невідома фігура: %T\n", v)
	}
}

func PrintAnything(v interface{}) {
	fmt.Printf("   Value: %v, Type: %T\n", v, v)
}
