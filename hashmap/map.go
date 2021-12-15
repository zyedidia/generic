package hashmap

// Hashable implements the Hash() function, and must support the == operator.
type Hashable interface {
	comparable
	Hash() uint64
}

type entry[K Hashable, V any] struct {
	key K
	filled bool
	value V
}

// A Map is a hashmap that supports copying via copy-on-write.
type Map[K Hashable, V any] struct {
	entries []entry[K, V]
	capacity uint64
	length uint64
	readonly bool
}

// NewCowMap constructs a new map with the given capacity.
func NewCowMap[K Hashable, V any](capacity uint64) *Map[K, V] {
	if capacity == 0 {
		capacity = 1
	}
	return &Map[K, V]{
		entries: make([]entry[K, V], capacity),
		capacity: capacity,
	}
}

// Get returns the value stored for this key, or false if there is no such
// value.
func (m *Map[K, V]) Get(key K) (V, bool) {
	hash := key.Hash()
	idx := hash & (m.capacity - 1)

	for m.entries[idx].filled {
		if (m.entries[idx].key == key) {
			return m.entries[idx].value, true
		}
		idx++
		if (idx >= m.capacity) {
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

func (m *Map[K, V]) expandto(newcap uint64) {
	newm := Map[K, V]{
		capacity: newcap,
		length: m.length,
		entries: make([]entry[K, V], newcap),
	}

	for _, ent := range m.entries {
		if ent.filled {
			newm.Set(ent.key, ent.value)
		}
	}
	m.capacity = newm.capacity
	m.entries = newm.entries
}

// Set maps the given key to the given value. If the key already exists its
// value will be overwritten with the new value.
func (m *Map[K, V]) Set(key K, val V) {
	if m.length >= m.capacity / 2 {
		m.expandto(m.capacity * 2)
	} else if m.readonly {
		entries := make([]entry[K, V], len(m.entries))
		copy(entries, m.entries)
		m.entries = entries
	}

	hash := key.Hash()
	idx := hash & (m.capacity - 1)

	for m.entries[idx].filled {
		if m.entries[idx].key == key {
			m.entries[idx].value = val
			return
		}
		idx++
		if idx >= m.capacity {
			idx = 0
		}
	}

	m.entries[idx].key = key;
	m.entries[idx].value = val;
	m.entries[idx].filled = true
	m.length++
}

// Copy returns a copy of this map. The copy will not allocate any memory until
// the first write, so any number of read-only copies can be made without any
// additional allocations.
func (m *Map[K, V]) Copy() *Map[K, V] {
	m.readonly = true
	return &Map[K, V]{
		entries: m.entries,
		capacity: m.capacity,
		length: m.length,
		readonly: true,
	}
}

// Keys returns the map's keys.
func (m *Map[K, V]) Keys() []K {
	keys := make([]K, m.length)

	for _, ent := range m.entries {
		if ent.filled {
			keys = append(keys, ent.key)
		}
	}

	return keys
}

// Range applies 'fn' to every value in the map.
func (m *Map[K, V]) Range(fn func (key K, val V)) {
	for _, ent := range m.entries {
		if ent.filled {
			fn(ent.key, ent.value)
		}
	}
}
