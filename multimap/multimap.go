// Package multimap provides an associative container that permits multiple entries with the same key.
//
// There are four implementations of the MultiMap data structure, identified by separate New* functions.
// They differ in the following ways:
//   - whether key type and value type must be comparable.
//   - whether duplicate entries (same key and same value) are permitted.
//   - whether keys and values are sorted or unsorted in Get, Each, and EachAssociation methods.
package multimap

// MultiMap is an associative container that contains a list of key-value pairs, while permitting multiple entries with the same key.
type MultiMap[K, V any] interface {
	// Dimension returns number of distinct keys.
	Dimension() int
	// Size returns total number of entries.
	Size() int

	// Count returns number of entries with a given key.
	Count(key K) int
	// Has determines whether at least one entry exists with a given key.
	Has(key K) bool
	// Get returns a list of values with a given key.
	Get(key K) []V

	// Put adds an entry.
	// Whether duplicate entries are allowed depends on the chosen implementation.
	Put(key K, value V)
	// Remove removes an entry.
	// If duplicate entries are allowed, this removes only one entry.
	// This is a no-op if the entry does not exist.
	Remove(key K, value V)
	// RemoveAll removes every entry with a given key.
	RemoveAll(key K)
	// Clear deletes all entries.
	Clear()

	// Each calls 'fn' on every entry.
	Each(fn func(key K, value V))
	// EachAssociation calls 'fn' on every key and list of values.
	EachAssociation(fn func(key K, values []V))
}

type baseMultiMap struct {
	size int
}

func (m baseMultiMap) Size() int {
	return m.size
}
