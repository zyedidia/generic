package hashset_test

import (
	g "github.com/zyedidia/generic"
	"github.com/zyedidia/generic/hashset"
	"testing"
)

func TestClearSetCopyIsolation(t *testing.T) {
	for _, original := range []bool{false, true} {
		name := "copy"
		if original {
			name = "original"
		}
		t.Run(name, func(t *testing.T) {
			s := hashset.New[int](4, g.Equals[int], func(i int) uint64 { return uint64(i) })
			s.Put(1)
			s.Put(3)
			c := s.Copy()
			target, kept := c, s
			if original {
				target, kept = s, c
			}
			target.Clear()
			if target.Size() != 0 || kept.Size() != 2 || !kept.Has(1) || !kept.Has(3) {
				t.Fatal("Clear changed sibling set")
			}
			target.Put(5)
			if kept.Has(5) {
				t.Fatal("reuse changed sibling")
			}
		})
	}
}
