// Package cache provides an implementation of a key-value store with a maximum
// size. Once the maximum size is reached, the cache uses a least-recently-used
// policy to evict old entries. The cache is implemented as a combined hashmap
// and linked list. This ensures all operations are constant-time.
package cache

import (
	"github.com/zyedidia/generic/list"
)

// A Cache is an LRU cache for keys and values. Each entry is
// put into the table with an associated key used for looking up the entry.
// The cache has a maximum size, and uses a least-recently-used eviction
// policy when there is not space for a new entry.
type Cache[K comparable, V any] struct {
	capacity int
	lru      list.List[KV[K, V]]
	table    map[K]*list.Node[KV[K, V]]
	evictCb  func(key K, val V)
}

type KV[K comparable, V any] struct {
	Key K
	Val V
}

// New returns a new Cache with the given capacity.
func New[K comparable, V any](capacity int) *Cache[K, V] {
	return &Cache[K, V]{
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

	if len(t.table) == t.capacity {
		t.evict()
	}
	n := &list.Node[KV[K, V]]{
		Value: KV[K, V]{
			Key: k,
			Val: e,
		},
	}
	t.lru.PushFrontNode(n)
	t.table[k] = n
}

func (t *Cache[K, V]) evict() {
	entry := t.lru.Back.Value
	if t.evictCb != nil {
		t.evictCb(entry.Key, entry.Val)
	}
	t.lru.Remove(t.lru.Back)
	delete(t.table, entry.Key)
}

// Remove causes the entry associated with the given key to be immediately
// evicted from the cache.
func (t *Cache[K, V]) Remove(k K) {
	if n, ok := t.table[k]; ok {
		t.lru.Remove(n)
		delete(t.table, k)
	}
}

// Resize changes the maximum capacity for this cache to 'capacity'.
func (t *Cache[K, V]) Resize(capacity int) {
	t.capacity = capacity
	for len(t.table) > capacity {
		t.evict()
	}
}

// Size returns the number of active elements in the cache.
func (t *Cache[K, V]) Size() int {
	return len(t.table)
}

// Capacity returns the maximum capacity of the cache.
func (t *Cache[K, V]) Capacity() int {
	return t.capacity
}

// Each calls 'fn' on every value in the cache, from most recently used to
// least recently used.
func (t *Cache[K, V]) Each(fn func(key K, val V)) {
	t.lru.Front.Each(func(kv KV[K, V]) {
		fn(kv.Key, kv.Val)
	})
}

// SetEvictCallback sets a callback to be invoked before an entry is evicted.
// This replaces any prior callback set by this method.
func (t *Cache[K, V]) SetEvictCallback(fn func(key K, val V)) {
	t.evictCb = fn
}
