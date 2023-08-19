// Package array2d contains an implementation of a 2D array.
package array2d

import (
	"fmt"
	"strings"
)

// New initializes a 2-dimensional array with all zero values.
func New[T any](width, height int) Array2D[T] {
	return Array2D[T]{
		width:  width,
		height: height,
		slice:  make([]T, width*height),
	}
}

// NewFilled initializes a 2-dimensional array with a value.
func NewFilled[T any](width, height int, value T) Array2D[T] {
	slice := make([]T, width*height)
	fill(slice, value)
	return Array2D[T]{
		width:  width,
		height: height,
		slice:  slice,
	}
}

// OfJagged initializes a 2-dimensional array based on a jagged
// slice of rows of values. Values from the jagged slice that are out of bounds
// are ignored.
func OfJagged[J ~[]S, S ~[]E, E any](width, height int, jagged J) Array2D[E] {
	arr := New[E](width, height)
	for y, row := range jagged {
		copy(arr.Row(y), row)
	}
	return arr
}

// Array2D is a 2-dimensional array.
type Array2D[T any] struct {
	width, height int
	slice         []T
}

// String returns a string representation of this array.
func (a Array2D[T]) String() string {
	var sb strings.Builder
	sb.WriteByte('[')
	for y := 0; y < a.height; y++ {
		if y > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteByte('[')
		for x := 0; x < a.width; x++ {
			if x > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprint(&sb, a.getUnchecked(x, y))
		}
		sb.WriteByte(']')
	}
	sb.WriteByte(']')
	return sb.String()
}

// Get returns a value from the array.
//
// The function will panic on out-of-bounds access.
func (a Array2D[T]) Get(x, y int) T {
	if x < 0 || x >= a.width {
		panic(fmt.Sprintf("array2d: x index out of range [%d] with width %d", x, a.width))
	}
	if y < 0 || y >= a.height {
		panic(fmt.Sprintf("array2d: y index out of range [%d] with height %d", y, a.height))
	}
	return a.getUnchecked(x, y)
}

func (a Array2D[T]) getUnchecked(x, y int) T {
	return a.slice[x+y*a.height]
}

// Set sets a value in the array.
//
// The function will panic on out-of-bounds access.
func (a Array2D[T]) Set(x, y int, value T) {
	if x < 0 || x >= a.width {
		panic(fmt.Sprintf("array2d: x index out of range [%d] with width %d", x, a.width))
	}
	if y < 0 || y >= a.height {
		panic(fmt.Sprintf("array2d: y index out of range [%d] with height %d", y, a.height))
	}
	a.setUnchecked(x, y, value)
}

func (a Array2D[T]) setUnchecked(x, y int, value T) {
	a.slice[x+y*a.height] = value
}

// Width returns the width of this array. The maximum x value is Width()-1.
func (a Array2D[T]) Width() int {
	return a.width
}

// Height returns the height of this array. The maximum y value is Height()-1.
func (a Array2D[T]) Height() int {
	return a.height
}

// Copy returns a shallow copy of this array.
func (a Array2D[T]) Copy() Array2D[T] {
	slice := make([]T, len(a.slice))
	copy(slice, a.slice)
	return Array2D[T]{
		width:  a.width,
		height: a.height,
		slice:  slice,
	}
}

// RowSpan returns a mutable slice for part of a row. Changing values in this
// slice will affect the array.
func (a Array2D[T]) RowSpan(x1, x2, y int) []T {
	if x1 < 0 || x1 >= a.width {
		panic(fmt.Sprintf("array2d: x1 index out of range [%d] with width %d", x1, a.width))
	}
	if y < 0 || y >= a.height {
		panic(fmt.Sprintf("array2d: y index out of range [%d] with height %d", y, a.height))
	}
	if x2 < 0 || x2 >= a.width {
		panic(fmt.Sprintf("array2d: x2 index out of range [%d] with width %d", x2, a.width))
	}
	if x2 < x1 {
		x1, x2 = x2, x1
	}
	return a.slice[x1+y*a.height : 1+x2+y*a.height]
}

// Row returns a mutable slice for an entire row. Changing values in this slice
// will affect the array.
func (a Array2D[T]) Row(y int) []T {
	if y < 0 || y >= a.height {
		panic(fmt.Sprintf("array2d: y index out of range [%d] with height %d", y, a.height))
	}
	return a.slice[y*a.height : a.width+y*a.height]
}

// Fill will assign all values inside the region to the specified value.
// The coordinates are inclusive, meaning all values from [x1,y1] including
// [x1,y1] to [x2,y2] including [x2,y2] are set.
//
// The method sorts the arguments, so x2 may be lower than x1 and y2 may be
// lower than y1.
func (a Array2D[T]) Fill(x1, y1, x2, y2 int, value T) {
	if x1 < 0 || x1 >= a.width {
		panic(fmt.Sprintf("array2d: x1 index out of range [%d] with width %d", x1, a.width))
	}
	if y1 < 0 || y1 >= a.height {
		panic(fmt.Sprintf("array2d: y1 index out of range [%d] with height %d", y1, a.height))
	}
	if x2 < 0 || x2 >= a.width {
		panic(fmt.Sprintf("array2d: x2 index out of range [%d] with width %d", x2, a.width))
	}
	if y2 < 0 || y2 >= a.height {
		panic(fmt.Sprintf("array2d: y2 index out of range [%d] with height %d", y2, a.height))
	}
	if x2 < x1 {
		x1, x2 = x2, x1
	}
	if y2 < y1 {
		y1, y2 = y2, y1
	}
	firstRow := a.slice[x1+y1*a.height : 1+x2+y1*a.height]
	fill(firstRow, value)
	for y := y1 + 1; y <= y2; y++ {
		copy(a.slice[x1+y*a.height:1+x2+y*a.height], firstRow)
	}
}

func fill[E any](slice []E, value E) {
	if len(slice) == 0 {
		return
	}
	// Exponential copy to fill a slice
	slice[0] = value
	for i := 1; i < len(slice); i += i {
		copy(slice[i:], slice[:i])
	}
}
