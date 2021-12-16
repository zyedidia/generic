package iter_test

import (
	"testing"

	"github.com/zyedidia/generic/iter"
)

func TestSliceIter(t *testing.T) {
	slice := []string{"foo", "bar", "baz", "quux"}
	it := iter.Slice(slice)
	var i int
	for val, ok := it(); ok; val, ok = it() {
		if slice[i] != val {
			t.Fatal("incorrect value")
		}
		i++
	}
}

func TestMapIter(t *testing.T) {
	m := map[string]int{
		"foo": 0,
		"bar": 1,
		"baz": 2,
	}

	it := iter.Map(m)
	it.For(func(i iter.KV[string, int]) {
		if i.Val != m[i.Key] {
			t.Fatal("incorrect")
		}
	})
}
