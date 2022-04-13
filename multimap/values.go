package multimap

import (
	"github.com/zyedidia/generic/avl"
	"golang.org/x/exp/slices"
)

type valuesContainer[V any] interface {
	Empty() bool
	Size() int
	Put(value V) int
	Remove(value V) int
	List() []V
	Each(fn func(value V))
}

var (
	_ valuesContainer[int] = valuesSet[int]{}
	_ valuesContainer[int] = (*valuesSlice[int])(nil)
)

type valuesSet[V any] struct {
	t *avl.Tree[V, struct{}]
}

func (vs valuesSet[V]) Empty() bool {
	return vs.t.Height() == 0
}

func (vs valuesSet[V]) Size() int {
	return vs.t.Size()
}

func (vs valuesSet[V]) has(value V) bool {
	_, ok := vs.t.Get(value)
	return ok
}

func (vs valuesSet[V]) Put(value V) int {
	if vs.has(value) {
		return 0
	}
	vs.t.Put(value, struct{}{})
	return 1
}

func (vs valuesSet[V]) Remove(value V) int {
	if !vs.has(value) {
		return 0
	}
	vs.t.Remove(value)
	return 1
}

func (vs valuesSet[V]) List() (values []V) {
	vs.Each(func(value V) {
		values = append(values, value)
	})
	return
}

func (vs valuesSet[V]) Each(fn func(value V)) {
	vs.t.Each(func(value V, _ struct{}) {
		fn(value)
	})
}

type valuesSlice[V comparable] []V

func (vs *valuesSlice[V]) Empty() bool {
	return len(*vs) == 0
}

func (vs *valuesSlice[V]) Size() int {
	return len(*vs)
}

func (vs *valuesSlice[V]) Put(value V) int {
	*vs = append(*vs, value)
	return 1
}

func (vs *valuesSlice[V]) Remove(value V) int {
	i := slices.Index(*vs, value)
	if i < 0 {
		return 0
	}
	(*vs)[i] = (*vs)[len(*vs)-1]
	*vs = (*vs)[:len(*vs)-1]
	return 1
}

func (vs *valuesSlice[V]) List() []V {
	return *vs
}

func (vs *valuesSlice[V]) Each(fn func(value V)) {
	for _, value := range *vs {
		fn(value)
	}
}
