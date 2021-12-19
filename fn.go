package generic

import "constraints"

func Max[T constraints.Ordered](a, b T) T {
	if a > b {
		return a
	}
	return b
}

func Min[T constraints.Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}

func MaxFunc[T any](a, b T, less func(a, b T) bool) T {
	if less(b, a) {
		return a
	}
	return b
}

func MinFunc[T any](a, b T, less func(a, b T) bool) T {
	if less(a, b) {
		return a
	}
	return b
}
