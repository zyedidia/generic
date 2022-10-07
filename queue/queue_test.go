package queue

import (
	"fmt"
	"testing"

	"github.com/zyedidia/generic/list"
)

func TestQueueEmpty(t *testing.T) {
	cases := []struct {
		name  string
		queue *Queue[int]
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

func TestQueuePeek(t *testing.T) {
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

func TestQueueTryPeek(t *testing.T) {
	t.Run("panics on empty queue", func(t *testing.T) {
		value, got := emptyQueue().TryPeek()
		want := false

		if got != want {
			t.Errorf("got %v, want %v; unexpected value: %v", got, want, value)
		}
	})

	t.Run("non-empty queue", func(t *testing.T) {
		gotValue, gotOk := nonEmptyQueue().TryPeek()
		wantValue := 1
		wantOk := true

		if gotOk != gotOk {
			t.Errorf("got ok %v, want ok %v", gotOk, wantOk)
		}
		if gotValue != wantValue {
			t.Errorf("got value %v, want value %v", gotValue, wantValue)
		}
	})
}

func TestQueueEnqueue(t *testing.T) {
	t.Run("empty queue", func(t *testing.T) {
		q := emptyQueue()

		q.Enqueue(1)

		if q.Len() != 1 {
			t.Errorf("want len 1, got len %d", q.Len())
		}

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

		if q.Len() != 3 {
			t.Errorf("want len 3, got len %d", q.Len())
		}

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

func TestQueueDequeue(t *testing.T) {
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

		if q.Len() != 1 {
			t.Errorf("want len 1 after dequeue, got %d", q.Len())
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

		if q.Len() != 0 {
			t.Errorf("want len 0 after empty, got %d", q.Len())
		}
	})
}

func TestQueueTryDequeue(t *testing.T) {
	t.Run("false on empty queue", func(t *testing.T) {
		value, got := emptyQueue().TryDequeue()
		want := false

		if got != want {
			t.Errorf("got %v, want %v; unexpected value: %v", got, want, value)
		}
	})

	t.Run("non-empty queue", func(t *testing.T) {
		q := nonEmptyQueue()

		// Non-empty after dequeue
		gotValue, gotOk := q.TryDequeue()

		if !gotOk {
			t.Errorf("got ok false, want ok true")
		} else {
			if gotValue != 1 {
				t.Errorf("got %v, want %v", gotValue, 1)
			}

			if q.list.Front.Value == 1 {
				t.Error("front of queue is still 1 after dequeue")
			}

			if q.Len() != 1 {
				t.Errorf("want len 1 after dequeue, got %d", q.Len())
			}

			if q.Empty() {
				t.Error("queue is empty")
			}
		}

		// Empty after dequeue
		gotValue, gotOk = q.TryDequeue()

		if !gotOk {
			t.Errorf("got ok false, want ok true")
		} else {
			if gotValue != 2 {
				t.Errorf("got %v, want %v", gotValue, 2)
			}

			if !q.Empty() {
				t.Error("queue is not empty")
			}

			if q.Len() != 0 {
				t.Errorf("want len 0 after empty, got %d", q.Len())
			}
		}
	})
}

func TestQueueEach(t *testing.T) {
	q := nonEmptyQueue()

	i := 1
	q.Each(func(item int) {
		if item != i {
			t.Errorf("got %v, want %v", item, i)
		}
		i++
	})
}

func TestQueueClear(t *testing.T) {
	cases := []struct {
		name  string
		queue *Queue[int]
	}{
		{
			name:  "empty queue",
			queue: emptyQueue(),
		},
		{
			name:  "non-empty queue",
			queue: nonEmptyQueue(),
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			c.queue.Clear()
			length := c.queue.Len()
			empty := c.queue.Empty()

			if length != 0 {
				t.Errorf("got len %v, want len %v", length, 0)
			}
			if empty != true {
				t.Errorf("got empty %v, want empty %v", empty, true)
			}
		})
	}
}

func ExampleQueue_Enqueue() {
	q := New[int]()
	q.Enqueue(1)
}

func ExampleQueue_Peek() {
	q := New[int]()
	q.Enqueue(1)

	fmt.Println(q.Peek())
	// Output: 1
}

func ExampleQueue_Dequeue() {
	q := New[int]()
	q.Enqueue(1)

	fmt.Println(q.Dequeue())
	// Output: 1
}

func Example() {
	q := New[int]()
	q.Enqueue(1)
	q.Enqueue(2)

	q.Each(func(i int) {
		fmt.Println(i)
	})
	// Output:
	// 1
	// 2
}

func ExampleQueue_Empty_empty() {
	q := New[int]()

	fmt.Println(q.Empty())
	// Output: true
}

func ExampleQueue_Empty_nonempty() {
	q := New[int]()
	q.Enqueue(1)

	fmt.Println(q.Empty())
	// Output: false
}

func emptyQueue() *Queue[int] {
	return New[int]()
}

func nonEmptyQueue() *Queue[int] {
	q := New[int]()
	q.list.Front = &list.Node[int]{Value: 1}
	q.list.Front.Next = &list.Node[int]{Value: 2, Prev: q.list.Front}
	q.list.Back = q.list.Front.Next
	q.length = 2
	return q
}
