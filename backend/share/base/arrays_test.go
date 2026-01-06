package base

import (
	"sort"
	"strings"
	"testing"
)

func TestBinarySearchArray(t *testing.T) {
	arr := []string{"1/2", "1/22", "a/b"}
	sort.Strings(arr)
	ks := []string{"1/2/33", "1/22/333", "a/b/C2", "1/", "11/", "DD1/2"}
	for _, k := range ks {
		i := BinarySearchArray(arr, k, func(a, b string) int {
			if strings.HasPrefix(b, a) {
				return 0
			} else {
				return strings.Compare(a, b)
			}
		})
		t.Logf("BinarySearchArray(%v) is %v, arr=%v\n", k, i, arr)
	}
}
