package base

func Concat[T any](a, b []T) (rs []T) {
	sz1, sz2 := len(a), len(b)
	if sz1 > sz2 {
		rs = append(a, b...)
	} else if sz2 > sz1 {
		rs = append(b, a...)
	} else if sz1+sz2 == 0 {
	}
	return rs
}
func ConcatMap[K comparable, V any](a, b map[K]V) (rs map[K]V) {
	sz1, sz2 := len(a), len(b)
	if sz2 == 0 {
		return a
	} else if sz1 == 0 {
		return b
	}
	for k, v := range b {
		a[k] = v
	}
	return rs
}
