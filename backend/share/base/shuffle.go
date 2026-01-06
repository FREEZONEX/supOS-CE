package base

import (
	"math/rand"
	"time"
)

func Shuffle[T any](list []T, rnd *rand.Rand) {
	r := rnd
	if r == nil {
		r = rand.New(rand.NewSource(time.Now().UnixNano()))
	}
	r.Shuffle(len(list), func(i, j int) {
		list[i], list[j] = list[j], list[i]
	})
}
