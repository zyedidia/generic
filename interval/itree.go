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
	"fmt"

	"github.com/zyedidia/generic"
	"golang.org/x/exp/constraints"
)

type KV[I constraints.Ordered, V any] struct {
	Low, High I
	Val       V
}

func newKV[I constraints.Ordered, V any](n *node[I, V]) KV[I, V] {
	return KV[I, V]{
		Low:  n.key.low,
		High: n.key.high,
		Val:  n.value,
	}
}

// intrvl represents an interval over [low, high).
type intrvl[I constraints.Ordered] struct {
	low, high I
}

func newIntrvl[I constraints.Ordered](low, high I) intrvl[I] {
	if low > high {
		panic(fmt.Sprintf("low cannot be greater than high: %v > %v", low, high))
	}
	return intrvl[I]{low, high}
}

func overlaps[I constraints.Ordered](i1 intrvl[I], i2 intrvl[I]) bool {
	return i1.low < i2.high && i1.high > i2.low
}

// Tree implements an interval tree. All intervals must have unique starting
// positions. Every low bound if an interval is inclusive, while high is
// exclusive.
type Tree[I constraints.Ordered, V any] struct {
	root *node[I, V]
}

// New returns an empty interval tree.
func New[I constraints.Ordered, V any]() *Tree[I, V] {
	return &Tree[I, V]{}
}

// Add associates the interval [low, high) with value.
//
// If an interval starting at low already exists in t, this method doesn't
// perform any change of the tree, but returns the conflicting interval.
func (t *Tree[I, V]) Add(low, high I, value V) (KV[I, V], bool) {
	newRoot, kv, ok := t.root.insert(newIntrvl(low, high), value, false)
	t.root = newRoot
	return kv, ok
}

// Put associates the interval [low, high) with value.
//
// If an interval starting at low already exists, this method will replace it.
// In such a case the conflicting (replaced) interval is returned.
func (t *Tree[I, V]) Put(low, high I, value V) (KV[I, V], bool) {
	newRoot, kv, ok := t.root.insert(newIntrvl(low, high), value, true)
	t.root = newRoot
	return kv, ok
}

// Overlaps returns all values that overlap with the given range. List returned
// is sorted by low positions of intervals.
func (t *Tree[I, V]) Overlaps(low, high I) []KV[I, V] {
	return t.root.overlaps(newIntrvl(low, high), nil)
}

// Remove deletes the interval starting at low. The removed interval is
// returned. If no such interval existed in a tree, the returned value is false.
func (t *Tree[I, V]) Remove(low I) (KV[I, V], bool) {
	newRoot, kv, ok := t.root.remove(low)
	t.root = newRoot
	return kv, ok
}

// Get returns the interval and value associated with the interval starting at
// low, or false if no such value exists.
func (t *Tree[I, V]) Get(low I) (KV[I, V], bool) {
	n := t.root.search(low)
	if n == nil {
		return KV[I, V]{}, false
	}
	return newKV(n), true
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

	height int
	left   *node[I, V]
	right  *node[I, V]

	// max is highest upper bound of all intervals stored in subtree which
	// node as its root.
	max I
}

// insert inserts interval key associated with value value to the tree.
//
// If interval starting at key.low already exists in a tree, behaviour of this
// method is defined by overwrite parameter. If it's true, the value is
// replaced. Otherwise whole subtree is left unchanged.
//
// This method returns new root node of a subtree rooted in n after insertion,
// an interval starting at key.low which already exists in the subtree and a
// flag if such an interval exists.
func (n *node[I, V]) insert(
	key intrvl[I],
	value V,
	overwrite bool,
) (*node[I, V], KV[I, V], bool) {
	if n == nil {
		return &node[I, V]{
			key:    key,
			value:  value,
			max:    key.high,
			height: 1,
			left:   nil,
			right:  nil,
		}, KV[I, V]{}, false
	}

	var kv KV[I, V]
	var evicted bool
	if key.low < n.key.low {
		n.left, kv, evicted = n.left.insert(key, value, overwrite)
	} else if key.low > n.key.low {
		n.right, kv, evicted = n.right.insert(key, value, overwrite)
	} else {
		if !overwrite {
			return n, newKV(n), true
		}

		kv, evicted = newKV(n), true

		n.key = key
		n.value = value
	}

	return n.rebalanceTree(), kv, evicted
}

func (n *node[I, V]) updateMax() {
	if n == nil {
		return
	}

	n.max = n.key.high
	if n.right != nil {
		n.max = generic.Max(n.max, n.right.max)
	}
	if n.left != nil {
		n.max = generic.Max(n.max, n.left.max)
	}
}

// remove removes interval starting at pos from a subtree. This function returns
// the new root of subtree rooted in n after pos removal, the KV removed and an
// information if any deletion happened (i.e. if interval starting at pos
// exists).
func (n *node[I, V]) remove(low I) (*node[I, V], KV[I, V], bool) {
	if n == nil {
		return nil, KV[I, V]{}, false
	}

	var kv KV[I, V]
	var removed bool
	if low < n.key.low {
		n.left, kv, removed = n.left.remove(low)
	} else if low > n.key.low {
		n.right, kv, removed = n.right.remove(low)
	} else {
		kv, removed = newKV(n), true
		n = n.removeThis()
	}

	return n.rebalanceTree(), kv, removed
}

// removeThis deletes n from subtree rooted in n and returns new root of the
// subtree.
func (n *node[I, V]) removeThis() *node[I, V] {
	// This can return nil if n has no children (n.right == nil).
	if n.left == nil {
		return n.right
	}
	if n.right == nil {
		return n.left
	}

	rightMinNode := n.right.findSmallest()
	n.key = rightMinNode.key
	n.value = rightMinNode.value
	n.right, _, _ = n.right.remove(rightMinNode.key.low)

	return n
}

func (n *node[I, V]) search(low I) *node[I, V] {
	if n == nil {
		return nil
	}

	if low < n.key.low {
		return n.left.search(low)
	} else if low > n.key.low {
		return n.right.search(low)
	} else {
		return n
	}
}

func (n *node[I, V]) overlaps(key intrvl[I], result []KV[I, V]) []KV[I, V] {
	if n == nil {
		return result
	}

	if key.low >= n.max {
		return result
	}

	result = n.left.overlaps(key, result)

	if overlaps(n.key, key) {
		result = append(result, newKV(n))
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
	n.height = 1 + generic.Max(n.left.getHeight(), n.right.getHeight())
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
