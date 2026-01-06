package base

type Either[A any, B any] struct {
	left   A
	right  B
	isLeft bool
}

func EitherLeft[A any, B any](a A) Either[A, B] {
	return Either[A, B]{left: a, isLeft: false}
}
func EitherRight[A any, B any](b B) Either[A, B] {
	return Either[A, B]{right: b, isLeft: true}
}
func (e Either[A, B]) Get() (a A, b B, isLeft bool) {
	return e.left, e.right, e.isLeft
}
