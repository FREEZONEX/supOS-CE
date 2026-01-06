package base

func Contains[T comparable](arr []T, k T) bool {
	if len(arr) > 0 {
		for _, a := range arr {
			if a == k {
				return true
			}
		}
	}
	return false
}
