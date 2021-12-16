package list_test

import (
	"testing"

	"github.com/zyedidia/generic/list"
)

func TestList(t *testing.T) {
	l := list.New[int]()
	l.PushBack(0)
	l.PushBack(1)
	l.PushBack(2)
	l.PushBack(3)

	var s int
	l.Front.Iter().For(func(i int) {
		s += i
	})

	if s != 6 {
		t.Fatal("incorrect sum")
	}
}
