package ulist

import (
	"github.com/zyedidia/generic/list"
)

// A UListIter points to an element in the UList.
type UListIter[V any] struct {
	node  *list.Node[blockPtr[V]]
	index int
}

// newIterFront returns a UListIter pointing to the first entry in 'ul'.
// If 'ul' is empty, an invalid iterator is returned.
func newIterFront[V any](ul *UList[V]) *UListIter[V] {
	return &UListIter[V]{
		node:  ul.ll.Front,
		index: 0,
	}
}

// newIterBack returns a UListIter pointing to the last entry in 'ul'.
// If 'ul' is empty, an invalid iterator is returned.
func newIterBack[V any](ul *UList[V]) *UListIter[V] {
	iter := UListIter[V]{
		node:  ul.ll.Back,
		index: 0,
	}
	if iter.node != nil {
		blk := *iter.node.Value
		iter.index = len(blk) - 1
	}
	return &iter
}

// IsValid returns true if the iterator points to a valid entry in the UList.
func (iter *UListIter[V]) IsValid() bool {
	if iter.node == nil {
		return false
	}
	blkPtr := iter.node.Value
	return iter.index >= 0 && iter.index < len(*blkPtr)
}

// Get returns the entry in the UList that the 'iter' is pointing to.
// This call should only ever be made when iter.IsValid() is true.
func (iter *UListIter[V]) Get() V {
	return (*iter.node.Value)[iter.index]
}

// Next moves the iterator one step forward and returns true if the iterator is valid.
func (iter *UListIter[V]) Next() bool {
	iter.index++
	blkPtr := iter.node.Value
	if iter.index >= len(*blkPtr) {
		if iter.node.Next != nil {
			iter.node = iter.node.Next
			iter.index = 0
		} else {
			// By not going past len, we can recover to the end using Prev().
			iter.index = len(*blkPtr)
		}
	}
	return iter.IsValid()
}

// Prev moves the iterator one step back and returns true if the iterator is valid.
func (iter *UListIter[V]) Prev() bool {
	iter.index--
	if iter.index < 0 {
		if iter.node.Prev != nil {
			iter.node = iter.node.Prev
			blkPtr := iter.node.Value
			iter.index = len(*blkPtr) - 1
		} else {
			// By not going further past -1, we can recover to the begin using Next().
			iter.index = -1
		}
	}
	return iter.IsValid()
}
