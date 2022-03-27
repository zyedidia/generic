// Package queue provides an implementation of a First In First Out (FIFO)
// queue. The FIFO queue is implemented using the doubly-linked list from the
// 'list' package.
package queue

import (
	"github.com/zyedidia/generic/list"
)

// Queue is a simple First In First Out (FIFO) queue.
type Queue[T any] struct {
	list *list.List[T]
}

// New returns an empty First In First Out (FIFO) queue.
func New[T any]() *Queue[T] {
	return &Queue[T]{
		list: list.New[T](),
	}
}

// Enqueue inserts 'value' to the end of the queue.
func (q *Queue[T]) Enqueue(value T) {
	q.list.PushBack(value)
}

// Dequeue removes and returns the item at the front of the queue.
//
// A panic occurs if the queue is Empty.
func (q *Queue[T]) Dequeue() T {
	value := q.list.Front.Value
	q.list.Remove(q.list.Front)

	return value
}

// Peek returns the item at the front of the queue without removing it.
//
// A panic occurs if the queue is Empty.
func (q *Queue[T]) Peek() T {
	return q.list.Front.Value
}

// Empty returns true if the queue is empty.
func (q *Queue[T]) Empty() bool {
	return q.list.Front == nil
}

// Each calls 'fn' on every item in the queue, starting with the least
// recently pushed element.
func (q *Queue[T]) Each(fn func(t T)) {
	q.list.Front.Each(fn)
}
