package base

import "fmt"

type HashSet[T comparable] struct {
	arr []T
	set map[T]int
}

func NewSet[T comparable](data []T) *HashSet[T] {
	var set = make(map[T]int, min(len(data), 4))
	initData := make([]T, 0, len(data))
	initData = addAll(initData, data, set)
	return &HashSet[T]{arr: initData, set: set}
}
func addAll[T comparable](initData, data []T, set map[T]int) []T {
	if len(data) > 0 {
		for i, d := range data {
			if _, has := set[d]; !has {
				set[d] = i
				initData = append(initData, d)
			}
		}
	}
	return initData
}
func NewEmptySet[T comparable](initCap int) *HashSet[T] {
	return &HashSet[T]{set: make(map[T]int, initCap)}
}
func (si *HashSet[T]) Size() int {
	if si == nil {
		return -1
	}
	return len(si.arr)
}
func (si *HashSet[T]) Clear() {
	if len(si.set) > 0 {
		for k := range si.set {
			delete(si.set, k)
		}
	}
	si.set = nil
	if len(si.arr) > 0 {
		si.arr = si.arr[:0]
	}
}
func (si *HashSet[T]) IsEmpty() bool {
	return si == nil || len(si.arr) == 0
}
func (si *HashSet[T]) Contains(n T) (has bool) {
	if si != nil {
		m := si.set
		if m != nil {
			_, has = m[n]
		}
	}
	return has
}

func (si *HashSet[T]) Add(n T) bool {
	m := si.set
	if m == nil {
		si.set = make(map[T]int)
		m = si.set
	}
	if _, has := m[n]; !has {
		m[n] = len(si.arr)
		si.arr = append(si.arr, n)
		return true
	}
	return false
}
func (si *HashSet[T]) AddAll(arr []T) {
	si.arr = addAll(si.arr, arr, si.set)
}
func (si *HashSet[T]) Remove(a T) bool {
	i, has := si.set[a]
	if has {
		delete(si.set, a)
		si.arr = append(si.arr[:i], si.arr[i+1:]...)
	}
	return has
}
func (si *HashSet[T]) RemoveAll(arr []T) {
	for _, a := range arr {
		delete(si.set, a)
	}
	si.arr = make([]T, 0, len(si.set))
	for _, a := range arr {
		si.arr = append(si.arr, a)
	}
}
func (si *HashSet[T]) AddCloneAndToList(n T) []T {
	if si.Contains(n) {
		return si.arr
	} else {
		return append(si.arr, n)
	}
}
func (si *HashSet[T]) Values() []T {
	return si.arr
}
func (si *HashSet[T]) Diff(n *HashSet[T]) *HashSet[T] {
	return NewSet[T](si.DiffArray(n))
}
func (si *HashSet[T]) DiffArray(n *HashSet[T]) []T {
	// Create a slice to store the different elements
	diff := make([]T, 0, si.Size())
	for _, num := range si.arr {
		// Check if the element is present in b
		if _, ok := n.set[num]; !ok {
			diff = append(diff, num)
		}
	}
	return diff
}
func (si *HashSet[T]) String() string {
	if si == nil {
		return "[]"
	}
	return fmt.Sprintf("%+v", si.arr)
}
