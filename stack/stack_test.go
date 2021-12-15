package stack_test

import (
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
}
