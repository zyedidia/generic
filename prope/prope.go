// Package prope provides an implementation of a persistent rope data structure.
// It is similar to the base rope data structure, but the changes
// are saved separately without modifying the original data structure by
// sharing data between multiple versions. The time complexity of operations
// stay the same, but they are generally a bit slower:
//
// * Remove: O(lg n).
//
// * Insert: O(lg n).
//
// * Slice: O(lg n + m), where m is the size of the slice.
//
// * At: O(lg n).
//
// The main difference is in space complexity, as the persistent data structure
// allows creating a copy with an insertion or removal in O(lg n) space,
// instead of duplicating the entire rope for each change.
// This also prevents the O(n) time complexity of cloning the rope
// to save a version, as this is done inside the operations in a more efficient
// manner.
package prope

import g "github.com/zyedidia/generic"

var (
	// SplitLength is the threshold above which slices will be split into
	// separate nodes. Larger values will take make operations take more
	// memory.
	SplitLength = 256
	// JoinLength is the threshold below which nodes will be merged into
	// slices.
	JoinLength = SplitLength / 2
	// RebalanceRatio is the threshold used to trigger a rebuild during a
	// rebalance operation.
	RebalanceRatio = 1.5
)

type nodeType byte

const (
	tLeaf nodeType = iota
	tNode
)

// A Node in the rope structure. If the kind is tLeaf, only the value and
// length are valid, and if the kind is tNode, only length, left, right are
// valid.
type Node[V any] struct {
	kind        nodeType
	value       []V
	length      int
	left, right *Node[V]
}

// New returns a new rope node from the given byte slice. The underlying
// data is not copied so the user should ensure that the slice will
// not be modified after the rope is created.
func New[V any](b []V) *Node[V] {
	n := &Node[V]{
		kind:   tLeaf,
		value:  b,
		length: len(b),
	}
	n.adjust()
	return n
}

// Len returns the number of elements stored in the rope.
func (n *Node[V]) Len() int {
	return n.length
}

// Value returns the elements of this node concatenated into a slice.
func (n *Node[V]) Value() []V {
	value := make([]V, n.length)
	n.copy(value)
	return value
}

// Slice returns the range of the rope from [start:end).
func (n *Node[V]) Slice(start, end int) []V {
	slice := make([]V, end-start)
	n.copySlice(slice, start, end)
	return slice
}

// At returns the element at the given position.
func (n *Node[V]) At(pos int) V {
	holder := make([]V, 1)
	n.copySlice(holder, pos, pos+1)
	return holder[0]
}

// Insert returns a new version of the rope with the given
// value inserted at pos.
func (n *Node[V]) Insert(pos int, value []V) *Node[V] {
	if n.kind == tLeaf {
		return New(insert(n.value, pos, value)) // Adjusting is done here
	}

	changedNode := &Node[V]{
		kind:   tNode,
		length: n.length + len(value),
		left:   n.left,
		right:  n.right,
	}

	if pos < n.left.length {
		changedNode.left = n.left.Insert(pos, value)
	} else {
		changedNode.right = n.right.Insert(pos-n.left.length, value)
	}
	return changedNode
}

// Remove returns a new version of the rope with the elements
// in the [start:end) range removed.
func (n *Node[V]) Remove(start, end int) *Node[V] {
	if n.kind == tLeaf {
		return New(remove(n.value, start, end))
	}

	changedNode := &Node[V]{
		kind: tNode,
	}

	leftStart, leftEnd := bound(start, end, n.left.length)
	changedNode.left = n.left.Remove(leftStart, leftEnd)

	rightStart, rightEnd := bound(start-n.left.length, end-n.left.length, n.right.length)
	changedNode.right = n.right.Remove(rightStart, rightEnd)

	changedNode.length = changedNode.right.length + changedNode.left.length
	changedNode.adjust()
	return changedNode
}

// SplitAt splits the node at the given index and returns two new ropes
// corresponding to the left and right portions of the split.
func (n *Node[V]) SplitAt(i int) (*Node[V], *Node[V]) {
	if n.kind == tLeaf {
		return New(n.value[:i]), New(n.value[i:])
	}
	if i == n.left.length {
		return n.left, n.right
	} else if i < n.left.length {
		l, r := n.left.SplitAt(i)
		return l, Join(r, n.right)
	} else {
		l, r := n.right.SplitAt(i - n.left.length)
		return Join(n.left, l), r
	}
}

// Rebalance finds unbalanced nodes and rebuilds them.
// Rebuilded nodes does not share memory with their old versions,
// so sometimes this operation will take up a lot of memory.
func (n *Node[V]) Rebalance() {
	if n.kind == tLeaf {
		return
	}
	lratio := float64(n.left.length) / float64(n.right.length)
	rratio := float64(n.right.length) / float64(n.left.length)
	if g.Max(lratio, rratio) > RebalanceRatio {
		n.Rebuild()
	} else {
		n.left.Rebalance()
		n.right.Rebalance()
	}
}

// Rebuild rebuilds the entire rope structure, resulting in a balanced tree.
// The rebuilded node does not share memory with its old versions,
// so this operation will take the same space as creating the node from scratch.
func (n *Node[V]) Rebuild() {
	if n.kind == tLeaf {
		return
	}
	*n = *New(n.Value())
}

// Join creates a merged version of all of the ropes.
func Join[V any](nodes ...*Node[V]) *Node[V] {
	if len(nodes) == 0 {
		return New([]V{})
	}
	accum := nodes[0]
	for _, node := range nodes[1:] {
		accum = &Node[V]{
			kind:   tNode,
			left:   accum,
			right:  node,
			length: accum.length + node.length,
		}
		accum.adjust()
	}
	return accum
}

func (n *Node[V]) copy(dst []V) {
	if n.kind == tLeaf {
		copy(dst, n.value)
		return
	}
	n.left.copy(dst)
	n.right.copy(dst[n.left.length:])
}

func (n *Node[V]) copySlice(dst []V, start, end int) {
	if start >= end {
		return
	}

	if n.kind == tLeaf {
		copy(dst, n.value[start:end])
		return
	}

	leftStart, leftEnd := bound(start, end, n.left.length)
	n.left.copySlice(dst, leftStart, leftEnd)

	rightStart, rightEnd := bound(start-n.left.length, end-n.left.length, n.right.length)
	n.right.copySlice(dst[leftEnd-leftStart:], rightStart, rightEnd)
}

func (n *Node[V]) adjust() {
	if n.kind == tLeaf && n.length > SplitLength {
		pivot := n.length / 2
		n.left = New(n.value[:pivot])
		n.right = New(n.value[pivot:])
		n.value = nil
		n.kind = tNode
	} else if n.kind == tNode && n.length < JoinLength {
		n.value = n.Value()
		n.left = nil
		n.right = nil
		n.kind = tLeaf
	}
}

// Bounds the start and end indices to a given length.
func bound(start, end int, length int) (newStart, newEnd int) {
	if start < 0 {
		start = 0
	} else if start > length {
		start = length
	}

	if end < 0 {
		end = 0
	} else if end > length {
		end = length
	}

	if start >= end {
		return 0, 0
	}
	return start, end
}

// Cannot modify slices, so no tricks are possible
func insert[V any](slice []V, k int, insertion []V) []V {
	result := make([]V, len(slice)+len(insertion))
	copy(result, slice[:k])
	copy(result[k:], insertion)
	copy(result[k+len(insertion):], slice[k:])
	return result
}

func remove[V any](slice []V, start, end int) []V {
	result := make([]V, len(slice)-(end-start))
	copy(result, slice[:start])
	copy(result[start:], slice[end:])
	return result
}
