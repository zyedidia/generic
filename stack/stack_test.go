package stack_test

import (
	"fmt"
	"testing"

	"github.com/zyedidia/generic/stack"
)

func assert(t *testing.T, fn func() bool) {
	if !fn() {
		t.Fatal("assert failed")
	}
}

func TestSimple(t *testing.T) {
	st := stack.New[int]()
	st.Push(0)
	assert(t, func() bool { return st.Peek() == 0 })
	st.Push(42)
	assert(t, func() bool { return st.Pop() == 42 })
	assert(t, func() bool { return st.Pop() == 0 })
	assert(t, func() bool { return st.Size() == 0 })
	assert(t, func() bool { return st.Pop() == 0 })
	assert(t, func() bool { return st.Peek() == 0 })

	_, okPeek := st.TryPeek()
	assert(t, func() bool { return !okPeek })

	_, okPop := st.TryPop()
	assert(t, func() bool { return !okPop })

	st.Push(99)
	valPeek, okPeek2 := st.TryPeek()
	assert(t, func() bool { return okPeek2 && valPeek == 99 })

	valPop, okPop2 := st.TryPop()
	assert(t, func() bool { return okPop2 && valPop == 99 })

	st.Push(0)
	peekZero, okPeekZero := st.TryPeek()
	assert(t, func() bool { return okPeekZero && peekZero == 0 })

	popZero, okPopZero := st.TryPop()
	assert(t, func() bool { return okPopZero && popZero == 0 })

	_, okEmptyPeek := st.TryPeek()
	assert(t, func() bool { return !okEmptyPeek })

	_, okEmptyPop := st.TryPop()
	assert(t, func() bool { return !okEmptyPop })
}

func Example() {
	st := stack.New[string]()
	st.Push("foo")
	st.Push("bar")

	fmt.Println(st.Pop())
	fmt.Println(st.Peek())

	st.Push("baz")
	fmt.Println(st.Size())
	// Output:
	// bar
	// foo
	// 2
}
