// Package hashset provides an implementation of a hashset.
package hashset

import (
	g "github.com/zyedidia/generic"
	"github.com/zyedidia/generic/hashmap"
	"math/rand"
)

// Set implements a hashset, using the hashmap as the underlying storage.
type Set[K any] struct {
	m *hashmap.Map[K, struct{}]
}

// New returns an empty hashset.
func New[K any](capacity uint64, equals g.EqualsFn[K], hash g.HashFn[K]) *Set[K] {
	return &Set[K]{
		m: hashmap.New[K, struct{}](capacity, equals, hash),
	}
}

// Of returns a new hashset initialized with the given 'vals'
func Of[K any](capacity uint64, equals g.EqualsFn[K], hash g.HashFn[K], vals ...K) *Set[K] {
	s := New[K](capacity, equals, hash)
	for _, val := range vals {
		s.Put(val)
	}
	return s
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

// Clear removes all elements from the set.
func (s *Set[K]) Clear() {
	s.m.Clear()
}

// Size returns the number of elements in the set.
func (s *Set[K]) Size() int {
	return s.m.Size()
}

// Each calls 'fn' on every item in the set in no particular order.
func (s *Set[K]) Each(fn func(key K)) {
	s.m.Each(func(key K, v struct{}) {
		fn(key)
	})
}

// Copy returns a copy of this set.
func (s *Set[K]) Copy() *Set[K] {
	return &Set[K]{
		m: s.m.Copy(),
	}
}

func (s *Set[K]) Values() []K {
	return s.m.Keys()
}

func (s *Set[K]) Random() K {
	return s.Values()[rand.Intn(s.Size())]
}
