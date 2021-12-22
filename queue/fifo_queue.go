// Package queue provides an implementation of a First In First Out (FIFO)
// queue. The FIFO queue is implemented using a doubly-linked list found in the
// 'list' package.
package queue

import (
	"github.com/zyedidia/generic/iter"
	"github.com/zyedidia/generic/list"
)

// FIFOQueue is a simple First In First Out (FIFO) queue.
type FIFOQueue[T any] struct {
	list *list.List[T]
}

// New returns an empty First In First Out (FIFO) queue.
func New[T any]() *FIFOQueue[T] {
	return &FIFOQueue[T]{
		list: list.New[T](),
	}
}

// Enqueue inserts 'value' to the end of the queue.
func (q *FIFOQueue[T]) Enqueue(value T) {
	q.list.PushBack(value)
}

// Dequeue removes and returns the item at the front of the queue.
//
// A panic occurs if the queue is Empty.
func (q *FIFOQueue[T]) Dequeue() T {
	value := q.list.Front.Value
	q.list.Remove(q.list.Front)

	return value
}

// Peek returns the item at the front of the queue without removing it.
//
// A panic occurs if the queue is Empty.
func (q *FIFOQueue[T]) Peek() T {
	return q.list.Front.Value
}

// Empty returns true if the queue is empty.
func (q *FIFOQueue[T]) Empty() bool {
	return q.list.Front == nil
}

// Iter returns a forward iterator, returning items starting from the front to
// the back of the queue.
func (q *FIFOQueue[T]) Iter() iter.Iter[T] {
	return q.list.Front.Iter()
}
