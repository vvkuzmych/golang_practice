package main

import (
	"fmt"
	"strings"
)

// ============= Example 1: Payment Strategy =============

// PaymentStrategy інтерфейс для методів оплати
type PaymentStrategy interface {
	Pay(amount float64) string
	GetFee(amount float64) float64
}

// CreditCardPayment стратегія
type CreditCardPayment struct {
	CardNumber string
}

func (c *CreditCardPayment) Pay(amount float64) string {
	return fmt.Sprintf("💳 Paid $%.2f with Credit Card (****%s)",
		amount, c.CardNumber[len(c.CardNumber)-4:])
}

func (c *CreditCardPayment) GetFee(amount float64) float64 {
	return amount * 0.02 // 2%
}

// PayPalPayment стратегія
type PayPalPayment struct {
	Email string
}

func (p *PayPalPayment) Pay(amount float64) string {
	return fmt.Sprintf("💰 Paid $%.2f via PayPal (%s)", amount, p.Email)
}

func (p *PayPalPayment) GetFee(amount float64) float64 {
	return amount * 0.04 // 4%
}

// CryptoPayment стратегія
type CryptoPayment struct {
	Currency string
}

func (c *CryptoPayment) Pay(amount float64) string {
	return fmt.Sprintf("₿  Paid $%.2f via %s", amount, c.Currency)
}

func (c *CryptoPayment) GetFee(amount float64) float64 {
	return 1.0 // flat fee
}

// ShoppingCart контекст
type ShoppingCart struct {
	paymentStrategy PaymentStrategy
	items           []string
	total           float64
}

func NewShoppingCart() *ShoppingCart {
	return &ShoppingCart{items: []string{}}
}

func (s *ShoppingCart) AddItem(item string, price float64) {
	s.items = append(s.items, item)
	s.total += price
}

func (s *ShoppingCart) SetPaymentStrategy(strategy PaymentStrategy) {
	s.paymentStrategy = strategy
}

func (s *ShoppingCart) Checkout() string {
	if s.paymentStrategy == nil {
		return "❌ No payment method selected"
	}

	fee := s.paymentStrategy.GetFee(s.total)
	total := s.total + fee

	result := fmt.Sprintf("🛒 Cart: %d items, Subtotal: $%.2f\n",
		len(s.items), s.total)
	result += fmt.Sprintf("   Fee: $%.2f, Total: $%.2f\n", fee, total)
	result += fmt.Sprintf("   %s", s.paymentStrategy.Pay(total))

	return result
}

// ============= Example 2: Compression Strategy =============

// CompressionStrategy інтерфейс
type CompressionStrategy interface {
	Compress(data string) string
	GetRatio() float64
}

// ZIPCompression стратегія
type ZIPCompression struct{}

func (z *ZIPCompression) Compress(data string) string {
	return fmt.Sprintf("[ZIP compressed: %s...]", data[:min(10, len(data))])
}

func (z *ZIPCompression) GetRatio() float64 {
	return 0.5 // 50% compression
}

// GZIPCompression стратегія
type GZIPCompression struct{}

func (g *GZIPCompression) Compress(data string) string {
	return fmt.Sprintf("[GZIP compressed: %s...]", data[:min(8, len(data))])
}

func (g *GZIPCompression) GetRatio() float64 {
	return 0.6 // 60% compression
}

// NoCompression стратегія
type NoCompression struct{}

func (n *NoCompression) Compress(data string) string {
	return data
}

func (n *NoCompression) GetRatio() float64 {
	return 1.0 // no compression
}

// FileManager контекст
type FileManager struct {
	compression CompressionStrategy
}

func NewFileManager() *FileManager {
	return &FileManager{compression: &NoCompression{}}
}

func (f *FileManager) SetCompression(strategy CompressionStrategy) {
	f.compression = strategy
}

func (f *FileManager) SaveFile(filename, data string) string {
	compressed := f.compression.Compress(data)
	originalSize := len(data)
	compressedSize := int(float64(originalSize) * f.compression.GetRatio())

	return fmt.Sprintf("💾 Saving %s\n"+
		"   Original: %d bytes\n"+
		"   Compressed: %d bytes (%.0f%%)\n"+
		"   Data: %s",
		filename, originalSize, compressedSize,
		f.compression.GetRatio()*100, compressed)
}

// ============= Example 3: Sorting Strategy =============

// SortStrategy інтерфейс
type SortStrategy interface {
	Sort(data []int) []int
	Name() string
}

// BubbleSort стратегія
type BubbleSort struct{}

func (b *BubbleSort) Sort(data []int) []int {
	result := make([]int, len(data))
	copy(result, data)

	for i := 0; i < len(result); i++ {
		for j := 0; j < len(result)-1-i; j++ {
			if result[j] > result[j+1] {
				result[j], result[j+1] = result[j+1], result[j]
			}
		}
	}
	return result
}

func (b *BubbleSort) Name() string {
	return "Bubble Sort"
}

// QuickSort стратегія
type QuickSort struct{}

func (q *QuickSort) Sort(data []int) []int {
	result := make([]int, len(data))
	copy(result, data)
	q.quicksort(result, 0, len(result)-1)
	return result
}

