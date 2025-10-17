package base

import "testing"

func TestLinkedHashMap(t *testing.T) {
	hm := NewLinkedHashMap[int, string]()
	t.Log(hm.Remove(1))
	hm.Put(1, "11")

	for it := hm.Iterator(); it.HasNext(); {
		k, v := it.Next()
		t.Logf("see: %v=%v\n", k, v)
		it.Remove()
	}
	t.Log(hm.Remove(1))

	hm.Put(9, "999")
	hm.Put(5, "555")
	hm.Put(7, "17")
	hm.Put(4, "四儿子")
	hm.Put(2, "x2")
	t.Log(hm.String(), hm.hashMap)

	itr := hm.Iterator()
	for itr.HasNext() {
		k, v := itr.Next()
		t.Logf("visit: %v=%v\n", k, v)
		if k == 4 {
			itr.Remove()
		}
	}
	hm.Put(2, "x22")
	t.Log(hm.String(), hm.hashMap)
}
