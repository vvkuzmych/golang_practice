package main

import (
	"fmt"
	"strings"
)

// ============= Example 1: Transport Factory =============

// Transport інтерфейс для всіх видів транспорту
type Transport interface {
	Deliver(destination string) string
	GetCost(distance int) float64
}

// Ship - доставка морем
type Ship struct {
	Name string
}

func (s *Ship) Deliver(destination string) string {
	return fmt.Sprintf("🚢 %s: Delivering to %s by sea", s.Name, destination)
}

func (s *Ship) GetCost(distance int) float64 {
	return float64(distance) * 0.5 // $0.5 per km
}

// Truck - доставка дорогою
type Truck struct {
	LicensePlate string
}

func (t *Truck) Deliver(destination string) string {
	return fmt.Sprintf("🚚 %s: Delivering to %s by road", t.LicensePlate, destination)
}

func (t *Truck) GetCost(distance int) float64 {
	return float64(distance) * 1.0 // $1.0 per km
}

// Plane - доставка повітрям
type Plane struct {
	FlightNumber string
}

func (p *Plane) Deliver(destination string) string {
	return fmt.Sprintf("✈️  %s: Delivering to %s by air", p.FlightNumber, destination)
}

func (p *Plane) GetCost(distance int) float64 {
	return float64(distance) * 3.0 // $3.0 per km
}

// NewTransport - factory method для створення транспорту
func NewTransport(transportType string) Transport {
	switch strings.ToLower(transportType) {
	case "sea", "ship":
		return &Ship{Name: "Cargo Ship S-123"}
	case "road", "truck":
		return &Truck{LicensePlate: "AA1234BB"}
	case "air", "plane":
		return &Plane{FlightNumber: "UA555"}
	default:
		// Default: truck
		return &Truck{LicensePlate: "DEFAULT-01"}
	}
}

// ============= Example 2: Payment Factory =============

// PaymentMethod інтерфейс для методів оплати
type PaymentMethod interface {
	ProcessPayment(amount float64) string
	GetFee(amount float64) float64
}

// CreditCard - оплата карткою
type CreditCard struct {
	CardNumber string
}

func (c *CreditCard) ProcessPayment(amount float64) string {
	return fmt.Sprintf("💳 Processing $%.2f via Credit Card ****%s",
		amount, c.CardNumber[len(c.CardNumber)-4:])
}

func (c *CreditCard) GetFee(amount float64) float64 {
	return amount * 0.03 // 3% fee
}

// PayPal - оплата через PayPal
type PayPal struct {
	Email string
}

func (p *PayPal) ProcessPayment(amount float64) string {
	return fmt.Sprintf("💰 Processing $%.2f via PayPal (%s)", amount, p.Email)
}

func (p *PayPal) GetFee(amount float64) float64 {
	return amount * 0.05 // 5% fee
}

// Crypto - оплата криптовалютою
type Crypto struct {
	Wallet string
}

func (c *Crypto) ProcessPayment(amount float64) string {
	return fmt.Sprintf("₿  Processing $%.2f via Crypto (wallet: %s)", amount, c.Wallet[:8]+"...")
}

func (c *Crypto) GetFee(amount float64) float64 {
	return 2.0 // fixed $2 fee
}

// NewPaymentMethod - factory method для створення методу оплати
func NewPaymentMethod(method string) PaymentMethod {
	switch strings.ToLower(method) {
	case "card", "creditcard":
		return &CreditCard{CardNumber: "1234567890123456"}
	case "paypal":
		return &PayPal{Email: "user@example.com"}
	case "crypto", "bitcoin":
		return &Crypto{Wallet: "1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa"}
	default:
		return &CreditCard{CardNumber: "0000000000000000"}
	}
}

// ============= Example 3: Logger Factory =============

// Logger інтерфейс для логерів
type Logger interface {
	Log(message string)
	LogLevel() string
}

// ConsoleLogger - лог в консоль
type ConsoleLogger struct {
	prefix string
}

func (c *ConsoleLogger) Log(message string) {
	fmt.Printf("🖥️  [CONSOLE] %s %s\n", c.prefix, message)
}

func (c *ConsoleLogger) LogLevel() string {
	return "INFO"
}

// FileLogger - лог у файл
type FileLogger struct {
	filename string
}

func (f *FileLogger) Log(message string) {
	fmt.Printf("📁 [FILE: %s] %s\n", f.filename, message)
}

func (f *FileLogger) LogLevel() string {
	return "DEBUG"
}

// RemoteLogger - лог на віддалений сервер
type RemoteLogger struct {
	endpoint string
}

func (r *RemoteLogger) Log(message string) {
	fmt.Printf("🌐 [REMOTE: %s] %s\n", r.endpoint, message)
}

func (r *RemoteLogger) LogLevel() string {
	return "ERROR"
}

// NewLogger - factory method для створення логера
func NewLogger(loggerType string) Logger {
	switch strings.ToLower(loggerType) {
	case "console":
		return &ConsoleLogger{prefix: "[APP]"}
	case "file":
		return &FileLogger{filename: "app.log"}
	case "remote":
		return &RemoteLogger{endpoint: "https://logs.example.com"}
	default:
		return &ConsoleLogger{prefix: "[DEFAULT]"}
	}
}

