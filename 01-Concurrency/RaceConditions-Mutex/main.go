// RACE CONDITIONS
// A race condition occurs when multiple goroutines access shared resources concurrently
// To detect race conditions, we can use the -race flag when running the program:
// go run -race main.go

package main

import (
	"fmt"
	"sync"
)

// balance represents our shared resource that multiple goroutines will access
var balance int = 100

// Deposit adds the given amount to the balance
// It takes a WaitGroup to coordinate goroutines and a RWMutex to prevent race conditions
// Only one goroutine can modify the balance at a time using Lock() for write operations
func Deposit(amount int, wg *sync.WaitGroup, mux *sync.RWMutex) {
	// Notify the WaitGroup that this goroutine is done when the function returns
	defer wg.Done()
	// Lock the mutex to ensure exclusive access to the balance
	// Lock() is used for write operations, blocking all other read and write operations
	mux.Lock()
	balance += amount
	// Unlock the mutex to allow other goroutines to access the balance
	mux.Unlock()
}

// Withdraw subtracts the given amount from the balance
// It takes a WaitGroup to coordinate goroutines and a RWMutex to prevent race conditions
// Similar to Deposit, it requires exclusive write access using Lock()
func Withdraw(amount int, wg *sync.WaitGroup, mux *sync.RWMutex) {
	defer wg.Done()
	// Lock() acquires exclusive write access, preventing any concurrent reads or writes
	mux.Lock()
	balance -= amount
	mux.Unlock()
}

// Balance returns the current balance
// It uses RLock() (Read Lock) since it only needs read access
// Multiple goroutines can read the balance simultaneously using RLock()
// However, if any goroutine has a write lock, read operations will be blocked
func Balance(mux *sync.RWMutex) int {
	// RLock() allows multiple readers to access the balance concurrently
	mux.RLock()
	b := balance
	// RUnlock() releases the read lock
	mux.RUnlock()
	return b
}

func main() {
	// Create a WaitGroup to wait for all goroutines to complete
	var wg sync.WaitGroup
	// Create a RWMutex (Read-Write Mutex) to prevent race conditions
	// RWMutex allows multiple readers but only one writer at a time
	// This is more efficient than a regular Mutex when we have more read operations than write operations
	var mux sync.RWMutex

	// Print initial balance using the read-only Balance function
	fmt.Printf("Initial balance: %d\n", Balance(&mux))

	// Launch 5 goroutines that deposit increasing amounts
	// Each goroutine will need exclusive write access using Lock()
	for index := 1; index <= 5; index++ {
		wg.Add(1) // Increment WaitGroup counter
		go Deposit(index*100, &wg, &mux)
	}

	// Launch 2 more goroutines with fixed deposit amounts
	// These also require exclusive write access
	wg.Add(2)
	go Deposit(100, &wg, &mux)
	go Deposit(200, &wg, &mux)

	// Launch 3 withdrawal goroutines
	// Each withdrawal needs exclusive write access to modify the balance
	wg.Add(3)
	go Withdraw(300, &wg, &mux)
	go Withdraw(200, &wg, &mux)
	go Withdraw(100, &wg, &mux)

	// Add a small delay to see intermediate balance
	// This demonstrates that we can read the balance while operations are in progress
	// Wait for all goroutines to complete their operations
	wg.Wait()
	// Print the final balance using the read-only Balance function
	fmt.Printf("Final balance: %d\n", Balance(&mux))
}
