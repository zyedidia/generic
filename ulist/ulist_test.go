package ulist

import (
	"fmt"
	"reflect"
	"runtime/debug"
	"testing"
	"unsafe"
)

func TestUList(t *testing.T) {
	entriesPerBlock := int(64 / unsafe.Sizeof(int(1)))
	ul := New[int](entriesPerBlock)

	// Constructor sanity test.
	checkEq(t, ul.entriesPerBlock, entriesPerBlock)
	checkEq(t, ul.ll.Front, nil)
	checkEq(t, ul.ll.Back, nil)

	expectedNumEntries := 0
	expectedNumBlocks := 0
	checkEq(t, ul.Size(), expectedNumEntries)
	checkEq(t, getNumUListEntries(ul), expectedNumEntries)
	checkEq(t, getNumUListBlocks(ul), expectedNumBlocks)

	// PushBack.
	for i := 0; i < entriesPerBlock+2; i++ {
		ul.PushBack(i)
		expectedNumEntries++
		if i%entriesPerBlock == 0 {
			expectedNumBlocks++
		}
		checkEq(t, ul.Size(), expectedNumEntries)
		checkEq(t, getNumUListEntries(ul), expectedNumEntries)
		checkEq(t, getNumUListBlocks(ul), expectedNumBlocks)
	}

	// PushFront.
	for i := 0; i < entriesPerBlock+2; i++ {
		ul.PushFront(i)
		expectedNumEntries++
		if i%entriesPerBlock == 0 {
			expectedNumBlocks++
		}
		checkEq(t, ul.Size(), expectedNumEntries)
		checkEq(t, getNumUListEntries(ul), expectedNumEntries)
		checkEq(t, getNumUListBlocks(ul), expectedNumBlocks)
	}

	// Validate entries.
	checkEq(t, getSlice(ul), []int{9, 8, 7, 6, 5, 4, 3, 2, 1, 0, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9})
	validateBlockCapacities(t, ul)
}

func checkEq[V any](t *testing.T, a V, b V) {
	//if a != b {
	if !reflect.DeepEqual(a, b) {
		t.Fatalf("got:%v, want:%v \n%s", a, b, debug.Stack())
	}
}

// Helper function that returns the number of entries in the unrolled list.
func getNumUListEntries[V any](ul *UList[V]) int {
	ret := int(0)
	mapper := func(val blockPtr[V]) {
		ret += len(*val)
	}
	ul.ll.Front.Each(mapper)
	return ret
}

// Helper function that returns the number of blocks in the unrolled list.
func getNumUListBlocks[V any](ul *UList[V]) int {
	ret := int(0)
	mapper := func(val blockPtr[V]) {
		ret += 1
	}
	ul.ll.Front.Each(mapper)
	return ret
}

// Helper function to print a debug string of 'ul'.
func getDebugString[V any](ul *UList[V]) string {
	ret := ""
	mapper := func(val blockPtr[V]) {
		ret += fmt.Sprintf("%v ->", *val)
	}
	ul.ll.Front.Each(mapper)
	return ret
}

// Helper function to return ulist as a slice.
func getSlice[V any](ul *UList[V]) []V {
	ret := make([]V, 0)
	mapper := func(val blockPtr[V]) {
		ret = append(ret, *val...)
	}
	ul.ll.Front.Each(mapper)
	return ret
}

// Helper function to check if all blocks in 'ul' are of the expected size.
func validateBlockCapacities[V any](t *testing.T, ul *UList[V]) {
	mapper := func(val blockPtr[V]) {
		checkEq(t, cap(*val), ul.entriesPerBlock)
	}
	ul.ll.Front.Each(mapper)
}
