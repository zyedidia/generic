package hashmap

import (
	"math/rand"
	"testing"

	g "github.com/zyedidia/generic"
)

func checkeq[K g.Hashable[K], V comparable](cm *Map[K, V], get func(k K) (V, bool), t *testing.T) {
	cm.Iter().For(func(kv KV[K, V]) {
		if ov, ok := get(kv.Key); !ok {
			t.Fatalf("key %v should exist", kv.Key)
		} else if kv.Val != ov {
			t.Fatalf("value mismatch: %v != %v", kv.Val, ov)
		}
	})
}

func TestCrossCheck(t *testing.T) {
	stdm := make(map[g.Uint64]uint32)
	cowm := NewMap[g.Uint64, uint32](1)

	const nops = 1000

	for i := 0; i < nops; i++ {
		key := g.Uint64(rand.Uint64())
		val := rand.Uint32()
		op := rand.Intn(2)

		switch op {
		case 0:
			stdm[key] = val
			cowm.Put(key, val)
		case 1:
			var del g.Uint64
			for k := range stdm {
				del = k
				break
			}
			delete(stdm, del)
			cowm.Remove(del)
		}

		checkeq(cowm, func(k g.Uint64) (uint32, bool) {
			v, ok := stdm[k]
			return v, ok
		}, t)
	}
}

func TestCopy(t *testing.T) {
	orig := NewMap[g.Uint64, uint32](1)

	for i := uint32(0); i < 10; i++ {
		orig.Put(g.Uint64(i), i)
	}

	cpy := orig.Copy()

	checkeq(cpy, orig.Get, t)

	cpy.Put(0, 42)

	if cpy.GetZ(0) != 42 {
		t.Fatal("didn't get 42")
	}
}