// ============= Helper Functions =============

func printSeparator(title string) {
	fmt.Printf("\n🔹 %s\n", title)
	fmt.Println(strings.Repeat("─", 50))
}

// ============= Main =============

func main() {
	fmt.Println("╔════════════════════════════════════════════════╗")
	fmt.Println("║       Factory Method Pattern Demo             ║")
	fmt.Println("╚════════════════════════════════════════════════╝")

	// ===== Example 1: Transport Factory =====
	printSeparator("Example 1: Transport Factory")

	destinations := []string{"Odesa", "Kyiv", "Lviv"}
	transportTypes := []string{"ship", "truck", "plane"}
	distance := 500

	for i, tType := range transportTypes {
		transport := NewTransport(tType)
		fmt.Println(transport.Deliver(destinations[i]))
		cost := transport.GetCost(distance)
		fmt.Printf("   Cost for %d km: $%.2f\n", distance, cost)
	}

	// ===== Example 2: Payment Factory =====
	printSeparator("Example 2: Payment Methods Factory")

	amount := 100.0
	paymentMethods := []string{"card", "paypal", "crypto"}

	for _, method := range paymentMethods {
		payment := NewPaymentMethod(method)
		fmt.Println(payment.ProcessPayment(amount))
		fee := payment.GetFee(amount)
		fmt.Printf("   Fee: $%.2f, Total: $%.2f\n", fee, amount+fee)
	}

	// ===== Example 3: Logger Factory =====
	printSeparator("Example 3: Logger Factory")

	loggers := []string{"console", "file", "remote"}

	for _, logType := range loggers {
		logger := NewLogger(logType)
		logger.Log(fmt.Sprintf("Application started (Level: %s)", logger.LogLevel()))
	}

	// ===== Real-World Scenario =====
	printSeparator("Real-World Scenario: E-commerce Order")

	// Замовлення
	order := struct {
		id          string
		destination string
		distance    int
		amount      float64
	}{
		id:          "ORD-12345",
		destination: "Kyiv",
		distance:    300,
		amount:      250.0,
	}

	fmt.Printf("\n📦 Processing Order: %s\n", order.id)
	fmt.Printf("   Destination: %s (%d km)\n", order.destination, order.distance)
	fmt.Printf("   Amount: $%.2f\n", order.amount)

	// Вибір транспорту (залежить від відстані)
	var transportType string
	if order.distance < 100 {
		transportType = "truck"
	} else if order.distance < 500 {
		transportType = "plane"
	} else {
		transportType = "ship"
	}

	transport := NewTransport(transportType)
	fmt.Printf("\n%s\n", transport.Deliver(order.destination))
	deliveryCost := transport.GetCost(order.distance)
	fmt.Printf("   Delivery cost: $%.2f\n", deliveryCost)

	// Вибір методу оплати (залежить від суми)
	var paymentType string
	if order.amount < 50 {
		paymentType = "card"
	} else if order.amount < 200 {
		paymentType = "paypal"
	} else {
		paymentType = "card" // для великих сум - картка
	}

	payment := NewPaymentMethod(paymentType)
	totalAmount := order.amount + deliveryCost
	fmt.Printf("\n%s\n", payment.ProcessPayment(totalAmount))
	paymentFee := payment.GetFee(totalAmount)
	fmt.Printf("   Payment fee: $%.2f\n", paymentFee)
	fmt.Printf("   Grand Total: $%.2f\n", totalAmount+paymentFee)

	// Логування
	logger := NewLogger("console")
	logger.Log(fmt.Sprintf("Order %s completed successfully", order.id))

	// ===== Flexibility Demo =====
	printSeparator("Flexibility: Easy to Add New Types")

	fmt.Println("\n💡 Adding new transport type is easy:")
	fmt.Println("   1. Create new struct implementing Transport interface")
	fmt.Println("   2. Add case to factory method")
	fmt.Println("   3. Done! No changes to existing code")

	fmt.Println("\n💡 Example: Adding Drone delivery:")
	fmt.Println("   type Drone struct { ID string }")
	fmt.Println("   case \"drone\": return &Drone{}")

	// ===== Summary =====
	fmt.Println("\n\n📝 ВИСНОВКИ")
	fmt.Println(strings.Repeat("─", 50))
	fmt.Println("✅ Централізована логіка створення об'єктів")
	fmt.Println("✅ Легко додавати нові типи")
	fmt.Println("✅ Слабка зв'язаність (loose coupling)")
	fmt.Println("✅ Клієнт працює через інтерфейс")
	fmt.Println("✅ Ідеально для: Transport, Payment, Logger, Database")

	fmt.Println("\n💡 ВИКОРИСТАННЯ:")
	fmt.Println("   - Тип залежить від умов runtime")
	fmt.Println("   - Потрібна гнучкість при розширенні")
	fmt.Println("   - Централізація створення об'єктів")
}
