// Package hashmap provides an implementation of a hashmap. The map uses linear
// probing and automatically resizes. The map can also be efficiently copied,
// and will perform copies lazily, using copy-on-write. However, the
// copy-on-write will copy the entire map after the first write. One can imagine
// a more efficient implementation that would split the map into chunks and use
// copy-on-write selectively for each chunk.
package hashmap

import (
	g "github.com/zyedidia/generic"
	"math/rand"
)

// A Map is a hashmap that supports copying via copy-on-write.
type Map[K, V any] struct {
	keys     []K
	values   []V
	filled   []bool
	capacity uint64
	length   uint64
	readonly bool

	ops ops[K]
}

type ops[T any] struct {
	equals func(a, b T) bool
	hash   func(t T) uint64
}

func pow2ceil(num uint64) uint64 {
	power := uint64(1)
	for power < num {
		power *= 2
	}
	return power
}

// New constructs a new map with the given capacity.
func New[K, V any](capacity uint64, equals g.EqualsFn[K], hash g.HashFn[K]) *Map[K, V] {
	if capacity == 0 {
		capacity = 1
	}
	capacity = pow2ceil(capacity)
	return &Map[K, V]{
		keys:     make([]K, capacity),
		values:   make([]V, capacity),
		filled:   make([]bool, capacity),
		capacity: capacity,
		ops: ops[K]{
			equals: equals,
			hash:   hash,
		},
	}
}

// Get returns the value stored for this key, or false if there is no such
// value.
func (m *Map[K, V]) Get(key K) (V, bool) {
	hash := m.ops.hash(key)
	idx := hash & (m.capacity - 1)

	for m.filled[idx] {
		if m.ops.equals(m.keys[idx], key) {
			return m.values[idx], true
		}
		idx++
		if idx >= m.capacity {
			idx = 0
		}
	}

	var v V
	return v, false
}

func (m *Map[K, V]) resize(newcap uint64) {
	newm := Map[K, V]{
		capacity: newcap,
		length:   m.length,
		keys:     make([]K, newcap),
		values:   make([]V, newcap),
		filled:   make([]bool, newcap),
		ops:      m.ops,
	}

	for i := range m.keys {
		if m.filled[i] {
			newm.Put(m.keys[i], m.values[i])
		}
	}

	m.capacity = newm.capacity
	m.keys = newm.keys
	m.values = newm.values
	m.filled = newm.filled
}

// Put maps the given key to the given value. If the key already exists its
// value will be overwritten with the new value.
func (m *Map[K, V]) Put(key K, val V) {
	if m.length >= m.capacity/2 {
		m.resize(m.capacity * 2)
	} else if m.readonly {
		keys := make([]K, len(m.keys), cap(m.keys))
		values := make([]V, len(m.values), cap(m.values))
		filled := make([]bool, len(m.filled), cap(m.filled))
		copy(keys, m.keys)
		copy(values, m.values)
		copy(filled, m.filled)
		m.keys = keys
		m.values = values
		m.filled = filled
		m.readonly = false
	}

	hash := m.ops.hash(key)
	idx := hash & (m.capacity - 1)

	for m.filled[idx] {
		if m.ops.equals(m.keys[idx], key) {
			m.values[idx] = val
			return
		}
		idx++
		if idx >= m.capacity {
			idx = 0
		}
	}

	m.keys[idx] = key
	m.values[idx] = val
	m.filled[idx] = true
	m.length++
}

func (m *Map[K, V]) remove(idx uint64) {
	var k K
	var v V
	m.filled[idx] = false
	m.keys[idx] = k
	m.values[idx] = v
	m.length--
}

// Remove removes the specified key-value pair from the map.
func (m *Map[K, V]) Remove(key K) {
	hash := m.ops.hash(key)
	idx := hash & (m.capacity - 1)

	for m.filled[idx] && !m.ops.equals(m.keys[idx], key) {
		idx = (idx + 1) & (m.capacity - 1)
	}

	if !m.filled[idx] {
		return
	}

	if m.readonly {
		keys := make([]K, len(m.keys), cap(m.keys))
		values := make([]V, len(m.values), cap(m.values))
		filled := make([]bool, len(m.filled), cap(m.filled))
		copy(keys, m.keys)
		copy(values, m.values)
		copy(filled, m.filled)
		m.keys = keys
		m.values = values
		m.filled = filled
		m.readonly = false
	}

	m.remove(idx)

	idx = (idx + 1) & (m.capacity - 1)
	for m.filled[idx] {
		krehash := m.keys[idx]
		vrehash := m.values[idx]
		m.remove(idx)
		m.Put(krehash, vrehash)
		idx = (idx + 1) & (m.capacity - 1)
	}

	// halves the array if it is 12.5% full or less
	if m.length > 0 && m.length <= m.capacity/8 {
		m.resize(m.capacity / 2)
	}
}

// Clear removes all key-value pairs from the map.
func (m *Map[K, V]) Clear() {
	for idx := range m.keys {
		if m.filled[idx] {
			m.remove(uint64(idx))
		}
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
		keys:     m.keys,
		values:   m.values,
		filled:   m.filled,
		capacity: m.capacity,
		length:   m.length,
		readonly: true,
		ops:      m.ops,
	}
}

// Each calls 'fn' on every key-value pair in the hashmap in no particular
// order.
func (m *Map[K, V]) Each(fn func(key K, val V)) {
	for idx := range m.keys {
		if m.filled[idx] {
			fn(m.keys[idx], m.values[idx])
		}
	}
}

// Keys returns the key of the hashmap
func (m *Map[K, V]) Keys() []K {
	keys := make([]K, len(m.keys))

	for idx := range m.keys {
		if m.filled[idx] {
			keys = append(keys, m.keys[idx])
		}
	}
	return keys[:m.Size()]
}

// Values returns the values of the hashmap
func (m *Map[K, V]) Values() []V {
	values := make([]V, len(m.keys))

	for idx := range m.keys {
		if m.filled[idx] {
			values = append(values, m.values[idx])
		}
	}
	return values[:m.Size()]
}

func (m *Map[K, V]) Random() (K, V) {
	randomIndex := rand.Intn(m.Size())
	return m.Keys()[randomIndex], m.Values()[randomIndex]
}
