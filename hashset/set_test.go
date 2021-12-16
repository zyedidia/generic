package hashset_test

import (
	"math/rand"
	"testing"

	g "github.com/zyedidia/generic"
	"github.com/zyedidia/generic/hashset"
)

func checkeq[K g.Hashable[K]](set *hashset.Set[K], get func(k K) bool, t *testing.T) {
	set.Iter().For(func(key K) {
		if !get(key) {
			t.Fatalf("value %v should be in the set", key)
		}
	})
}

func TestCrossCheck(t *testing.T) {
	stdm := make(map[int]bool)
	set := hashset.New[g.Int](1)

	const nops = 1000
	for i := 0; i < nops; i++ {
		op := rand.Intn(2)
		switch op {
		case 0:
			key := rand.Int()
			stdm[key] = true
			set.Put(g.Int(key))
		case 1:
			var del int
			for k := range stdm {
				del = k
				break
			}
			delete(stdm, del)
			set.Remove(g.Int(del))
		}

		checkeq(set, func(k g.Int) bool {
			_, ok := stdm[int(k)]
			return ok
		}, t)
	}
}
