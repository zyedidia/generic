// Package hashmap provides an implementation of a hashmap. The map uses linear
// probing and automatically resizes. The map can also be efficiently copied,
// and will perform copies lazily, using copy-on-write. However, the
// copy-on-write will copy the entire map after the first write. One can imagine
// a more efficient implementation that would split the map into chunks and use
// copy-on-write selectively for each chunk.
package hashmap

import (
	g "github.com/zyedidia/generic"
	"github.com/zyedidia/generic/iter"
)

type entry[K g.Hashable[K], V any] struct {
	key    K
	filled bool
	value  V
}

// A Map is a hashmap that supports copying via copy-on-write.
type Map[K g.Hashable[K], V any] struct {
	entries  []entry[K, V]
	capacity uint64
	length   uint64
	readonly bool
}

// NewMap constructs a new map with the given capacity.
func NewMap[K g.Hashable[K], V any](capacity uint64) *Map[K, V] {
	if capacity == 0 {
		capacity = 1
	}
	return &Map[K, V]{
		entries:  make([]entry[K, V], capacity),
		capacity: capacity,
	}
}

// Get returns the value stored for this key, or false if there is no such
// value.
func (m *Map[K, V]) Get(key K) (V, bool) {
	hash := key.Hash()
	idx := hash & (m.capacity - 1)

	for m.entries[idx].filled {
		if m.entries[idx].key.Equals(key) {
			return m.entries[idx].value, true
		}
		idx++
		if idx >= m.capacity {
			idx = 0
		}
	}

	var v V
	return v, false
}

// GetZ returns the value stored for this key, or its zero value if there is no
// such value.
func (m *Map[K, V]) GetZ(key K) V {
	v, _ := m.Get(key)
	return v
}

func (m *Map[K, V]) resize(newcap uint64) {
	newm := Map[K, V]{
		capacity: newcap,
		length:   m.length,
		entries:  make([]entry[K, V], newcap),
	}

	for _, ent := range m.entries {
		if ent.filled {
			newm.Put(ent.key, ent.value)
		}
	}
	m.capacity = newm.capacity
	m.entries = newm.entries
}

// Put maps the given key to the given value. If the key already exists its
// value will be overwritten with the new value.
func (m *Map[K, V]) Put(key K, val V) {
	if m.length >= m.capacity/2 {
		m.resize(m.capacity * 2)
	} else if m.readonly {
		entries := make([]entry[K, V], len(m.entries))
		copy(entries, m.entries)
		m.entries = entries
	}

	hash := key.Hash()
	idx := hash & (m.capacity - 1)

	for m.entries[idx].filled {
		if m.entries[idx].key.Equals(key) {
			m.entries[idx].value = val
			return
		}
		idx++
		if idx >= m.capacity {
			idx = 0
		}
	}

	m.entries[idx].key = key
	m.entries[idx].value = val
	m.entries[idx].filled = true
	m.length++
}

func (m *Map[K, V]) remove(idx uint64) {
	var k K
	var v V
	m.entries[idx].filled = false
	m.entries[idx].key = k
	m.entries[idx].value = v
	m.length--
}

// Remove removes the specified key-value pair from the map.
func (m *Map[K, V]) Remove(key K) {
	hash := key.Hash()
	idx := hash & (m.capacity - 1)

	for m.entries[idx].filled && !m.entries[idx].key.Equals(key) {
		idx = (idx + 1) % m.capacity
	}

	if !m.entries[idx].filled {
		return
	}

	m.remove(idx)

	idx = (idx + 1) % m.capacity
	for m.entries[idx].filled {
		krehash := m.entries[idx].key
		vrehash := m.entries[idx].value
		m.remove(idx)
		m.Put(krehash, vrehash)
		idx = (idx + 1) % m.capacity
	}

	// halves the array if it is 12.5% full or less
	if m.length > 0 && m.length <= m.capacity/8 {
		m.resize(m.capacity / 2)
	}
}

// Size returns the number of items in the map.
func (m *Map[K, V]) Size() int {
	return int(m.length)
}

// Copy returns a copy of this map. The copy will not allocate any memory until
// the first write, so any number of read-only copies can be made without any
// additional allocations.
func (m *Map[K, V]) Copy() *Map[K, V] {
	m.readonly = true
	return &Map[K, V]{
		entries:  m.entries,
		capacity: m.capacity,
		length:   m.length,
		readonly: true,
	}
}

type KV[K g.Hashable[K], V any] struct {
	Key K
	Val V
}

// Iter returns an iterator over all key-value pairs in the map.
func (m *Map[K, V]) Iter() iter.Iter[KV[K, V]] {
	kvs := make([]KV[K, V], 0, m.length)
	for _, ent := range m.entries {
		if ent.filled {
			kvs = append(kvs, KV[K, V]{
				Key: ent.key,
				Val: ent.value,
			})
		}
	}
	return iter.Slice(kvs)
}