func (q *QuickSort) quicksort(arr []int, low, high int) {
	if low < high {
		pivot := arr[high]
		i := low - 1

		for j := low; j < high; j++ {
			if arr[j] < pivot {
				i++
				arr[i], arr[j] = arr[j], arr[i]
			}
		}
		arr[i+1], arr[high] = arr[high], arr[i+1]

		pi := i + 1
		q.quicksort(arr, low, pi-1)
		q.quicksort(arr, pi+1, high)
	}
}

func (q *QuickSort) Name() string {
	return "Quick Sort"
}

// Sorter контекст
type Sorter struct {
	strategy SortStrategy
}

func (s *Sorter) SetStrategy(strategy SortStrategy) {
	s.strategy = strategy
}

func (s *Sorter) Sort(data []int) []int {
	return s.strategy.Sort(data)
}

// ============= Helper Functions =============

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// ============= Main =============

func main() {
	fmt.Println("╔════════════════════════════════════════════════╗")
	fmt.Println("║         Strategy Pattern Demo                  ║")
	fmt.Println("╚════════════════════════════════════════════════╝")

	// ===== Example 1: Payment Strategy =====
	fmt.Println("\n🔹 Example 1: Payment Strategies")
	fmt.Println(strings.Repeat("─", 50))

	cart := NewShoppingCart()
	cart.AddItem("Laptop", 1200.00)
	cart.AddItem("Mouse", 25.00)
	cart.AddItem("Keyboard", 75.00)

	// Оплата карткою
	fmt.Println("\n💳 Paying with Credit Card:")
	cart.SetPaymentStrategy(&CreditCardPayment{CardNumber: "1234567890123456"})
	fmt.Println(cart.Checkout())

	// Оплата PayPal
	cart2 := NewShoppingCart()
	cart2.AddItem("Phone", 800.00)
	fmt.Println("\n💰 Paying with PayPal:")
	cart2.SetPaymentStrategy(&PayPalPayment{Email: "user@example.com"})
	fmt.Println(cart2.Checkout())

	// Оплата криптою
	cart3 := NewShoppingCart()
	cart3.AddItem("Tablet", 400.00)
	fmt.Println("\n₿  Paying with Crypto:")
	cart3.SetPaymentStrategy(&CryptoPayment{Currency: "Bitcoin"})
	fmt.Println(cart3.Checkout())

	// ===== Example 2: Compression Strategy =====
	fmt.Println("\n\n🔹 Example 2: Compression Strategies")
	fmt.Println(strings.Repeat("─", 50))

	fileManager := NewFileManager()
	data := "This is a large file with lots of data that needs compression"

	// No compression
	fmt.Println("\n📄 No Compression:")
	fileManager.SetCompression(&NoCompression{})
	fmt.Println(fileManager.SaveFile("document.txt", data))

	// ZIP compression
	fmt.Println("\n📦 ZIP Compression:")
	fileManager.SetCompression(&ZIPCompression{})
	fmt.Println(fileManager.SaveFile("document.txt", data))

	// GZIP compression
	fmt.Println("\n🗜️  GZIP Compression:")
	fileManager.SetCompression(&GZIPCompression{})
	fmt.Println(fileManager.SaveFile("document.txt", data))

	// ===== Example 3: Sorting Strategy =====
	fmt.Println("\n\n🔹 Example 3: Sorting Strategies")
	fmt.Println(strings.Repeat("─", 50))

	data2 := []int{64, 34, 25, 12, 22, 11, 90}
	fmt.Printf("\nOriginal: %v\n", data2)

	sorter := &Sorter{}

	// Bubble Sort
	sorter.SetStrategy(&BubbleSort{})
	sorted1 := sorter.Sort(data2)
	fmt.Printf("Bubble Sort: %v\n", sorted1)

	// Quick Sort
	sorter.SetStrategy(&QuickSort{})
	sorted2 := sorter.Sort(data2)
	fmt.Printf("Quick Sort: %v\n", sorted2)

	// ===== Runtime Strategy Change =====
	fmt.Println("\n\n🔹 Runtime Strategy Change")
	fmt.Println(strings.Repeat("─", 50))

	cart4 := NewShoppingCart()
	cart4.AddItem("Item", 100.00)

	fmt.Println("\n1️⃣ First try with Card:")
	cart4.SetPaymentStrategy(&CreditCardPayment{CardNumber: "1111222233334444"})
	fmt.Printf("   Fee: $%.2f\n", cart4.paymentStrategy.GetFee(100.00))

	fmt.Println("\n2️⃣ Changed mind, use Crypto:")
	cart4.SetPaymentStrategy(&CryptoPayment{Currency: "Bitcoin"})
	fmt.Printf("   Fee: $%.2f\n", cart4.paymentStrategy.GetFee(100.00))
	fmt.Println("   ✅ Strategy changed at runtime!")

	// ===== Summary =====
	fmt.Println("\n\n📝 ВИСНОВКИ")
	fmt.Println(strings.Repeat("─", 50))
	fmt.Println("✅ Взаємозамінні алгоритми")
	fmt.Println("✅ Легко додавати нові стратегії")
	fmt.Println("✅ Зміна поведінки в runtime")
	fmt.Println("✅ Чистий код без if/switch")
	fmt.Println("✅ Open/Closed Principle")

	fmt.Println("\n💡 ВИКОРИСТАННЯ:")
	fmt.Println("   - Payment methods")
	fmt.Println("   - Compression algorithms")
	fmt.Println("   - Sorting algorithms")
	fmt.Println("   - Route calculation")
	fmt.Println("   - Authentication methods")
}
