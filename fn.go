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

func MaxG[T Lesser[T]](a, b T) T {
	if b.Less(a) {
		return a
	}
	return b
}

func MinG[T Lesser[T]](a, b T) T {
	if a.Less(b) {
		return a
	}
	return b
}
