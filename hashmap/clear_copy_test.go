package hashmap_test

import (
	g "github.com/zyedidia/generic"
	"github.com/zyedidia/generic/hashmap"
	"testing"
)

func TestClearCopyIsolation(t *testing.T) {
	for _, clearOriginal := range []bool{false, true} {
		name := "copy"
		if clearOriginal {
			name = "original"
		}
		t.Run(name, func(t *testing.T) {
			m := hashmap.New[int, string](4, g.Equals[int], func(i int) uint64 { return uint64(i % 2) })
			m.Put(1, "one")
			m.Put(3, "three")
			c := m.Copy()
			target, kept := c, m
			if clearOriginal {
				target, kept = m, c
			}
			target.Clear()
			if target.Size() != 0 {
				t.Fatal("cleared map not empty")
			}
			for k, want := range map[int]string{1: "one", 3: "three"} {
				if got, ok := kept.Get(k); !ok || got != want {
					t.Fatalf("Clear erased another map: key%d got%q/%v", k, got, ok)
				}
			}
			if kept.Size() != 2 {
				t.Fatal("copy size changed")
			}
			target.Put(5, "new")
			kept.Put(7, "other")
			if _, ok := kept.Get(5); ok {
				t.Fatal("cleared map reused shared storage")
			}
			if _, ok := target.Get(7); ok {
				t.Fatal("sibling reused shared storage")
			}
		})
	}
}
func TestClearNestedCopyIsolation(t *testing.T) {
	m := hashmap.New[int, int](4, g.Equals[int], func(i int) uint64 { return uint64(i) })
	m.Put(1, 2)
	a := m.Copy()
	b := a.Copy()
	a.Clear()
	for _, kept := range []*hashmap.Map[int, int]{m, b} {
		if v, ok := kept.Get(1); !ok || v != 2 || kept.Size() != 1 {
			t.Fatal("Clear erased nested copy")
		}
	}
}
func TestClearUnsharedControl(t *testing.T) {
	m := hashmap.New[int, int](4, g.Equals[int], func(i int) uint64 { return uint64(i) })
	m.Put(1, 2)
	m.Clear()
	if m.Size() != 0 {
		t.Fatal("clear")
	}
	m.Put(3, 4)
	if v, ok := m.Get(3); !ok || v != 4 {
		t.Fatal("reuse")
	}
}
