// This program scans ports on a given website to check which ones are open
// RUN PROGRAM WITH FLAGS
// go run port.go --webSite=scanme.webscantest.com
package main

import (
	"flag"
	"fmt"
	"net"
	"sync"
)

// Define command line flag for the website to scan
// Default value is scanme.nmap.org which is a site specifically for testing port scanning
var webSite = flag.String("site", "scanme.nmap.org", "url to scan ports")

func main() {
	// Parse command line flags
	flag.Parse()

	// Create a WaitGroup to synchronize all goroutines
	var wg sync.WaitGroup

	// Iterate through first 3000 ports
	for i := range 3000 {
		// Increment WaitGroup counter before launching goroutine
		wg.Add(1)

		// Launch goroutine for each port scan
		go func(port int) {
			// Ensure WaitGroup is decremented when goroutine completes
			defer wg.Done()

			// Attempt to establish TCP connection to host:port
			conn, err := net.Dial("tcp", net.JoinHostPort(*webSite, fmt.Sprintf("%d", port)))
			if err != nil {
				// If connection fails, port is closed or filtered
				return
			}
			// Close connection immediately after successful connection
			conn.Close()
			// Print message for open ports
			fmt.Printf("Port %d is open\n", port)
		}(i)
	}

	// Wait for all port scanning goroutines to complete
	wg.Wait()
}
