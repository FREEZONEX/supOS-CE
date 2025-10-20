package base

func FoldLeft[A any, B any](arr []A, init B, op func(b B, a A) B) (rs B) {
	rs = init
	for _, v := range arr {
		rs = op(rs, v)
	}
	return rs
}
func FoldRight[A any, B any](arr []A, init B, op func(a A, b B) B) (rs B) {
	rs = init
	for _, v := range arr {
		rs = op(v, rs)
	}
	return rs
}
