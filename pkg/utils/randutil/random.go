package randutil

import (
	"math/rand"
	"time"
)

var (
	r *rand.Rand
)

func init() {
	r = rand.New(rand.NewSource(time.Now().Unix()))
}

func Intn(n int) int {
	if r == nil {
		return -1
	}
	return r.Intn(n)
}
