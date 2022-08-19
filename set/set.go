package set

import (
	"fmt"
	"sort"

	"github.com/zyedidia/generic"
	"github.com/zyedidia/generic/hashset"
	"github.com/zyedidia/generic/mapset"
)

func NewMapset[K comparable](in ...K) Set[K] {
	con := func() SetOf[K] { return mapset.New[K]() }
	set := NewSet(con, in...)
	return set
}

func NewHashset[K comparable](cap uint64, equals generic.EqualsFn[K], hash generic.HashFn[K], in ...K) Set[K] {
	con := func() SetOf[K] { return hashset.New(cap, equals, hash) }
	set := NewSet(con, in...)
	return set
}

func NewSet[K comparable, S func() SetOf[K]](con S, in ...K) Set[K] {
	set := con()
	for _, v := range in {
		set.Put(v)
	}
	return Set[K]{
		new:   con,
		SetOf: set,
	}
}

type SetOf[K comparable] interface {
	Put(val K)
	Has(val K) bool
	Remove(val K)
	Size() int
	Each(fn func(key K))
}

type Set[K comparable] struct {
	SetOf[K]
	new func() SetOf[K]
}

// Intersection returns the intersection of the current set with `others`. A new set is returned. Non-mutating action.
func (s Set[K]) Intersection(others ...SetOf[K]) Set[K] {
	return s.Clone().InPlaceIntersection(others...)
}

// Difference returns the elements of the current set that are not in `others`. A new set is returned. Non-mutating action.
func (s Set[K]) Difference(others ...SetOf[K]) Set[K] {
	return s.Clone().InPlaceDifference(others...)
}

// Union returns the elements of teh current set that are common with `others`. A new set is returned. Non-mutating action.
func (s Set[K]) Union(others ...SetOf[K]) Set[K] {
	return s.Clone().InPlaceUnion(others...)
}

// ConstSymmetricDifference returns the elements of the current set
func (s Set[K]) ConstSymmetricDifference(with ...K) Set[K] {
	return s.SymmetricDifference(NewSet(s.new, with...))
}

func (s Set[K]) ConstIntersection(with ...K) Set[K] {
	return s.Clone().InPlaceIntersection(NewSet(s.new, with...))
}

func (s Set[K]) ConstDifference(with ...K) Set[K] {
	return s.Clone().InPlaceDifference(NewSet(s.new, with...))
}

func (s Set[K]) ConstUnion(with ...K) Set[K] {
	return s.Clone().InPlaceUnion(NewSet(s.new, with...))
}

// Clone returns a copy of the current set
func (s Set[K]) Clone() Set[K] {
	new := NewSet(s.new)
	s.Each(func(key K) { new.Put(key) })
	return new
}

// String provides a sorted list of strings representing the set.
func (s Set[K]) String() string {
	out := make([]string, 0, s.Size())
	s.Each(func(key K) { out = append(out, fmt.Sprintf(`%v`, key)) })
	sort.Strings(out)
	return fmt.Sprintf("%v", out)
}

// Map converts the current set to a map where each value in the set is a key of the map. This is useful for converting sets to standard go types.
func (s Set[K]) Map() map[K]struct{} {
	out := make(map[K]struct{}, s.Size())
	s.Each(func(key K) {
		out[key] = struct{}{}
	})
	return out
}

// SymmetricDifference returns a Set of elements that are not in the current set or any of the comparator `others` sets. This does not mutate the current set.
func (s Set[K]) SymmetricDifference(others ...SetOf[K]) Set[K] {
	new := s.Clone()
	seen := new.Clone()
	for _, other := range others {
		other.Each(func(key K) {
			if seen.Has(key) {
				new.Remove(key)
				return
			}
			new.Put(key)
			seen.Put(key)
		})
	}
	return new
}

// InPlaceIntersection removes any elements from the current set that match any elements from provided other sets. This mutates the current set.
func (s Set[K]) InPlaceIntersection(others ...SetOf[K]) Set[K] {
	for _, other := range others {
		s.Each(func(key K) {
			if !other.Has(key) {
				s.Remove(key)
			}
		})
	}
	return s
}

// InPlaceDifference removes the matching elements of N provided sets from the current set. This is mutates the current set.
func (s Set[K]) InPlaceDifference(others ...SetOf[K]) Set[K] {
	for _, other := range others {
		other.Each(func(key K) {
			s.Remove(key)
		})
	}
	return s
}

// InPlaceUnion adds the matching elements of N provided sets to the current set. This is mutates the current set.
func (s Set[K]) InPlaceUnion(others ...SetOf[K]) Set[K] {
	for _, other := range others {
		other.Each(func(key K) {
			s.Put(key)
		})
	}
	return s
}

// Keys provides the keys of the set.
func (s Set[K]) Keys() []K {
	out := make([]K, 0, s.Size())
	s.Each(func(key K) {
		out = append(out, key)
	})
	return out
}

func (s Set[K]) IsDisjoint(other SetOf[K]) bool {
	// TODO: maybe optimize?
	return s.Intersection(other).Size() > 0
}

// IsSubset returns if a set is a subset
func (s Set[K]) IsSubset(of SetOf[K]) bool {
	subset := true
	s.Each(func(key K) {
		if !of.Has(key) {
			subset = false
		}
	})
	return subset
}

// IsSuperset returns if a set is a superset of the provided set. Equal sets are considered supersets of each other.
func (s Set[K]) IsSuperset(of SetOf[K]) bool {
	superset := true
	of.Each(func(key K) {
		if !s.Has(key) {
			superset = false
		}
	})
	return superset
}

// Equal compares two sets
func (s Set[K]) Equal(to SetOf[K]) bool {
	if s.Size() != to.Size() {
		return false
	}
	return s.Union(to).Size() == s.Size()
}

// IsProperSubset returns true if `to` is a subset of `s` but is not equal
func (s Set[K]) IsProperSubset(to SetOf[K]) bool {
	if s.Equal(to) {
		return false
	}
	return s.IsSubset(to)
}

// IsProperSuperset returns true if `to` is a superset of `s` but is not equal
func (s Set[K]) IsProperSuperset(to SetOf[K]) bool {
	if s.Equal(to) {
		return false
	}
	return s.IsSuperset(to)
}
