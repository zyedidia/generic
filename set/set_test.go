package set_test

import (
	"math/rand"
	"testing"

	g "github.com/zyedidia/generic"
	"github.com/zyedidia/generic/set"
)

func TestCrossCheck(t *testing.T) {
	stdm := make(map[int]bool)
	set := set.New[g.Int](1)

	const nops = 1000
	for i := 0; i < nops; i++ {
		op := rand.Intn(2)
		switch op {
		case 0:
			key := rand.Int()
			stdm[key] = true
			set.Put(g.Int(key))
		case 1:
		}
	}
}
