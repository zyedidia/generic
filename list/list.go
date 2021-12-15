package list

type List[V any] struct {
	Head, Tail *Node[V]
}

type Node[V any] struct {
	Value        V
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
