package ulist

import (
	"testing"
	"unsafe"
)

func TestUListIter(t *testing.T) {
	entriesPerBlock := int(64 / unsafe.Sizeof(int(1)))
	ul := New[int](entriesPerBlock)

	// Constructor sanity test.
	iter := ul.Begin()
	checkEq(t, iter.node, nil)
	checkEq(t, iter.index, 0)
	checkEq(t, iter.IsValid(), false)

	iter = ul.End()
	checkEq(t, iter.node, nil)
	checkEq(t, iter.index, 0)
	checkEq(t, iter.IsValid(), false)

	ul.PushBack(1)
	ul.PushBack(2)
	ul.PushBack(3)

	// Sanity test iterator of a non empty ulist.
	iter = ul.Begin()
	checkEq(t, iter.node, ul.ll.Front)
	checkEq(t, iter.index, 0)
	checkEq(t, iter.IsValid(), true)

	iter = ul.End()
	checkEq(t, iter.node, ul.ll.Back)
	checkEq(t, iter.index, 2)
	checkEq(t, iter.IsValid(), true)
}

func TestUListIterIteration(t *testing.T) {
	entriesPerBlock := int(64 / unsafe.Sizeof(int(1)))
	ul := New[int](entriesPerBlock)
	for i := 0; i < 10; i++ {
		ul.PushBack(i)
	}

	// Iterate begin to end.
	ret := make([]int, 0)
	for iter := ul.Begin(); iter.IsValid(); iter.Next() {
		v := iter.Get()
		ret = append(ret, v)
	}
	checkEq(t, ret, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9})

	// Iterate end to begin.
	ret = make([]int, 0)
	for iter := ul.End(); iter.IsValid(); iter.Prev() {
		v := iter.Get()
		ret = append(ret, v)
	}
	checkEq(t, ret, []int{9, 8, 7, 6, 5, 4, 3, 2, 1, 0})

	// Bi-directional iteration.
	iter := ul.Begin()
	checkEq(t, iter.IsValid(), true)
	checkEq(t, iter.Get(), 0)
	checkEq(t, iter.Next(), true)
	checkEq(t, iter.Get(), 1)
	checkEq(t, iter.Next(), true)
	checkEq(t, iter.Get(), 2)
	checkEq(t, iter.Prev(), true)
	checkEq(t, iter.Get(), 1)
	checkEq(t, iter.Prev(), true)
	checkEq(t, iter.Get(), 0)
	checkEq(t, iter.Prev(), false)

	// Seek to the end of the block.
	for i := 0; i < entriesPerBlock; i++ {
		checkEq(t, iter.Next(), true)
	}
	checkEq(t, iter.Get(), 7)

	// Straddle back and forth across block boundaries.
	checkEq(t, iter.Next(), true)
	checkEq(t, iter.Get(), 8)
	checkEq(t, iter.Prev(), true)
	checkEq(t, iter.Get(), 7)

	validateBlockCapacities(t, ul)
}

func TestUListIterRecovery(t *testing.T) {
	entriesPerBlock := int(64 / unsafe.Sizeof(int(1)))
	ul := New[int](entriesPerBlock)
	for i := 0; i < 5; i++ {
		ul.PushBack(i)
	}

	// No matter how far we iterate forward, we should be able to jump back
	// to the end element by calling Prev().
	iter := ul.Begin()
	for i := 0; i < ul.Size()*2; i++ {
		iter.Next()
	}
	checkEq(t, iter.IsValid(), false)
	checkEq(t, iter.Prev(), true)
	checkEq(t, iter.Get(), 4)

	// Test the same when iterating in reverse.
	iter = ul.End()
	for i := 0; i < ul.Size()*2; i++ {
		iter.Prev()
	}
	checkEq(t, iter.IsValid(), false)
	checkEq(t, iter.Next(), true)
	checkEq(t, iter.Get(), 0)
}

