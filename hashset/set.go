package hashset

import (
	g "github.com/zyedidia/generic"
	"github.com/zyedidia/generic/hashmap"
	"github.com/zyedidia/generic/iter"
)

type Set[K g.Hashable[K]] struct {
	m *hashmap.Map[K, struct{}]
}

func New[K g.Hashable[K]](capacity uint64) *Set[K] {
	return &Set[K]{
		m: hashmap.NewMap[K, struct{}](capacity),
	}
}

func (s *Set[K]) Put(val K) {
	s.m.Set(val, struct{}{})
}

func (s *Set[K]) Has(val K) bool {
	_, ok := s.m.Get(val)
	return ok
}

func (s *Set[K]) Remove(val K) {
	s.m.Delete(val)
}

func (s *Set[K]) Iter() iter.Iter[K] {
	it := s.m.Iter()
	return func() (K, bool) {
		kv, ok := it()
		return kv.Key, ok
	}
}
