package btree_test

import (
	g "github.com/zyedidia/generic"
	"github.com/zyedidia/generic/btree"
	"testing"
)

func TestOverwriteSize(t *testing.T) {
	for _, n := range []int{1, 64, 256} {
		t.Run(string(rune('a'+n/64)), func(t *testing.T) {
			tr := btree.New[int, int](g.Less[int])
			for i := 0; i < n; i++ {
				tr.Put(i, i)
			}
			for i := 0; i < n; i++ {
				tr.Put(i, -i)
			}
			if tr.Size() != n {
				t.Fatalf("overwrite grew size: got%d want%d", tr.Size(), n)
			}
			seen := 0
			tr.Each(func(k, v int) {
				seen++
				if v != -k {
					t.Fatal("value not updated")
				}
			})
			if seen != n {
				t.Fatal("duplicate enumeration")
			}
			for i := 0; i < n; i++ {
				tr.Remove(i)
				tr.Remove(i)
			}
			if tr.Size() != 0 {
				t.Fatal("not empty after removing every key")
			}
		})
	}
}
func TestOverwriteRevivedTombstone(t *testing.T) {
	tr := btree.New[int, int](g.Less[int])
	tr.Put(1, 1)
	tr.Remove(1)
	tr.Put(1, 2)
	tr.Put(1, 3)
	if tr.Size() != 1 {
		t.Fatal("revived key counted twice")
	}
	tr.Remove(1)
	if tr.Size() != 0 {
		t.Fatal("revived key removal leaves inflated size")
	}
}
func TestDistinctPutSizeControl(t *testing.T) {
	tr := btree.New[int, int](g.Less[int])
	for i := 0; i < 100; i++ {
		tr.Put(i, i)
	}
	if tr.Size() != 100 {
		t.Fatal("distinct keys count")
	}
}
