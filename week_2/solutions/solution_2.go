package main

import (
	"fmt"
	"math"
	"sort"
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
	return "Rectangle"
}

func (r Rectangle) String() string {
	return fmt.Sprintf("Rectangle{Width: %.2f, Height: %.2f}", r.Width, r.Height)
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
	return "Circle"
}

func (c Circle) String() string {
	return fmt.Sprintf("Circle{Radius: %.2f}", c.Radius)
}

// ============= Square =============

type Square struct {
	Side float64
}

func (s Square) Area() float64 {
	return s.Side * s.Side
}

func (s Square) Perimeter() float64 {
	return 4 * s.Side
}

func (s Square) Name() string {
	return "Square"
}

func (s Square) String() string {
	return fmt.Sprintf("Square{Side: %.2f}", s.Side)
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
	return "Triangle"
}

func (t Triangle) String() string {
	return fmt.Sprintf("Triangle{A: %.2f, B: %.2f, C: %.2f}", t.A, t.B, t.C)
}

// ============= Helper Functions =============

// PrintShapeInfo виводить інформацію про фігуру
func PrintShapeInfo(s Shape) {
	fmt.Printf("📐 %s\n", s.Name())
	fmt.Printf("   Area: %.2f\n", s.Area())
	fmt.Printf("   Perimeter: %.2f\n", s.Perimeter())
}

// TotalArea обчислює загальну площу всіх фігур
func TotalArea(shapes []Shape) float64 {
	total := 0.0
	for _, shape := range shapes {
		total += shape.Area()
	}
	return total
}

// AverageArea обчислює середню площу
func AverageArea(shapes []Shape) float64 {
	if len(shapes) == 0 {
		return 0
	}
	return TotalArea(shapes) / float64(len(shapes))
}

// LargestShape знаходить найбільшу фігуру за площею
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

// SmallestShape знаходить найменшу фігуру за площею
func SmallestShape(shapes []Shape) Shape {
	if len(shapes) == 0 {
		return nil
	}

	smallest := shapes[0]
	for _, shape := range shapes[1:] {
		if shape.Area() < smallest.Area() {
			smallest = shape
		}
	}
	return smallest
}

// FilterByMinArea фільтрує фігури за мінімальною площею
func FilterByMinArea(shapes []Shape, minArea float64) []Shape {
	result := []Shape{}
	for _, shape := range shapes {
		if shape.Area() >= minArea {
			result = append(result, shape)
		}
	}
	return result
}

// SortByArea сортує фігури за площею
func SortByArea(shapes []Shape) {
	sort.Slice(shapes, func(i, j int) bool {
		return shapes[i].Area() < shapes[j].Area()
	})
}

// CountByType підраховує кількість фігур кожного типу
func CountByType(shapes []Shape) map[string]int {
	counts := make(map[string]int)
	for _, shape := range shapes {
		counts[shape.Name()]++
	}
	return counts
}

// ============= Main =============

