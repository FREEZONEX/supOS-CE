package base

func SlidingOne[T any](a []T, size int) [][]T {
	return Sliding(a, size, 1)
}

func Sliding[T any](array []T, size, step int) [][]T {
	if step < 1 {
		step = 1
	}
	if size < 1 {
		size = 1
	}
	N := len(array)
	if size >= N || step >= N {
		return [][]T{array}
	}
	rs := make([][]T, 0, len(array))
	for i := 0; i < N; i += step {
		slice := make([]T, 0, size)
		for j := 0; j < size; j++ {
			index := i + j
			if index < N {
				slice = append(slice, array[index])
			} else {
				break
			}
		}
		rs = append(rs, slice)
		if i+size >= N {
			break
		}
	}
	return rs
}
