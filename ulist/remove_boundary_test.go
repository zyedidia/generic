package ulist_test

import (
	"reflect"
	"testing"

	"github.com/zyedidia/generic/ulist"
)

func TestRemoveContinuesAcrossBlockBoundary(t *testing.T) {
	for _, tc := range []struct {
		name                   string
		capacity, count, index int
		expected               int
		hasNext                bool
	}{
		{"full first block tail", 2, 4, 1, 2, true},
		{"full middle block tail", 3, 8, 5, 6, true},
		{"single-entry block control", 1, 3, 1, 2, true},
		{"within block control", 4, 4, 1, 2, true},
		{"last entry control", 3, 8, 7, 0, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			l := ulist.New[int](tc.capacity)
			for i := 0; i < tc.count; i++ {
				l.PushBack(i)
			}
			it := l.Begin()
			for i := 0; i < tc.index; i++ {
				it.Next()
			}
			l.Remove(it)
			if it.IsValid() != tc.hasNext {
				t.Fatalf("iterator valid=%v want%v", it.IsValid(), tc.hasNext)
			}
			if tc.hasNext && it.Get() != tc.expected {
				t.Fatalf("next=%d want%d", it.Get(), tc.expected)
			}
			if l.Size() != tc.count-1 {
				t.Fatal("size mismatch")
			}
		})
	}
}

func TestRemovePartiallyFilledBlockTail(t *testing.T) {
	l := ulist.New[int](3)
	for i := 0; i < 7; i++ {
		l.PushBack(i)
	}
	it := l.Begin()
	l.Remove(it)
	it.Next()
	l.Remove(it)
	if !it.IsValid() || it.Get() != 3 {
		t.Fatal("lost next block after deleting tail of a partially filled block")
	}
}

func TestRemoveWhileFilteringAcrossBlocks(t *testing.T) {
	l := ulist.New[int](2)
	for i := 0; i < 8; i++ {
		l.PushBack(i)
	}
	for it := l.Begin(); it.IsValid(); {
		if it.Get()%2 == 1 {
			l.Remove(it)
		} else {
			it.Next()
		}
	}
	got := []int{}
	for it := l.Begin(); it.IsValid(); it.Next() {
		got = append(got, it.Get())
	}
	if !reflect.DeepEqual(got, []int{0, 2, 4, 6}) || l.Size() != 4 {
		t.Fatalf("filter stopped early: %v size%d", got, l.Size())
	}
}
