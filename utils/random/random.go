package random

import "math/rand"

func Number(min int64, max int64) int64 {
	return rand.Int63n(max+1-min) + min
}
