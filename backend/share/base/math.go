package base

import (
	"errors"
	"math"
)

func Abs(n int) int {
	if n == math.MinInt {
		return 0
	} else if n < 0 {
		return -n
	}
	return n
}

var ArithmeticOverflow = errors.New("overflow")

func AddExact[N ~int | ~int8 | ~int16 | ~int32 | ~int64 |
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr](x, y N) (rs N, err error) {
	rs = x + y
	if (x^rs)&(y^rs) < 0 {
		err = ArithmeticOverflow
	}
	return rs, err
}
func ToIntExact(n int64) bool {
	i := int32(n)
	N := int64(i)
	return N == n
}
