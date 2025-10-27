package base

import "cmp"

func SanYuan[T any](exp bool, a, b T) T {
	if exp {
		return a
	} else {
		return b
	}
}
func SanF[T any](exp bool, a, b func() T) T {
	if exp {
		return a()
	} else {
		return b()
	}
}
func SanA[T any](exp bool, a T, b func() T) T {
	if exp {
		return a
	} else {
		return b()
	}
}
func GetOrElse[T cmp.Ordered](p *T, val T) T {
	if p != nil {
		return *p
	} else {
		return val
	}
}
