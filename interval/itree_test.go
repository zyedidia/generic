package interval

import (
	"testing"
)

func TestOverlaps(t *testing.T) {
	tests := []struct {
		l1, h1 int
		l2, h2 int
		expect bool
	}{
		{0, 5, 5, 10, false},
		{0, 5, 4, 5, true},
	}

	for _, tt := range tests {
		t.Run("test", func(t *testing.T) {
			if overlaps(intrvl{tt.l1, tt.h1}, tt.l2, tt.h2) != tt.expect {
				t.Fatalf("[%d, %d) vs [%d, %d): expected %v, got %v", tt.l1, tt.h1, tt.l2, tt.h2, tt.expect, !tt.expect)
			}
		})
	}
}

func TestSimple(t *testing.T) {
	tree := New[string]()
	tree.Put(0, 10, "foo")
	tree.Put(5, 9, "bar")
	tree.Put(10, 11, "baz")
	tree.Put(-10, -5, "quux")

	if len(tree.Overlaps(9, 15)) != 2 {
		t.Fatal("wrong number of overlaps")
	}
}
