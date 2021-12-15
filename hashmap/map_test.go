package hashmap

import (
	"testing"
	"math/rand"
)

func checkeq[K Hashable, V comparable](cm *Map[K, V], get func(k K) (V, bool), t *testing.T) {
	cm.Range(func(k K, v V) {
		if ov, ok := get(k); !ok {
			t.Fatalf("key %v should exist", k)
		} else if v != ov {
			t.Fatalf("value mismatch: %v != %v", v, ov)
		}
	})
}

func TestLookupMap(t *testing.T) {
	stdm := make(map[Uint64]uint32)
	cowm := NewCowMap[Uint64, uint32](1)

	const nops = 1000

	for i := 0; i < nops; i++ {
		key := Uint64(rand.Uint64())
		val := rand.Uint32()

		stdm[key] = val
		cowm.Set(key, val)

		checkeq(cowm, func(k Uint64) (uint32, bool) {
			v, ok := stdm[k]
			return v, ok
		}, t)
	}
}

func TestCopy(t *testing.T) {
	orig := NewCowMap[Uint64, uint32](1)

	for i := uint32(0); i < 10; i++ {
		orig.Set(Uint64(i), i)
	}

	cpy := orig.Copy()

	checkeq(cpy, orig.Get, t)

	cpy.Set(0, 42)

	if cpy.GetZ(0) != 42 {
		t.Fatal("didn't get 42")
	}
}
