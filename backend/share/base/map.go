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
func MapDistinct[E any, K comparable](arr []E, op func(e E) K) (rs []K) {
	if len(arr) > 0 {
		rsMap := NewLinkedHashMap[K, int]()
		for i, v := range arr {
			k := op(v)
			rsMap.Put(k, i)
		}
		rs = rsMap.Keys()
	}
	return rs
}
func MapDistinctFree[E any, K comparable](arr []E, op func(e E) K) (rs []K) {
	if len(arr) > 0 {
		rsMap := make(map[K]int)
		for i, v := range arr {
			k := op(v)
			rsMap[k] = i
		}
		rs = make([]K, len(rsMap))
		i := 0
		for k := range rsMap {
			rs[i] = k
			i++
		}
	}
	return rs
}
func MapFilter[K comparable, V any](m map[K]V, f func(V) bool) (rs map[K]V) {
	rs = make(map[K]V, len(m))
	if len(m) > 0 {
		for k, v := range m {
			if f(v) {
				rs[k] = v
			}
		}
	}
	return rs
}
func MapKeys[K comparable, V any](m map[K]V) (rs []K) {
	rs = make([]K, len(m))
	if len(m) > 0 {
		i := 0
		for k := range m {
			rs[i] = k
			i++
		}
	}
	return rs
}
func MapValues[K comparable, V any](m map[K]V) (rs []V) {
	rs = make([]V, len(m))
	if len(m) > 0 {
		i := 0
		for _, v := range m {
			rs[i] = v
			i++
		}
	}
	return rs
}
func MapMapValues[K comparable, V, V2 any](m map[K]V, op func(e V) V2) (rs []V2) {
	rs = make([]V2, len(m))
	if len(m) > 0 {
		i := 0
		for _, v := range m {
			rs[i] = op(v)
			i++
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
func PutIfAbsent[K comparable, V any](m map[K]V, key K, val V) bool {
	if _, has := m[key]; !has {
		m[key] = val
		return true
	}
	return false
}
