package avl_test

import (
	"fmt"
	"math/rand"
	"testing"

	g "github.com/zyedidia/generic"
	"github.com/zyedidia/generic/avl"
)

func checkeq[K any, V comparable](cm *avl.Tree[K, V], n int, get func(k K) (V, bool), t *testing.T) {
	if sz := cm.Size(); sz != n {
		t.Fatalf("size mismatch: %d != %d", sz, n)
	}
	cm.Each(func(key K, val V) {
		if ov, ok := get(key); !ok {
			t.Fatalf("key %v should exist", key)
		} else if val != ov {
			t.Fatalf("value mismatch: %v != %v", val, ov)
		}
	})
}

func TestCrossCheck(t *testing.T) {
	stdm := make(map[int]int)
	get := func(k int) (int, bool) {
		v, ok := stdm[int(k)]
		return v, ok
	}
	tree := avl.New[int, int](g.Less[int])
	checkeq(tree, len(stdm), get, t)

	const nops = 1000
	for i := 0; i < nops; i++ {
		key := rand.Intn(100)
		val := rand.Int()
		op := rand.Intn(2)

		switch op {
		case 0:
			stdm[key] = val
			tree.Put(key, val)
		case 1:
			var del int
			for k := range stdm {
				del = k
				break
			}
			delete(stdm, del)
			tree.Remove(del)
		}

		checkeq(tree, len(stdm), get, t)
	}
}

func Example() {
	tree := avl.New[int, string](g.Less[int])

	tree.Put(42, "foo")
	tree.Put(-10, "bar")
	tree.Put(0, "baz")
	tree.Put(10, "quux")
	tree.Remove(10)

	tree.Each(func(key int, val string) {
		fmt.Println(key, val)
	})

	// Output:
	// -10 bar
	// 0 baz
	// 42 foo
}
