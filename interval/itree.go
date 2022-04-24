// Package interval provides an implementation of an interval tree built using
// an augmented AVL tree. An interval tree stores values associated with
// intervals, and can efficiently determine which intervals overlap with
// others. All intervals must have a unique starting position. It supports the
// following operations, where 'n' is the number of
// intervals in the tree:
//
// * Put: add an interval to the tree. Complexity: O(lg n).
//
// * Get: find an interval with a given starting position. Complexity O(lg n).
//
// * Overlaps: find all intervals that overlap with a given interval. Complexity:
//   O(lg n + m), where 'm' is the size of the result (number of overlapping
//   intervals found).
//
// * Remove: remove the interval at a given position. Complexity: O(lg n).
package interval

import (
	g "github.com/zyedidia/generic"
	"golang.org/x/exp/constraints"
)

type KV[V any] struct {
	Low, High int
	Val       V
}

// intrvl represents an interval over [low, high).
type intrvl[I constraints.Ordered] struct {
	low, high I
}

func newIntrvl[I constraints.Ordered](low, high I) intrvl[I] {
	return intrvl[I]{low, high}
}

func overlaps[I constraints.Ordered](i1 intrvl[I], low, high I) bool {
	return i1.low < high && i1.high > low
}

// Tree implements an interval tree. All intervals must have unique starting
// positions.
type Tree[I constraints.Ordered, V any] struct {
	root *node[I, V]
}

// New returns an empty interval tree.
func New[I constraints.Ordered, V any]() *Tree[I, V] {
	return &Tree[I, V]{}
}

// Put associates the interval 'key' with 'value'.
func (t *Tree[I, V]) Put(low, high I, value V) {
	t.root = t.root.add(newIntrvl(low, high), value)
}

// Overlaps returns all values that overlap with the given range.
func (t *Tree[I, V]) Overlaps(low, high I) []V {
	return t.root.overlaps(newIntrvl(low, high), nil)
}

// Remove deletes the interval starting at 'pos'.
func (t *Tree[I, V]) Remove(pos I) {
	t.root = t.root.remove(pos)
}

// Get returns the value associated with the interval starting at 'pos', or
// 'false' if no such value exists.
func (t *Tree[I, V]) Get(pos I) (V, bool) {
	n := t.root.search(pos)
	if n == nil {
		var v V
		return v, false
	}
	return n.value, true
}

// Each calls 'fn' on every element in the tree, and its corresponding
// interval, in order sorted by starting position.
func (t *Tree[I, V]) Each(fn func(low, high I, val V)) {
	t.root.each(fn)
}

// Height returns the height of the tree.
func (t *Tree[I, V]) Height() int {
	return t.root.getHeight()
}

// Size returns the number of elements in the tree.
func (t *Tree[I, V]) Size() int {
	return t.root.size()
}

type node[I constraints.Ordered, V any] struct {
	key   intrvl[I]
	value V

	max    I
	height int
	left   *node[I, V]
	right  *node[I, V]
}

func (n *node[I, V]) add(key intrvl[I], value V) *node[I, V] {
	if n == nil {
		return &node[I, V]{
			key:    key,
			value:  value,
			max:    key.high,
			height: 1,
			left:   nil,
			right:  nil,
		}
	}

	if key.low < n.key.low {
		n.left = n.left.add(key, value)
	} else if key.low > n.key.low {
		n.right = n.right.add(key, value)
	} else {
		n.value = value
	}
	return n.rebalanceTree()
}

func (n *node[I, V]) updateMax() {
	if n != nil {
		if n.right != nil {
			n.max = g.Max(n.max, n.right.max)
		}
		if n.left != nil {
			n.max = g.Max(n.max, n.left.max)
		}
	}
}

func (n *node[I, V]) remove(pos I) *node[I, V] {
	if n == nil {
		return nil
	}
	if pos < n.key.low {
		n.left = n.left.remove(pos)
	} else if pos > n.key.low {
		n.right = n.right.remove(pos)
	} else {
		if n.left != nil && n.right != nil {
			rightMinNode := n.right.findSmallest()
			n.key = rightMinNode.key
			n.value = rightMinNode.value
			n.right = n.right.remove(rightMinNode.key.low)
		} else if n.left != nil {
			n = n.left
		} else if n.right != nil {
			n = n.right
		} else {
			n = nil
			return n
		}

	}
	return n.rebalanceTree()
}

func (n *node[I, V]) search(pos I) *node[I, V] {
	if n == nil {
		return nil
	}
	if pos < n.key.low {
		return n.left.search(pos)
	} else if pos > n.key.low {
		return n.right.search(pos)
	} else {
		return n
	}
}

func (n *node[I, V]) overlaps(key intrvl[I], result []V) []V {
	if n == nil {
		return result
	}

	if key.low >= n.max {
		return result
	}

	result = n.left.overlaps(key, result)

	if overlaps(n.key, key.low, key.high) {
		result = append(result, n.value)
	}

	if key.high <= n.key.low {
		return result
	}

	result = n.right.overlaps(key, result)
	return result
}

func (n *node[I, V]) each(fn func(low, high I, val V)) {
	if n == nil {
		return
	}
	n.left.each(fn)
	fn(n.key.low, n.key.high, n.value)
	n.right.each(fn)
}

func (n *node[I, V]) getHeight() int {
	if n == nil {
		return 0
	}
	return n.height
}

func (n *node[I, V]) recalculateHeight() {
	n.height = 1 + g.Max(n.left.getHeight(), n.right.getHeight())
}

func (n *node[I, V]) rebalanceTree() *node[I, V] {
	if n == nil {
		return n
	}
	n.recalculateHeight()
	n.updateMax()

	balanceFactor := n.left.getHeight() - n.right.getHeight()
	if balanceFactor <= -2 {
		if n.right.left.getHeight() > n.right.right.getHeight() {
			n.right = n.right.rotateRight()
		}
		return n.rotateLeft()
	} else if balanceFactor >= 2 {
		if n.left.right.getHeight() > n.left.left.getHeight() {
			n.left = n.left.rotateLeft()
		}
		return n.rotateRight()
	}
	return n
}

func (n *node[I, V]) rotateLeft() *node[I, V] {
	newRoot := n.right
	n.right = newRoot.left
	newRoot.left = n

	n.recalculateHeight()
	n.updateMax()
	newRoot.recalculateHeight()
	newRoot.updateMax()
	return newRoot
}

func (n *node[I, V]) rotateRight() *node[I, V] {
	newRoot := n.left
	n.left = newRoot.right
	newRoot.right = n

	n.recalculateHeight()
	n.updateMax()
	newRoot.recalculateHeight()
	newRoot.updateMax()
	return newRoot
}

func (n *node[I, V]) findSmallest() *node[I, V] {
	if n.left != nil {
		return n.left.findSmallest()
	} else {
		return n
	}
}

func (n *node[I, V]) size() int {
	s := 1
	if n.left != nil {
		s += n.left.size()
	}
	if n.right != nil {
		s += n.right.size()
	}
	return s
}
