package base

func Exists[T any](arr []T, op func(e T) bool) bool {
	if len(arr) > 0 {
		for _, a := range arr {
			if op(a) {
				return true
			}
		}
	}
	return false
}
