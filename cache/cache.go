// Package cache provides an implementation of a key-value store with a maximum
// size. Once the maximum size is reached, the cache uses a least-recently-used
// policy to evict old entries. The cache is implemented as a combined hashmap
// and linked list. This ensures all operations are constant-time.
package cache

import (
	"github.com/zyedidia/generic/iter"
	"github.com/zyedidia/generic/list"
)

// A Cache is an LRU cache for keys and values. Each entry is
// put into the table with an associated key used for looking up the entry.
// The cache has a maximum size, and uses a least-recently-used eviction
// policy when there is not space for a new entry.
type Cache[K comparable, V any] struct {
	size     int
	capacity int
	lru      list.List[KV[K, V]]
	table    map[K]*list.Node[KV[K, V]]
}

type KV[K comparable, V any] struct {
	Key K
	Val V
}

// New returns a new Cache with the given capacity.
func New[K comparable, V any](capacity int) *Cache[K, V] {
	return &Cache[K, V]{
		size:     0,
		capacity: capacity,
		lru:      list.List[KV[K, V]]{},
		table:    make(map[K]*list.Node[KV[K, V]]),
	}
}

// Get returns the entry associated with a given key, and a boolean indicating
// whether the key exists in the table.
func (t *Cache[K, V]) Get(k K) (V, bool) {
	if n, ok := t.table[k]; ok {
		t.lru.Remove(n)
		t.lru.PushFrontNode(n)
		return n.Value.Val, true
	}
	var v V
	return v, false
}

// Put adds a new key-entry pair to the table.
func (t *Cache[K, V]) Put(k K, e V) {
	if n, ok := t.table[k]; ok {
		n.Value.Val = e
		t.lru.Remove(n)
		t.lru.PushFrontNode(n)
		return
	}

	if t.size == t.capacity {
		t.evict()
	}
	n := &list.Node[KV[K, V]]{
		Value: KV[K, V]{
			Key: k,
			Val: e,
		},
	}
	t.lru.PushFrontNode(n)
	t.size++
	t.table[k] = n
}

func (t *Cache[K, V]) evict() {
	key := t.lru.Back.Value.Key
	t.lru.Remove(t.lru.Back)
	t.size--
	delete(t.table, key)
}

// Remove causes the entry associated with the given key to be immediately
// evicted from the cache.
func (t *Cache[K, V]) Remove(k K) {
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

	for i := 0; i < t.capacity-size; i++ {
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

// Iter returns an iterator over all key-value pairs in the cache. It iterates
// in order of most recently used to least recently used.
func (t *Cache[K, V]) Iter() iter.Iter[KV[K, V]] {
	return t.lru.Front.Iter()
}
