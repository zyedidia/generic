package btree_test

import (
	"math/rand"
	"testing"

	g "github.com/zyedidia/generic"
	"github.com/zyedidia/generic/btree"
)

func checkeq[K g.Lesser[K], V comparable](cm *btree.Tree[K, V], get func(k K) (V, bool), t *testing.T) {
	cm.Iter().For(func(kv btree.KV[K, V]) {
		if ov, ok := get(kv.Key); !ok {
			t.Fatalf("key %v should exist", kv.Key)
		} else if kv.Val != ov {
			t.Fatalf("value mismatch: %v != %v", kv.Val, ov)
		}
	})
}

func TestCrossCheck(t *testing.T) {
	stdm := make(map[int]int)
	tree := btree.New[g.Int, int]()

	const nops = 1000

	for i := 0; i < nops; i++ {
		key := rand.Int()
		val := rand.Int()
		op := rand.Intn(1)

		switch op {
		case 0:
			stdm[key] = val
			tree.Put(g.Int(key), val)
		}

		checkeq(tree, func(k g.Int) (int, bool) {
			v, ok := stdm[int(k)]
			return v, ok
		}, t)
	}
}
