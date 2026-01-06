package spring_test

import (
	"fmt"
)

func (u *UserService) OnEventUserCreatedEvent_2000(userData UserCreatedEvent) error {
	u.Logger.Log(fmt.Sprintf("UserEventHandler: Processing UserCreatedEvent %s", userData.GetType()))

	u.Logger.Log(fmt.Sprintf("User created: ID=%s, Username=%s, Email=%s",
		userData.UserID, userData.Username, userData.Email))

	return nil
}
func (u *UserService) OnEventOrderCreatedEvent(orderData OrderCreatedEvent) error {
	u.Logger.Log(fmt.Sprintf("UserEventHandler: Processing OrderCreatedEvent %s", orderData.GetType()))
	u.Logger.Log(fmt.Sprintf("Order created: ID=%s, UserID=%s, Amount=%.2f",
		orderData.OrderID, orderData.UserID, orderData.Amount))

	return nil
}

// OrderEventHandler handles order-related events.

// PaymentEventHandler handles payment-related events.
type PaymentEventHandler struct {
	Logger *Logger
}

func (h *PaymentEventHandler) OnEventPaymentProcessedEvent(paymentData PaymentProcessedEvent) error {
	h.Logger.Log(fmt.Sprintf("PaymentEventHandler: Processing event %s", paymentData.GetType()))

	h.Logger.Log(fmt.Sprintf("Payment processed: ID=%s, OrderID=%s, Amount=%.2f, Status=%s",
		paymentData.PaymentID, paymentData.OrderID, paymentData.Amount, paymentData.Status))

	return nil
}
func (u *PaymentEventHandler) OnEventUserCreatedEvent_100(userData UserCreatedEvent) error {
	u.Logger.Log(fmt.Sprintf("PaymentEventHandler: Processing UserCreatedEvent %s", userData.GetType()))
	return nil
}
func (h *PaymentEventHandler) GetEventType() string {
	return "payment.processed"
}
