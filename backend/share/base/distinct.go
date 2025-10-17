package base

func Distinct[T comparable](arr []T) (rs []T) {
	if len(arr) > 0 {
		hashMap := make(map[T]int, len(arr))
		rs = make([]T, 0, len(arr))
		for i, a := range arr {
			if _, has := hashMap[a]; !has {
				hashMap[a] = i
				rs = append(rs, a)
			}
		}
	}
	return rs
}
