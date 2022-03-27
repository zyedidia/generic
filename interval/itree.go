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
)

type KV[V any] struct {
	Low, High int
	Val       V
}

// intrvl represents an interval over [low, high).
type intrvl struct {
	Low, High int
}

func overlaps(i1 intrvl, low, high int) bool {
	return i1.Low < high && i1.High > low
}

// Tree implements an interval tree. All intervals must have unique starting
// positions.
type Tree[V any] struct {
	root *node[V]
}

// New returns an empty interval tree.
func New[V any]() *Tree[V] {
	return &Tree[V]{}
}

// Put associates the interval 'key' with 'value'.
func (t *Tree[V]) Put(low, high int, value V) {
	t.root = t.root.add(intrvl{low, high}, value)
}

// Overlaps returns all values that overlap with the given range.
func (t *Tree[V]) Overlaps(low, high int) []V {
	var result []V
	return t.root.overlaps(intrvl{low, high}, result)
}

// Remove deletes the interval starting at 'pos'.
func (t *Tree[V]) Remove(pos int) {
	t.root = t.root.remove(pos)
}

// Get returns the value associated with the interval starting at 'pos', or
// 'false' if no such value exists.
func (t *Tree[V]) Get(pos int) (V, bool) {
	n := t.root.search(pos)
	if n == nil {
		var v V
		return v, false
	}
	return n.value, true
}

// Each calls 'fn' on every element in the tree, and its corresponding
// interval, in order sorted by starting position.
func (t *Tree[V]) Each(fn func(low, high int, val V)) {
	t.root.each(fn)
}

// Height returns the height of the tree.
func (t *Tree[V]) Height() int {
	return t.root.getHeight()
}

// Size returns the number of elements in the tree.
func (t *Tree[V]) Size() int {
	return t.root.size()
}

type node[V any] struct {
	key   intrvl
	value V

	max    int
	height int
	left   *node[V]
	right  *node[V]
}

func (n *node[V]) add(key intrvl, value V) *node[V] {
	if n == nil {
		return &node[V]{
			key:    key,
			value:  value,
			max:    key.High,
			height: 1,
			left:   nil,
			right:  nil,
		}
	}

	if key.Low < n.key.Low {
		n.left = n.left.add(key, value)
	} else if key.Low > n.key.Low {
		n.right = n.right.add(key, value)
	} else {
		n.value = value
	}
	return n.rebalanceTree()
}

func (n *node[V]) updateMax() {
	if n != nil {
		if n.right != nil {
			n.max = g.Max(n.max, n.right.max)
		}
		if n.left != nil {
			n.max = g.Max(n.max, n.left.max)
		}
	}
}

func (n *node[V]) remove(pos int) *node[V] {
	if n == nil {
		return nil
	}
	if pos < n.key.Low {
		n.left = n.left.remove(pos)
	} else if pos > n.key.Low {
		n.right = n.right.remove(pos)
	} else {
		if n.left != nil && n.right != nil {
			rightMinNode := n.right.findSmallest()
			n.key = rightMinNode.key
			n.value = rightMinNode.value
			n.right = n.right.remove(rightMinNode.key.Low)
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

func (n *node[V]) search(pos int) *node[V] {
	if n == nil {
		return nil
	}
	if pos < n.key.Low {
		return n.left.search(pos)
	} else if pos > n.key.Low {
		return n.right.search(pos)
	} else {
		return n
	}
}

func (n *node[V]) overlaps(key intrvl, result []V) []V {
	if n == nil {
		return result
	}

	if key.Low >= n.max {
		return result
	}

	result = n.left.overlaps(key, result)

	if overlaps(n.key, key.Low, key.High) {
		result = append(result, n.value)
	}

	if key.High <= n.key.Low {
		return result
	}

	result = n.right.overlaps(key, result)
	return result
}

func (n *node[V]) each(fn func(low, high int, val V)) {
	if n == nil {
		return
	}
	n.left.each(fn)
	fn(n.key.Low, n.key.High, n.value)
	n.right.each(fn)
}

func (n *node[V]) getHeight() int {
	if n == nil {
		return 0
	}
	return n.height
}

func (n *node[V]) recalculateHeight() {
	n.height = 1 + g.Max(n.left.getHeight(), n.right.getHeight())
}

func (n *node[V]) rebalanceTree() *node[V] {
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

func (n *node[V]) rotateLeft() *node[V] {
	newRoot := n.right
	n.right = newRoot.left
	newRoot.left = n

	n.recalculateHeight()
	n.updateMax()
	newRoot.recalculateHeight()
	newRoot.updateMax()
	return newRoot
}

func (n *node[V]) rotateRight() *node[V] {
	newRoot := n.left
	n.left = newRoot.right
	newRoot.right = n

	n.recalculateHeight()
	n.updateMax()
	newRoot.recalculateHeight()
	newRoot.updateMax()
	return newRoot
}

func (n *node[V]) findSmallest() *node[V] {
	if n.left != nil {
		return n.left.findSmallest()
	} else {
		return n
	}
}

func (n *node[V]) size() int {
	s := 1
	if n.left != nil {
		s += n.left.size()
	}
	if n.right != nil {
		s += n.right.size()
	}
	return s
}
