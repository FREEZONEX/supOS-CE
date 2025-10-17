// Package treemap provides a generic key-sorted map.
// It uses red-black tree under the hood.
// Iterators are designed after C++.
//
// Example:
//
//	package main
//
//	import (
//	    "fmt"
//	)
//
//	func main() {
//	    tree := treemap.New[int, string]()
//	    tree.Set(1, "World")
//	    tree.Set(0, "Hello")
//	    for it := tree.Iterator(); it.HasNext(); it.Next() {
//	        fmt.Println(it.Key(), it.Value())
//	    }
//	}
//
//	// Output:
//	// 0 Hello
//	// 1 World
package base

import (
	"cmp"
)

// TreeMap is the generic red-black tree based map
type TreeMap[Key, Value any] struct {
	endNode    *node[Key, Value]
	beginNode  *node[Key, Value]
	count      int
	keyCompare func(a Key, b Key) bool
}

func (t *TreeMap[Key, Value]) Clone() *TreeMap[Key, Value] {
	rs := &TreeMap[Key, Value]{
		count:      t.count,
		keyCompare: t.keyCompare,
		beginNode:  t.beginNode.clone(),
		endNode:    t.endNode.clone(),
	}
	return rs
}

type node[Key, Value any] struct {
	right   *node[Key, Value]
	left    *node[Key, Value]
	parent  *node[Key, Value]
	isBlack bool
	key     Key
	value   Value
}

func (n *node[Key, Value]) clone() (rs *node[Key, Value]) {
	if n != nil {
		rs = &node[Key, Value]{
			isBlack: n.isBlack,
			key:     n.key,
			value:   n.value,
			left:    n.left.clone(),
			right:   n.right.clone(),
		}
	}
	return rs
}

// New creates and returns new TreeMap.
func New[Key cmp.Ordered, Value any]() *TreeMap[Key, Value] {
	endNode := &node[Key, Value]{isBlack: true}
	return &TreeMap[Key, Value]{beginNode: endNode, endNode: endNode, keyCompare: defaultKeyCompare[Key]}
}

// NewWithKeyCompare creates and returns new TreeMap with the specified key compare function.
// Parameter keyCompare is a function returning a < b.
func NewWithKeyCompare[Key, Value any](
	keyCompare func(a, b Key) bool,
) *TreeMap[Key, Value] {
	endNode := &node[Key, Value]{isBlack: true}
	return &TreeMap[Key, Value]{beginNode: endNode, endNode: endNode, keyCompare: keyCompare}
}

// Size returns total count of elements in a map.
// Complexity: O(1).
func (t *TreeMap[Key, Value]) Size() int { return t.count }
func (t *TreeMap[Key, Value]) IsEmpty() bool {
	return t.Size() == 0
}
func (t *TreeMap[Key, Value]) PutIfAbsent(key Key, value Value) (oldValue Value) {
	oldValue, has := t.Get(key)
	if !has {
		t.put(key, value)
	}
	return oldValue
}
func (t *TreeMap[Key, Value]) Put(key Key, value Value) (oldValue Value) {
	oldValue, _ = t.Get(key)
	t.put(key, value)
	return oldValue
}

