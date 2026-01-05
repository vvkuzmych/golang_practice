package main

import (
	"errors"
	"fmt"
	"math"
)

// ============= Calculator =============

type Calculator struct {
	result  float64
	history []string
}

// NewCalculator створює новий калькулятор
func NewCalculator() *Calculator {
	return &Calculator{
		result:  0,
		history: []string{},
	}
}

// Add додає значення
func (c *Calculator) Add(value float64) *Calculator {
	c.result += value
	c.addToHistory(fmt.Sprintf("+ %.2f = %.2f", value, c.result))
	return c
}

// Subtract віднімає значення
func (c *Calculator) Subtract(value float64) *Calculator {
	c.result -= value
	c.addToHistory(fmt.Sprintf("- %.2f = %.2f", value, c.result))
	return c
}

// Multiply множить на значення
func (c *Calculator) Multiply(value float64) *Calculator {
	c.result *= value
	c.addToHistory(fmt.Sprintf("× %.2f = %.2f", value, c.result))
	return c
}

// Divide ділить на значення
func (c *Calculator) Divide(value float64) error {
	if value == 0 {
		return errors.New("division by zero")
	}
	c.result /= value
	c.addToHistory(fmt.Sprintf("÷ %.2f = %.2f", value, c.result))
	return nil
}

// Sqrt обчислює квадратний корінь
func (c *Calculator) Sqrt() error {
	if c.result < 0 {
		return errors.New("cannot take square root of negative number")
	}
	oldResult := c.result
	c.result = math.Sqrt(c.result)
	c.addToHistory(fmt.Sprintf("√%.2f = %.2f", oldResult, c.result))
	return nil
}

// Power піднесення до степеня
func (c *Calculator) Power(exp float64) *Calculator {
	oldResult := c.result
	c.result = math.Pow(c.result, exp)
	c.addToHistory(fmt.Sprintf("%.2f^%.2f = %.2f", oldResult, exp, c.result))
	return c
}

// Result повертає поточний результат
func (c Calculator) Result() float64 {
	return c.result
}

// Reset скидає результат до 0
func (c *Calculator) Reset() {
	c.result = 0
	c.history = []string{}
	c.addToHistory("Reset to 0")
}

// History повертає історію операцій
func (c Calculator) History() []string {
	return c.history
}

// String текстове представлення
func (c Calculator) String() string {
	return fmt.Sprintf("Calculator: %.2f", c.result)
}

func (c *Calculator) addToHistory(operation string) {
	c.history = append(c.history, operation)
}

// ============= Main =============

func main() {
	fmt.Println("╔══════════════════════════════════════════╗")
	fmt.Println("║        Calculator Solution               ║")
	fmt.Println("╚══════════════════════════════════════════╝")

	// ===== Basic Operations =====
	fmt.Println("\n🔹 Базові операції")
	fmt.Println("─────────────────────────────────────────")

	calc := NewCalculator()
	fmt.Printf("Початкове значення: %.2f\n\n", calc.Result())

	calc.Add(10)
	fmt.Printf("Після Add(10): %.2f\n", calc.Result())

	calc.Add(5)
	fmt.Printf("Після Add(5): %.2f\n", calc.Result())

	calc.Multiply(2)
	fmt.Printf("Після Multiply(2): %.2f\n", calc.Result())

	calc.Subtract(10)
	fmt.Printf("Після Subtract(10): %.2f\n", calc.Result())

	calc.Divide(4)
	fmt.Printf("Після Divide(4): %.2f\n", calc.Result())

	fmt.Printf("\n%s\n", calc)

	// ===== Error Handling =====
	fmt.Println("\n\n🔹 Обробка помилок")
	fmt.Println("─────────────────────────────────────────")

	fmt.Printf("Поточне значення: %.2f\n", calc.Result())

	err := calc.Divide(0)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
	}
	fmt.Printf("Значення після помилки: %.2f (не змінилось)\n", calc.Result())

	// Негативний корінь
	calc.Reset()
	calc.Add(-25)
	err = calc.Sqrt()
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
	}

	// ===== Advanced Operations =====
	fmt.Println("\n\n🔹 Додаткові операції")
	fmt.Println("─────────────────────────────────────────")

	calc.Reset()
	calc.Add(16)
	fmt.Printf("Значення: %.2f\n", calc.Result())

	calc.Sqrt()
	fmt.Printf("Після Sqrt(): %.2f\n", calc.Result())

	calc.Power(3)
	fmt.Printf("Після Power(3): %.2f\n", calc.Result())

	// ===== Chaining =====
	fmt.Println("\n\n🔹 Chainable методи")
	fmt.Println("─────────────────────────────────────────")

	calc.Reset()
	calc.Add(10).Multiply(2).Subtract(5).Add(3)
	fmt.Printf("Результат ланцюга: %.2f\n", calc.Result())

	// ===== History =====
	fmt.Println("\n\n🔹 Історія операцій")
	fmt.Println("─────────────────────────────────────────")

	for i, op := range calc.History() {
		fmt.Printf("%d. %s\n", i+1, op)
	}

	// ===== Complex Example =====
	fmt.Println("\n\n🔹 Складний приклад")
	fmt.Println("─────────────────────────────────────────")

	calc2 := NewCalculator()

	// (5 + 3) × 2 - 10 / 2
	calc2.Add(5).
		Add(3).
		Multiply(2).
		Subtract(10)

	calc2.Divide(2)

	fmt.Printf("Результат: %.2f\n", calc2.Result())
	fmt.Println("\nІсторія:")
	for _, op := range calc2.History() {
		fmt.Printf("  %s\n", op)
	}

	// ===== Multiple Calculators =====
	fmt.Println("\n\n🔹 Кілька калькуляторів")
	fmt.Println("─────────────────────────────────────────")

	calc_a := NewCalculator()
	calc_b := NewCalculator()

	calc_a.Add(100).Divide(4)
	calc_b.Add(50).Multiply(2)

	fmt.Printf("Калькулятор A: %.2f\n", calc_a.Result())
	fmt.Printf("Калькулятор B: %.2f\n", calc_b.Result())

	// Сума результатів
	calc_sum := NewCalculator()
	calc_sum.Add(calc_a.Result()).Add(calc_b.Result())
	fmt.Printf("Сума: %.2f\n", calc_sum.Result())

	// ===== Summary =====
	fmt.Println("\n\n📝 Висновки")
	fmt.Println("─────────────────────────────────────────")
	fmt.Println("✅ Реалізовано:")
	fmt.Println("   • Базові операції (+, -, ×, ÷)")
	fmt.Println("   • Додаткові операції (√, ^)")
	fmt.Println("   • Обробка помилок")
	fmt.Println("   • Історія операцій")
	fmt.Println("   • Chainable методи")
	fmt.Println("   • Pointer receivers для зміни стану")
	fmt.Println("   • Value receiver для читання")
}
