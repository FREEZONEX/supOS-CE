package spring_test

import (
	"backend/share/spring"
	"fmt"
	"time"
)

// Logger represents a logging service.
type Logger struct {
	Level string
}

func (l *Logger) Log(message string) {
	fmt.Printf("[%s] %s\n", l.Level, message)
}

// UserService represents a user service.
type UserService struct {
	Logger *Logger
}

func (u *UserService) CreateUser(username, email string) error {
	userID := fmt.Sprintf("user-%d", time.Now().Unix())

	u.Logger.Log(fmt.Sprintf("Creating user: %s (%s)", username, email))

	// Publish user created event
	event := UserCreatedEvent{
		Event:    Event{Type: "user.created", Timestamp: time.Now()},
		UserID:   userID,
		Username: username,
		Email:    email,
	}

	return spring.PublishEvent(event)
}

// OrderService represents an order service.
type OrderService struct {
	Logger *Logger
}

func (o *OrderService) CreateOrder(userID string, amount float64) error {
	orderID := fmt.Sprintf("order-%d", time.Now().Unix())

	o.Logger.Log(fmt.Sprintf("Creating order for user: %s, amount: %.2f", userID, amount))

	// Publish order created event
	event := OrderCreatedEvent{
		Event:   Event{Type: "order.created", Timestamp: time.Now()},
		OrderID: orderID,
		UserID:  userID,
		Amount:  amount,
	}

	return spring.PublishEvent(event)
}

// PaymentService represents a payment service.
type PaymentService struct {
	Logger *Logger
}

func (p *PaymentService) ProcessPayment(orderID string, amount float64) error {
	paymentID := fmt.Sprintf("payment-%d", time.Now().Unix())

	p.Logger.Log(fmt.Sprintf("Processing payment for order: %s, amount: %.2f", orderID, amount))

	// Publish payment processed event
	event := PaymentProcessedEvent{
		Event:     Event{Type: "payment.processed", Timestamp: time.Now()},
		PaymentID: paymentID,
		OrderID:   orderID,
		Amount:    amount,
		Status:    "completed",
	}

	return spring.PublishEvent(event)
}
