package interval

import (
	"fmt"
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
			if overlaps(newIntrvl(tt.l1, tt.h1), newIntrvl(tt.l2, tt.h2)) != tt.expect {
				t.Fatalf("[%d, %d) vs [%d, %d): expected %v, got %v", tt.l1, tt.h1, tt.l2, tt.h2, tt.expect, !tt.expect)
			}
		})
	}
}

func TestPut(t *testing.T) {
	tree := New[int, string]()
	tree.Put(5, 7, "foo1")
	tree.Put(5, 9, "foo2")
	tree.Put(2, 4, "foo3")
	tree.Put(8, 9, "foo4")

	tests := []struct {
		low, high int
		vals      []string
	}{{
		low:  6,
		high: 7,
		vals: []string{"foo2"},
	}, {
		low:  7,
		high: 8,
		vals: []string{"foo2"},
	}, {
		low:  8,
		high: 9,
		vals: []string{"foo2", "foo4"},
	}, {
		low:  3,
		high: 6,
		vals: []string{"foo3", "foo2"},
	}}

	for i, tt := range tests {
		tt := tt
		t.Run(fmt.Sprintf("test_%d", i), func(t *testing.T) {
			ov := tree.Overlaps(tt.low, tt.high)
			if len(ov) != len(tt.vals) {
				t.Fatalf("Len missmatch: expected %d, got %d",
					len(tt.vals), len(ov))
			}

			for i, v := range tt.vals {
				if ov[i].Val == v {
					continue
				}

				t.Fatalf("Value mismatch at position %d: expected %q, got %q",
					i, v, ov[i].Val)
			}
		})
	}
}

func Example() {
	tree := New[int, string]()
	tree.Put(0, 10, "foo")
	tree.Put(5, 9, "bar")
	tree.Put(10, 11, "baz")
	tree.Put(-10, 4, "quux")

	vals := tree.Overlaps(4, 10)
	for _, v := range vals {
		fmt.Println(v.Val)
	}
	// Output:
	// foo
	// bar
}