// Put Set sets the value and silently overrides previous value if it exists.
// Complexity: O(log N).
func (t *TreeMap[Key, Value]) put(key Key, value Value) {
	parent := t.endNode
	current := parent.left
	less := true
	for current != nil {
		parent = current
		switch {
		case t.keyCompare(key, current.key):
			current = current.left
			less = true
		case t.keyCompare(current.key, key):
			current = current.right
			less = false
		default:
			current.value = value
			return
		}
	}
	x := &node[Key, Value]{parent: parent, value: value, key: key}
	if less {
		parent.left = x
	} else {
		parent.right = x
	}
	if t.beginNode.left != nil {
		t.beginNode = t.beginNode.left
	}
	t.insertFixup(x)
	t.count++
}
func (t *TreeMap[Key, Value]) FirstEntry() (k Key, v Value, ok bool) {
	if nd := t.beginNode; t.count > 0 && nd != nil {
		k, v, ok = nd.key, nd.value, true
	}
	return k, v, ok
}
func (t *TreeMap[Key, Value]) ForEach(visitor func(k Key, v Value)) {
	for itr := t.Iterator(); itr.HasNext(); {
		k, v := itr.Next()
		visitor(k, v)
	}
}
func (t *TreeMap[Key, Value]) ForEachValues(visitor func(v Value)) {
	for itr := t.Iterator(); itr.HasNext(); {
		_, v := itr.Next()
		visitor(v)
	}
}
func (t *TreeMap[Key, Value]) Keys() (ks []Key) {
	ks = make([]Key, 0, t.Size())
	for itr := t.Iterator(); itr.HasNext(); {
		k, _ := itr.Next()
		ks = append(ks, k)
	}
	return ks
}
func (t *TreeMap[Key, Value]) Values() (vs []Value) {
	vs = make([]Value, 0, t.Size())
	for itr := t.Iterator(); itr.HasNext(); {
		_, v := itr.Next()
		vs = append(vs, v)
	}
	return vs
}
func (t *TreeMap[Key, Value]) RemoveAll(ks []Key) {
	if len(ks) > 0 {
		for _, k := range ks {
			t.Remove(k)
		}
	}
}
func (t *TreeMap[Key, Value]) LastEntry() (k Key, v Value, ok bool) {
	if t.count > 0 {
		nd := t.endNode.left
		if nd != nil {
			nd = mostRight(nd)
			k, v, ok = nd.key, nd.value, true
		}
	}
	return k, v, ok
}

// Remove deletes the value.
// Complexity: O(log N).
func (t *TreeMap[Key, Value]) Remove(key Key) (v Value) {
	z := t.findNode(key)
	if z == nil {
		return v
	}
	if t.beginNode == z {
		if z.right != nil {
			t.beginNode = z.right
		} else {
			t.beginNode = z.parent
		}
	}
	t.count--
	v = z.value
	removeNode(t.endNode.left, z)
	return v
}
func (t *TreeMap[Key, Value]) RemoveWhen(key Key, v Value, valueEquals func(a, b Value) bool) bool {
	cur, has := t.Get(key)
	if !has || !valueEquals(cur, v) {
		return false
	}
	t.Remove(key)
	return true
}

// Get retrieves a value from a map for specified key and reports if it exists.
// Complexity: O(log N).
func (t *TreeMap[Key, Value]) Get(id Key) (Value, bool) {
	node := t.findNode(id)
	if node == nil {
		node = t.endNode
	}
	return node.value, node != t.endNode
}

// Clear clears the map.
// Complexity: O(1).
func (t *TreeMap[Key, Value]) Clear() {
	t.count = 0
	t.beginNode = t.endNode
	t.endNode.left = nil
}

// Contains checks if key exists in a map.
// Complexity: O(log N)
func (t *TreeMap[Key, Value]) Contains(id Key) bool { return t.findNode(id) != nil }

// Range returns a pair of iterators that you can use to go through all the keys in the range [from, to].
// More specifically it returns iterators pointing to lower bound and upper bound.
// Complexity: O(log N).
func (t *TreeMap[Key, Value]) Range(from, to Key) (*ForwardIterator[Key, Value], *ReverseIterator[Key, Value]) {
	return t.LowerBound(from), t.UpperBound(to)
}

// LowerBound returns an iterator pointing to the first element that is >= key and iterate ASC.
// Complexity: O(log N).
func (t *TreeMap[Key, Value]) LowerBound(key Key) *ForwardIterator[Key, Value] {
	var result *node[Key, Value]
	node := t.endNode.left
	if node == nil {
		return &ForwardIterator[Key, Value]{tree: t}
	}
	for {
		if t.keyCompare(node.key, key) {
			if node.right != nil {
				node = node.right
			} else {
				return &ForwardIterator[Key, Value]{tree: t, node: result}
			}
		} else if t.keyCompare(key, node.key) {
			result = node
			if node.left != nil {
				node = node.left
			} else {
				return &ForwardIterator[Key, Value]{tree: t, node: result}
			}
		} else {
			return &ForwardIterator[Key, Value]{tree: t, node: node}
		}
	}
}

