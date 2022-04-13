package multimap

import (
	g "github.com/zyedidia/generic"
	"github.com/zyedidia/generic/avl"
)

type avlMultiMap[K, V any, C valuesContainer[V]] struct {
	baseMultiMap
	keyLess    g.LessFn[K]
	keys       *avl.Tree[K, C]
	makeValues func() C
}

func (m *avlMultiMap[K, V, C]) Dimension() int {
	return m.keys.Size()
}

func (m *avlMultiMap[K, V, C]) Count(key K) int {
	values, ok := m.keys.Get(key)
	if !ok {
		return 0
	}
	return values.Size()
}

func (m *avlMultiMap[K, V, C]) Has(key K) bool {
	_, ok := m.keys.Get(key)
	return ok
}

func (m *avlMultiMap[K, V, C]) Get(key K) []V {
	values, ok := m.keys.Get(key)
	if !ok {
		return nil
	}
	return values.List()
}

func (m *avlMultiMap[K, V, C]) Put(key K, value V) {
	values, ok := m.keys.Get(key)
	if !ok {
		values = m.makeValues()
		m.keys.Put(key, values)
	}

	m.size += values.Put(value)
}

func (m *avlMultiMap[K, V, C]) Remove(key K, value V) {
	values, ok := m.keys.Get(key)
	if !ok {
		return
	}

	m.size -= values.Remove(value)
	if values.Empty() {
		m.keys.Remove(key)
	}
}

func (m *avlMultiMap[K, V, C]) RemoveAll(key K) {
	values, ok := m.keys.Get(key)
	if !ok {
		return
	}

	m.size -= values.Size()
	m.keys.Remove(key)
}

func (m *avlMultiMap[K, V, C]) Clear() {
	m.size = 0
	m.keys = avl.New[K, C](m.keyLess)
}

func (m *avlMultiMap[K, V, C]) Each(fn func(key K, value V)) {
	m.keys.Each(func(key K, values C) {
		values.Each(func(value V) {
			fn(key, value)
		})
	})
}

func (m *avlMultiMap[K, V, C]) EachAssociation(fn func(key K, values []V)) {
	m.keys.Each(func(key K, values C) {
		fn(key, values.List())
	})
}

// NewAvlSlice creates a MultiMap using AVL tree and builtin slice.
//  - Value type must be comparable.
//  - Duplicate entries are permitted.
//  - Keys are sorted, but values are unsorted.
func NewAvlSlice[K any, V comparable](keyLess g.LessFn[K]) MultiMap[K, V] {
	m := &avlMultiMap[K, V, *valuesSlice[V]]{
		keyLess: keyLess,
		makeValues: func() *valuesSlice[V] {
			return &valuesSlice[V]{}
		},
	}
	m.Clear()
	return m
}

// NewAvlSet creates a MultiMap using AVL tree and AVL set.
//  - Duplicate entries are not permitted.
//  - Both keys and values are sorted.
func NewAvlSet[K, V any](keyLess g.LessFn[K], valueLess g.LessFn[V]) MultiMap[K, V] {
	m := &avlMultiMap[K, V, valuesSet[V]]{
		keyLess: keyLess,
		makeValues: func() valuesSet[V] {
			return valuesSet[V]{
				t: avl.New[V, struct{}](valueLess),
			}
		},
	}
	m.Clear()
	return m
}
