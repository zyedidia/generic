package iter

// Iter is a generic iterator. When called it returns the next item, along with
// an ok indicator. If no remaining item exists the bool will be false.
type Iter[T any] func() (T, bool)

// For calls 'fn' on every value returned by the iterator.
func (it Iter[T]) For(fn func(t T)) {
	for val, ok := it(); ok; val, ok = it() {
		fn(val)
	}
}

// ForBreak is the same as 'For' but if 'fn' returns false it breaks from the
// loop early.
func (it Iter[T]) ForBreak(fn func(t T) bool) {
	for val, ok := it(); ok; val, ok = it() {
		if !fn(val) {
			break
		}
	}
}

// Slice returns an iterator for a slice.
func Slice[T any](slice []T) Iter[T] {
	var i int
	return func() (t T, done bool) {
		if i >= len(slice) {
			return t, false
		}

		r := slice[i]
		i++
		return r, true
	}
}

// KV is a key-value pair used in a standard map.
type KV[K comparable, V any] struct {
	Key K
	Val V
}

// Map returns an iterator a map.
func Map[K comparable, V any](m map[K]V) Iter[KV[K, V]] {
	keys := make([]KV[K, V], 0, len(m))
	for k, v := range m {
		keys = append(keys, KV[K, V]{
			Key: k,
			Val: v,
		})
	}

	return Slice[KV[K, V]](keys)
}