// UpperBound returns an iterator pointing to the first element that is <= key and iterate DESC.
// Complexity: O(log N).
func (t *TreeMap[Key, Value]) UpperBound(key Key) *ReverseIterator[Key, Value] {
	var result *node[Key, Value]
	node := t.endNode.left
	if node == nil {
		return &ReverseIterator[Key, Value]{tree: t}
	}
	for {
		if t.keyCompare(node.key, key) {
			result = node
			if node.right != nil {
				node = node.right
			} else {
				return &ReverseIterator[Key, Value]{tree: t, node: result}
			}
		} else if t.keyCompare(key, node.key) {
			if node.left != nil {
				//result = node
				node = node.left
			} else {
				return &ReverseIterator[Key, Value]{tree: t, node: result}
			}
		} else {
			return &ReverseIterator[Key, Value]{tree: t, node: node}
		}
	}
}

// Iterator returns an iterator for tree map.
// It starts at the first element and goes to the one-past-the-end position.
// You can iterate a map at O(N) complexity.
// Method complexity: O(1)
func (t *TreeMap[Key, Value]) Iterator() *ForwardIterator[Key, Value] {
	return &ForwardIterator[Key, Value]{tree: t, node: t.beginNode}
}

// Reverse returns a reverse iterator for tree map.
// It starts at the last element and goes to the one-before-the-start position.
// You can iterate a map at O(N) complexity.
// Method complexity: O(log N)
func (t *TreeMap[Key, Value]) Reverse() *ReverseIterator[Key, Value] {
	node := t.endNode.left
	if node != nil {
		node = mostRight(node)
	}
	return &ReverseIterator[Key, Value]{tree: t, node: node}
}

func defaultKeyCompare[Key cmp.Ordered](
	a, b Key,
) bool {
	return a < b
}

func (t *TreeMap[Key, Value]) findNode(id Key) *node[Key, Value] {
	current := t.endNode.left
	for current != nil {
		switch {
		case t.keyCompare(id, current.key):
			current = current.left
		case t.keyCompare(current.key, id):
			current = current.right
		default:
			return current
		}
	}
	return nil
}

func mostLeft[Key, Value any](
	x *node[Key, Value],
) *node[Key, Value] {
	for x.left != nil {
		x = x.left
	}
	return x
}

func mostRight[Key, Value any](
	x *node[Key, Value],
) *node[Key, Value] {
	for x.right != nil {
		x = x.right
	}
	return x
}

func successor[Key, Value any](
	x *node[Key, Value],
) *node[Key, Value] {
	if x.right != nil {
		return mostLeft(x.right)
	}
	for x != x.parent.left {
		x = x.parent
	}
	return x.parent
}

func predecessor[Key, Value any](
	x *node[Key, Value],
) *node[Key, Value] {
	if x.left != nil {
		return mostRight(x.left)
	}
	for x.parent != nil && x != x.parent.right {
		x = x.parent
	}
	return x.parent
}

func rotateLeft[Key, Value any](
	x *node[Key, Value],
) {
	y := x.right
	x.right = y.left
	if x.right != nil {
		x.right.parent = x
	}
	y.parent = x.parent
	if x == x.parent.left {
		x.parent.left = y
	} else {
		x.parent.right = y
	}
	y.left = x
	x.parent = y
}

func rotateRight[Key, Value any](
	x *node[Key, Value],
) {
	y := x.left
	x.left = y.right
	if x.left != nil {
		x.left.parent = x
	}
	y.parent = x.parent
	if x == x.parent.left {
		x.parent.left = y
	} else {
		x.parent.right = y
	}
	y.right = x
	x.parent = y
}

