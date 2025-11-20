package spring

import (
	"fmt"
	"reflect"
	"sync"

	do "gitee.com/supos-community-edition/di/v2"
)

var beanFactory = do.New()
var lazyCallbacks = make([]func(), 0, 128)
var beanLock sync.RWMutex

func RegisterLazy[Component any](beanProvider func() Component) bool {
	serviceName := do.NameOf[Component]()
	if beanFactory.ServiceExists(serviceName) {
		return false
	}
	beanLock.Lock()
	lazyCallbacks = append(lazyCallbacks, func() {
		do.MustInvoke[Component](beanFactory)
	})
	beanLock.Unlock()
	do.Provide[Component](beanFactory, func(injector do.Injector) (Component, error) {
		bean := beanProvider()
		registerEventListener(bean)
		return bean, nil
	})
	return true
}
func RegisterLazyNamed[Component any](name string, beanProvider func() Component) bool {
	serviceName := do.NameOf[Component]()
	if beanFactory.ServiceExists(serviceName) {
		return false
	}

	beanLock.Lock()
	lazyCallbacks = append(lazyCallbacks, func() {
		do.MustInvokeNamed[Component](beanFactory, name)
	})
	beanLock.Unlock()
	do.ProvideNamed[Component](beanFactory, name, func(injector do.Injector) (Component, error) {
		bean := beanProvider()
		registerEventListener(bean)
		return bean, nil
	})
	return true
}
func RegisterBean[Component any](bean Component) bool {
	serviceName := do.NameOf[Component]()
	if beanFactory.ServiceExists(serviceName) {
		return false
	}
	registerEventListener(bean)
	do.ProvideValue[Component](beanFactory, bean)
	return true
}
func RegisterBeanNamed[Component any](name string, bean Component) bool {
	serviceName := do.NameOf[Component]()
	if beanFactory.ServiceExists(serviceName) {
		return false
	}
	registerEventListener(bean)
	do.ProvideNamedValue[Component](beanFactory, name, bean)
	return true
}

func GetBean[Component any]() Component {
	rs, er := GetBeanOrErr[Component]()
	if er != nil {
		panic(er)
	}
	return rs
}
func GetBeansOfType[Component any]() (list []Component) {
	list = make([]Component, 0, 4)
	beanFactory.InstanceForEach(func(s string, scope *do.Scope, a any) bool {
		if im, ok := a.(Component); ok && !reflect.ValueOf(a).IsNil() {
			list = append(list, im)
		}
		return true
	})
	return
}

var notFoundError = fmt.Errorf("beanNotFound")

func GetBeanOrErr[Component any]() (Component, error) {
	rs, er := do.Invoke[Component](beanFactory)
	if er != nil {
		var com Component
		t := reflect.TypeOf(&com)
		if t != nil && t.Kind() == reflect.Pointer && t.Elem().Kind() == reflect.Interface {
			list := GetBeansOfType[Component]()
			if sz := len(list); sz > 0 {
				if sz == 1 {
					er = nil
					rs = list[0]
				} else {
					er = fmt.Errorf(`too many beans of type: %s`, do.NameOf[Component]())
				}
			}
		}
	}
	var a any = rs
	if reflect.ValueOf(a).IsNil() {
		er = notFoundError
	}
	return rs, er
}

// RefreshBeanContext 项目初始化完成之后调用，相当于 spring onApplicationRefreshed
func RefreshBeanContext() {
	beanLock.Lock()
	for _, call := range lazyCallbacks {
		call()
	}
	lazyCallbacks = lazyCallbacks[:]
	beanLock.Unlock()

	onRefreshBeanContext()
}
