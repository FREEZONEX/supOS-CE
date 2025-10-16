package spring

import (
	"reflect"
	"strings"
	"sync"

	typetostring "github.com/samber/go-type-to-string"
)

var eventHandlers = make(map[string][]func(any))
var eventLock sync.RWMutex

func registerEventListener(obj any) {
	t := reflect.TypeOf(obj)
	for i := 0; i < t.NumMethod(); i++ {
		method := t.Method(i)
		if strings.HasPrefix(method.Name, "OnEvent") && method.Type.NumIn() == 2 {
			eventType := typetostring.GetReflectType(method.Type.In(1))
			eventLock.Lock()
			eventHandlers[eventType] = append(eventHandlers[eventType], func(event any) {
				method.Func.Call([]reflect.Value{reflect.ValueOf(obj), reflect.ValueOf(event)})
			})
			eventLock.Unlock()
		}
	}
}
func PublishEvent(event any) error {
	eventType := typetostring.GetReflectType(reflect.TypeOf(event))
	eventLock.RLock()
	if hs, has := eventHandlers[eventType]; has {
		eventLock.RUnlock()
		for _, handler := range hs {
			handler(event)
		}
	} else {
		eventLock.RUnlock()
	}
	return nil
}
