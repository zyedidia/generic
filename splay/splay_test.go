package splay_test

import (
	"fmt"
	"math/rand"
	"testing"

	g "github.com/zyedidia/generic"
	agg "github.com/zyedidia/generic/aggregator"
	"github.com/zyedidia/generic/splay"
)

func checkeq[K, V comparable, A, R any](cm *splay.Tree[K, V, A, R], stdm map[K]V, t *testing.T) {
	n := len(stdm)
	if sz := cm.Size(); sz != n {
		t.Fatalf("size mismatch: %d != %d", sz, n)
	}
	for key, ov := range stdm {
		val, ok := cm.Get(key)
		if !ok {
			t.Fatalf("key %v should exist", key)
		} else if val != ov {
			t.Fatalf("value mismatch: %v != %v", val, ov)
		}
	}
}

func TestCrossCheck1(t *testing.T) {
	stdm := make(map[int]int)
	tree := splay.New(g.Less[int], agg.NewValueAggregator[int]())
	checkeq(tree, stdm, t)

	const nops = 3000
	for i := 0; i < nops; i++ {
		key := rand.Intn(100)
		val := rand.Int()
		op := rand.Intn(3)

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
		case 2:
			val, ok := tree.Get(key)
			mapV, mapOk := stdm[key]
			if ok != mapOk {
				t.Fatalf("key %v exists in one implementation but missing in another", key)
			} else if ok == true && val != mapV {
				t.Fatalf("value mismatch: %v != %v", val, mapV)
			}
		}
	}
}

func TestCrossCheck2(t *testing.T) {
	stdm := make(map[int]int)
	tree := splay.New(g.Less[int], agg.NewValueAggregator[int]())
	checkeq(tree, stdm, t)

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
		checkeq(tree, stdm, t)
	}
}

func TestPopUp(t *testing.T) {
	stdm := [100]int{}
	tree := splay.New(g.Less[int], agg.NewMinMaxAggregator(g.Less[int]))

	for i := 0; i < len(stdm); i++ {
		stdm[i] = -1
	}

	const nops = 3000
	for i := 0; i < nops; i++ {
		key := rand.Intn(100)
		val := rand.Intn(0xffff)
		op := rand.Intn(3)

		switch op {
		case 0:
			stdm[key] = val
			tree.Put(key, val)
		case 1:
			stdm[key] = -1
			tree.Remove(key)
		case 2:
			l := key
			r := rand.Intn(100)
			if l > r {
				l, r = r, l
			} else if l == r {
				r += 1
			}
			agg := tree.Range(l, r)
			if agg == nil {
				continue
			}
			treeMin := agg.Min()
			treeMax := agg.Max()
			listMin := 0x10000
			listMax := 0
			for i := l; i < r; i++ {
				if stdm[i] >= 0 {
					listMin = g.Min(listMin, stdm[i])
					listMax = g.Max(listMax, stdm[i])
				}
			}
			if treeMin != listMin || treeMax != listMax {
				t.Fatalf("value mismatch: (%v, %v) != (%v, %v)", listMin, listMax, treeMin, treeMax)
			}
		}
	}
}

func TestPushDown(t *testing.T) {
	stdm := make(map[int]int)
	tree := splay.New(g.Less[int], agg.NewRangeAssignAggregator[int]())
	checkeq(tree, stdm, t)

	const nops = 1000
	for i := 0; i < nops; i++ {
		key := rand.Intn(100)
		val := rand.Int()
		op := rand.Intn(3)

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
		case 2:
			l := key
			r := rand.Intn(100)
			if l > r {
				l, r = r, l
			} else if l == r {
				r += 1
			}
			tree.Range(l, r).Assign(val)
			for k := range stdm {
				if l <= k && k < r {
					stdm[k] = val
				}
			}
		}
		checkeq(tree, stdm, t)
	}
}

func Example() {
	tree := splay.New(g.Less[int], agg.NewMinMaxAggregator(g.Less[string]))

	tree.Put(42, "foo")
	tree.Put(-10, "bar")
	tree.Put(0, "baz")
	tree.Put(10, "quux")
	tree.Remove(10)

	tree.Each(func(key int, val string) {
		fmt.Println(key, val)
	})

	fmt.Println(tree.Range(-10, 10).Min())

	// Output:
	// -10 bar
	// 0 baz
	// 42 foo
	// bar
}