func main() {
	fmt.Println("╔══════════════════════════════════════════╗")
	fmt.Println("║        Shape Interface Solution          ║")
	fmt.Println("╚══════════════════════════════════════════╝")

	// ===== Create Shapes =====
	fmt.Println("\n🔹 Створення фігур")
	fmt.Println("─────────────────────────────────────────")

	rect := Rectangle{Width: 10, Height: 5}
	circle := Circle{Radius: 7}
	square := Square{Side: 6}
	triangle := Triangle{A: 3, B: 4, C: 5}

	shapes := []Shape{rect, circle, square, triangle}

	// ===== Individual Info =====
	fmt.Println("\n🔹 Інформація про кожну фігуру")
	fmt.Println("─────────────────────────────────────────")

	for _, shape := range shapes {
		PrintShapeInfo(shape)
		fmt.Println()
	}

	// ===== Statistics =====
	fmt.Println("🔹 Статистика")
	fmt.Println("─────────────────────────────────────────")

	fmt.Printf("Кількість фігур: %d\n", len(shapes))
	fmt.Printf("Загальна площа: %.2f\n", TotalArea(shapes))
	fmt.Printf("Середня площа: %.2f\n", AverageArea(shapes))

	largest := LargestShape(shapes)
	fmt.Printf("Найбільша: %s (%.2f)\n", largest.Name(), largest.Area())

	smallest := SmallestShape(shapes)
	fmt.Printf("Найменша: %s (%.2f)\n", smallest.Name(), smallest.Area())

	// ===== Filter =====
	fmt.Println("\n\n🔹 Фільтрація: площа > 40")
	fmt.Println("─────────────────────────────────────────")

	filtered := FilterByMinArea(shapes, 40)
	for _, shape := range filtered {
		fmt.Printf("  %s: %.2f\n", shape.Name(), shape.Area())
	}

	// ===== Sorting =====
	fmt.Println("\n\n🔹 Сортування за площею")
	fmt.Println("─────────────────────────────────────────")

	sorted := make([]Shape, len(shapes))
	copy(sorted, shapes)
	SortByArea(sorted)

	for i, shape := range sorted {
		fmt.Printf("%d. %s: %.2f\n", i+1, shape.Name(), shape.Area())
	}

	// ===== Count by Type =====
	fmt.Println("\n\n🔹 Підрахунок за типами")
	fmt.Println("─────────────────────────────────────────")

	// Додамо ще фігури
	allShapes := []Shape{
		rect, circle, square, triangle,
		Rectangle{Width: 8, Height: 4},
		Circle{Radius: 5},
		Square{Side: 10},
	}

	counts := CountByType(allShapes)
	for shapeType, count := range counts {
		fmt.Printf("  %s: %d\n", shapeType, count)
	}

	// ===== Comparison =====
	fmt.Println("\n\n🔹 Порівняння фігур")
	fmt.Println("─────────────────────────────────────────")

	rect1 := Rectangle{Width: 10, Height: 10}
	circle1 := Circle{Radius: math.Sqrt(100 / math.Pi)} // площа = 100

	fmt.Printf("Rectangle 10×10: %.2f\n", rect1.Area())
	fmt.Printf("Circle r=%.2f: %.2f\n", circle1.Radius, circle1.Area())
	fmt.Printf("Приблизно однакова площа? %v\n",
		math.Abs(rect1.Area()-circle1.Area()) < 0.01)

	// ===== More Shapes =====
	fmt.Println("\n\n🔹 Колекція різних фігур")
	fmt.Println("─────────────────────────────────────────")

	collection := []Shape{
		Rectangle{Width: 5, Height: 3},
		Circle{Radius: 4},
		Square{Side: 4},
		Triangle{A: 5, B: 5, C: 5}, // рівносторонній
		Rectangle{Width: 20, Height: 2},
		Circle{Radius: 10},
	}

	fmt.Printf("Всього фігур: %d\n", len(collection))
	fmt.Printf("Загальна площа: %.2f\n", TotalArea(collection))
	fmt.Printf("Середня площа: %.2f\n", AverageArea(collection))

	largest2 := LargestShape(collection)
	fmt.Printf("\nНайбільша фігура:\n")
	PrintShapeInfo(largest2)

	// ===== Type Assertion =====
	fmt.Println("\n\n🔹 Type Assertion")
	fmt.Println("─────────────────────────────────────────")

	var s Shape = Circle{Radius: 5}

	if c, ok := s.(Circle); ok {
		fmt.Printf("✅ Це Circle з радіусом %.2f\n", c.Radius)
	}

	if _, ok := s.(Rectangle); ok {
		fmt.Println("Це Rectangle")
	} else {
		fmt.Println("❌ Це НЕ Rectangle")
	}

	// Type Switch
	fmt.Println("\n🔀 Type Switch:")
	for _, shape := range shapes {
		DescribeShape(shape)
	}

	// ===== Summary =====
	fmt.Println("\n\n📝 Висновки")
	fmt.Println("─────────────────────────────────────────")
	fmt.Println("✅ Реалізовано:")
	fmt.Println("   • Interface Shape")
	fmt.Println("   • 4 типи фігур (Rectangle, Circle, Square, Triangle)")
	fmt.Println("   • Polymorphism через інтерфейс")
	fmt.Println("   • Фільтрація, сортування, статистика")
	fmt.Println("   • Type assertions і type switch")
	fmt.Println()
	fmt.Println("💡 Всі фігури реалізують Shape неявно!")
	fmt.Println("   Різні типи працюють через єдиний інтерфейс")
}

func DescribeShape(s Shape) {
	switch v := s.(type) {
	case Rectangle:
		fmt.Printf("  Rectangle: %.2f × %.2f\n", v.Width, v.Height)
	case Circle:
		fmt.Printf("  Circle: r = %.2f\n", v.Radius)
	case Square:
		fmt.Printf("  Square: side = %.2f\n", v.Side)
	case Triangle:
		fmt.Printf("  Triangle: %.2f, %.2f, %.2f\n", v.A, v.B, v.C)
	default:
		fmt.Printf("  Unknown shape: %T\n", v)
	}
}
