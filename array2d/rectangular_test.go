package array2d

import (
	"fmt"
	"reflect"
	"testing"
)

func TestRectangularArray(t *testing.T) {
	for _, size := range [][2]int{{4, 2}, {2, 4}, {1, 3}, {3, 1}} {
		t.Run(fmt.Sprint(size), func(t *testing.T) {
			w, h := size[0], size[1]
			rows := make([][]int, h)
			for y := range rows {
				rows[y] = make([]int, w)
				for x := range rows[y] {
					rows[y][x] = 1 + y*w + x
				}
			}
			check := func(a Array2D[int]) {
				t.Helper()
				for y := range rows {
					if !reflect.DeepEqual(a.Row(y), rows[y]) {
						t.Errorf("row %d = %v, want %v", y, a.Row(y), rows[y])
					}
					if !reflect.DeepEqual(a.RowSpan(w-1, 0, y), rows[y]) {
						t.Errorf("span %d = %v, want %v", y, a.RowSpan(w-1, 0, y), rows[y])
					}
					for x, want := range rows[y] {
						if got := a.Get(x, y); got != want {
							t.Errorf("Get(%d,%d)=%d, want %d", x, y, got, want)
						}
					}
				}
				if got, want := a.String(), fmt.Sprint(rows); got != want {
					t.Errorf("String=%s, want %s", got, want)
				}
			}
			a := New[int](w, h)
			for y := range rows {
				for x, v := range rows[y] {
					a.Set(x, y, v)
				}
			}
			check(a)
			check(a.Copy())
			check(OfJagged(w, h, rows))
			a = NewFilled(w, h, -1)
			firstY := 1
			if h == 1 {
				firstY = 0
			}
			a.Fill(w-1, h-1, 0, firstY, 7)
			for y := range rows {
				for x := range rows[y] {
					want := -1
					if y >= firstY {
						want = 7
					}
					if a.Get(x, y) != want {
						t.Errorf("Fill at %d,%d=%d, want %d", x, y, a.Get(x, y), want)
					}
				}
			}
		})
	}
}
