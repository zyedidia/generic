package aggregator

// Array is an array designed to test aggregators.
// The first element of the array is the root node,
// and the rest of the elements are the children of the root node.
// For example, if the size is 4, the structure will be
//
//     [0]
//     ├─ [1]
//     ├─ [2]
//     └─ [3]
type Array[V, A, R any] struct {
	values     []A
	aggregator Aggregator[V, A, R]
}

// Size returns the number of elements in the array.
func (a *Array[V, A, R]) Size() int {
	return len(a.values)
}

func (a *Array[V, A, R]) leaves() []*A {
	leaves := make([]*A, 0, len(a.values)-1)
	for i := 1; i < len(a.values); i++ {
		leaves = append(leaves, &a.values[i])
	}
	return leaves
}

// Get returns the value of index 'key'.
func (a *Array[V, A, R]) Get(key int) V {
	if key > 0 {
		a.aggregator.PushDown(&a.values[0], a.leaves())
	}
	return a.aggregator.Value(&a.values[key])
}

// Put updates the 'value' of index 'key'.
func (a *Array[V, A, R]) Put(key int, value V) {
	a.aggregator.PushDown(&a.values[0], a.leaves())
	a.values[key] = a.aggregator.FromValue(value)
	a.aggregator.PopUp(&a.values[0], a.leaves())
}

// Range returns the aggregator associated with key range [l, r),
// which can be used to obtain statistics or do range-based update.
// Note that the range is only valid before next operation.
func (a *Array[V, A, R]) Range(l, r int) R {
	if l <= 0 && r >= len(a.values) {
		return a.aggregator.RangeView([]*A{&a.values[0]}, []*A{})
	}
	a.aggregator.PushDown(&a.values[0], a.leaves())
	values := make([]*A, 0, r-l)
	for i := l; i < r; i++ {
		values = append(values, &a.values[i])
	}
	return a.aggregator.RangeView([]*A{}, values)
}

// NewArray returns an empty array for test with the given size and aggregator.
func NewArray[V, A, R any](size int, aggregator Aggregator[V, A, R]) *Array[V, A, R] {
	if size < 1 {
		return nil
	}
	return &Array[V, A, R]{
		values:     make([]A, size),
		aggregator: aggregator,
	}
}
