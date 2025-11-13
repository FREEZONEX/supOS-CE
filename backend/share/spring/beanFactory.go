package spring

import (
	"sync"

	"github.com/samber/do/v2"
)

var beanFactory = do.New()
var lazyCallbacks = make([]func(), 0, 128)
var beanLock sync.RWMutex

func RegisterLazy[Component any](beanProvider func() Component) {
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
}
func RegisterLazyNamed[Component any](name string, beanProvider func() Component) {
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
}
func RegisterBean[Component any](bean Component) {
	registerEventListener(bean)
	do.ProvideValue[Component](beanFactory, bean)
}

func RegisterBeanNamed[Component any](name string, bean Component) {
	registerEventListener(bean)
	do.ProvideNamedValue[Component](beanFactory, name, bean)
}

func GetBean[Component any]() Component {
	return do.MustInvoke[Component](beanFactory)
}

func GetBeanOrErr[Component any]() (Component, error) {
	return do.Invoke[Component](beanFactory)
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
