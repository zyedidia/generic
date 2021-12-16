package avl_test

import (
	"math/rand"
	"testing"

	g "github.com/zyedidia/generic"
	"github.com/zyedidia/generic/avl"
)

func checkeq[K g.Lesser[K], V comparable](cm *avl.Tree[K, V], get func(k K) (V, bool), t *testing.T) {
	cm.Iter().For(func(kv avl.KV[K, V]) {
		if ov, ok := get(kv.Key); !ok {
			t.Fatalf("key %v should exist", kv.Key)
		} else if kv.Val != ov {
			t.Fatalf("value mismatch: %v != %v", kv.Val, ov)
		}
	})
}

func TestCrossCheck(t *testing.T) {
	stdm := make(map[int]int)
	tree := avl.New[g.Int, int]()

	const nops = 1000

	for i := 0; i < nops; i++ {
		key := rand.Int()
		val := rand.Int()
		op := rand.Intn(2)

		switch op {
		case 0:
			stdm[key] = val
			tree.Add(g.Int(key), val)
		case 1:
			var del int
			for k := range stdm {
				del = k
				break
			}
			delete(stdm, del)
			tree.Remove(g.Int(del))
		}

		checkeq(tree, func(k g.Int) (int, bool) {
			v, ok := stdm[int(k)]
			return v, ok
		}, t)
	}
}
