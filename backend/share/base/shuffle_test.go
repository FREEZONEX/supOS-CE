package base

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestShuffle(t *testing.T) {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	list := []interface{}{1, 2, 3, 4, 5}
	Shuffle(list, rnd)
	fmt.Println(list)

	ns := []int{0x10, 0x20, 01}
	for _, n := range ns {
		t.Logf("%v >> 4 = %d", n, n>>4)
	}
}
