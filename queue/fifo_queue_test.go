package queue

import (
	"fmt"
	"testing"

	"github.com/zyedidia/generic/list"
)

func TestFIFOQueueEmpty(t *testing.T) {
	cases := []struct {
		name  string
		queue *FIFOQueue[int]
		want  bool
	}{
		{
			name:  "empty queue",
			queue: emptyQueue(),
			want:  true,
		},
		{
			name:  "non-empty queue",
			queue: nonEmptyQueue(),
			want:  false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := c.queue.Empty()

			if got != c.want {
				t.Errorf("got %v, want %v", got, c.want)
			}
		})
	}
}

func TestFIFOQueuePeek(t *testing.T) {
	t.Run("panics on empty queue", func(t *testing.T) {
		defer func() {
			if err := recover(); err == nil {
				t.Error("peeking on empty queue did not panic")
			}
		}()

		emptyQueue().Peek()
	})

	t.Run("non-empty queue", func(t *testing.T) {
		got := nonEmptyQueue().Peek()
		want := 1

		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})
}

func TestFIFOQueueEnqueue(t *testing.T) {
	t.Run("empty queue", func(t *testing.T) {
		q := emptyQueue()

		q.Enqueue(1)

		if q.list.Front == nil {
			t.Error("front is nil")
		}

		if q.list.Front.Value != 1 {
			t.Errorf("got %v, want %v for front value", q.list.Front.Value, 1)
		}

		if q.list.Back == nil {
			t.Error("back is nil")
		}

		if q.list.Front != q.list.Back {
			t.Error(("front and back are not the same"))
		}
	})

	t.Run("non-empty queue", func(t *testing.T) {
		q := nonEmptyQueue()

		q.Enqueue(3)

		if q.list.Front.Value != 1 {
			t.Errorf("got %v, want %v for front value", q.list.Front.Value, 1)
		}

		if q.list.Back == nil {
			t.Error("back is nil")
		}

		if q.list.Back.Value != 3 {
			t.Errorf("got %v, want %v for back value", q.list.Back.Value, 3)
		}
	})
}

func TestFIFOQueueDequeue(t *testing.T) {
	t.Run("panics on empty queue", func(t *testing.T) {
		defer func() {
			if err := recover(); err == nil {
				t.Error("dequeue on empty queue did not panic")
			}
		}()

		emptyQueue().Dequeue()
	})

	t.Run("non-empty queue", func(t *testing.T) {
		q := nonEmptyQueue()

		// Non-empty after dequeue
		got := q.Dequeue()

		if got != 1 {
			t.Errorf("got %v, want %v", got, 1)
		}

		if q.list.Front.Value == 1 {
			t.Error("front of queue is still 1 after dequeue")
		}

		if q.Empty() {
			t.Error("queue is empty")
		}

		// Empty after dequeue
		got = q.Dequeue()

		if got != 2 {
			t.Errorf("got %v, want %v", got, 2)
		}

		if !q.Empty() {
			t.Error("queue is not empty")
		}
	})
}

func TestFIFOQueueIter(t *testing.T) {
	q := nonEmptyQueue()

	i := 1
	q.Iter().For(func(item int) {
		if item != i {
			t.Errorf("got %v, want %v", item, i)
		}
		i++
	})
}

func ExampleFIFOQueue_Enqueue() {
	q := New[int]()
	q.Enqueue(1)
}

func ExampleFIFOQueue_Peek() {
	q := New[int]()
	q.Enqueue(1)

	fmt.Println(q.Peek())
	// Output: 1
}

func ExampleFIFOQueue_Dequeue() {
	q := New[int]()
	q.Enqueue(1)

	fmt.Println(q.Dequeue())
	// Output: 1
}

func ExampleFIFOQueue_Iter() {
	q := New[int]()
	q.Enqueue(1)
	q.Enqueue(2)

	q.Iter().For(func(i int) {
		fmt.Println(i)
	})
	// Output:
	// 1
	// 2
}

func ExampleFIFOQueue_Empty_empty() {
	q := New[int]()

	fmt.Println(q.Empty())
	// Output: true
}

func ExampleFIFOQueue_Empty_nonempty() {
	q := New[int]()
	q.Enqueue(1)

	fmt.Println(q.Empty())
	// Output: false
}

func emptyQueue() *FIFOQueue[int] {
	return New[int]()
}

func nonEmptyQueue() *FIFOQueue[int] {
	q := New[int]()
	q.list.Front = &list.Node[int]{Value: 1}
	q.list.Front.Next = &list.Node[int]{Value: 2, Prev: q.list.Front}
	q.list.Back = q.list.Front.Next
	return q
}
