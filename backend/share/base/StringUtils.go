package base

import (
	"fmt"
	"strings"
)

func HashCode(s string) int {
	var h = 0
	for _, c := range s {
		h = 31*h + int(c)
	}
	return h
}

const maxInt = int(^uint(0) >> 1)

// Join concatenates the elements of its first argument to create a single string. The separator
// string sep is placed between elements in the resulting string.
func Join[T fmt.Stringer](elems []T, sep string) string {
	switch len(elems) {
	case 0:
		return ""
	case 1:
		return elems[0].String()
	}

	var n int
	if len(sep) > 0 {
		if len(sep) >= maxInt/(len(elems)-1) {
			panic("strings: Join output length overflow")
		}
		n += len(sep) * (len(elems) - 1)
	}
	for _, elem := range elems {
		str := elem.String()
		if len(str) > maxInt-n {
			panic("strings: Join output length overflow")
		}
		n += len(str)
	}

	var b strings.Builder
	b.Grow(n)
	b.WriteString(elems[0].String())
	for _, s := range elems[1:] {
		b.WriteString(sep)
		b.WriteString(s.String())
	}
	return b.String()
}
