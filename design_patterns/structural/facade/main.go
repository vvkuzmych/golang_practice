package main

import (
	"fmt"
	"strings"
)

// ============= Complex Subsystems =============

// Inventory system
type Inventory struct{}

func (i *Inventory) CheckAvailability(product string) bool {
	fmt.Printf("   📦 Checking inventory for %s... OK\n", product)
	return true
}

// Payment system
type Payment struct{}

func (p *Payment) ProcessPayment(amount float64) bool {
	fmt.Printf("   💳 Processing payment $%.2f... OK\n", amount)
	return true
}

// Shipping system
type Shipping struct{}

func (s *Shipping) CalculateShipping(address string) float64 {
	fmt.Printf("   🚚 Calculating shipping to %s... $10.00\n", address)
	return 10.00
}

func (s *Shipping) ShipProduct(product, address string) bool {
	fmt.Printf("   📮 Shipping %s to %s... OK\n", product, address)
	return true
}

// Notification system
type Notification struct{}

func (n *Notification) SendEmail(email, message string) {
	fmt.Printf("   📧 Email sent to %s: %s\n", email, message)
}

func (n *Notification) SendSMS(phone, message string) {
	fmt.Printf("   📱 SMS sent to %s: %s\n", phone, message)
}

// ============= Facade =============

type OrderFacade struct {
	inventory    *Inventory
	payment      *Payment
	shipping     *Shipping
	notification *Notification
}

func NewOrderFacade() *OrderFacade {
	return &OrderFacade{
		inventory:    &Inventory{},
		payment:      &Payment{},
		shipping:     &Shipping{},
		notification: &Notification{},
	}
}

func (o *OrderFacade) PlaceOrder(product string, amount float64, address, email string) bool {
	fmt.Println("\n🛒 Processing order...")

	// Step 1: Check inventory
	if !o.inventory.CheckAvailability(product) {
		fmt.Println("   ❌ Product not available")
		return false
	}

	// Step 2: Calculate shipping
	shippingCost := o.shipping.CalculateShipping(address)
	total := amount + shippingCost

	// Step 3: Process payment
	if !o.payment.ProcessPayment(total) {
		fmt.Println("   ❌ Payment failed")
		return false
	}

	// Step 4: Ship product
	if !o.shipping.ShipProduct(product, address) {
		fmt.Println("   ❌ Shipping failed")
		return false
	}

	// Step 5: Send notifications
	o.notification.SendEmail(email, "Your order has been shipped!")

	fmt.Printf("\n   ✅ Order placed successfully! Total: $%.2f\n", total)
	return true
}

// ============= Example 2: Computer Builder Facade =============

type CPU struct{}

func (c *CPU) Initialize() {
	fmt.Println("   🔧 CPU initialized")
}

type Memory struct{}

func (m *Memory) Load() {
	fmt.Println("   💾 Memory loaded")
}

type HardDrive struct{}

func (h *HardDrive) Read() {
	fmt.Println("   💿 Hard drive ready")
}

// ComputerFacade
type ComputerFacade struct {
	cpu       *CPU
	memory    *Memory
	hardDrive *HardDrive
}

func NewComputer() *ComputerFacade {
	return &ComputerFacade{
		cpu:       &CPU{},
		memory:    &Memory{},
		hardDrive: &HardDrive{},
	}
}

func (c *ComputerFacade) Start() {
	fmt.Println("\n💻 Starting computer...")
	c.cpu.Initialize()
	c.memory.Load()
	c.hardDrive.Read()
	fmt.Println("   ✅ Computer started!")
}

// ============= Main =============

func main() {
	fmt.Println("╔════════════════════════════════════════════════╗")
	fmt.Println("║           Facade Pattern Demo                  ║")
	fmt.Println("╚════════════════════════════════════════════════╝")

	// ===== Example 1: E-commerce Order =====
	fmt.Println("\n🔹 Example 1: E-commerce Order Processing")
	fmt.Println(strings.Repeat("─", 50))

	fmt.Println("\n❌ Without Facade (complex):")
	fmt.Println("   inventory := &Inventory{}")
	fmt.Println("   if !inventory.CheckAvailability(...) { return }")
	fmt.Println("   payment := &Payment{}")
	fmt.Println("   if !payment.ProcessPayment(...) { return }")
	fmt.Println("   shipping := &Shipping{}")
	fmt.Println("   cost := shipping.CalculateShipping(...)")
	fmt.Println("   shipping.ShipProduct(...)")
	fmt.Println("   notification := &Notification{}")
	fmt.Println("   notification.SendEmail(...)")
	fmt.Println("   → Багато кроків!")

	fmt.Println("\n✅ With Facade (simple):")
	orderFacade := NewOrderFacade()
	orderFacade.PlaceOrder(
		"Laptop",
		1299.99,
		"123 Main St, Kyiv",
		"user@example.com",
	)

	// ===== Example 2: Computer Boot =====
	fmt.Println("\n\n🔹 Example 2: Computer Boot Process")
	fmt.Println(strings.Repeat("─", 50))

	fmt.Println("\n❌ Without Facade:")
	fmt.Println("   cpu := &CPU{}")
	fmt.Println("   cpu.Initialize()")
	fmt.Println("   memory := &Memory{}")
	fmt.Println("   memory.Load()")
	fmt.Println("   hd := &HardDrive{}")
	fmt.Println("   hd.Read()")
	fmt.Println("   → Складно!")

	fmt.Println("\n✅ With Facade:")
	computer := NewComputer()
	computer.Start()

	// ===== Example 3: Multiple Orders =====
	fmt.Println("\n\n🔹 Example 3: Processing Multiple Orders")
	fmt.Println(strings.Repeat("─", 50))

	orders := []struct {
		product string
		price   float64
		address string
		email   string
	}{
		{"Phone", 799.00, "Lviv, Ukraine", "john@example.com"},
		{"Tablet", 499.00, "Odesa, Ukraine", "jane@example.com"},
	}

	facade := NewOrderFacade()
	for i, order := range orders {
		fmt.Printf("\n📦 Order #%d:\n", i+1)
		facade.PlaceOrder(order.product, order.price, order.address, order.email)
	}

	// ===== Summary =====
	fmt.Println("\n\n📝 ВИСНОВКИ")
	fmt.Println(strings.Repeat("─", 50))
	fmt.Println("✅ Спрощує складну систему")
	fmt.Println("✅ Один метод замість багатьох викликів")
	fmt.Println("✅ Приховує деталі реалізації")
	fmt.Println("✅ Зменшує coupling")

	fmt.Println("\n💡 ВИКОРИСТАННЯ:")
	fmt.Println("   - E-commerce checkout process")
	fmt.Println("   - Computer boot sequence")
	fmt.Println("   - Authentication systems")
	fmt.Println("   - Complex API wrappers")

	fmt.Println("\n📚 Go stdlib приклади:")
	fmt.Println("   - http.ListenAndServe()")
	fmt.Println("   - database/sql")
}
