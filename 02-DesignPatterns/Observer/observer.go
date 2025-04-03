// The Observer Pattern is a behavioral design pattern that defines a subscription mechanism
// to notify multiple objects about any events happening to the object they're observing.
//
// In this example, we implement a product availability notification system where:
// 1. We have a Topic (Item) that can be observed
// 2. We have Observers (EmailClient) that subscribe to receive updates
// 3. When the Item becomes available, it notifies all its observers
// 4. The observers receive the notification and execute their logic (send email)

package main

import "fmt"

// Topic defines the interface for objects that can be observed
type Topic interface {
	// Register adds a new observer to receive updates
	Register(observer Observer)
	// Broadcast notifies all registered observers
	Broadcast()
}

// Observer defines the interface for objects that want to receive updates
type Observer interface {
	// getId returns the unique identifier of the observer
	getId() string
	// updateValue receives updates from the Topic
	updateValue(string)
}

// Item represents a product that can be available or not
// Implements the Topic interface to be observable
type Item struct {
	observers []Observer // List of subscribed observers
	name      string     // Product name
	price     int        // Product price
}

// NewItem creates a new Item instance with the specified name
func NewItem(name string) *Item {
	return &Item{
		name: name,
	}
}

// UpdateAvailable marks the item as available, updates its price and notifies observers
func (i *Item) UpdateAvailable() {
	fmt.Printf("The item %s is now available\n", i.name)
	i.price = 100
	i.Broadcast()
}

// Register adds a new observer to the item's list of observers
func (i *Item) Register(observer Observer) {
	i.observers = append(i.observers, observer)
}

// Broadcast notifies all registered observers about changes in the item
func (i *Item) Broadcast() {
	for _, observer := range i.observers {
		observer.updateValue(i.name)
	}
}

// EmailClient represents a client that will receive email notifications
// Implements the Observer interface
type EmailClient struct {
	id string // Client's email
}

// updateValue implements the email notification logic when an item becomes available
func (e *EmailClient) updateValue(itemName string) {
	fmt.Printf("Sending email - %s is now available for client %s\n", itemName, e.id)
}

// getId returns the client's email
func (e EmailClient) getId() string {
	return e.id
}

// SmsClient represents a client that will receive SMS notifications
// Implements the Observer interface
type SmsClient struct {
	id string // Client's phone number
}

// updateValue implements the SMS notification logic when an item becomes available
func (s *SmsClient) updateValue(itemName string) {
	fmt.Printf("Sending SMS - %s is now available for client %s\n", itemName, s.id)
}

// getId returns the client's phone number
func (s SmsClient) getId() string {
	return s.id
}

func main() {
	// Example of Observer pattern usage
	item := NewItem("RTX 5090")
	emailClient := &EmailClient{id: "test@test.com"}
	secondEmailClient := &EmailClient{id: "test2@test.com"}
	smsClient := &SmsClient{id: "+525555555555"}
	// Register two clients to receive notifications
	item.Register(emailClient)
	item.Register(secondEmailClient)
	// Register a client to receive SMS notifications
	item.Register(smsClient)
	// Update item availability, which will notify both clients
	item.UpdateAvailable()
}
