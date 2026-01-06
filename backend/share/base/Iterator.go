package base

type Iterator[T any] interface {
	HasNext() bool
	Next() T
}

type EntryIterator[K any, V any] interface {
	HasNext() bool
	Next() (K, V)
}
type RemoveAbleIterator[T any] interface {
	HasNext() bool
	Next() T
	Remove()
}

type EntryRemoveAbleIterator[K any, V any] interface {
	HasNext() bool
	Next() (K, V)
	Remove()
}
type _emptyItr[T any] struct {
}

func (e _emptyItr[T]) HasNext() bool {
	return false
}
func (e _emptyItr[T]) Next() (rs T) {
	return rs
}
func EmptyIterator[T any]() Iterator[T] {
	return _emptyItr[T]{}
}
func EmptyEntryIterator[K, V any]() EntryIterator[K, V] {
	return _emptyEntryItr[K, V]{}
}

type _emptyEntryItr[K, V any] struct {
}

func (e _emptyEntryItr[K, V]) HasNext() bool {
	return false
}
func (e _emptyEntryItr[K, V]) Next() (k K, v V) {
	return k, v
}

type arrayIterator[T any] struct {
	array []T
	i     int
}

func (a *arrayIterator[T]) HasNext() bool {
	return a.i < len(a.array)
}
func (a *arrayIterator[T]) Next() (rs T) {
	rs = a.array[a.i]
	a.i++
	return rs
}
func ArrayIterator[T any](a []T) Iterator[T] {
	return &arrayIterator[T]{array: a}
}
