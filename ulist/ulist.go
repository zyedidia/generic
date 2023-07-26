package ulist

import (
	"github.com/zyedidia/generic/list"
)

//          ---------    ---------
//  UList:  | Block | <-> | Block | <-> ...
//          ---------    ---------
// A UList is represented internally as a list.List (doubly linked list)
// of pointers to blocks of entries.
//
// Block:  []V{ ... }
// A Block is a slice of entries and forms the Node in the list.List.

// Type alias for a pointer to a block of entries.
type blockPtr[V any] *[]V

// UList implements a doubly-linked unolled list.
type UList[V any] struct {
	ll              list.List[blockPtr[V]]
	entriesPerBlock int
	size            int
}

// New returns an empty unrolled linked list.
// 'entriesPerBlock' is the number of entries to store in each block.
// This value should ideally be the size of a cache-line or multiples there-of.
// See: https://en.wikipedia.org/wiki/Unrolled_linked_list
func New[V any](entriesPerBlock int) *UList[V] {
	return &UList[V]{
		ll:              *list.New[blockPtr[V]](),
		entriesPerBlock: entriesPerBlock,
		size:            0,
	}
}

// Size returns the number of entries in 'ul'.
func (ul *UList[V]) Size() int {
	return ul.size
}

// PushBack adds 'v to the end of the ulist.
func (ul *UList[V]) PushBack(v V) {
	if !hasCapacity[V](ul.ll.Back) {
		ul.ll.PushBack(ul.newBlock())
	}
	blkPtr := ul.ll.Back.Value
	*blkPtr = append(*blkPtr, v)
	ul.size++
}

// PushFront adds 'v to the beginning of the ulist.
func (ul *UList[V]) PushFront(v V) {
	if !hasCapacity[V](ul.ll.Front) {
		ul.ll.PushFront(ul.newBlock())
	}
	ul.prependToBlock(v, ul.ll.Front.Value)
	ul.size++
}

// Begin returns an UListIter pointing to the first entry in the UList.
func (ul *UList[V]) Begin() *UListIter[V] {
	return newIterFront(ul)
}

// End returns an UListIter pointing to the last entry in the UList.
func (ul *UList[V]) End() *UListIter[V] {
	return newIterBack(ul)
}

// AddAfter adds 'v' to 'ul' after the entry pointed to by 'iter'.
// 'iter' is expected to be valid, i.e. iter->IsValid() == true.
// 'iter' is updated to now point to the new entry added, such that
// iter->Get() == 'v'.
func (ul *UList[V]) AddAfter(iter *UListIter[V], v V) {
	ul.size++
	blkPtr := iter.node.Value
	// Adding to a block with spare capacity.
	if hasCapacity(iter.node) {
		iter.index++
		*blkPtr = append((*blkPtr)[:iter.index+1], (*blkPtr)[iter.index:]...)
		(*blkPtr)[iter.index] = v
		return
	}
	// Adding to an already full block.
	if iter.index == len(*blkPtr)-1 {
		// When adding to the end of a block, 'v' is the overflow.
		iter.addOverflowToNextBlock(ul, v)
		iter.Next()
		return
	}
	// When adding 'v' in the middle, the last entry in the block is the overflow.
	overflow := (*blkPtr)[len(*blkPtr)-1]
	iter.addOverflowToNextBlock(ul, overflow)
	iter.index++
	// Slide entries beyond the write index right by one spot and write the value.
	*blkPtr = append((*blkPtr)[:iter.index+1], (*blkPtr)[iter.index:len(*blkPtr)-1]...)
	(*blkPtr)[iter.index] = v
}

// AddBefore adds 'v' to 'ul' before the entry pointed to by 'iter'.
// 'iter' is expected to be valid, i.e. iter->IsValid() == true.
// 'iter' is updated to now point to the new entry added, such that
// iter->Get() == 'v'.
func (ul *UList[V]) AddBefore(iter *UListIter[V], v V) {
	writeIter := *iter
	hasPrev := writeIter.Prev()
	if !hasPrev {
		ul.PushFront(v)
		*iter = *ul.Begin()
		return
	}
	*iter = writeIter
	ul.AddAfter(iter, v)
}

// Remove deletes the entry in 'ul' pointed to by 'iter'.
// 'iter' is moved forward in the process. i.e. iter.Get() returns the element in 'ul'
// that occurs after the deleted entry.
func (ul *UList[V]) Remove(iter *UListIter[V]) {
	ul.size--
	blkPtr := iter.node.Value
	*blkPtr = append((*blkPtr)[:iter.index], (*blkPtr)[iter.index+1:]...)
	if len(*blkPtr) == 0 {
		// Block got emptied.
		ul.ll.Remove(iter.node)
		iter.Next()
		return
	}
}

func hasCapacity[V any](llNode *list.Node[blockPtr[V]]) bool {
	if llNode == nil {
		return false
	}
	blkPtr := llNode.Value
	return len(*blkPtr) < cap(*blkPtr)
}

func (ul *UList[V]) newBlock() blockPtr[V] {
	blk := make([]V, 0, ul.entriesPerBlock)
	return &blk
}

func (ul *UList[V]) prependToBlock(v V, blkPtr blockPtr[V]) {
	tmp := *ul.newBlock()
	tmp = append(tmp, v)
	// 'append' returns a slice with capacity of the first variable.
	// To maintain the propoer capacity, we use 'tmp' with an explicitly defined capacity.
	*blkPtr = append(tmp, *blkPtr...)
}

func (iter *UListIter[V]) addOverflowToNextBlock(ul *UList[V], v V) {
	if hasCapacity(iter.node.Next) {
		ul.prependToBlock(v, iter.node.Next.Value)
	} else {
		newBlk := make([]V, 0, ul.entriesPerBlock)
		newBlk = append(newBlk, v)
		ul.ll.InsertAfter(iter.node, &list.Node[blockPtr[V]]{
			Value: &newBlk,
		})
	}
}
