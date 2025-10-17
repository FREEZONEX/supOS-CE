package base

func Sum[E int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | float32 | float64](arr []E) (rs E) {
	for _, v := range arr {
		rs += v
	}
	return rs
}
