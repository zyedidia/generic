// Package heap provides an implementation of a binary heap.
// A binary heap (binary min-heap) is a tree with the property that each node
// is the minimum-valued node in its subtree.
package heap

import (
	g "github.com/zyedidia/generic"
)

// Heap implements a binary heap.
type Heap[T any] struct {
	data []T
	less func(a, b T) bool
}

// New returns a new heap with the given less function.
func New[T any](less g.LessFn[T]) *Heap[T] {
	return &Heap[T]{
		data: make([]T, 0),
		less: less,
	}
}

// From returns a new heap with the given less function and initial data.
func From[T any](less g.LessFn[T], t ...T) *Heap[T] {
	return FromSlice(less, t)
}

// FromSlice returns a new heap with the given less function and initial data.
// The `data` is not copied and used as the inside array.
func FromSlice[T any](less g.LessFn[T], data []T) *Heap[T] {
	n := len(data)
	for i := n/2 - 1; i >= 0; i-- {
		down(data, i, less)
	}

	return &Heap[T]{
		data: data,
		less: less,
	}
}

// Push pushes the given element onto the heap.
func (h *Heap[T]) Push(x T) {
	h.data = append(h.data, x)
	up(h.data, len(h.data)-1, h.less)
}

// Pop removes and returns the minimum element from the heap. If the heap is
// empty, it returns zero value and false.
func (h *Heap[T]) Pop() (T, bool) {
	var x T
	if h.Size() == 0 {
		return x, false
	}

	x = h.data[0]

	h.data[0] = h.data[len(h.data)-1]
	h.data = h.data[:len(h.data)-1]
	down(h.data, 0, h.less)

	return x, true
}

// Peek returns the minimum element from the heap without removing it. if the
// heap is empty, it returns zero value and false.
func (h *Heap[T]) Peek() (T, bool) {
	if h.Size() == 0 {
		var x T
		return x, false
	}

	return h.data[0], true
}

// Size returns the number of elements in the heap.
func (h *Heap[T]) Size() int {
	return len(h.data)
}

func down[T any](h []T, i int, less g.LessFn[T]) {
	for {
		left, right := 2*i+1, 2*i+2
		if left >= len(h) || left < 0 { // `left < 0` in case of overflow
			break
		}

		// find the smallest child
		j := left
		if right < len(h) && less(h[right], h[left]) {
			j = right
		}

		if !less(h[j], h[i]) {
			break
		}

		h[i], h[j] = h[j], h[i]
		i = j
	}
}

func up[T any](h []T, i int, less g.LessFn[T]) {
	for {
		parent := (i - 1) / 2
		if i == 0 || !less(h[i], h[parent]) {
			break
		}

		h[i], h[parent] = h[parent], h[i]
		i = parent
	}
}
