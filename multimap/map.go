package multimap

import (
	g "github.com/zyedidia/generic"
	"github.com/zyedidia/generic/avl"
)

type mapMultiMap[K comparable, V any, C valuesContainer[V]] struct {
	baseMultiMap
	keys       map[K]C
	makeValues func() C
}

func (m *mapMultiMap[K, V, C]) Dimension() int {
	return len(m.keys)
}

func (m *mapMultiMap[K, V, C]) Count(key K) int {
	values, ok := m.keys[key]
	if !ok {
		return 0
	}
	return values.Size()
}

func (m *mapMultiMap[K, V, C]) Has(key K) bool {
	_, ok := m.keys[key]
	return ok
}

func (m *mapMultiMap[K, V, C]) Get(key K) []V {
	values, ok := m.keys[key]
	if !ok {
		return nil
	}
	return values.List()
}

func (m *mapMultiMap[K, V, C]) Put(key K, value V) {
	values, ok := m.keys[key]
	if !ok {
		values = m.makeValues()
		m.keys[key] = values
	}

	m.size += values.Put(value)
}

func (m *mapMultiMap[K, V, C]) Remove(key K, value V) {
	values, ok := m.keys[key]
	if !ok {
		return
	}

	m.size -= values.Remove(value)
	if values.Empty() {
		delete(m.keys, key)
	}
}

func (m *mapMultiMap[K, V, C]) RemoveAll(key K) {
	values, ok := m.keys[key]
	if !ok {
		return
	}

	m.size -= values.Size()
	delete(m.keys, key)
}

func (m *mapMultiMap[K, V, C]) Clear() {
	m.size = 0
	m.keys = map[K]C{}
}

func (m *mapMultiMap[K, V, C]) Each(fn func(key K, value V)) {
	for key, values := range m.keys {
		values.Each(func(value V) {
			fn(key, value)
		})
	}
}

func (m *mapMultiMap[K, V, C]) EachAssociation(fn func(key K, values []V)) {
	for key, values := range m.keys {
		fn(key, values.List())
	}
}

// NewMapSlice creates a MultiMap using builtin map and builtin slice.
//  - Both key type and value type must be comparable.
//  - Duplicate entries are permitted.
//  - Both keys and values are unsorted.
func NewMapSlice[K, V comparable]() MultiMap[K, V] {
	m := &mapMultiMap[K, V, *valuesSlice[V]]{
		makeValues: func() *valuesSlice[V] {
			return &valuesSlice[V]{}
		},
	}
	m.Clear()
	return m
}

// NewMapSet creates a MultiMap using builtin map and AVL set.
//  - Key type must be comparable.
//  - Duplicate entries are not permitted.
//  - Values are sorted, but keys are unsorted.
func NewMapSet[K comparable, V any](valueLess g.LessFn[V]) MultiMap[K, V] {
	m := &mapMultiMap[K, V, valuesSet[V]]{
		makeValues: func() valuesSet[V] {
			return valuesSet[V]{
				t: avl.New[V, struct{}](valueLess),
			}
		},
	}
	m.Clear()
	return m
}
