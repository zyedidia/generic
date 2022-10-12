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
	Clear()
	Size() int
	Each(fn func(key K))
}

type Set[K comparable] struct {
	SetOf[K]
	new func() SetOf[K]
}

func (s Set[K]) Intersection(others ...SetOf[K]) Set[K] {
	return s.Clone().InPlaceIntersection(others...)
}
func (s Set[K]) Difference(others ...SetOf[K]) Set[K] {
	return s.Clone().InPlaceDifference(others...)
}
func (s Set[K]) Union(others ...SetOf[K]) Set[K] {
	return s.Clone().InPlaceUnion(others...)
}

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

func (s Set[K]) Clone() Set[K] {
	new := NewSet(s.new)
	s.Each(func(key K) { new.Put(key) })
	return new
}

func (s Set[K]) String() string {
	out := make([]string, 0, s.Size())
	s.Each(func(key K) { out = append(out, fmt.Sprintf(`%v`, key)) })
	sort.Strings(out)
	return fmt.Sprintf("%v", out)
}

func (s Set[K]) Map() map[K]struct{} {
	out := make(map[K]struct{}, s.Size())
	s.Each(func(key K) {
		out[key] = struct{}{}
	})
	return out
}

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

func (s Set[K]) InPlaceDifference(others ...SetOf[K]) Set[K] {
	for _, other := range others {
		other.Each(func(key K) {
			s.Remove(key)
		})
	}
	return s
}

func (s Set[K]) InPlaceUnion(others ...SetOf[K]) Set[K] {
	for _, other := range others {
		other.Each(func(key K) {
			s.Put(key)
		})
	}
	return s
}

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

func (s Set[K]) IsSubset(of SetOf[K]) bool {
	subset := true
	s.Each(func(key K) {
		if !of.Has(key) {
			subset = false
		}
	})
	return subset
}

func (s Set[K]) IsSuperset(of SetOf[K]) bool {
	superset := true
	of.Each(func(key K) {
		if !s.Has(key) {
			superset = false
		}
	})
	return superset
}

func (s Set[K]) Equal(to SetOf[K]) bool {
	if s.Size() != to.Size() {
		return false
	}
	return s.Union(to).Size() == s.Size()
}

func (s Set[K]) IsProperSubset(to SetOf[K]) bool {
	if s.Equal(to) {
		return false
	}
	return s.IsSubset(to)
}

func (s Set[K]) IsProperSuperset(to SetOf[K]) bool {
	if s.Equal(to) {
		return false
	}
	return s.IsSuperset(to)
}
