// The Adapter Pattern is a structural design pattern that allows objects with incompatible interfaces
// to collaborate. It acts as a wrapper between two objects, catching calls for one object and
// transforming them to format and interface recognizable by the second object.
//
// Key benefits of the Adapter Pattern:
// - Allows integration of incompatible interfaces
// - Promotes code reusability
// - Increases flexibility by decoupling client code from implementation
// - Makes existing code work with new code without modification
//
// Common use cases:
// - When you need to use an existing class but its interface isn't compatible
// - When you want to create a reusable class that cooperates with classes that don't have compatible interfaces
// - When you need to integrate third-party code without modifying it
//
// In this example, we implement a payment system that:
// 1. Defines a Payment interface for standard payment operations
// 2. Has a CashPayment implementation that follows the interface
// 3. Has a BankPayment implementation with a different interface
// 4. Uses an Adapter to make BankPayment compatible with the Payment interface
// 5. Demonstrates how both payment types can be processed uniformly

package main

import "fmt"

// Payment defines the standard interface for all payment methods
type Payment interface {
	Pay()
}

// CashPayment implements the Payment interface directly
type CashPayment struct{}

func (c *CashPayment) Pay() {
	fmt.Println("Paying with cash")
}

// ProcessPayment handles any payment method that implements the Payment interface
func ProcessPayment(p Payment) {
	p.Pay()
}

// BankPayment represents a payment system with an incompatible interface
type BankPayment struct{}

func (b *BankPayment) Pay(amount int) {
	fmt.Printf("Paying %d with bank transfer\n", amount)
}

// BankPaymentAdapter adapts BankPayment to match the Payment interface
type BankPaymentAdapter struct {
	bankPayment *BankPayment
	bankAccount int
}

// Pay implements the Payment interface for BankPaymentAdapter
func (b *BankPaymentAdapter) Pay() {
	b.bankPayment.Pay(b.bankAccount)
}

func main() {
	// Example of direct Payment interface usage
	cash := &CashPayment{}
	ProcessPayment(cash)

	// Example of adapted payment method usage
	bankAdapter := &BankPaymentAdapter{
		bankPayment: &BankPayment{},
		bankAccount: 5,
	}
	ProcessPayment(bankAdapter)
}