func TestUListIterDeletes(t *testing.T) {
	entriesPerBlock := int(64 / unsafe.Sizeof(int(1)))
	ul := New[int](entriesPerBlock)
	for i := 0; i < 10; i++ {
		ul.PushBack(i)
	}

	// ul: [0,1,...,7] -> [8,9]
	expectedNumEntries := 10
	expectedNumBlocks := 2
	checkEq(t, ul.Size(), expectedNumEntries)
	checkEq(t, getNumUListEntries(ul), expectedNumEntries)
	checkEq(t, getNumUListBlocks(ul), expectedNumBlocks)
	validateBlockCapacities(t, ul)

	// Deleting the last element should invalidate the iterator.
	// ul: [0,1,...,7] -> [8,_]
	//                       ^
	iter := ul.End()
	checkEq(t, iter.Get(), 9)
	ul.Remove(iter)
	checkEq(t, iter.IsValid(), false)
	expectedNumEntries--
	validateBlockCapacities(t, ul)

	// Deleting the last element should get rid of the block.
	// ul: [0,1,...,7]
	checkEq(t, iter.Prev(), true)
	checkEq(t, iter.Get(), 8)
	ul.Remove(iter)
	checkEq(t, iter.IsValid(), false)
	checkEq(t, iter.Prev(), true)
	checkEq(t, iter.Get(), 7)

	expectedNumEntries--
	expectedNumBlocks--
	checkEq(t, ul.Size(), expectedNumEntries)
	checkEq(t, getNumUListEntries(ul), expectedNumEntries)
	checkEq(t, getNumUListBlocks(ul), expectedNumBlocks)
	validateBlockCapacities(t, ul)

	// Empty the ulist from begin.
	for iter := ul.Begin(); iter.IsValid(); ul.Remove(iter) {
	}
	expectedNumEntries = 0
	expectedNumBlocks = 0
	checkEq(t, ul.Size(), expectedNumEntries)
	checkEq(t, getNumUListEntries(ul), expectedNumEntries)
	checkEq(t, getNumUListBlocks(ul), expectedNumBlocks)
	validateBlockCapacities(t, ul)
}

func TestUListIterAdd(t *testing.T) {
	entriesPerBlock := int(64 / unsafe.Sizeof(int(1)))
	ul := New[int](entriesPerBlock)
	for i := 0; i < entriesPerBlock-1; i++ {
		ul.PushBack(i)
	}

	// ul: [0,1,...,6]
	iter := ul.End()
	checkEq(t, iter.Get(), 6)

	expectedNumEntries := 7
	expectedNumBlocks := 1
	checkEq(t, ul.Size(), expectedNumEntries)
	checkEq(t, getNumUListEntries(ul), expectedNumEntries)
	checkEq(t, getNumUListBlocks(ul), expectedNumBlocks)
	checkEq(t, getSlice(ul), []int{0, 1, 2, 3, 4, 5, 6})

	// Add entries.
	ul.AddAfter(iter, 7)
	checkEq(t, iter.Get(), 7)
	ul.AddAfter(iter, 8)
	checkEq(t, iter.Get(), 8)

	// ul: [0,1,...,6,7] -> [8]
	expectedNumEntries = 9
	expectedNumBlocks = 2
	checkEq(t, ul.Size(), expectedNumEntries)
	checkEq(t, getNumUListEntries(ul), expectedNumEntries)
	checkEq(t, getNumUListBlocks(ul), expectedNumBlocks)
	checkEq(t, getSlice(ul), []int{0, 1, 2, 3, 4, 5, 6, 7, 8})

	// Add to an already full block.
	// ul: [0,-1,1,...,6] -> [7,8]
	iter = ul.Begin()
	ul.AddAfter(iter, -1)
	checkEq(t, iter.Get(), -1)
	expectedNumEntries++
	checkEq(t, ul.Size(), expectedNumEntries)
	checkEq(t, getNumUListEntries(ul), expectedNumEntries)
	checkEq(t, getNumUListBlocks(ul), expectedNumBlocks)
	checkEq(t, getSlice(ul), []int{0, -1, 1, 2, 3, 4, 5, 6, 7, 8})

	// Add to an already full block, such that the new element overflows.
	// ul: [0,-1,1,...,6] -> [-6,7,8]
	iter = ul.End()
	iter.Prev()
	iter.Prev()
	checkEq(t, iter.Get(), 6)
	ul.AddAfter(iter, -6)
	checkEq(t, iter.Get(), -6)
	expectedNumEntries++
	checkEq(t, ul.Size(), expectedNumEntries)
	checkEq(t, getNumUListEntries(ul), expectedNumEntries)
	checkEq(t, getNumUListBlocks(ul), expectedNumBlocks)
	checkEq(t, getSlice(ul), []int{0, -1, 1, 2, 3, 4, 5, 6, -6, 7, 8})

	validateBlockCapacities(t, ul)

	// Test AddBefore on the first element.
	iter = ul.Begin()
	ul.AddBefore(iter, 100)
	checkEq(t, iter.Get(), 100)
	expectedNumEntries++
	expectedNumBlocks++

	// Test AddBefore in the middle.
	checkEq(t, iter.Next(), true)
	checkEq(t, iter.Next(), true)
	checkEq(t, iter.Next(), true)
	checkEq(t, iter.Get(), 1)
	ul.AddBefore(iter, 111)
	checkEq(t, iter.Get(), 111)
	expectedNumEntries++

	// ul: [100] -> [0,-1,111,1,...,5] -> [6,-6,7,8]
	checkEq(t, ul.Size(), expectedNumEntries)
	checkEq(t, getNumUListEntries(ul), expectedNumEntries)
	checkEq(t, getNumUListBlocks(ul), expectedNumBlocks)
	checkEq(t, getSlice(ul), []int{100, 0, -1, 111, 1, 2, 3, 4, 5, 6, -6, 7, 8})

	validateBlockCapacities(t, ul)
}
