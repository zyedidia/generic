// Package bimap provides an implementation of a bi-directional map.
//
// It is implemented by using two Go maps, which keeps the lookup speed
// identical for both forward and reverse lookups, however it also doubles the
// memory usage of the map.
package bimap

// Of returns a new [Bimap] initiated with the keys and values
// from the given map.
func Of[K, V comparable](m map[K]V) Bimap[K, V] {
	bm := Bimap[K, V]{}
	for k, v := range m {
		bm.Add(k, v)
	}
	return bm
}

// Bimap is a bi-directional map where both the keys and values are indexed
// against each other, allowing performant lookup on both keys and values,
// at the cost of double the memory usage.
type Bimap[K, V comparable] struct {
	forward map[K]V
	reverse map[V]K
}

// Len returns the number of key-value pairs in this map.
func (b *Bimap[K, V]) Len() int {
	if b == nil {
		return 0
	}
	return len(b.forward)
}

// Add another key-value pair to be indexed inside this map. Both the key
// and the value is indexed, to allow performant lookups on both key and value.
//
// On collisions, the old values will be overwritten.
func (b *Bimap[K, V]) Add(key K, value V) {
	if oldVal, ok := b.GetForward(key); ok {
		delete(b.reverse, oldVal)
	}
	if oldKey, ok := b.GetReverse(value); ok {
		delete(b.forward, oldKey)
	}
	if b.forward == nil {
		b.forward = make(map[K]V)
		b.reverse = make(map[V]K)
	}
	b.forward[key] = value
	b.reverse[value] = key
}

// RemoveForward removes a key-value pair from this map based on the key.
func (b *Bimap[K, V]) RemoveForward(key K) {
	if value, ok := b.forward[key]; ok {
		delete(b.reverse, value)
		delete(b.forward, key)
	}
}

// RemoveReverse removes a key-value pair from this map based on the value.
func (b *Bimap[K, V]) RemoveReverse(value V) {
	if key, ok := b.reverse[value]; ok {
		delete(b.reverse, value)
		delete(b.forward, key)
	}
}

// Each loops over all the values in this map.
func (b *Bimap[K, V]) Each(f func(key K, value V)) {
	for k, v := range b.forward {
		f(k, v)
	}
}

// ContainsForward checks if the given key exists.
func (b *Bimap[K, V]) ContainsForward(key K) bool {
	_, ok := b.forward[key]
	return ok
}

// GetForward performs a lookup on the key to get the value.
func (b *Bimap[K, V]) GetForward(key K) (V, bool) {
	value, ok := b.forward[key]
	return value, ok
}

// ContainsReverse checks if the given value exists.
func (b *Bimap[K, V]) ContainsReverse(value V) bool {
	_, ok := b.reverse[value]
	return ok
}

// GetReverse performs a lookup on the value to get the key.
func (b *Bimap[K, V]) GetReverse(value V) (K, bool) {
	key, ok := b.reverse[value]
	return key, ok
}

// Clear empties this bidirectional map, removing all items.
func (b *Bimap[K, V]) Clear() {
	clear(b.forward)
	clear(b.reverse)
}

// Copy creates a shallow copy of this bidirectional map.
func (b *Bimap[K, V]) Copy() Bimap[K, V] {
	return Bimap[K, V]{
		forward: shallowCopy(b.forward),
		reverse: shallowCopy(b.reverse),
	}
}

func clear[M ~map[K]V, K comparable, V any](m M) {
	// Relies on the compiler optimization introduced in Go v1.11
	// https://go.dev/doc/go1.11#performance-compiler
	for k := range m {
		delete(m, k)
	}
}

func shallowCopy[M ~map[K]V, K comparable, V any](m M) M {
	newMap := make(M, len(m))
	for k, v := range m {
		newMap[k] = v
	}
	return newMap
}
