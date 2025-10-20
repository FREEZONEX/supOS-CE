package base

import (
	"strconv"
	"strings"
)

type StringBuilder struct {
	strings.Builder
}

func (s *StringBuilder) Append(str string) *StringBuilder {
	s.WriteString(str)
	return s
}
func (s *StringBuilder) Int(n int) *StringBuilder {
	s.WriteString(strconv.Itoa(n))
	return s
}
func (s *StringBuilder) Long(n int64) *StringBuilder {
	s.WriteString(strconv.FormatInt(n, 10))
	return s
}
