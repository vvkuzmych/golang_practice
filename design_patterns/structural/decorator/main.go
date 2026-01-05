package main

import (
	"fmt"
	"strings"
	"time"
)

// ============= Example 1: Coffee Decorator =============

// Coffee базовий інтерфейс
type Coffee interface {
	GetDescription() string
	GetCost() float64
}

// Espresso базова кава
type Espresso struct{}

func (e *Espresso) GetDescription() string {
	return "Espresso"
}

func (e *Espresso) GetCost() float64 {
	return 2.00
}

// MilkDecorator додає молоко
type MilkDecorator struct {
	coffee Coffee
}

func (m *MilkDecorator) GetDescription() string {
	return m.coffee.GetDescription() + ", Milk"
}

func (m *MilkDecorator) GetCost() float64 {
	return m.coffee.GetCost() + 0.50
}

// SugarDecorator додає цукор
type SugarDecorator struct {
	coffee Coffee
}

func (s *SugarDecorator) GetDescription() string {
	return s.coffee.GetDescription() + ", Sugar"
}

func (s *SugarDecorator) GetCost() float64 {
	return s.coffee.GetCost() + 0.25
}

// WhippedCreamDecorator додає вершки
type WhippedCreamDecorator struct {
	coffee Coffee
}

func (w *WhippedCreamDecorator) GetDescription() string {
	return w.coffee.GetDescription() + ", Whipped Cream"
}

func (w *WhippedCreamDecorator) GetCost() float64 {
	return w.coffee.GetCost() + 0.75
}

// ============= Example 2: Text Decorator =============

// TextProcessor інтерфейс
type TextProcessor interface {
	Process(text string) string
}

// PlainText базова реалізація
type PlainText struct{}

func (p *PlainText) Process(text string) string {
	return text
}

// UpperCaseDecorator
type UpperCaseDecorator struct {
	processor TextProcessor
}

func (u *UpperCaseDecorator) Process(text string) string {
	return strings.ToUpper(u.processor.Process(text))
}

// TrimDecorator
type TrimDecorator struct {
	processor TextProcessor
}

func (t *TrimDecorator) Process(text string) string {
	return strings.TrimSpace(t.processor.Process(text))
}

// QuoteDecorator
type QuoteDecorator struct {
	processor TextProcessor
}

func (q *QuoteDecorator) Process(text string) string {
	return `"` + q.processor.Process(text) + `"`
}

// ============= Example 3: HTTP Middleware (Decorator) =============

// Handler інтерфейс (як http.Handler)
type Handler interface {
	Handle(request string) string
}

// BasicHandler базовий обробник
type BasicHandler struct{}

func (b *BasicHandler) Handle(request string) string {
	return fmt.Sprintf("Response: %s", request)
}

// LoggingMiddleware логує запити
type LoggingMiddleware struct {
	handler Handler
}

func (l *LoggingMiddleware) Handle(request string) string {
	fmt.Printf("📝 [LOG] Request: %s\n", request)
	response := l.handler.Handle(request)
	fmt.Printf("📝 [LOG] Response sent\n")
	return response
}

// AuthMiddleware перевіряє авторизацію
type AuthMiddleware struct {
	handler Handler
	token   string
}

func (a *AuthMiddleware) Handle(request string) string {
	if a.token == "" {
		fmt.Println("🔒 [AUTH] No token, access denied")
		return "401 Unauthorized"
	}
	fmt.Printf("🔒 [AUTH] Token valid: %s\n", a.token[:8]+"...")
	return a.handler.Handle(request)
}

// TimingMiddleware вимірює час
type TimingMiddleware struct {
	handler Handler
}

func (t *TimingMiddleware) Handle(request string) string {
	start := time.Now()
	response := t.handler.Handle(request)
	duration := time.Since(start)
	fmt.Printf("⏱️  [TIMING] Took %v\n", duration)
	return response
}

// CompressionMiddleware "стискає" відповідь
type CompressionMiddleware struct {
	handler Handler
}

func (c *CompressionMiddleware) Handle(request string) string {
	response := c.handler.Handle(request)
	fmt.Println("🗜️  [COMPRESSION] Response compressed")
	return "[COMPRESSED] " + response
}

// ============= Main =============

