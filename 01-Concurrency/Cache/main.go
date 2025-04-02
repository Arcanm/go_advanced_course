package main

import (
	"fmt"
	"sync"
	"time"
)

// FibonacciCached calculates the Fibonacci number for a given value 'n'
// It uses a caching system to store previously calculated results
// and avoid redundant calculations
// Parameters:
//   - n: The position in the Fibonacci sequence to calculate
//   - m: Pointer to the Memory cache system
//
// Returns: The Fibonacci number at position n
func FibonacciCached(n int, m *Memory) int {
	if n <= 1 {
		return n
	}
	// Gets the previous values from cache and adds them
	return m.Get(n-1) + m.Get(n-2)
}

// Function is a type that defines the signature of functions that can be cached
// It takes a key and a pointer to the cache memory system
type Function func(key int, m *Memory) int

// Memory implements a thread-safe caching system
// This structure ensures safe concurrent access to cached values
type Memory struct {
	f     Function    // The function to be cached
	cache map[int]int // Map that stores cached results
	mux   sync.Mutex  // Mutex to ensure thread-safe access to the cache
}

// NewCache creates a new instance of the caching system
// Parameters:
//   - f: The function to be cached
//
// Returns: A pointer to a new Memory instance
func NewCache(f Function) *Memory {
	return &Memory{
		f:     f,
		cache: make(map[int]int),
	}
}

// Get retrieves a value from the cache. If it doesn't exist, calculates and stores it
// This method is thread-safe thanks to mutex implementation
// Parameters:
//   - key: The input value for which we want to cache the result
//
// Returns: The cached or newly calculated result
func (m *Memory) Get(key int) int {
	// First attempt to read from cache, protected by mutex
	m.mux.Lock()
	result, exists := m.cache[key]
	m.mux.Unlock()

	// If the value doesn't exist in cache, we calculate it
	if !exists {
		m.mux.Lock()
		// Calculate the result using the stored function
		result = m.f(key, m)
		// Store the result in cache
		m.cache[key] = result
		m.mux.Unlock()
	}
	return result
}

func main() {
	// Create a new cache instance for the Fibonacci function
	cache := NewCache(FibonacciCached)

	// List of Fibonacci numbers we want to calculate
	// Note that some numbers are repeated to demonstrate cache effectiveness
	tasks := []int{42, 40, 41, 42, 38, 1000}

	// Calculate each number and measure the time taken
	// This demonstrates how subsequent calculations of the same number
	// are much faster due to caching
	for _, n := range tasks {
		start := time.Now()
		value := cache.Get(n)
		// Print: calculated number, elapsed time, and result
		fmt.Printf(" %d, %s, %d\n", n, time.Since(start), value)
	}
}
