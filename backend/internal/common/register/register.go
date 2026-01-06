package register

import "sync"

// Registry defines a generic and concurrent-safe registry.
// It uses a sync.Map to store instances of any type, keyed by a string.
// This provides a centralized and safe way to manage shared instances across different parts of an application.
type Registry[T any] struct {
	instances sync.Map
}

// NewRegistry creates a new registry.
func NewRegistry[T any]() *Registry[T] {
	return &Registry[T]{}
}

// Register adds a new instance to the registry with a specified name.
// If the name already exists, it will be overwritten.
func (r *Registry[T]) Register(name string, instance T) {
	r.instances.Store(name, instance)
}

// Get retrieves an instance from the registry by its name.
// It returns the instance and a boolean indicating if the instance was found.
func (r *Registry[T]) Get(name string) (T, bool) {
	if val, ok := r.instances.Load(name); ok {
		return val.(T), true
	}
	var zero T
	return zero, false
}

// GetAll retrieves all instances from the registry as a map.
func (r *Registry[T]) GetAll() map[string]T {
	all := make(map[string]T)
	r.instances.Range(func(key, value any) bool {
		all[key.(string)] = value.(T)
		return true
	})
	return all
}

// Remove deletes an instance from the registry by its name.
func (r *Registry[T]) Remove(name string) {
	r.instances.Delete(name)
}

// InstanceScanner defines an interface for scanning and registering instances.
type InstanceScanner[T any] interface {
	// ScanAndRegister scans for instances and registers them into the provided registry.
	ScanAndRegister(registry *Registry[T])
}
