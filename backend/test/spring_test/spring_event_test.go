package spring_test

import (
	"backend/share/spring"
	"fmt"
	"reflect"
	"testing"

	typetostring "github.com/samber/go-type-to-string"
)

func TestEventName(t *testing.T) {
	vs := []any{UserCreatedEvent{}, &UserCreatedEvent{}}
	for i, v := range vs {
		eventType := typetostring.GetReflectType(reflect.TypeOf(v))
		t.Logf("eventType[%d]=%#v\n", i, eventType)
	}
}
func TestSpringEvent(t *testing.T) {

	fmt.Println("=== Event-Driven Architecture Example ===")
	fmt.Println("This example demonstrates how to build an event-driven system")
	fmt.Printf("using the do library for dependency injection.\n\n")

	// Step 1: RegisterLazy core infrastructure services
	fmt.Println("Step 1: Registering core infrastructure services")
	spring.RegisterBean[*Logger](&Logger{Level: "INFO"})

	// Step 2: RegisterLazy event handlers
	fmt.Println("Step 2: Registering event handlers")
	spring.RegisterLazy(func() *PaymentEventHandler {
		return &PaymentEventHandler{
			Logger: spring.GetBean[*Logger](),
		}
	})

	// Step 3: RegisterLazy business services that publish events
	fmt.Println("Step 3: Registering business services that publish events")
	spring.RegisterLazy(func() *UserService {
		return &UserService{
			Logger: spring.GetBean[*Logger](),
		}
	})

	spring.RegisterLazy(func() *OrderService {
		return &OrderService{
			Logger: spring.GetBean[*Logger](),
		}
	})

	spring.RegisterLazy(func() *PaymentService {
		return &PaymentService{
			Logger: spring.GetBean[*Logger](),
		}
	})

	// Step 4: RegisterLazy main application
	fmt.Println("Step 4: Registering main application")
	spring.RegisterLazy(func() *Application {
		return &Application{
			UserService:    spring.GetBean[*UserService](),
			OrderService:   spring.GetBean[*OrderService](),
			PaymentService: spring.GetBean[*PaymentService](),
			Logger:         spring.GetBean[*Logger](),
		}
	})
	app := spring.GetBean[*Application]()
	spring.RefreshBeanContext()
	app.Run()
}
