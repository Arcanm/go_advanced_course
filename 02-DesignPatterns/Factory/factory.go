// The Factory Pattern is a creational design pattern that provides an interface for creating objects
// without explicitly specifying their exact classes. It encapsulates object creation logic in a
// separate component, allowing the system to be independent of how its objects are created.
//
// Key benefits of the Factory Pattern:
// - Decouples object creation from the code that uses the objects
// - Makes the system more flexible and extensible
// - Centralizes complex creation logic
// - Promotes consistency in object creation
//
// Common use cases:
// - When you need to create different but related objects
// - When you want to delegate object creation to subclasses
// - When you want to hide complex creation logic from client code
//
// In this example, we implement a computer factory that:
// 1. Defines an IProduct interface that all products must implement
// 2. Has a base Computer struct with common attributes and behaviors
// 3. Creates s pecific types (Laptop, Desktop) that inherit from Computer
// 4. Uses a Factory function to centralize and encapsulate creation logic
// 5. Allows clients to create objects without knowing implementation details

package main

import (
	"fmt"
)

// IProduct defines the interface that all products must implement
type IProduct interface {
	setStock(stock int)
	getStock() int
	setName(name string)
	getName() string
}

// Computer is the base structure that contains common attributes
type Computer struct {
	stock int
	name  string
}

// Methods of the Computer structure that implement the IProduct interface
func (c *Computer) setStock(stock int) {
	c.stock = stock
}

func (c Computer) getStock() int {
	return c.stock
}

func (c *Computer) setName(name string) {
	c.name = name
}

func (c Computer) getName() string {
	return c.name
}

// Laptop is a specific type of Computer
type Laptop struct {
	Computer
}

// newLaptop is a constructor function that creates a new Laptop instance
func newLaptop() IProduct {
	return &Laptop{
		Computer: Computer{
			stock: 11,
			name:  "Laptop",
		},
	}
}

// Desktop is another specific type of Computer
type Desktop struct {
	Computer
}

// newDesktop is a constructor function that creates a new Desktop instance
func newDesktop() IProduct {
	return &Desktop{
		Computer: Computer{
			stock: 66,
			name:  "Desktop",
		},
	}
}

// ComputerFactory is the factory function that centralizes the creation of different computer types.
// It receives the type of computer to create and returns a new instance of the corresponding type.
// This encapsulates all object creation logic in one place, making it easier to modify or extend.
func ComputerFactory(computerType IProduct) (IProduct, error) {
	switch computerType.(type) {
	case *Desktop:
		return newDesktop(), nil
	case *Laptop:
		return newLaptop(), nil
	default:
		return nil, fmt.Errorf("Unknown computer type: %T", computerType)
	}
}

// String implements the Stringer interface to display product information
func (c *Computer) String() string {
	s := fmt.Sprintf("Product: %s, with Stock: %d", c.name, c.stock)
	return s
}

func main() {
	// Example of Factory Pattern usage
	// The client only needs to know about the IProduct interface and the Factory function
	// It doesn't need to know the implementation details of each computer type
	laptop, _ := ComputerFactory(&Laptop{})
	desktop, _ := ComputerFactory(&Desktop{})
	fmt.Println(laptop)
	fmt.Println(desktop)
}