func main() {
	fmt.Println("╔════════════════════════════════════════════════╗")
	fmt.Println("║         Decorator Pattern Demo                 ║")
	fmt.Println("╚════════════════════════════════════════════════╝")

	// ===== Example 1: Coffee Shop =====
	fmt.Println("\n🔹 Example 1: Coffee Shop")
	fmt.Println(strings.Repeat("─", 50))

	// Простий еспресо
	var coffee Coffee = &Espresso{}
	fmt.Printf("☕ %s: $%.2f\n", coffee.GetDescription(), coffee.GetCost())

	// Еспресо з молоком
	coffee = &MilkDecorator{coffee: &Espresso{}}
	fmt.Printf("☕ %s: $%.2f\n", coffee.GetDescription(), coffee.GetCost())

	// Еспресо з молоком і цукром
	coffee = &Espresso{}
	coffee = &MilkDecorator{coffee: coffee}
	coffee = &SugarDecorator{coffee: coffee}
	fmt.Printf("☕ %s: $%.2f\n", coffee.GetDescription(), coffee.GetCost())

	// Максимальний набір
	coffee = &Espresso{}
	coffee = &MilkDecorator{coffee: coffee}
	coffee = &SugarDecorator{coffee: coffee}
	coffee = &WhippedCreamDecorator{coffee: coffee}
	fmt.Printf("☕ %s: $%.2f\n", coffee.GetDescription(), coffee.GetCost())

	// ===== Example 2: Text Processing =====
	fmt.Println("\n\n🔹 Example 2: Text Processing Pipeline")
	fmt.Println(strings.Repeat("─", 50))

	originalText := "  hello world  "
	fmt.Printf("\nOriginal: '%s'\n", originalText)

	// Plain
	var processor TextProcessor = &PlainText{}
	fmt.Printf("Plain: '%s'\n", processor.Process(originalText))

	// Trim
	processor = &TrimDecorator{processor: &PlainText{}}
	fmt.Printf("Trim: '%s'\n", processor.Process(originalText))

	// Trim + UpperCase
	processor = &PlainText{}
	processor = &TrimDecorator{processor: processor}
	processor = &UpperCaseDecorator{processor: processor}
	fmt.Printf("Trim + Upper: '%s'\n", processor.Process(originalText))

	// Trim + UpperCase + Quote
	processor = &PlainText{}
	processor = &TrimDecorator{processor: processor}
	processor = &UpperCaseDecorator{processor: processor}
	processor = &QuoteDecorator{processor: processor}
	fmt.Printf("Trim + Upper + Quote: '%s'\n", processor.Process(originalText))

	// ===== Example 3: HTTP Middleware =====
	fmt.Println("\n\n🔹 Example 3: HTTP Middleware Stack")
	fmt.Println(strings.Repeat("─", 50))

	// Базовий handler
	fmt.Println("\n1️⃣ No middleware:")
	var handler Handler = &BasicHandler{}
	response := handler.Handle("GET /api/users")
	fmt.Printf("   Result: %s\n", response)

	// З логуванням
	fmt.Println("\n2️⃣ With Logging:")
	handler = &LoggingMiddleware{handler: &BasicHandler{}}
	response = handler.Handle("GET /api/users")
	fmt.Printf("   Result: %s\n", response)

	// З авторизацією (без токена)
	fmt.Println("\n3️⃣ With Auth (no token):")
	handler = &AuthMiddleware{handler: &BasicHandler{}, token: ""}
	response = handler.Handle("GET /api/users")
	fmt.Printf("   Result: %s\n", response)

	// З авторизацією (з токеном)
	fmt.Println("\n4️⃣ With Auth (with token):")
	handler = &AuthMiddleware{handler: &BasicHandler{}, token: "abc123token456"}
	response = handler.Handle("GET /api/users")
	fmt.Printf("   Result: %s\n", response)

	// Повний стек middleware
	fmt.Println("\n5️⃣ Full Middleware Stack:")
	fmt.Println("   (Timing → Logging → Auth → Compression → Handler)")

	handler = &BasicHandler{}
	handler = &CompressionMiddleware{handler: handler}
	handler = &AuthMiddleware{handler: handler, token: "valid-token-123"}
	handler = &LoggingMiddleware{handler: handler}
	handler = &TimingMiddleware{handler: handler}

	response = handler.Handle("GET /api/users")
	fmt.Printf("   Final Result: %s\n", response)

	// ===== Order Matters =====
	fmt.Println("\n\n🔹 Order of Decorators Matters!")
	fmt.Println(strings.Repeat("─", 50))

	text := "  hello  "

	// Order 1: Trim → Upper → Quote
	fmt.Println("\n📝 Order 1: Trim → Upper → Quote")
	processor = &PlainText{}
	processor = &TrimDecorator{processor: processor}
	processor = &UpperCaseDecorator{processor: processor}
	processor = &QuoteDecorator{processor: processor}
	result1 := processor.Process(text)
	fmt.Printf("   Result: %s\n", result1)

	// Order 2: Quote → Trim → Upper
	fmt.Println("\n📝 Order 2: Quote → Trim → Upper")
	processor = &PlainText{}
	processor = &QuoteDecorator{processor: processor}
	processor = &TrimDecorator{processor: processor}
	processor = &UpperCaseDecorator{processor: processor}
	result2 := processor.Process(text)
	fmt.Printf("   Result: %s\n", result2)

	fmt.Println("\n⚠️  Different results! Order is important!")

	// ===== Summary =====
	fmt.Println("\n\n📝 ВИСНОВКИ")
	fmt.Println(strings.Repeat("─", 50))
	fmt.Println("✅ Динамічне додавання функціональності")
	fmt.Println("✅ Комбінування декораторів")
	fmt.Println("✅ Немає зміни базового класу")
	fmt.Println("✅ Single Responsibility Principle")
	fmt.Println("✅ Open/Closed Principle")

	fmt.Println("\n💡 ВИКОРИСТАННЯ:")
	fmt.Println("   - HTTP middleware (logging, auth, metrics)")
	fmt.Println("   - Stream processing (buffering, encryption)")
	fmt.Println("   - Text processing pipelines")
	fmt.Println("   - UI component enhancement")

	fmt.Println("\n📚 Go stdlib приклади:")
	fmt.Println("   - io.Reader wrappers (bufio, gzip, etc)")
	fmt.Println("   - http.Handler middleware")
	fmt.Println("   - context.Context wrapping")

	fmt.Println("\n⚠️  ВАЖЛИВО:")
	fmt.Println("   - Порядок декораторів має значення!")
	fmt.Println("   - Кожен декоратор - одна відповідальність")
}
