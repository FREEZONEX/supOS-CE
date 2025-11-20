package spring

import (
	"backend/internal/common/event"
	"fmt"
	"math"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"sync"
	"unicode"

	typetostring "github.com/samber/go-type-to-string"
)

var eventHandlers = make(map[string][]eventListener)
var eventLock sync.RWMutex

type eventListener struct {
	name     string
	listener func(any) error
	order    int64
}

func (e eventListener) String() string {
	return e.name
}
func registerEventListener(obj any) {
	if obj == nil || reflect.ValueOf(obj).IsNil() {
		return
	}
	t := reflect.TypeOf(obj)
	beanTypeName := typetostring.GetReflectType(t)
	if x := strings.Index(beanTypeName, "/"); x > 0 {
		beanTypeName = strings.Replace(beanTypeName[x+1:], "/", ".", -1)
	}
	errorType := reflect.TypeOf((*error)(nil)).Elem()
	for i := 0; i < t.NumMethod(); i++ {
		method := t.Method(i)
		if strings.HasPrefix(method.Name, "OnEvent") && method.Type.NumIn() == 2 {
			eventType := typetostring.GetReflectType(method.Type.In(1))
			mName := method.Name
			order, pareErr := extractTailNumbers(mName) //方法名最后一段约定为优先级数值（数字越小越优先）
			if pareErr != nil {
				order = math.MaxInt64
			}
			eh := eventListener{
				name: beanTypeName + "." + mName,
				listener: func(event any) (er error) {
					vs := method.Func.Call([]reflect.Value{reflect.ValueOf(obj), reflect.ValueOf(event)})
					if sz := len(vs); sz > 0 {
						if vLast := vs[sz-1]; !vLast.IsNil() && vLast.Type().Implements(errorType) {
							er = vLast.Interface().(error)
						}
					}
					return er
				},
				order: order,
			}
			eventLock.Lock()
			eventHandlers[eventType] = append(eventHandlers[eventType], eh)
			eventLock.Unlock()
		}
	}
}

var noTailNumber = fmt.Errorf("no trailing numbers found")

func extractTailNumbers(s string) (int64, error) {
	numbers := ""
	for i := len(s) - 1; i >= 0; i-- {
		if unicode.IsDigit(rune(s[i])) {
			numbers = string(s[i]) + numbers
		} else {
			break
		}
	}

	if numbers == "" {
		return 0, noTailNumber
	}

	return strconv.ParseInt(numbers, 10, 64)
}
func PublishEvent(eventObj any) error {
	eventType := typetostring.GetReflectType(reflect.TypeOf(eventObj))
	eventLock.RLock()
	listeners, has := eventHandlers[eventType]
	eventLock.RUnlock()
	if !has {
		return nil
	}
	if eventStatusAware, isAware := eventObj.(event.EventStatusAware); isAware {
		TOTAL := len(listeners)
		for i, handler := range listeners {
			eventStatusAware.BeforeEvent(TOTAL, i, handler.name)
			err := handler.listener(eventObj)
			if err != nil {
				eventStatusAware.AfterEvent(TOTAL, i, handler.name, err)
				return err
			} else {
				eventStatusAware.AfterEvent(TOTAL, i, handler.name, nil)
			}
		}
	} else {
		for _, handler := range listeners {
			err := handler.listener(eventObj)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
func onRefreshBeanContext() {
	eventLock.Lock()
	for _, listeners := range eventHandlers {
		if len(listeners) > 1 {
			sort.Sort(eventListenerArray(listeners))
		}
	}
	eventLock.Unlock()
}

type eventListenerArray []eventListener

func (x eventListenerArray) Len() int           { return len(x) }
func (x eventListenerArray) Less(i, j int) bool { return x[i].order < x[j].order }
func (x eventListenerArray) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }
