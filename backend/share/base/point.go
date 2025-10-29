package base

func P2v[T int | int32 | float32 | float64 | string](p *T) (rs T) {
	if p != nil {
		rs = *p
	}
	return
}
func V2p[T int | int32 | float32 | float64 | string](p T) (rs *T) {
	return &p
}
