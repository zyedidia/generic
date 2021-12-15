package cache

import "github.com/zyedidia/generic/list"

// A Cache is an LRU cache for keys and values. Each entry is
// put into the table with an associated key used for looking up the entry.
// The cache has a maximum size, and uses a least-recently-used eviction
// policy when there is not space for a new entry.
type Cache[K comparable, V any] struct {
	size     int
	capacity int
	lru      list.List[kv[K, V]]
	table    map[K]*list.Node[kv[K, V]]
}

type kv[K comparable, V any] struct {
	key K
	val V
}

// NewCache returns a new Cache with the given capacity.
func NewCache[K comparable, V any](capacity int) *Cache[K, V] {
	return &Cache[K, V]{
		size:     0,
		capacity: capacity,
		lru:      list.List[kv[K, V]]{},
		table:    make(map[K]*list.Node[kv[K, V]]),
	}
}

// Get returns the entry associated with a given key, and a boolean indicating
// whether the key exists in the table.
func (t *Cache[K, V]) Get(k K) (V, bool) {
	if n, ok := t.table[k]; ok {
		t.lru.Remove(n)
		t.lru.Append(n)
		return n.Value.val, true
	}
	var v V
	return v, false
}

// GetZ is the same as Get but returns the zero-value if k is not found.
func (t *Cache[K, V]) GetZ(k K) V {
	v, _ := t.Get(k)
	return v
}

// Put adds a new key-entry pair to the table.
func (t *Cache[K, V]) Put(k K, e V) {
	if n, ok := t.table[k]; ok {
		n.Value.val = e
		t.lru.Remove(n)
		t.lru.Append(n)
		return
	}

	if t.size == t.capacity {
		t.evict()
	}
	n := &list.Node[kv[K, V]]{
		Value: kv[K, V]{
			key: k,
			val: e,
		},
	}
	t.lru.Append(n)
	t.size++
	t.table[k] = n
}

func (t *Cache[K, V]) evict() {
	key := t.lru.Tail.Value.key
	t.lru.Remove(t.lru.Tail)
	t.size--
	delete(t.table, key)
}

// Delete causes the entry associated with the given key to be immediately
// evicted from the cache.
func (t *Cache[K, V]) Delete(k K) {
	if n, ok := t.table[k]; ok {
		t.lru.Remove(n)
		t.size--
		delete(t.table, k)
	}
}

// Resize changes the maximum capacity for this cache to 'size'.
func (t *Cache[K, V]) Resize(size int) {
	if t.capacity == size {
		return
	} else if t.capacity < size {
		t.capacity = size
		return
	}

	for i := 0; i < t.capacity - size; i++ {
		t.evict()
	}

	t.capacity = size
}

// Size returns the number of active elements in the cache.
func (t *Cache[K, V]) Size() int {
	return t.size
}

// Capacity returns the maximum capacity of the cache.
func (t *Cache[K, V]) Capacity() int {
	return t.capacity
}
