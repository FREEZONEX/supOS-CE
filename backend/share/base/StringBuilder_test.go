package base

import "testing"

func TestAppendAndSetLast(t *testing.T) {
	bs := &StringBuilder{}
	bs.Append("{ )")
	bs.SetLast('}')

	t.Log(bs.String())
}
