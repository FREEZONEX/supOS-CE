package base

type Comparable[T any] interface {
	CompareTo(other T) int
}
type Int int
type Int8 int8
type Int16 int16
type Int32 int32
type Int64 int64
type Uint uint
type Uint8 uint8
type Uint16 uint16
type Uint32 uint32
type Uint64 uint64
type UintPtr uintptr
type Float32 float32
type Float64 float64
type String string

func (x Int) CompareTo(y Int) int {
	if x < y {
		return -1
	}
	if x > y {
		return +1
	}
	return 0
}
func (x Int8) CompareTo(y Int8) int {
	if x < y {
		return -1
	}
	if x > y {
		return +1
	}
	return 0
}
func (x Int16) CompareTo(y Int16) int {
	if x < y {
		return -1
	}
	if x > y {
		return +1
	}
	return 0
}
func (x Int32) CompareTo(y Int32) int {
	if x < y {
		return -1
	}
	if x > y {
		return +1
	}
	return 0
}
func (x Int64) CompareTo(y Int64) int {
	if x < y {
		return -1
	}
	if x > y {
		return +1
	}
	return 0
}
func (x Uint) CompareTo(y Uint) int {
	if x < y {
		return -1
	}
	if x > y {
		return +1
	}
	return 0
}
func (x Uint8) CompareTo(y Uint8) int {
	if x < y {
		return -1
	}
	if x > y {
		return +1
	}
	return 0
}
func (x Uint16) CompareTo(y Uint16) int {
	if x < y {
		return -1
	}
	if x > y {
		return +1
	}
	return 0
}
func (x Uint32) CompareTo(y Uint32) int {
	if x < y {
		return -1
	}
	if x > y {
		return +1
	}
	return 0
}
func (x Uint64) CompareTo(y Uint64) int {
	if x < y {
		return -1
	}
	if x > y {
		return +1
	}
	return 0
}
func (x UintPtr) CompareTo(y UintPtr) int {
	if x < y {
		return -1
	}
	if x > y {
		return +1
	}
	return 0
}
func (x Float32) CompareTo(y Float32) int {
	if x < y {
		return -1
	}
	if x > y {
		return +1
	}
	return 0
}
func (x Float64) CompareTo(y Float64) int {
	if x < y {
		return -1
	}
	if x > y {
		return +1
	}
	return 0
}
func (x String) CompareTo(y String) int {
	if x < y {
		return -1
	}
	if x > y {
		return +1
	}
	return 0
}
