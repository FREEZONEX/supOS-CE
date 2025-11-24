package base

import (
	"fmt"
	"strings"
)

type LinkedHashMap[K comparable, V any] struct {
	hashMap map[K]*hNode[K, V]
	head    *hNode[K, V]
	tail    *hNode[K, V]
}
type hNode[K comparable, V any] struct {
	key   K
	value V
	prev  *hNode[K, V]
	next  *hNode[K, V]
}

func (n *hNode[K, V]) String() string {
	return fmt.Sprint(n.value)
}
func NewLinkedHashMap[K comparable, V any]() *LinkedHashMap[K, V] {
	return &LinkedHashMap[K, V]{hashMap: make(map[K]*hNode[K, V]), head: &hNode[K, V]{}}
}
func (h *LinkedHashMap[K, V]) IsEmpty() bool {
	return len(h.hashMap) == 0
}
func (h *LinkedHashMap[K, V]) Size() int {
	return len(h.hashMap)
}
func (h *LinkedHashMap[K, V]) Put(k K, v V) {
	if old, has := h.hashMap[k]; !has {
		node := &hNode[K, V]{key: k, value: v}
		if h.tail != nil {
			h.tail.next = node
			node.prev = h.tail
		} else {
			h.head.next = node
			node.prev = h.head
		}
		h.tail = node
		h.hashMap[k] = node
	} else {
		old.value = v
	}
}
func (h *LinkedHashMap[K, V]) Get(k K) (v V) {
	val, has := h.hashMap[k]
	if has {
		v = val.value
	}
	return v
}
func (h *LinkedHashMap[K, V]) Remove(k K) (oldValue V) {
	if old, has := h.hashMap[k]; has {
		delete(h.hashMap, k)
		oldValue = old.value
		old.prev.next = old.next
		if old.next != nil {
			old.next.prev = old.prev
		}
		old.prev, old.next = nil, nil
	}
	return oldValue
}
func (h *LinkedHashMap[K, V]) Iterator() EntryRemoveAbleIterator[K, V] {
	return &linkMapItr[K, V]{lm: h, p: h.head.next}
}
func (h *LinkedHashMap[K, V]) Range(visit func(K, V)) {
	for p := h.head; p.next != nil; p = p.next {
		visit(p.next.key, p.next.value)
	}
}
func (h *LinkedHashMap[K, V]) Keys() []K {
	vs := make([]K, 0, h.Size())
	for p := h.head; p.next != nil; p = p.next {
		vs = append(vs, p.next.key)
	}
	return vs
}
func (h *LinkedHashMap[K, V]) Values() []V {
	vs := make([]V, 0, h.Size())
	for p := h.head; p.next != nil; p = p.next {
		vs = append(vs, p.next.value)
	}
	return vs
}
func (h *LinkedHashMap[K, V]) String() string {
	sz := h.Size()
	if sz == 0 {
		return "{}"
	}
	b := &strings.Builder{}
	b.Grow(sz * 32)
	b.WriteByte('{')
	for p := h.head; p.next != nil; p = p.next {
		k, v := p.next.key, p.next.value
		b.WriteString(fmt.Sprintf("%v:%+v,", k, v))
	}
	str := s2b(b.String())
	str[len(str)-1] = '}'
	return b2s(str)
}

type linkMapItr[K comparable, V any] struct {
	lm     *LinkedHashMap[K, V]
	p      *hNode[K, V]
	curKey K
}

func (lmItr *linkMapItr[K, V]) HasNext() bool {
	return lmItr.p != nil
}
func (lmItr *linkMapItr[K, V]) Next() (K, V) {
	p := lmItr.p
	lmItr.curKey = p.key
	lmItr.p = p.next
	return p.key, p.value
}
func (lmItr *linkMapItr[K, V]) Remove() {
	lmItr.lm.Remove(lmItr.curKey)
}
