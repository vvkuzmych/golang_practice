package main

import (
	"fmt"
	"strings"
)

// ============= Prototype Interface =============

type Cloneable interface {
	Clone() Cloneable
	GetInfo() string
}

// ============= Document Prototype =============

type Document struct {
	Title    string
	Content  string
	Author   string
	Pages    int
	Metadata map[string]string
}

func (d *Document) Clone() Cloneable {
	// Deep copy of metadata
	metadata := make(map[string]string)
	for k, v := range d.Metadata {
		metadata[k] = v
	}

	return &Document{
		Title:    d.Title,
		Content:  d.Content,
		Author:   d.Author,
		Pages:    d.Pages,
		Metadata: metadata,
	}
}

func (d *Document) GetInfo() string {
	return fmt.Sprintf("📄 '%s' by %s (%d pages)", d.Title, d.Author, d.Pages)
}

// ============= Shape Prototype =============

type Shape interface {
	Cloneable
	Draw() string
}

type Circle struct {
	X      int
	Y      int
	Radius int
	Color  string
}

func (c *Circle) Clone() Cloneable {
	return &Circle{
		X:      c.X,
		Y:      c.Y,
		Radius: c.Radius,
		Color:  c.Color,
	}
}

func (c *Circle) GetInfo() string {
	return fmt.Sprintf("⭕ Circle at (%d,%d), radius=%d, color=%s",
		c.X, c.Y, c.Radius, c.Color)
}

func (c *Circle) Draw() string {
	return fmt.Sprintf("Drawing circle: %s", c.GetInfo())
}

type Rectangle struct {
	X      int
	Y      int
	Width  int
	Height int
	Color  string
}

func (r *Rectangle) Clone() Cloneable {
	return &Rectangle{
		X:      r.X,
		Y:      r.Y,
		Width:  r.Width,
		Height: r.Height,
		Color:  r.Color,
	}
}

func (r *Rectangle) GetInfo() string {
	return fmt.Sprintf("▭  Rectangle at (%d,%d), %dx%d, color=%s",
		r.X, r.Y, r.Width, r.Height, r.Color)
}

func (r *Rectangle) Draw() string {
	return fmt.Sprintf("Drawing rectangle: %s", r.GetInfo())
}

// ============= Registry (Prototype Manager) =============

type PrototypeRegistry struct {
	prototypes map[string]Cloneable
}

func NewRegistry() *PrototypeRegistry {
	return &PrototypeRegistry{
		prototypes: make(map[string]Cloneable),
	}
}

func (r *PrototypeRegistry) Register(name string, prototype Cloneable) {
	r.prototypes[name] = prototype
}

func (r *PrototypeRegistry) Create(name string) Cloneable {
	if prototype, exists := r.prototypes[name]; exists {
		return prototype.Clone()
	}
	return nil
}

// ============= Main =============

