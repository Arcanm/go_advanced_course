package main

import (
	"fmt"
	"sync"
	"time"
)

// ExpensiveFibonacci simulates an expensive calculation by adding a delay
// Parameters:
//   - n: The input number for which we want to calculate fibonacci
//
// Returns: For demonstration purposes, it just returns the input number
func ExpensiveFibonacci(n int) int {
	fmt.Printf("Calculating expensive fibonacci for %d\n", n)
	// Simulate heavy processing by sleeping for 5 seconds
	time.Sleep(5 * time.Second)
	return n
}

// Service implements a concurrent job processing system with deduplication
// It ensures that identical jobs are only processed once, while multiple
// requesters wait for the same result
type Service struct {
	InProgress map[int]bool       // Tracks which jobs are currently being calculated
	IsPending  map[int][]chan int // Stores channels for workers waiting for results
	Lock       sync.RWMutex       // RWMutex allows multiple simultaneous reads but exclusive writes
}

// Work processes a job request, implementing deduplication logic
// If the same job is already in progress, it waits for the result
// If it's a new job, it performs the calculation and notifies any waiting workers
// Parameters:
//   - job: The number for which we want to calculate fibonacci
func (s *Service) Work(job int) {
	// First, check if this job is already being processed
	// Using RLock() for better performance as multiple goroutines can read simultaneously
	s.Lock.RLock()
	exists := s.InProgress[job]
	if exists {
		// If job exists, release read lock before proceeding
		s.Lock.RUnlock()

		// Create a channel to receive the result when ready
		response := make(chan int)
		defer close(response) // Ensure channel cleanup when function returns

		// Switch to exclusive lock to modify the pending workers list
		s.Lock.Lock()
		s.IsPending[job] = append(s.IsPending[job], response)
		s.Lock.Unlock()

		fmt.Printf("Waiting for response for job %d\n", job)
		// Block until result is received through the channel
		resp := <-response
		fmt.Printf("Job %d finished with result %d\n", job, resp)
		return
	}
	s.Lock.RUnlock()

	// If job doesn't exist, acquire exclusive lock to mark it as in progress
	s.Lock.Lock()
	s.InProgress[job] = true
	s.Lock.Unlock()

	fmt.Printf("Calculating fibonacci for %d\n", job)
	// Perform the expensive calculation
	result := ExpensiveFibonacci(job)

	// Check if there are any workers waiting for this result
	s.Lock.RLock()
	pendingWorkers, exists := s.IsPending[job]
	s.Lock.RUnlock()

	if exists {
		// Notify all waiting workers by sending the result through their channels
		for _, response := range pendingWorkers {
			response <- result
		}
		fmt.Printf("Result %d sent to %d workers\n", result, len(pendingWorkers))
	}

	// Clean up: mark job as complete and clear pending workers list
	s.Lock.Lock()
	s.InProgress[job] = false
	s.IsPending[job] = make([]chan int, 0)
	s.Lock.Unlock()
}

// NewService creates and initializes a new Service instance
// Returns: A pointer to the new Service with initialized maps
func NewService() *Service {
	return &Service{
		InProgress: make(map[int]bool),
		IsPending:  make(map[int][]chan int),
	}
}

// massiveOperations demonstrates the Service usage with concurrent job processing
// It creates multiple goroutines that process both unique and duplicate jobs,
// showing how the service handles deduplication and concurrent access
func massiveOperations() {
	service := NewService()
	// Define a list of jobs with some duplicates to demonstrate deduplication
	jobs := []int{3, 4, 5, 5, 4, 3, 2, 1, 0}

	// Use WaitGroup to ensure all goroutines complete before function returns
	var wg sync.WaitGroup
	wg.Add(len(jobs))

	// Launch a goroutine for each job
	for _, job := range jobs {
		go func(job int) {
			defer wg.Done() // Ensure WaitGroup is decremented when goroutine completes
			service.Work(job)
		}(job)
	}

	// Wait for all goroutines to complete
	wg.Wait()
}
