package main

import (
	"fmt"
	"sync"
	"time"
)

func ExpensiveFibonacci(n int) int {
	fmt.Printf("Calculating expensive fibonacci for %d\n", n)
	time.Sleep(5 * time.Second)
	return n
}

type Service struct {
	InProgress map[int]bool
	IsPending  map[int][]chan int
	Cache      map[int]int
	Lock       sync.Mutex
}

func (s *Service) Work(job int) {
	s.Lock.Lock()

	// Check cache first
	if result, exists := s.Cache[job]; exists {
		fmt.Printf("Found in cache! Job %d: %d\n", job, result)
		fmt.Printf("Job %d completed with cached result %d\n", job, result)
		s.Lock.Unlock()
		return
	}

	// If job is in progress, wait for result
	if s.InProgress[job] {
		response := make(chan int)
		s.IsPending[job] = append(s.IsPending[job], response)
		s.Lock.Unlock()

		fmt.Printf("Waiting for response for job %d\n", job)
		resp := <-response
		fmt.Printf("Job %d finished with result %d\n", job, resp)
		return
	}

	// Mark job as in progress
	s.InProgress[job] = true
	s.Lock.Unlock()

	// Calculate result
	fmt.Printf("Calculating fibonacci for %d\n", job)
	result := ExpensiveFibonacci(job)

	// Update cache and notify pending workers
	s.Lock.Lock()
	s.Cache[job] = result
	s.InProgress[job] = false

	// Notify pending workers
	if pendingWorkers, exists := s.IsPending[job]; exists {
		for _, response := range pendingWorkers {
			response <- result
		}
		fmt.Printf("Result %d sent to %d workers\n", result, len(pendingWorkers))
		delete(s.IsPending, job)
	}
	s.Lock.Unlock()

	fmt.Printf("Job %d finished with result %d\n", job, result)
}

func NewService() *Service {
	return &Service{
		InProgress: make(map[int]bool),
		IsPending:  make(map[int][]chan int),
		Cache:      make(map[int]int),
	}
}

func main() {
	service := NewService()
	jobs := []int{44, 46, 47, 49, 44, 34, 47, 41, 41, 33,
		43, 40, 46, 43, 45, 49, 34, 47, 36, 43, 48, 41,
		35, 46, 33, 42, 49, 47, 32, 30, 50, 50, 30, 40,
		44, 30, 49, 34, 48, 43, 50, 42, 48, 31, 35, 30,
		33, 40, 40, 50, 49, 47, 36, 43, 48, 41, 35, 46,
	}
	var wg sync.WaitGroup
	wg.Add(len(jobs))

	for _, job := range jobs {
		go func(job int) {
			defer wg.Done()
			service.Work(job)
		}(job)
	}

	wg.Wait()
	wg.Add(len(jobs))

	for _, job := range jobs {
		go func(job int) {
			defer wg.Done()
			service.Work(job)
		}(job)
	}

	wg.Wait()
}
