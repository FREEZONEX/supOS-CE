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
func (s *StringBuilder) AppendI(n int) *StringBuilder {
	s.WriteString(strconv.Itoa(n))
	return s
}
