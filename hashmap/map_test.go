package hashmap

import (
	"fmt"
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
		key := g.Uint64(rand.Intn(100))
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

	if v, _ := cpy.Get(0); v != 42 {
		t.Fatal("didn't get 42")
	}
}

func Example() {
	m := NewMap[g.String, g.Int](1)
	m.Put("foo", 42)
	m.Put("bar", 13)

	fmt.Println(m.Get("foo"))
	fmt.Println(m.Get("baz"))

	m.Remove("foo")

	fmt.Println(m.Get("foo"))

	// Output:
	// 42 true
	// 0 false
	// 0 false
}
