package main

import "fmt"

// 12. Strategy Pattern - Шаблон Стратегія

type PaymentStrategy interface {
	Pay(amount float64) string
}

type CreditCard struct{ CardNumber string }

func (c CreditCard) Pay(amount float64) string {
	return fmt.Sprintf("Paid %.2f with Credit Card %s", amount, c.CardNumber)
}

type PayPal struct{ Email string }

func (p PayPal) Pay(amount float64) string {
	return fmt.Sprintf("Paid %.2f with PayPal %s", amount, p.Email)
}

type Cash struct{}

func (c Cash) Pay(amount float64) string {
	return fmt.Sprintf("Paid %.2f in Cash", amount)
}

type ShoppingCart struct {
	paymentMethod PaymentStrategy
}

func (s *ShoppingCart) SetPaymentMethod(pm PaymentStrategy) {
	s.paymentMethod = pm
}

func (s *ShoppingCart) Checkout(amount float64) {
	result := s.paymentMethod.Pay(amount)
	fmt.Println(result)
}

func main() {
	cart := &ShoppingCart{}

	cart.SetPaymentMethod(CreditCard{"1234-5678"})
	cart.Checkout(100.50)

	cart.SetPaymentMethod(PayPal{"user@example.com"})
	cart.Checkout(75.25)

	cart.SetPaymentMethod(Cash{})
	cart.Checkout(50.00)
}
