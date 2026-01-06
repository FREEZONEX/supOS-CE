package base

func P2v[T int | int8 | int16 | int32 | int64 | float32 | float64 | string | bool | byte](p *T) (rs T) {
	if p != nil {
		rs = *p
	}
	return
}
func P2vWithDefault[T int | int8 | int16 | int32 | int64 | float32 | float64 | string | bool | byte](p *T, def T) T {
	if p != nil {
		return *p
	}
	return def
}
func V2p[T int | int8 | int16 | int32 | int64 | float32 | float64 | string | bool | byte](p T) (rs *T) {
	return &p
}
