package iter

import "github.com/zyedidia/generic"

func Slice[T any](slice []T) generic.Iter[T] {
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

type KV[K comparable, V any] struct {
	Key K
	Val V
}

func Map[K comparable, V any](m map[K]V) generic.Iter[KV[K, V]] {
	keys := make([]KV[K, V], 0, len(m))
	for k, v := range m {
		keys = append(keys, KV[K, V]{
			Key: k,
			Val: v,
		})
	}

	return Slice[KV[K, V]](keys)
}
