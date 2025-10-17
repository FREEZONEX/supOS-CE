package base

import (
	"fmt"
	"math"
	"testing"
)

func TestTreeMap_FirstEntry(t *testing.T) {
	tree := New[int, string]()
	if fwi := tree.LowerBound(0); fwi.HasNext() {
		t.Fatal("LowerBound 空树不能有下一步")
	}
	if fwi := tree.UpperBound(math.MaxInt); fwi.HasNext() {
		t.Fatal("UpperBound 空树不能有下一步")
	}
	tree.Put(1, "World")
	tree.Put(0, "Hello")
	tree.Put(7, "b")
	tree.Put(3, "a")
	tree.Put(9, "c")

	k, v, ok := tree.FirstEntry()
	t.Log("FirstKey:", k, v, ok)

	k, v, ok = tree.LastEntry()
	t.Log("LastEntry:", k, v, ok)
}
func TestTreeMap_Iterator(t *testing.T) {
	tree := New[int, string]()
	if fwi := tree.LowerBound(0); fwi.HasNext() {
		t.Fatal("LowerBound 空树不能有下一步")
	}
	if fwi := tree.UpperBound(math.MaxInt); fwi.HasNext() {
		t.Fatal("UpperBound 空树不能有下一步")
	}
	tree.Put(1, "World")
	fmt.Println("Iterator only one.")
	for it := tree.Iterator(); it.HasNext(); {
		k, v := it.Next()
		fmt.Println(k, v)
	}
	fmt.Println("search >=Max when only one element")
	if fwi := tree.LowerBound(math.MaxInt); fwi.HasNext() {
		t.Fatal(">=Max　不能有值")
	}
	if fwi := tree.UpperBound(math.MinInt); fwi.HasNext() {
		t.Fatal("<=Min　不能有值")
	}
	//
	tree.Put(0, "Hello")
	tree.Put(7, "b")
	tree.Put(3, "a")
	tree.Put(9, "c")
	assertItr(t, "Iterator All:", tree.Iterator(), []int{0, 1, 3, 7, 9})
	assertItr(t, "Iterator Reverse:", tree.Reverse(), []int{9, 7, 3, 1, 0})
	assertTreeMap(t, tree, true, math.MinInt, []int{0, 1, 3, 7, 9})
	assertTreeMap(t, tree, true, 2, []int{3, 7, 9})
	assertTreeMap(t, tree, true, 5, []int{7, 9})

	assertTreeMap(t, tree, false, 0, []int{0})
	assertTreeMap(t, tree, false, 1, []int{1, 0})
	assertTreeMap(t, tree, false, 5, []int{3, 1, 0})
	assertTreeMap(t, tree, false, math.MaxInt, []int{9, 7, 3, 1, 0})
	assertTreeMap(t, tree, false, math.MinInt, []int{})
}
func assertTreeMap(t *testing.T, tm *TreeMap[int, string], biggerThan bool, key int, expect []int) {
	var itr EntryIterator[int, string]
	tag := ""
	if biggerThan {
		itr = tm.LowerBound(key)
		tag = " >= "
	} else {
		itr = tm.UpperBound(key)
		tag = " <= "
	}
	assertItr(t, fmt.Sprintf("%s %d", tag, key), itr, expect)
}
func assertItr(t *testing.T, tag string, itr EntryIterator[int, string], expect []int) {
	rs := make([]int, 0, 16)
	for itr.HasNext() {
		k, _ := itr.Next()
		rs = append(rs, k)
	}
	if !Equals(rs, expect) {
		t.Fatalf("%s: Expect %+v but got: %+v\n", tag, expect, rs)
	} else {
		t.Logf("%s: got: %+v\n", tag, rs)
	}
}
func TestTreeMap_IteratorRemove(t *testing.T) {
	tree := New[int, string]()
	tree.Put(1, "World")
	tree.Put(0, "Hello")
	tree.Put(7, "b")
	tree.Put(3, "a")
	tree.Put(9, "c")

	{
		k, v, ok := tree.LastEntry()
		t.Log("LastEntry: ", k, v, ok)
	}
	{
		itr := tree.Iterator()
		for itr.HasNext() {
			k, v := itr.Next()
			t.Log(k, "=", v)
			if k == 1 {
				itr.Remove()
			}
		}
	}
	fmt.Println()
	t.Log("After remove World: ")
	{
		itr := tree.Iterator()
		for itr.HasNext() {
			k, v := itr.Next()
			t.Log(k, "=", v)
		}
	}
}