func (t *TreeMap[Key, Value]) insertFixup(x *node[Key, Value]) {
	root := t.endNode.left
	x.isBlack = x == root
	for x != root && !x.parent.isBlack {
		if x.parent == x.parent.parent.left {
			y := x.parent.parent.right
			if y != nil && !y.isBlack {
				x = x.parent
				x.isBlack = true
				x = x.parent
				x.isBlack = x == root
				y.isBlack = true
			} else {
				if x != x.parent.left {
					x = x.parent
					rotateLeft(x)
				}
				x = x.parent
				x.isBlack = true
				x = x.parent
				x.isBlack = false
				rotateRight(x)
				break
			}
		} else {
			y := x.parent.parent.left
			if y != nil && !y.isBlack {
				x = x.parent
				x.isBlack = true
				x = x.parent
				x.isBlack = x == root
				y.isBlack = true
			} else {
				if x == x.parent.left {
					x = x.parent
					rotateRight(x)
				}
				x = x.parent
				x.isBlack = true
				x = x.parent
				x.isBlack = false
				rotateLeft(x)
				break
			}
		}
	}
}

// noinspection GoNilness
//
//nolint:gocyclo
func removeNode[Key, Value any](
	root, z *node[Key, Value],
) {
	var y *node[Key, Value]
	if z.left == nil || z.right == nil {
		y = z
	} else {
		y = successor(z)
	}
	var x *node[Key, Value]
	if y.left != nil {
		x = y.left
	} else {
		x = y.right
	}
	var w *node[Key, Value]
	if x != nil {
		x.parent = y.parent
	}
	if y == y.parent.left {
		y.parent.left = x
		if y != root {
			w = y.parent.right
		} else {
			root = x // w == nil
		}
	} else {
		y.parent.right = x
		w = y.parent.left
	}
	removedBlack := y.isBlack
	if y != z {
		y.parent = z.parent
		if z == z.parent.left {
			y.parent.left = y
		} else {
			y.parent.right = y
		}
		y.left = z.left
		y.left.parent = y
		y.right = z.right
		if y.right != nil {
			y.right.parent = y
		}
		y.isBlack = z.isBlack
		if root == z {
			root = y
		}
	}
	if removedBlack && root != nil {
		if x != nil {
			x.isBlack = true
		} else {
			for {
				if w != w.parent.left {
					if !w.isBlack {
						w.isBlack = true
						w.parent.isBlack = false
						rotateLeft(w.parent)
						if root == w.left {
							root = w
						}
						w = w.left.right
					}
					if (w.left == nil || w.left.isBlack) && (w.right == nil || w.right.isBlack) {
						w.isBlack = false
						x = w.parent
						if x == root || !x.isBlack {
							x.isBlack = true
							break
						}
						if x == x.parent.left {
							w = x.parent.right
						} else {
							w = x.parent.left
						}
					} else {
						if w.right == nil || w.right.isBlack {
							w.left.isBlack = true
							w.isBlack = false
							rotateRight(w)
							w = w.parent
						}
						w.isBlack = w.parent.isBlack
						w.parent.isBlack = true
						w.right.isBlack = true
						rotateLeft(w.parent)
						break
					}
				} else {
					if !w.isBlack {
						w.isBlack = true
						w.parent.isBlack = false
						rotateRight(w.parent)
						if root == w.right {
							root = w
						}
						w = w.right.left
					}
					if (w.left == nil || w.left.isBlack) && (w.right == nil || w.right.isBlack) {
						w.isBlack = false
						x = w.parent
						if !x.isBlack || x == root {
							x.isBlack = true
							break
						}
						if x == x.parent.left {
							w = x.parent.right
						} else {
							w = x.parent.left
						}
					} else {
						if w.left == nil || w.left.isBlack {
							w.right.isBlack = true
							w.isBlack = false
							rotateLeft(w)
							w = w.parent
						}
						w.isBlack = w.parent.isBlack
						w.parent.isBlack = true
						w.left.isBlack = true
						rotateRight(w.parent)
						break
					}
				}
			}
		}
	}
}

// ForwardIterator represents a position in a tree map.
// It is designed to iterate a map in a forward order.
// It can point to any position from the first element to the one-past-the-end element.
type ForwardIterator[Key, Value any] struct {
	tree *TreeMap[Key, Value]
	node *node[Key, Value]
	cur  Key
}

// HasNext reports if the iterator position is valid.
// In other words it returns true if an iterator is not at the one-past-the-end position.
func (i *ForwardIterator[Key, Value]) HasNext() bool {
	return i.node != nil && i.node != i.tree.endNode
}

