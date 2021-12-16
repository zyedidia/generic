package list

import "github.com/zyedidia/generic/iter"

// List implements a doubly-linked list.
type List[V any] struct {
	Front, Back *Node[V]
}

// Node is a node in the linked list.
type Node[V any] struct {
	Value      V
	Prev, Next *Node[V]
}

// New returns an empty linked list.
func New[V any]() *List[V] {
	return &List[V]{}
}

// PushBack adds 'v' to the end of the list.
func (l *List[V]) PushBack(v V) {
	l.PushBackNode(&Node[V]{
		Value: v,
	})
}

// PushFront adds 'v' to the beginning of the list.
func (l *List[V]) PushFront(v V) {
	l.PushFrontNode(&Node[V]{
		Value: v,
	})
}

// PushBackNode adds the node 'n' to the back of the list.
func (l *List[V]) PushBackNode(n *Node[V]) {
	n.Next = nil
	n.Prev = l.Back
	if l.Back != nil {
		l.Back.Next = n
	} else {
		l.Front = n
	}
	l.Back = n
}

// PushFrontNode adds the node 'n' to the front of the list.
func (l *List[V]) PushFrontNode(n *Node[V]) {
	n.Next = l.Front
	n.Prev = nil
	if l.Front != nil {
		l.Front.Prev = n
	} else {
		l.Back = n
	}
	l.Front = n
}

// Remove removes the node 'n' from the list.
func (l *List[V]) Remove(n *Node[V]) {
	if n.Next != nil {
		n.Next.Prev = n.Prev
	} else {
		l.Back = n.Prev
	}
	if n.Prev != nil {
		n.Prev.Next = n.Next
	} else {
		l.Front = n.Next
	}
}

// Iter returns a forward iterator, going from front to back starting at node
// 'n'.
func (n *Node[V]) Iter() iter.Iter[V] {
	node := n
	return func() (v V, ok bool) {
		if node == nil {
			return v, false
		}
		v = node.Value
		node = node.Next
		return v, true
	}
}

// Iter returns a reverse iterator, going from back to front starting at node
// 'n'.
func (n *Node[V]) ReverseIter() iter.Iter[V] {
	node := n
	return func() (v V, ok bool) {
		if node == nil {
			return v, false
		}
		v = node.Value
		node = node.Prev
		return v, true
	}
}
