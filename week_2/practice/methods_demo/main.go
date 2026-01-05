package main

import "fmt"

// ============= Rectangle with Value and Pointer Receivers =============

type Rectangle struct {
	Width  int
	Height int
}

// Value receiver - тільки читання
func (r Rectangle) Area() int {
	return r.Width * r.Height
}

func (r Rectangle) Perimeter() int {
	return 2 * (r.Width + r.Height)
}

// Value receiver - НЕ змінює оригінал
func (r Rectangle) ScaleValue(factor int) {
	r.Width *= factor
	r.Height *= factor
	fmt.Printf("  Всередині ScaleValue: %+v\n", r)
}

// Pointer receiver - змінює оригінал
func (r *Rectangle) ScalePointer(factor int) {
	r.Width *= factor
	r.Height *= factor
	fmt.Printf("  Всередині ScalePointer: %+v\n", r)
}

func (r Rectangle) String() string {
	return fmt.Sprintf("Rectangle{Width: %d, Height: %d}", r.Width, r.Height)
}

// ============= Counter =============

type Counter struct {
	count int
}

// Value receiver - не працює для зміни
func (c Counter) IncrementValue() {
	c.count++
}

// Pointer receiver - працює для зміни
func (c *Counter) IncrementPointer() {
	c.count++
}

func (c Counter) Value() int {
	return c.count
}

// ============= BankAccount =============

type BankAccount struct {
	owner   string
	balance float64
}

func NewBankAccount(owner string, initialBalance float64) *BankAccount {
	return &BankAccount{
		owner:   owner,
		balance: initialBalance,
	}
}

// Pointer receivers для зміни даних
func (b *BankAccount) Deposit(amount float64) {
	if amount > 0 {
		b.balance += amount
		fmt.Printf("  ✅ Депозит %s: +%.2f грн\n", b.owner, amount)
	}
}

func (b *BankAccount) Withdraw(amount float64) bool {
	if amount > 0 && b.balance >= amount {
		b.balance -= amount
		fmt.Printf("  ✅ Зняття %s: -%.2f грн\n", b.owner, amount)
		return true
	}
	fmt.Printf("  ❌ Зняття %s: недостатньо коштів (є %.2f, потрібно %.2f)\n", b.owner, b.balance, amount)
	return false
}

// Value receiver для читання
func (b BankAccount) Balance() float64 {
	return b.balance
}

func (b BankAccount) Info() string {
	return fmt.Sprintf("%s: %.2f грн", b.owner, b.balance)
}

// ============= Main =============

func main() {
	fmt.Println("╔══════════════════════════════════════════╗")
	fmt.Println("║  Methods Demo: Value vs Pointer         ║")
	fmt.Println("╚══════════════════════════════════════════╝")

	// ===== Rectangle Demo =====
	fmt.Println("\n📐 RECTANGLE DEMO")
	fmt.Println("─────────────────────────────────────────")

	rect := Rectangle{Width: 10, Height: 5}
	fmt.Printf("Початковий: %s\n", rect)
	fmt.Printf("Area: %d\n", rect.Area())
	fmt.Printf("Perimeter: %d\n", rect.Perimeter())

	fmt.Println("\n🔹 ScaleValue (value receiver):")
	rect.ScaleValue(2)
	fmt.Printf("Після ScaleValue: %s ← НЕ ЗМІНИВСЯ!\n", rect)

	fmt.Println("\n🔹 ScalePointer (pointer receiver):")
	rect.ScalePointer(2)
	fmt.Printf("Після ScalePointer: %s ← ЗМІНИВСЯ!\n", rect)

	// ===== Counter Demo =====
	fmt.Println("\n\n🔢 COUNTER DEMO")
	fmt.Println("─────────────────────────────────────────")

	counter := Counter{count: 0}
	fmt.Printf("Початкове значення: %d\n", counter.Value())

	fmt.Println("\n🔹 IncrementValue (value receiver):")
	counter.IncrementValue()
	fmt.Printf("Після IncrementValue: %d ← НЕ ЗБІЛЬШИВСЯ!\n", counter.Value())

	fmt.Println("\n🔹 IncrementPointer (pointer receiver):")
	counter.IncrementPointer()
	fmt.Printf("Після IncrementPointer: %d ← ЗБІЛЬШИВСЯ!\n", counter.Value())

	counter.IncrementPointer()
	counter.IncrementPointer()
	fmt.Printf("Після ще 2x IncrementPointer: %d\n", counter.Value())

	// ===== BankAccount Demo =====
	fmt.Println("\n\n💰 BANK ACCOUNT DEMO")
	fmt.Println("─────────────────────────────────────────")

	account := NewBankAccount("Іван Петренко", 1000.0)
	fmt.Printf("Створено: %s\n", account.Info())

	fmt.Println("\nОперації:")
	account.Deposit(500.0)
	account.Withdraw(300.0)
	account.Withdraw(2000.0) // недостатньо коштів
	account.Deposit(200.0)

	fmt.Printf("\nПідсумок: %s\n", account.Info())

	// ===== Multiple Accounts =====
	fmt.Println("\n\n👥 MULTIPLE ACCOUNTS")
	fmt.Println("─────────────────────────────────────────")

	accounts := []*BankAccount{
		NewBankAccount("Марія Іванова", 5000.0),
		NewBankAccount("Петро Сидоренко", 3000.0),
		NewBankAccount("Оксана Коваль", 7500.0),
	}

	fmt.Println("Всі рахунки:")
	for _, acc := range accounts {
		fmt.Printf("  • %s\n", acc.Info())
	}

	// Переказ
	fmt.Println("\n💸 Переказ 1000 грн від Оксани до Марії")
	if accounts[2].Withdraw(1000.0) {
		accounts[0].Deposit(1000.0)
	}

	fmt.Println("\nПісля переказу:")
	for _, acc := range accounts {
		fmt.Printf("  • %s\n", acc.Info())
	}

	// ===== Summary =====
	fmt.Println("\n\n📝 ВИСНОВКИ")
	fmt.Println("─────────────────────────────────────────")
	fmt.Println("✅ Value receiver:")
	fmt.Println("   • Отримує КОПІЮ")
	fmt.Println("   • НЕ може змінити оригінал")
	fmt.Println("   • Використовується для читання даних")
	fmt.Println()
	fmt.Println("✅ Pointer receiver:")
	fmt.Println("   • Отримує ВКАЗІВНИК")
	fmt.Println("   • МОЖЕ змінити оригінал")
	fmt.Println("   • Використовується для зміни даних")
	fmt.Println()
	fmt.Println("💡 Правило: Якщо метод змінює дані → pointer receiver")
}
