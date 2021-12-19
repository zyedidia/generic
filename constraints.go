package generic

func Compare[T any](a, b T, less func(a, b T) bool) int {
	if less(a, b) {
		return -1
	} else if less(b, a) {
		return 1
	}
	return 0
}
