package generic

type Lesser[T any] interface {
	Less(other T) bool
}

func Compare[T Lesser[T]](a, b T) int {
	if a.Less(b) {
		return -1
	} else if b.Less(a) {
		return 1
	}
	return 0
}

type Comparable[T any] interface {
	Equals(other T) bool
}

type Hashable[T any] interface {
	Comparable[T]
	Hash() uint64
}

type Sliceable[T any] interface {
	At(idx int) T
	Slice(low, high int) Sliceable[T]
	Append(s Sliceable[T]) Sliceable[T]
}
