package interval

import (
	g "github.com/zyedidia/generic"
	"github.com/zyedidia/generic/iter"
)

type KV[V any] struct {
	Key Range
	Val V
}

type Range struct {
	Low, High int
}

func overlaps(i1 Range, low, high int) bool {
	return i1.Low <= high && i1.High >= low
}

type Tree[V any] struct {
	root *node[V]
}

func (t *Tree[V]) Add(key Range, value V) {
	t.root = t.root.add(key, value)
}

func (t *Tree[V]) Overlaps(key Range) []V {
	var result []V
	return t.root.overlaps(key, result)
}

func (t *Tree[V]) Remove(key Range) {
	t.root = t.root.remove(key)
}

func (t *Tree[V]) Search(pos int) (V, bool) {
	n := t.root.search(pos)
	if n == nil {
		var v V
		return v, false
	}
	return n.value, true
}

func (t *Tree[V]) Iter() iter.Iter[KV[V]] {
	return t.root.iter()
}

func (t *Tree[V]) Height() int {
	return t.root.getHeight()
}

type node[V any] struct {
	key   Range
	value V

	max    int
	height int
	left   *node[V]
	right  *node[V]
}

func (n *node[V]) add(key Range, value V) *node[V] {
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

func (n *node[V]) remove(key Range) *node[V] {
	if n == nil {
		return nil
	}
	if key.Low < n.key.Low {
		n.left = n.left.remove(key)
	} else if key.Low > n.key.Low {
		n.right = n.right.remove(key)
	} else {
		if n.left != nil && n.right != nil {
			rightMinNode := n.right.findSmallest()
			n.key = rightMinNode.key
			n.value = rightMinNode.value
			n.right = n.right.remove(rightMinNode.key)
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

func (n *node[V]) overlaps(key Range, result []V) []V {
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

func (n *node[V]) iter() iter.Iter[KV[V]] {
	if n == nil {
		return func() (v KV[V], ok bool) {
			return v, false
		}
	}

	var didself bool
	left := n.left.iter()
	right := n.right.iter()
	return func() (KV[V], bool) {
		v, ok := left()
		if ok {
			return v, true
		} else if !didself {
			didself = true
			return KV[V]{
				Key: n.key,
				Val: n.value,
			}, true
		}
		return right()
	}
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
