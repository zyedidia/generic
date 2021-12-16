package list

import "github.com/zyedidia/generic/iter"

type List[V any] struct {
	Head, Tail *Node[V]
}

type Node[V any] struct {
	Value      V
	Prev, Next *Node[V]
}

func (l *List[V]) Append(n *Node[V]) {
	n.Next = l.Head
	n.Prev = nil
	if l.Head != nil {
		l.Head.Prev = n
	} else {
		l.Tail = n
	}
	l.Head = n
}

func (l *List[V]) Remove(n *Node[V]) {
	if n.Next != nil {
		n.Next.Prev = n.Prev
	} else {
		l.Tail = n.Prev
	}
	if n.Prev != nil {
		n.Prev.Next = n.Next
	} else {
		l.Head = n.Next
	}
}

func (n *Node[V]) Iter() iter.Iter[V] {
	node := n
	return func() (v V, ok bool) {
		if node == nil {
			return v, false
		}
		node = node.Next
		return node.Value, true
	}
}

func (n *Node[V]) ReverseIter() iter.Iter[V] {
	node := n
	return func() (v V, ok bool) {
		if node == nil {
			return v, false
		}
		node = node.Prev
		return node.Value, true
	}
}
