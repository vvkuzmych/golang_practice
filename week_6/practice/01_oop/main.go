package main

import "fmt"

// ===== 1. Інкапсуляція =====

type BankAccount struct {
	owner   string  // приватне поле
	balance float64 // приватне поле
}

func NewBankAccount(owner string, initialBalance float64) *BankAccount {
	return &BankAccount{
		owner:   owner,
		balance: initialBalance,
	}
}

func (a *BankAccount) Deposit(amount float64) {
	if amount > 0 {
		a.balance += amount
		fmt.Printf("Deposited $%.2f. New balance: $%.2f\n", amount, a.balance)
	}
}

func (a *BankAccount) Withdraw(amount float64) error {
	if amount > a.balance {
		return fmt.Errorf("insufficient funds")
	}
	a.balance -= amount
	fmt.Printf("Withdrawn $%.2f. New balance: $%.2f\n", amount, a.balance)
	return nil
}

func (a *BankAccount) GetBalance() float64 {
	return a.balance
}

// ===== 2. Поліморфізм =====

type Shape interface {
	Area() float64
	Name() string
}

type Circle struct {
	Radius float64
}

func (c Circle) Area() float64 {
	return 3.14 * c.Radius * c.Radius
}

func (c Circle) Name() string {
	return "Circle"
}

type Rectangle struct {
	Width, Height float64
}

func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

func (r Rectangle) Name() string {
	return "Rectangle"
}

func PrintShapeInfo(s Shape) {
	fmt.Printf("%s area: %.2f\n", s.Name(), s.Area())
}

// ===== 3. Композиція =====

type Engine struct {
	Horsepower int
}

func (e *Engine) Start() {
	fmt.Println("Engine started")
}

type Car struct {
	Brand  string
	Engine Engine // композиція
}

func (c *Car) Drive() {
	fmt.Printf("Driving %s\n", c.Brand)
	c.Engine.Start()
}

func main() {
	fmt.Println("=== 1. Інкапсуляція ===")
	account := NewBankAccount("John", 1000)
	account.Deposit(500)
	account.Withdraw(200)
	fmt.Printf("Final balance: $%.2f\n", account.GetBalance())

	fmt.Println("\n=== 2. Поліморфізм ===")
	shapes := []Shape{
		Circle{Radius: 5},
		Rectangle{Width: 10, Height: 20},
	}
	for _, shape := range shapes {
		PrintShapeInfo(shape)
	}

	fmt.Println("\n=== 3. Композиція ===")
	car := Car{
		Brand:  "Toyota",
		Engine: Engine{Horsepower: 150},
	}
	car.Drive()
}
