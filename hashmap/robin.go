package hashmap

import (
	g "github.com/zyedidia/generic"
)

const (
	emptyBucket  = -1
	resizeFactor = 2
)

type bucket[K comparable, V any] struct {
	key K
	// psl is the probe sequence length (PSL), which is the distance value from
	// the optimum insertion. -1 or `emptyBucket` signals a free slot.
	// inspired from:
	//  - https://programming.guide/robin-hood-hashing.html
	//  - https://cs.uwaterloo.ca/research/tr/1986/CS-86-14.pdf
	psl   int8
	value V
}

// RobinMap is a hashmap that uses linear probing in combination with
// robin hood hashing as collision strategy.
type RobinMap[K comparable, V any] struct {
	buckets []bucket[K, V]
	hasher  g.HashFn[K]
	// length stores the current inserted elements
	length uintptr
	// indexMask is used for a bitwise AND on the hash value,
	// because the size of the underlying array is a power of two value
	indexMask uintptr
	// log2Index is the number of extra reserved bytes at the end of the array,
	// to sparse the length check while probing.
	// Furthermore this value is the maximum possible PSL over the hashmap,
	// because a grow is forced if this value will raised during the insert operation.
	log2Index int8
}

// go:inline
func newBucketArray[K comparable, V any](capacity uintptr) []bucket[K, V] {
	buckets := make([]bucket[K, V], capacity)
	for i := range buckets {
		buckets[i].psl = emptyBucket
	}
	return buckets
}

// NewRobinMapWithHasher constructs a new map with the given hasher function.
func NewRobinMapWithHasher[K comparable, V any](hasher g.HashFn[K]) *RobinMap[K, V] {

	capacity := uintptr(4)
	log2Index := uintptr(2)

	return &RobinMap[K, V]{
		buckets:   newBucketArray[K, V](capacity + log2Index + 1),
		indexMask: capacity - 1,
		log2Index: int8(log2Index),
		hasher:    hasher,
	}
}

// getBucket return a pointer to the bucket if the key was found.
// Furthermore the index of the underlying array and the psl is returned.
//
// Note:
//   - There exists also other search strategies like organ-pipe search
//     or smart search, where searching starts at the end (tracks max PSL)
//     or around the mean value (mean, mean − 1, mean + 1, mean − 2, mean + 2, ...)
//   - Here it is used the simplest technic, which is more cache friendly and
//     does not track other metic values.
//
// go:inline
func (m *RobinMap[K, V]) getBucket(key K) (*bucket[K, V], uintptr, int8) {
	hash := uintptr(m.hasher(key))
	idx := hash & m.indexMask

	psl := int8(0)
	for ; psl <= m.buckets[idx].psl; psl++ {
		if m.buckets[idx].key == key {
			return &m.buckets[idx], idx, psl
		}
		idx++
	}
	return nil, idx, psl
}

// Get returns the value stored for this key, or false if there is no such value.
func (m *RobinMap[K, V]) Get(key K) (V, bool) {
	var v V
	e, _, _ := m.getBucket(key)
	if e == nil {
		return v, false
	}
	return e.value, true
}

// Reserve sets the number of entires in the container to the most appropriate
// to contain at least n elements. If n is lower than that, the function may have no effect.
func (m *RobinMap[K, V]) Reserve(n uintptr) {
	newCap := uintptr(g.NextPowerOf2(uint64(resizeFactor * n)))
	if (m.indexMask + 1) < newCap {
		m.resize(newCap)
	}
}

// go:inline
func (m *RobinMap[K, V]) grow() {
	capacity := m.indexMask + 1
	m.resize(capacity * resizeFactor)
}

// go:inline
func (m *RobinMap[K, V]) resize(n uintptr) {
	// extra space, that is at the same time the worse case lookup time
	log2Index := (3 * uintptr(g.Log2(uint64(n)))) / 2
	newm := RobinMap[K, V]{
		indexMask: n - 1,
		log2Index: int8(log2Index),
		length:    m.length,
		buckets:   newBucketArray[K, V](n + log2Index + 1),
		hasher:    m.hasher,
	}

	for _, ent := range m.buckets {
		if ent.psl != emptyBucket {
			newm.Put(ent.key, ent.value)
		}
	}
	m.indexMask = newm.indexMask
	m.log2Index = newm.log2Index
	m.buckets = newm.buckets
}

// Put maps the given key to the given value. If the key already exists its
// value will be overwritten with the new value.
func (m *RobinMap[K, V]) Put(key K, val V) {

	e, idx, psl := m.getBucket(key)

	// override old value
	if e != nil {
		e.value = val
		return
	}

	m.robinHoodEmplace(bucket[K, V]{key: key, value: val, psl: psl}, idx)
}

// robinHoodEmplace applies the Robin Hood creed to all following entires until a empty is found.
// Robin Hood creed: "takes from the rich and gives to the poor".
// rich means, low psl
// poor means, higher psl
//
// The result is a normal distribution of the PSL values,
// where the expected length of the longest PSL is O(log(n))
func (m *RobinMap[K, V]) robinHoodEmplace(e bucket[K, V], idx uintptr) {

	if m.length >= m.indexMask {
		m.grow()
		m.Put(e.key, e.value)
		return
	}

	for ; ; e.psl++ {
		if m.buckets[idx].psl == emptyBucket {
			// emplace the element, a valid bucket was found
			m.buckets[idx] = e
			m.length++
			return
		}
		// force resize to leave out overflow check of m.buckets
		if e.psl >= m.log2Index {
			m.grow()
			m.Put(e.key, e.value)
			return
		}
		if e.psl > m.buckets[idx].psl {
			g.Swap(&e, &m.buckets[idx])
		}

		idx++
	}
}

// Remove deletes the specified key-value pair from the map.
func (m *RobinMap[K, V]) Remove(key K) {
	current, idx, _ := m.getBucket(key)
	if current == nil {
		return
	}

	m.length--
	current.psl = emptyBucket // make as empty, because we want to remove it

	idx++
	next := &m.buckets[idx]

	// now, back shift all buckets until we found a optimum or empty one
	for next.psl > 0 {
		next.psl--
		g.Swap(current, next)
		current = next
		idx++
		next = &m.buckets[idx]
	}
}

// Clear removes all key-value pairs from the map.
func (m *RobinMap[K, V]) Clear() {
	for idx := range m.buckets {
		m.buckets[idx].psl = emptyBucket
	}
}

// Load return the current load of the hash map.
func (m *RobinMap[K, V]) Load() float64 {
	capacity := m.indexMask + 1
	return float64(m.length) / float64(capacity)
}

// Size returns the number of items in the map.
func (m *RobinMap[K, V]) Size() int {
	return int(m.length)
}

// Copy returns a copy of this map.
func (m *RobinMap[K, V]) Copy() *RobinMap[K, V] {
	return &RobinMap[K, V]{
		buckets:   m.buckets,
		indexMask: m.indexMask,
		log2Index: m.log2Index,
		length:    m.length,
		hasher:    m.hasher,
	}
}

// Each calls 'fn' on every key-value pair in the hashmap in no particular order.
func (m *RobinMap[K, V]) Each(fn func(key K, val V)) {
	for _, ent := range m.buckets {
		if ent.psl != emptyBucket {
			fn(ent.key, ent.value)
		}
	}
}
