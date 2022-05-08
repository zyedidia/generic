package heap_test

import (
	"fmt"
	"testing"

	"github.com/zyedidia/generic/heap"
)

func TestCreateHeapFromSlice(t *testing.T) {
	cases := []struct {
		name   string
		array  []int
		sorted []int
		less   func(int, int) bool
	}{
		{
			name:  "empty",
			array: []int{},
			less:  func(a, b int) bool { return a < b },
		},
		{
			name:   "non-empty (minheap)",
			array:  []int{5, 3, 6, 2, 4, 1},
			sorted: []int{1, 2, 3, 4, 5, 6},
			less:   func(a, b int) bool { return a < b },
		},
		{
			name:   "non-empty (maxheap)",
			array:  []int{5, 7, 9, 2, -1},
			sorted: []int{9, 7, 5, 2, -1},
			less:   func(a, b int) bool { return a > b },
		},
	}

	for _, c := range cases {
		t.Run(c.name+" From", func(t *testing.T) {
			// copy array
			array := make([]int, len(c.array))
			copy(array, c.array)

			heap := heap.From(c.less, array...)

			for i, v := range c.sorted {
				peek, ok := heap.Pop()
				if !ok {
					t.Errorf("pop not ok, idx: %v", i)
				}
				if peek != v {
					t.Errorf("peek not equal, idx: %v", i)
				}
			}
		})

		t.Run(c.name+" FromArray", func(t *testing.T) {
			heap := heap.FromSlice(c.less, c.array)

			for i, v := range c.sorted {
				peek, ok := heap.Pop()
				if !ok {
					t.Errorf("pop not ok, idx: %v", i)
				}
				if peek != v {
					t.Errorf("peek not equal, idx: %v", i)
				}
			}
		})
	}
}

func TestBasic(t *testing.T) {
	cases := []struct {
		name        string
		data        []int
		peeksInPush []int
		peeksInPop  []int
		less        func(int, int) bool
	}{
		{
			name:        "minheap",
			data:        []int{5, 7, 9, 2, -1},
			peeksInPush: []int{5, 5, 5, 2, -1},
			peeksInPop:  []int{-1, 2, 5, 7, 9},
			less:        func(a, b int) bool { return a < b },
		},
		{
			name:        "maxheap",
			data:        []int{5, 3, 6, 2, 4, 1},
			peeksInPush: []int{5, 5, 6, 6, 6, 6},
			peeksInPop:  []int{6, 5, 4, 3, 2, 1},
			less:        func(a, b int) bool { return a > b },
		},
		{
			name:        "minheap",
			data:        []int{9, 10, 8, 7, 9, 4, 8, 4, -1, -1, 2, 3, 5},
			peeksInPush: []int{9, 9, 8, 7, 7, 4, 4, 4, -1, -1, -1, -1, -1},
			peeksInPop:  []int{-1, -1, 2, 3, 4, 4, 5, 7, 8, 8, 9, 9, 10},
			less:        func(a, b int) bool { return a < b },
		},
		{
			name:        "maxheap",
			data:        []int{1, 5, 4, 7, 0, 8, 4, -1, 9, 2, 3, 5},
			peeksInPush: []int{1, 5, 5, 7, 7, 8, 8, 8, 9, 9, 9, 9},
			peeksInPop:  []int{9, 8, 7, 5, 5, 4, 4, 3, 2, 1, 0, -1},
			less:        func(a, b int) bool { return a >= b },
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			heap := heap.New(c.less)

			if heap.Size() != 0 {
				t.Errorf("heap len not 0")
			}

			// push all elements
			for i, v := range c.data {
				heap.Push(v)

				peek, ok := heap.Peek()
				if !ok {
					t.Errorf("peek not ok, idx: %v", i)
				}
				if peek != c.peeksInPush[i] {
					t.Errorf("peek not equal, idx: %v", i)
				}
			}

			if heap.Size() != len(c.data) {
				t.Errorf("heap len not equal to data len")
			}

			// pop all elements
			for idx, peek := range c.peeksInPop {
				v, ok := heap.Pop()
				if !ok {
					t.Errorf("pop not ok, idx: %v", idx)
				}
				if v != peek {
					t.Errorf("peek not equal, idx: %v", idx)
				}
			}

			if heap.Size() != 0 {
				t.Errorf("heap len not 0")
			}
		})
	}
}

func TestPopAndPeekOnEmpty(t *testing.T) {
	heap := heap.New(func(a, b int) bool { return a < b })

	var v int
	var ok bool

	_, ok = heap.Peek()
	if ok {
		t.Errorf("peek returns ok on empty heap")
	}

	_, ok = heap.Pop()
	if ok {
		t.Errorf("pop returns ok on empty heap")
	}

	heap.Push(1)

	v, ok = heap.Peek()
	if !ok {
		t.Errorf("peek not ok on non-empty heap")
	}
	if v != 1 {
		t.Errorf("expect peek %v, but got %v", 1, v)
	}

	v, ok = heap.Pop()
	if !ok {
		t.Errorf("pop not ok on non-empty heap")
	}
	if v != 1 {
		t.Errorf("expect pop %v, but got %v", 1, v)
	}
}

func Example() {
	heap := heap.New(func(a, b int) bool { return a < b })

	heap.Push(5)
	heap.Push(2)
	heap.Push(3)

	v, _ := heap.Pop()
	fmt.Println(v)

	v, _ = heap.Peek()
	fmt.Println(v)
	// Output:
	// 2
	// 3
}

func ExampleFrom() {
	heap := heap.From(func(a, b int) bool { return a < b }, 5, 2, 3)

	v, _ := heap.Pop()
	fmt.Println(v)

	v, _ = heap.Peek()
	fmt.Println(v)
	// Output:
	// 2
	// 3
}

func ExampleFromSlice() {
	heap := heap.FromSlice(func(a, b int) bool { return a > b }, []int{-1, 5, 2, 3})

	v, _ := heap.Pop()
	fmt.Println(v)

	v, _ = heap.Peek()
	fmt.Println(v)
	// Output:
	// 5
	// 3
}

func ExampleHeap_Pop() {
	heap := heap.New(func(a, b int) bool { return a < b })

	heap.Push(5)

	v, ok := heap.Pop()
	fmt.Println(v, ok)

	// pop on empty
	v, ok = heap.Pop()
	fmt.Println(v, ok)
	// Output:
	// 5 true
	// 0 false
}
