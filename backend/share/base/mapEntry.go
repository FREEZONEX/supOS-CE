package base

type Entry[K any, V any] interface {
	GetKey() K
	GetValue() V
	SetValue(V) V
}
type SimpleImmutableEntry[K any, V any] struct {
	key   K
	value V
}

func NewSimpleImmutableEntry[K any, V any](k K, v V) SimpleImmutableEntry[K, V] {
	return SimpleImmutableEntry[K, V]{key: k, value: v}
}
func (s SimpleImmutableEntry[K, V]) GetKey() K {
	return s.key
}
func (s SimpleImmutableEntry[K, V]) GetValue() V {
	return s.value
}
func (s SimpleImmutableEntry[K, V]) SetValue(V) V {
	return s.value
}
