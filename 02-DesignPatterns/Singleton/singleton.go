// The Singleton Pattern is a creational design pattern that ensures a class has only one
// instance and provides a global access point to it. It's useful when we need to coordinate
// actions through a single control point, such as:
// - Database connections
// - Global configurations
// - Shared caches
//
// In this example, we implement the Singleton for a database connection:
// 1. We have a global 'db' variable that holds the single instance
// 2. We use a mutex to ensure thread-safe access to instance creation
// 3. The getDatabaseInstance() method implements "lazy initialization" logic
// 4. The main() demo shows how multiple goroutines try to access the same instance

package main

import (
	"fmt"
	"sync"
	"time"
)

// Database represents our "connection" to the database
type Database struct{}

// mutex to ensure thread-safety in instance creation
var mutex = sync.Mutex{}

// CreateSingleConnection simulates creating a database connection
func (Database) CreateSingleConnection() {
	fmt.Println("Connection singleton for database")
	time.Sleep(3 * time.Second) // Simulate connection work
	fmt.Println("Connection created")
}

// The only Database instance that will exist
var db *Database

// getDatabaseInstance implements the Singleton pattern
// Returns the single instance, creating it if it doesn't exist
func getDatabaseInstance() *Database {
	mutex.Lock()
	defer mutex.Unlock()

	if db == nil {
		fmt.Println("Creating new database instance")
		db = &Database{}
		db.CreateSingleConnection()
	} else {
		fmt.Println("Database instance already created")
	}
	return db
}

func main() {
	// Demonstrate the Singleton with multiple goroutines
	var wg sync.WaitGroup
	for range 10 {
		wg.Add(1)

		go func() {
			defer wg.Done()
			getDatabaseInstance()
		}()
	}
	wg.Wait()
}
