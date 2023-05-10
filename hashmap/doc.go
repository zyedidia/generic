// Package hashmap provides several implementation of hashmaps.
package hashmap

// HashMap is the collection of the basic interface functions.
type HashMap[K comparable, V any] struct {
	Get     func(key K) (V, bool)
	Reserve func(n uintptr)
	Load    func() float64
	Put     func(key K, val V)
	Remove  func(key K)
	Clear   func()
	Size    func() int
	Each    func(fn func(key K, val V))
}
