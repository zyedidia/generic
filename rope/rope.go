package rope

import g "github.com/zyedidia/generic"

var (
	// SplitLength is the threshold above which slices will be split into
	// separate nodes.
	SplitLength = 4096 * 4
	// JoinLength is the threshold below which nodes will be merged into
	// slices.
	JoinLength = SplitLength / 2
	// RebalanceRatio is the threshold used to trigger a rebuild during a
	// rebalance operation.
	RebalanceRatio = 1.2
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
// data is not copied so the user should ensure that it is okay to insert and
// delete from the input slice.
func New[V any](b []V) *Node[V] {
	n := &Node[V]{
		kind:   tLeaf,
		value:  b[0:len(b):len(b)],
		length: len(b),
	}
	n.adjust()
	return n
}

// Len returns the number of elements stored in the rope.
func (n *Node[V]) Len() int {
	return n.length
}

func (n *Node[V]) adjust() {
	switch n.kind {
	case tLeaf:
		if n.length > SplitLength {
			divide := n.length / 2
			n.left = New(n.value[:divide])
			n.right = New(n.value[divide:])
			n.value = nil
			n.kind = tNode
			n.length = n.left.length + n.right.length
		}
	case tNode:
		if n.length < JoinLength {
			n.value = n.Value()
			n.left = nil
			n.right = nil
			n.kind = tLeaf
			n.length = len(n.value)
		}
	}
}

// Value returns the elements of this node concatenated into a slice. May
// return the underyling slice without copying, so do not modify the returned
// slice.
func (n *Node[V]) Value() []V {
	switch n.kind {
	case tLeaf:
		return n.value
	case tNode:
		return concat(n.left.Value(), n.right.Value())
	}
	panic("unreachable")
}

// Remove deletes the range [start:end) (exclusive bound) from the rope.
func (n *Node[V]) Remove(start, end int) {
	switch n.kind {
	case tLeaf:
		// slice tricks delete
		n.value = append(n.value[:start], n.value[end:]...)
		n.length = len(n.value)
	case tNode:
		leftLength := n.left.length
		leftStart := g.Min(start, leftLength)
		leftEnd := g.Min(end, leftLength)
		rightLength := n.right.length
		rightStart := g.Max(0, g.Min(start-leftLength, rightLength))
		rightEnd := g.Max(0, g.Min(end-leftLength, rightLength))
		if leftStart < leftLength {
			n.left.Remove(leftStart, leftEnd)
		}
		if rightEnd > 0 {
			n.right.Remove(rightStart, rightEnd)
		}
		n.length = n.left.length + n.right.length
	}
	n.adjust()
}

// Insert inserts the given value at pos.
func (n *Node[V]) Insert(pos int, value []V) {
	switch n.kind {
	case tLeaf:
		// slice tricks insert
		n.value = insert(n.value, pos, value)
		n.length = len(n.value)
	case tNode:
		leftLength := n.left.length
		if pos < leftLength {
			n.left.Insert(pos, value)
		} else {
			n.right.Insert(pos-leftLength, value)
		}
		n.length = n.left.length + n.right.length
	}
	n.adjust()
}

// Slice returns the range of the rope from [start:end). The returned slice
// is not copied.
func (n *Node[V]) Slice(start, end int) []V {
	if start >= end {
		return []V{}
	}

	switch n.kind {
	case tLeaf:
		return n.value[start:end]
	case tNode:
		leftLength := n.left.length
		leftStart := g.Min(start, leftLength)
		leftEnd := g.Min(end, leftLength)
		rightLength := n.right.length
		rightStart := g.Max(0, g.Min(start-leftLength, rightLength))
		rightEnd := g.Max(0, g.Min(end-leftLength, rightLength))

		if leftStart != leftEnd {
			if rightStart != rightEnd {
				return concat(n.left.Slice(leftStart, leftEnd), n.right.Slice(rightStart, rightEnd))
			} else {
				return n.left.Slice(leftStart, leftEnd)
			}
		} else {
			if rightStart != rightEnd {
				return n.right.Slice(rightStart, rightEnd)
			} else {
				return []V{}
			}
		}
	}
	panic("unreachable")
}

// At returns the element at the given position.
func (n *Node[V]) At(pos int) V {
	s := n.Slice(pos, pos+1)
	return s[0]
}

// SplitAt splits the node at the given index and returns two new ropes
// corresponding to the left and right portions of the split.
func (n *Node[V]) SplitAt(i int) (*Node[V], *Node[V]) {
	switch n.kind {
	case tLeaf:
		return New(n.value[:i]), New(n.value[i:])
	case tNode:
		m := n.left.length
		if i == m {
			return n.left, n.right
		} else if i < m {
			l, r := n.left.SplitAt(i)
			return l, join(r, n.right)
		}
		l, r := n.right.SplitAt(i - m)
		return join(n.left, l), r
	}
	panic("unreachable")
}

func join[V any](l, r *Node[V]) *Node[V] {
	n := &Node[V]{
		left:   l,
		right:  r,
		length: l.length + r.length,
		kind:   tNode,
	}
	n.adjust()
	return n
}

// Join merges all the given ropes together into one rope.
func Join[V any](a, b *Node[V], more ...*Node[V]) *Node[V] {
	s := join(a, b)
	for _, n := range more {
		s = join(s, n)
	}
	return s
}

// Rebuild rebuilds the entire rope structure, resulting in a balanced tree.
func (n *Node[V]) Rebuild() {
	switch n.kind {
	case tNode:
		n.value = concat(n.left.Value(), n.right.Value())
		n.left = nil
		n.right = nil
		n.adjust()
	}
}

// Rebalance finds unbalanced nodes and rebuilds them.
func (n *Node[V]) Rebalance() {
	switch n.kind {
	case tNode:
		lratio := float64(n.left.length) / float64(n.right.length)
		rratio := float64(n.right.length) / float64(n.left.length)
		if lratio > RebalanceRatio || rratio > RebalanceRatio {
			n.Rebuild()
		} else {
			n.left.Rebalance()
			n.right.Rebalance()
		}
	}
}

// Each applies the given function to every leaf node in order.
func (n *Node[V]) Each(fn func(n *Node[V]) bool) bool {
	switch n.kind {
	case tLeaf:
		return fn(n)
	default: // case tNode
		if n.left.Each(fn) {
			return true
		}
		return n.right.Each(fn)
	}
}

// from slice tricks
func insert[V any](s []V, k int, vs []V) []V {
	if n := len(s) + len(vs); n <= cap(s) {
		s2 := s[:n]
		copy(s2[k+len(vs):], s[k:])
		copy(s2[k:], vs)
		return s2
	}
	s2 := make([]V, len(s)+len(vs))
	copy(s2, s[:k])
	copy(s2[k:], vs)
	copy(s2[k+len(vs):], s[k:])
	return s2
}

func concat[V any](a, b []V) []V {
	c := make([]V, 0, len(a)+len(b))
	c = append(c, a...)
	c = append(c, b...)
	return c
}