// Next moves an iterator to the next element.
// It panics if it goes out of bounds.
func (i *ForwardIterator[Key, Value]) Next() (k Key, v Value) {
	if i.node == i.tree.endNode {
		panic("out of bound iteration")
	}
	k, v = i.node.key, i.node.value
	i.cur = k
	i.node = successor(i.node)
	return k, v
}
func (i *ForwardIterator[Key, Value]) Remove() {
	i.tree.Remove(i.cur)
}
func (i *ForwardIterator[Key, Value]) Clone() (c *TreeMap[Key, Value]) {
	cpItr := &ForwardIterator[Key, Value]{
		tree: i.tree,
		node: i.node,
		cur:  i.cur,
	}
	endNode := &node[Key, Value]{isBlack: true}
	t := &TreeMap[Key, Value]{
		count:      0,
		keyCompare: i.tree.keyCompare,
		beginNode:  endNode,
		endNode:    endNode,
	}
	for cpItr.HasNext() {
		k, v := cpItr.Next()
		t.Put(k, v)
	}
	return t
}

// Prev moves an iterator to the previous element.
// It panics if it goes out of bounds.
//func (i *ForwardIterator[Key, Value]) Prev() (k Key, v Value) {
//	if i.node != nil {
//		i.cur = i.node.key
//	}
//	i.node = predecessor(i.node)
//
//	if i.node == nil {
//		panic("out of bound iteration")
//	}
//	return i.node.key, i.node.value
//}

// Key returns a key at the iterator position
//func (i ForwardIterator[Key, Value]) Key() Key { return i.node.key }
//
//// Value returns a value at the iterator position
//func (i ForwardIterator[Key, Value]) Value() Value { return i.node.value }

// ReverseIterator represents a position in a tree map.
// It is designed to iterate a map in a reverse order.
// It can point to any position from the one-before-the-start element to the last element.
type ReverseIterator[Key, Value any] struct {
	tree *TreeMap[Key, Value]
	node *node[Key, Value]
	cur  Key
}

// HasNext reports if the iterator position is valid.
// In other words it returns true if an iterator is not at the one-before-the-start position.
func (i *ReverseIterator[Key, Value]) HasNext() bool {
	return i.node != nil
}

// Next moves an iterator to the next element in reverse order.
// It panics if it goes out of bounds.
func (i *ReverseIterator[Key, Value]) Next() (k Key, v Value) {
	if i.node == nil {
		panic("out of bound iteration")
	}
	k, v = i.node.key, i.node.value
	i.cur = k
	i.node = predecessor(i.node)
	return k, v
}
func (i *ReverseIterator[Key, Value]) Remove() {
	i.tree.Remove(i.cur)
}
func (i *ReverseIterator[Key, Value]) Clone() (c *TreeMap[Key, Value]) {
	cpItr := &ReverseIterator[Key, Value]{
		tree: i.tree,
		node: i.node,
		cur:  i.cur,
	}
	endNode := &node[Key, Value]{isBlack: true}
	t := &TreeMap[Key, Value]{
		count:      0,
		keyCompare: i.tree.keyCompare,
		beginNode:  endNode,
		endNode:    endNode,
	}
	for cpItr.HasNext() {
		k, v := cpItr.Next()
		t.Put(k, v)
	}
	return t
}

// Prev moves an iterator to the previous element in reverse order.
// It panics if it goes out of bounds.
//func (i *ReverseIterator[Key, Value]) Prev() (k Key, v Value) {
//	if i.node != nil {
//		i.cur = i.node.key
//		i.node = successor(i.node)
//	} else {
//		i.node = i.tree.beginNode
//		i.cur = i.node.key
//	}
//	k, v = i.node.key, i.node.value
//	if i.node == i.tree.endNode {
//		panic("out of bound iteration")
//	}
//	return k, v
//}

//// Key returns a key at the iterator position
//func (i ReverseIterator[Key, Value]) Key() Key { return i.node.key }
//
//// Value returns a value at the iterator position
//func (i ReverseIterator[Key, Value]) Value() Value { return i.node.value }
