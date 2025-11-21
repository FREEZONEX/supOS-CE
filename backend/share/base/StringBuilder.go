package base

import (
	"backend/share/base/buffer"
	"strconv"
)

type StringBuilder struct {
	bs buffer.ByteBuffer
}

func (s *StringBuilder) Grow(n int) *StringBuilder {
	s.bs.Grow(n)
	return s
}
func (s *StringBuilder) Append(str string) *StringBuilder {
	_, _ = s.bs.WriteString(str)
	return s
}
func (s *StringBuilder) Int(n int) *StringBuilder {
	_, _ = s.bs.WriteString(strconv.Itoa(n))
	return s
}
func (s *StringBuilder) Long(n int64) *StringBuilder {
	_, _ = s.bs.WriteString(strconv.FormatInt(n, 10))
	return s
}
func (s *StringBuilder) SetLast(d byte) *StringBuilder {
	s.bs.PutLast(d)
	return s
}
func (s *StringBuilder) String() string {
	return s.bs.String()
}
func (s *StringBuilder) Reset() {
	s.bs.Reset()
}
