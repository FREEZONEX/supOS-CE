package base

func Filter[T any](arr []T, match func(e T) bool) (rs []T) {
	if sz := len(arr); sz > 0 {
		rs = make([]T, 0, sz)
		for _, a := range arr {
			if match(a) {
				rs = append(rs, a)
			}
		}
	}
	return rs
}
func FilterAndMap[T, N any](arr []T, mf func(e T) (v N, ok bool)) (rs []N) {
	if sz := len(arr); sz > 0 {
		rs = make([]N, 0, sz)
		for _, a := range arr {
			if v, match := mf(a); match {
				rs = append(rs, v)
			}
		}
	}
	return rs
}
func FilterAndFlatMap[T, N any](arr []T, mf func(e T) (vs []N, ok bool)) (rs []N) {
	if sz := len(arr); sz > 0 {
		rs = make([]N, 0, sz)
		for _, a := range arr {
			if vs, match := mf(a); match {
				rs = append(rs, vs...)
			}
		}
	}
	return rs
}
func FilterSet[T comparable](arr []T, match func(e T) bool) (rs *HashSet[T]) {
	if sz := len(arr); sz > 0 {
		rs = NewEmptySet[T](sz)
		for _, a := range arr {
			if match(a) {
				rs.Add(a)
			}
		}
	}
	return rs
}
func IndexOf[T comparable](arr []T, k T) int {
	for i, e := range arr {
		if e == k {
			return i
		}
	}
	return -1
}
