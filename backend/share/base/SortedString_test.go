package base

import (
	"strings"
	"testing"
)

func TestSortedString_Put(t *testing.T) {
	s := SortedString("")
	kvs := [][2]string{{"city", "HangZhou"}, {"am", "09"}, {"biz", "1"}, {"name", "LiLei"}, {"00", "First"}}
	for _, e := range kvs {
		k, v := e[0], e[1]
		s.Put(k, v)
	}
	size := s.Size()
	if size != len(kvs) {
		t.Fatalf("expect Size:%d, but got: %d", len(kvs), size)
	}
	for _, e := range kvs {
		k, v := e[0], e[1]
		vGet, _ := s.Get(k)
		t.Logf("get(%s)=%v\n", k, vGet)
		if vGet != v {
			t.Fatalf("expect %s, but got: %s", v, vGet)
		}
	}
	{
		entries := s.Entries()
		t.Log("1 entries:", entries)
	}
	{
		k := kvs[0][0]
		v := strings.ToUpper(kvs[0][1])
		s.Put(k, v)
		vGet, _ := s.Get(k)
		if vGet != v {
			t.Fatalf("expect %s, but got: %s", v, vGet)
		}
		t.Log("AfterPut 1 entries:", s.Entries())
		v += "_BJ"
		s.Put(k, v)
		vGet, _ = s.Get(k)
		if vGet != v {
			t.Fatalf("expect %s, but got: %s", v, vGet)
		}
		t.Log("AfterPut 2 entries:", s.Entries())
		v = kvs[0][1]
		s.Put(k, v)
		vGet, _ = s.Get(k)
		if vGet != v {
			t.Fatalf("expect %s, but got: %s", v, vGet)
		}
		t.Log("AfterPut 3 entries:", s.Entries())
	}
	t.Log("Final entries:", s.Entries()) // Output: [[am 09] [city HangZhou]]
	v1, _ := s.Get("name")
	t.Logf("baseGet(%s)=%v\n", "name", v1)
	v2, _ := s.Get("none")
	t.Logf("baseGet(%s)=%v\n", "none", v2)

	s.Put("name", "")
	if name, _ := s.Get("name"); name != "" {
		t.Fatalf("expect %s, but got: %s", "", name)
	}
	if !s.ContainsKey("name") {
		t.Fatalf("expect %s, but got: %s", "ContainsKey('name')", "false")
	}
}
