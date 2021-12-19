// Package hashset provides an implementation of a hashset.
package hashset

import (
	g "github.com/zyedidia/generic"
	"github.com/zyedidia/generic/hashmap"
	"github.com/zyedidia/generic/iter"
)

// Set implements a hashset, using the hashmap as the underlying storage.
type Set[K any] struct {
	m *hashmap.Map[K, struct{}]
}

// New returns an empty hashset.
func New[K any](capacity uint64, equals g.Equaler[K], hash g.Hasher[K]) *Set[K] {
	return &Set[K]{
		m: hashmap.NewMap[K, struct{}](capacity, equals, hash),
	}
}

// Put adds 'val' to the set.
func (s *Set[K]) Put(val K) {
	s.m.Put(val, struct{}{})
}

// Has returns true only if 'val' is in the set.
func (s *Set[K]) Has(val K) bool {
	_, ok := s.m.Get(val)
	return ok
}

// Remove removes 'val' from the set.
func (s *Set[K]) Remove(val K) {
	s.m.Remove(val)
}

// Size returns the number of elements in the set.
func (s *Set[K]) Size() int {
	return s.m.Size()
}

// Iter returns an iterator over all values in the set.
func (s *Set[K]) Iter() iter.Iter[K] {
	it := s.m.Iter()
	return func() (K, bool) {
		kv, ok := it()
		return kv.Key, ok
	}
}
