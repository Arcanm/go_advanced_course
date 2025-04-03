package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

// Command line flags for client configuration
var (
	port = flag.Int("port", 3090, "port to connect to")
	host = flag.String("host", "localhost", "host to connect to")
)

// main is the entry point of the chat client application
// It establishes a connection to the chat server and handles bidirectional communication
func main() {
	// Parse command line flags
	flag.Parse()

	// Connect to the chat server
	conn, err := net.Dial("tcp", net.JoinHostPort(*host, fmt.Sprintf("%d", *port)))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// Channel to signal when either goroutine finishes
	done := make(chan struct{})

	// Goroutine to read from the server and write to stdout
	// This handles incoming messages from other clients
	go func() {
		// Copy all data from the connection to stdout
		io.Copy(os.Stdout, conn)
		// Log when the connection is closed
		log.Println("Connection closed by remote host")
		// Signal that this goroutine is done
		done <- struct{}{}
	}()

	// Goroutine to read from stdin and write to the server
	// This handles outgoing messages from this client
	go func() {
		// Copy all data from stdin to the connection
		io.Copy(conn, os.Stdin)
		// Signal that this goroutine is done
		done <- struct{}{}
	}()

	// Wait for either goroutine to finish
	// This blocks until the connection is closed or the user exits
	<-done
}