func main() {
	fmt.Println("╔════════════════════════════════════════════════╗")
	fmt.Println("║         Prototype Pattern Demo                 ║")
	fmt.Println("╚════════════════════════════════════════════════╝")

	// ===== Demo 1: Document Cloning =====
	fmt.Println("\n🔹 Demo 1: Document Cloning")
	fmt.Println(strings.Repeat("─", 50))

	original := &Document{
		Title:   "Design Patterns",
		Content: "This is a book about design patterns...",
		Author:  "Gang of Four",
		Pages:   395,
		Metadata: map[string]string{
			"ISBN":     "978-0201633610",
			"Language": "English",
		},
	}

	fmt.Println("📄 Original document:")
	fmt.Println("  ", original.GetInfo())
	fmt.Printf("   Metadata: %v\n", original.Metadata)

	// Clone
	copy := original.Clone().(*Document)
	copy.Title = "Design Patterns - Copy"
	copy.Metadata["ISBN"] = "978-XXXXXXXXXX"

	fmt.Println("\n📄 Cloned document (modified):")
	fmt.Println("  ", copy.GetInfo())
	fmt.Printf("   Metadata: %v\n", copy.Metadata)

	fmt.Println("\n📄 Original unchanged:")
	fmt.Println("  ", original.GetInfo())
	fmt.Printf("   Metadata: %v\n", original.Metadata)

	// ===== Demo 2: Shape Cloning =====
	fmt.Println("\n\n🔹 Demo 2: Shape Cloning")
	fmt.Println(strings.Repeat("─", 50))

	redCircle := &Circle{
		X:      100,
		Y:      100,
		Radius: 50,
		Color:  "red",
	}

	fmt.Println("\n⭕ Original red circle:")
	fmt.Println("  ", redCircle.GetInfo())

	// Clone and modify
	blueCircle := redCircle.Clone().(*Circle)
	blueCircle.X = 200
	blueCircle.Color = "blue"

	greenCircle := redCircle.Clone().(*Circle)
	greenCircle.X = 300
	greenCircle.Color = "green"

	fmt.Println("\n⭕ Cloned circles:")
	fmt.Println("  ", blueCircle.GetInfo())
	fmt.Println("  ", greenCircle.GetInfo())

	fmt.Println("\n⭕ Original still red:")
	fmt.Println("  ", redCircle.GetInfo())

	// ===== Demo 3: Prototype Registry =====
	fmt.Println("\n\n🔹 Demo 3: Prototype Registry")
	fmt.Println(strings.Repeat("─", 50))

	registry := NewRegistry()

	// Register prototypes
	registry.Register("default-circle", &Circle{
		X:      0,
		Y:      0,
		Radius: 25,
		Color:  "black",
	})

	registry.Register("default-rectangle", &Rectangle{
		X:      0,
		Y:      0,
		Width:  100,
		Height: 50,
		Color:  "gray",
	})

	registry.Register("report-template", &Document{
		Title:   "Monthly Report",
		Content: "Template content...",
		Author:  "System",
		Pages:   1,
		Metadata: map[string]string{
			"Type":     "Report",
			"Template": "Monthly",
		},
	})

	fmt.Println("\n📦 Creating objects from registry:")

	// Create from prototypes
	circle1 := registry.Create("default-circle").(*Circle)
	circle1.X = 50
	circle1.Y = 50
	circle1.Color = "yellow"
	fmt.Println("  ", circle1.GetInfo())

	circle2 := registry.Create("default-circle").(*Circle)
	circle2.X = 150
	circle2.Y = 150
	circle2.Color = "purple"
	fmt.Println("  ", circle2.GetInfo())

	rect := registry.Create("default-rectangle").(*Rectangle)
	rect.X = 200
	rect.Y = 200
	rect.Color = "orange"
	fmt.Println("  ", rect.GetInfo())

	report := registry.Create("report-template").(*Document)
	report.Title = "January 2024 Report"
	report.Author = "John Doe"
	fmt.Println("  ", report.GetInfo())

	// ===== Demo 4: Performance Comparison =====
	fmt.Println("\n\n🔹 Demo 4: Why Use Prototype?")
	fmt.Println(strings.Repeat("─", 50))

	fmt.Println("\n❌ Without Prototype:")
	fmt.Println("   1. Create new object")
	fmt.Println("   2. Set all fields manually")
	fmt.Println("   3. Load data from DB/file")
	fmt.Println("   4. Initialize complex structures")
	fmt.Println("   → Slow and repetitive!")

	fmt.Println("\n✅ With Prototype:")
	fmt.Println("   1. Clone existing object")
	fmt.Println("   2. Modify only what's needed")
	fmt.Println("   → Fast and simple!")

	// ===== Demo 5: Game Example =====
	fmt.Println("\n\n🔹 Demo 5: Game Enemy Spawning")
	fmt.Println(strings.Repeat("─", 50))

	enemyPrototype := &Circle{
		X:      0,
		Y:      0,
		Radius: 10,
		Color:  "red",
	}

	fmt.Println("\n🎮 Spawning enemies from prototype:")

	positions := []struct{ x, y int }{
		{100, 100},
		{200, 150},
		{300, 120},
		{400, 180},
	}

	for i, pos := range positions {
		enemy := enemyPrototype.Clone().(*Circle)
		enemy.X = pos.x
		enemy.Y = pos.y
		fmt.Printf("   Enemy %d: %s\n", i+1, enemy.GetInfo())
	}

	// ===== Summary =====
	fmt.Println("\n\n📝 ВИСНОВКИ")
	fmt.Println(strings.Repeat("─", 50))
	fmt.Println("✅ Швидке копіювання складних об'єктів")
	fmt.Println("✅ Незалежність від конкретних класів")
	fmt.Println("✅ Альтернатива підкласам для конфігурацій")
	fmt.Println("✅ Prototype Registry для управління шаблонами")

	fmt.Println("\n💡 ВИКОРИСТАННЯ:")
	fmt.Println("   - Copy-paste документів")
	fmt.Println("   - Клонування game objects")
	fmt.Println("   - Database record templates")
	fmt.Println("   - Configuration cloning")

	fmt.Println("\n⚠️  ВАЖЛИВО:")
	fmt.Println("   - Deep copy vs Shallow copy")
	fmt.Println("   - Обережно з circular references")
}
