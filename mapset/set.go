// Package mapset provides an implementation of a set using the built-in map.
package mapset

// Set implements a hashset, using the hashmap as the underlying storage.
type Set[K comparable] struct {
	m map[K]struct{}
}

// New returns an empty hashset.
func New[K comparable]() Set[K] {
	return Set[K]{
		m: make(map[K]struct{}),
	}
}

// Of returns a new hashset initialized with the given 'vals'
func Of[K comparable](vals ...K) Set[K] {
	s := New[K]()
	for _, val := range vals {
		s.Put(val)
	}
	return s
}

// Put adds 'val' to the set.
func (s Set[K]) Put(val K) {
	s.m[val] = struct{}{}
}

// Has returns true only if 'val' is in the set.
func (s Set[K]) Has(val K) bool {
	_, ok := s.m[val]
	return ok
}

// Remove removes 'val' from the set.
func (s Set[K]) Remove(val K) {
	delete(s.m, val)
}

// Size returns the number of elements in the set.
func (s Set[K]) Size() int {
	return len(s.m)
}

// Each calls 'fn' on every item in the set in no particular order.
func (s Set[K]) Each(fn func(key K)) {
	for k, _ := range s.m {
		fn(k)
	}
}
