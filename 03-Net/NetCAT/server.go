package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
)

// Client represents a connected user in the chat system.
// It's a channel that can only send strings (chan<- string)
type Client chan<- string

// Global variables for managing the chat system
var (
	// IncomingClients channel receives new clients when they connect
	IncomingClients = make(chan Client)
	// LeavingClients channel receives clients when they disconnect
	LeavingClients = make(chan Client)
	// ChatMessages channel receives all messages to be broadcasted
	ChatMessages = make(chan string)
	// Host and Port for the server configuration
	Host = flag.String("host", "localhost", "host to connect to")
	Port = flag.Int("port", 3090, "port to connect to")
)

// HandleConn manages a single client connection
// It creates a message channel for the client, sends welcome message,
// and handles incoming messages until the client disconnects
func HandleConn(conn net.Conn) {
	defer conn.Close()

	// Create a channel for this client's messages
	clientMessages := make(chan string)
	// Start a goroutine to write messages to this client
	go MessageWriter(conn, clientMessages)

	// Get client's address as their name
	clientName := conn.RemoteAddr().String()

	// Send welcome message to the new client
	clientMessages <- fmt.Sprintf("Welcome to the chat, %s!", clientName)
	// Broadcast that a new client has joined
	ChatMessages <- fmt.Sprintf("New client %s has joined", clientName)
	// Register this client in the system
	IncomingClients <- clientMessages

	// Create a scanner to read messages from the client
	inputMessage := bufio.NewScanner(conn)
	// Read messages until the client disconnects
	for inputMessage.Scan() {
		// Broadcast the message to all clients
		ChatMessages <- clientName + ": " + inputMessage.Text()
	}

	// Client has disconnected
	LeavingClients <- clientMessages
	// Broadcast that the client has left
	ChatMessages <- fmt.Sprintf("Client %s has left", clientName)
}

// MessageWriter continuously reads from the client's message channel
// and writes the messages to the client's connection
func MessageWriter(conn net.Conn, clientMessages <-chan string) {
	// Range over the channel until it's closed
	for msg := range clientMessages {
		// Write each message to the client's connection
		fmt.Fprintln(conn, msg)
	}
}

// Broadcast manages the distribution of messages to all connected clients
// It maintains a map of all connected clients and handles:
// - Broadcasting messages to all clients
// - Adding new clients
// - Removing disconnected clients
func Broadcast() {
	// Map to keep track of all connected clients
	clients := make(map[Client]bool)

	// Infinite loop to handle all channel events
	for {
		select {
		// When a new message arrives
		case message := <-ChatMessages:
			// Send the message to all connected clients
			for client := range clients {
				client <- message
			}
		// When a new client connects
		case client := <-IncomingClients:
			// Add the client to the map
			clients[client] = true
		// When a client disconnects
		case leavingClient := <-LeavingClients:
			// Remove the client from the map
			delete(clients, leavingClient)
			// Close the client's message channel
			close(leavingClient)
		}
	}
}

// StartServer initializes the chat server
// It sets up the TCP listener and handles incoming connections
func StartServer() {
	// Create a TCP listener on the specified host and port
	listener, err := net.Listen("tcp", net.JoinHostPort(*Host, fmt.Sprintf("%d", *Port)))
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	// Start the broadcast goroutine
	go Broadcast()

	// Accept incoming connections
	for {
		// Wait for a new connection
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		// Handle the connection in a new goroutine
		go HandleConn(conn)
	}
}

// main is the entry point of the chat server application
// It parses command line flags and starts the chat server
func main() {
	// Parse command line flags (host and port)
	flag.Parse()

	// Log that the server is starting
	log.Println("Starting chat server...")

	// Start the chat server
	StartServer()
}
