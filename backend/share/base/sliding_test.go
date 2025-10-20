package base

import (
	"reflect"
	"testing"
)

func TestSliding(t *testing.T) {
	arr := []int{1, 2, 3, 4, 5}
	args := [][]int{{3, 1}, {3, 2}, {3, 3}}
	expects := [][][]int{{{1, 2, 3}, {2, 3, 4}, {3, 4, 5}},
		{{1, 2, 3}, {3, 4, 5}},
		{{1, 2, 3}, {4, 5}},
	}
	for i, arg := range args {
		size, step := arg[0], arg[1]
		rs := Sliding(arr, size, step)
		expect := expects[i]
		if reflect.DeepEqual(rs, expect) {
			t.Logf("sliding(%d,%d) = %+v", size, step, rs)
		} else {
			t.Fatalf("sliding(%d,%d): expect: %+v, but got: %+v", size, step, expect, rs)
		}
	}

	arr = []int{1, 2, 3, 4, 5, 6, 7, 8}
	rs := Sliding(arr, 3, 2)
	expect := [][]int{{1, 2, 3}, {3, 4, 5}, {5, 6, 7}, {7, 8}}
	if reflect.DeepEqual(rs, expect) {
		t.Logf("[1~8].sliding(3,3) = %+v", rs)
	} else {
		t.Fatalf("[1~8].sliding(3,3): expect: %+v, but got: %+v", expect, rs)
	}
}
