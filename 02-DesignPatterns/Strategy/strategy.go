// The Strategy Pattern is a behavioral design pattern that enables selecting an algorithm's behavior at runtime.
// It lets you define a family of algorithms, encapsulate each one, and make them interchangeable.
//
// Key benefits of the Strategy Pattern:
// - Enables runtime switching between different algorithms
// - Promotes loose coupling between algorithms and code that uses them
// - Makes code more flexible and maintainable
// - Follows Open/Closed Principle by allowing new strategies without modifying existing code
//
// In this example, we implement a password protection system that:
// 1. Defines a HashAlgorithm interface for different hashing strategies
// 2. Has concrete implementations (SHA, MD5) of the hashing interface
// 3. Uses a PasswordProtector that can work with any hash algorithm
// 4. Allows switching between hash algorithms at runtime
// 5. Demonstrates how different strategies can be used interchangeably

package main

import "fmt"

// PasswordProtector holds user credentials and the selected hash algorithm
type PasswordProtector struct {
	user          string
	password      string
	hashAlgorithm HashAlgorithm
}

// HashAlgorithm defines the interface that all hash strategies must implement
type HashAlgorithm interface {
	Hash(p *PasswordProtector)
}

// NewPasswordProtector creates a new PasswordProtector instance with the specified hash algorithm
func NewPasswordProtector(user string, password string, hashAlgorithm HashAlgorithm) *PasswordProtector {
	return &PasswordProtector{
		user:          user,
		password:      password,
		hashAlgorithm: hashAlgorithm,
	}
}

// SetHashAlgorithm allows changing the hash strategy at runtime
func (p *PasswordProtector) SetHashAlgorithm(hashAlgorithm HashAlgorithm) {
	p.hashAlgorithm = hashAlgorithm
}

// Hash executes the selected hash algorithm on the password
func (p *PasswordProtector) Hash() {
	p.hashAlgorithm.Hash(p)
}

// SHA implements the HashAlgorithm interface using SHA strategy
type SHA struct{}

func (s *SHA) Hash(p *PasswordProtector) {
	fmt.Printf("Hashing password for %s using SHA\n", p.user)
}

// MD5 implements the HashAlgorithm interface using MD5 strategy
type MD5 struct{}

func (m *MD5) Hash(p *PasswordProtector) {
	fmt.Printf("Hashing password for %s using MD5\n", p.user)
}

func main() {
	// Create instances of different hash strategies
	sha := &SHA{}
	md5 := &MD5{}

	// Create password protector with initial SHA strategy
	passwordProtector := NewPasswordProtector("Andres", "password", sha)
	passwordProtector.Hash()

	// Switch to MD5 strategy at runtime
	passwordProtector.SetHashAlgorithm(md5)
	passwordProtector.Hash()
}
