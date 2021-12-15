package generic

type Orderable[T any] interface {
	Compare(other T) int
}
