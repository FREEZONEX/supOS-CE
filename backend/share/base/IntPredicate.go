package base

type IntPredicate func(n int) bool

func (I IntPredicate) And(other IntPredicate) IntPredicate {
	return func(n int) bool {
		return I(n) && other(n)
	}
}
func (I IntPredicate) Or(other IntPredicate) IntPredicate {
	return func(n int) bool {
		return I(n) || other(n)
	}
}
func (I IntPredicate) Negate() IntPredicate {
	return func(n int) bool {
		return !I(n)
	}
}
