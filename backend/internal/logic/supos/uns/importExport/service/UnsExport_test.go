package service

import "testing"

func Test_getLayAndIdsInner(t *testing.T) {
	layRec, ids := getLayAndIdsInner([]int64{3}, []int64{4, 5}, []string{"1/2/3", "1/2/3/4", "1/2/5"})
	t.Log("layRec = ", layRec)
	t.Log("ids = ", ids)
	//layRec =  [1/2/3]
	//ids =  [5]
}
