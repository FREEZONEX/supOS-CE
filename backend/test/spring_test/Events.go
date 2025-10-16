package spring_test

import (
	"time"
)

// Event represents a generic event.
type Event struct {
	Type      string
	Timestamp time.Time
}

func (e Event) GetType() string {
	return e.Type
}
func (e Event) GetTimestamp() time.Time {
	return e.Timestamp
}

// UserCreatedEvent represents a user created event.
type UserCreatedEvent struct {
	Event
	UserID   string
	Username string
	Email    string
}

// OrderCreatedEvent represents an order created event.
type OrderCreatedEvent struct {
	Event
	OrderID string
	UserID  string
	Amount  float64
}

// PaymentProcessedEvent represents a payment processed event.
type PaymentProcessedEvent struct {
	Event
	PaymentID string
	OrderID   string
	Amount    float64
	Status    string
}
