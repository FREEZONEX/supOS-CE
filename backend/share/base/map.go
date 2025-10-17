package base

func Map[E any, M any](arr []E, op func(e E) M) (rs []M) {
	rs = make([]M, len(arr))
	if len(arr) > 0 {
		for i, v := range arr {
			rs[i] = op(v)
		}
	}
	return rs
}
func MapArrayToMap[E any, K comparable, V any](arr []E, op func(e E) (ok bool, k K, v V)) (rs map[K]V) {
	rs = make(map[K]V, len(arr))
	if len(arr) > 0 {
		for _, v := range arr {
			ok, k, v2 := op(v)
			if ok {
				rs[k] = v2
			}
		}
	}
	return rs
}
func MapMapV[K comparable, V, V2 any](m map[K]V, op func(e V) V2) (rs map[K]V2) {
	rs = make(map[K]V2, len(m))
	if len(m) > 0 {
		for k, v := range m {
			rs[k] = op(v)
		}
	}
	return rs
}
func MapMap[K1 comparable, V1 any, K2 comparable, V2 any](m map[K1]V1, op func(k K1, v V1) (K2, V2)) (rs map[K2]V2) {
	rs = make(map[K2]V2, len(m))
	if len(m) > 0 {
		for k, v := range m {
			k2, v2 := op(k, v)
			rs[k2] = v2
		}
	}
	return rs
}
func MapContainsKey[K comparable, V any](m map[K]V, key K) bool {
	_, has := m[key]
	return has
}
