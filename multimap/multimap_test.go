package multimap_test

import (
	"testing"

	g "github.com/zyedidia/generic"
	"github.com/zyedidia/generic/multimap"
	"golang.org/x/exp/slices"
)

type entry struct {
	Key   string
	Value int
}

type association struct {
	Key    string
	Values []int
}

type Fixture struct {
	t testing.TB
	m multimap.MultiMap[string, int]

	allowDuplicate bool
	sortedKeys     bool
	sortedValues   bool
}

func (f *Fixture) checkEmpty(key string) {
	if f.m.Has(key) {
		f.t.Errorf("Has(%v) is true", key)
	}
	if n := f.m.Count(key); n != 0 {
		f.t.Errorf("%v count %d is not zero", key, n)
	}
	if list := f.m.Get(key); len(list) != 0 {
		f.t.Errorf("%v values %v is not empty", key, list)
	}
}

func (f *Fixture) checkEach(expected []entry) {
	if size := f.m.Size(); size != len(expected) {
		f.t.Errorf("size %d differs from expected %d", size, len(expected))
	}

	actual := []entry{}
	f.m.Each(func(key string, value int) {
		actual = append(actual, entry{key, value})
	})

	switch {
	case !f.sortedKeys && !f.sortedValues:
		slices.SortFunc(actual, func(a, b entry) bool { return a.Key < b.Key || (a.Key == b.Key && a.Value < b.Value) })
	case !f.sortedKeys:
		slices.SortStableFunc(actual, func(a, b entry) bool { return a.Key < b.Key })
	case !f.sortedValues:
		if !slices.IsSortedFunc(actual, func(a, b entry) bool { return a.Key < b.Key }) {
			f.t.Errorf("multimap %v keys are unsorted", actual)
		}
		slices.SortFunc(actual, func(a, b entry) bool { return a.Key < b.Key || (a.Key == b.Key && a.Value < b.Value) })
	}
	if !slices.Equal(actual, expected) {
		f.t.Errorf("multimap %v differs from expected %v", actual, expected)
	}
}

func (f *Fixture) checkAssociations(expected []association) {
	if dim := f.m.Dimension(); dim != len(expected) {
		f.t.Errorf("dimension %d differs from expected %d", dim, len(expected))
	}

	for _, a := range expected {
		if !f.m.Has(a.Key) {
			f.t.Errorf("Has(%v) is false", a.Key)
		}
		if n := f.m.Count(a.Key); n != len(a.Values) {
			f.t.Errorf("%v count %d differs from expected %d", a.Key, n, len(a.Values))
		}

		list := f.m.Get(a.Key)
		if !f.sortedValues {
			list = slices.Clone(list)
			slices.Sort(list)
		}
		if !slices.Equal(list, a.Values) {
			f.t.Errorf("%v values %v differs from expected %v", a.Key, list, a.Values)
		}
	}

	actual := []association{}
	f.m.EachAssociation(func(key string, values []int) {
		if !f.sortedValues {
			values = slices.Clone(values)
			slices.Sort(values)
		}
		actual = append(actual, association{key, values})
	})

	if !f.sortedKeys {
		slices.SortFunc(actual, func(a, b association) bool { return a.Key < b.Key })
	}
	if !slices.EqualFunc(actual, expected, func(a, b association) bool {
		return a.Key == b.Key && slices.Equal(a.Values, b.Values)
	}) {
		f.t.Errorf("multimap %v differs from expected %v", actual, expected)
	}
}

func testMultiMap(t testing.TB, m multimap.MultiMap[string, int], allowDuplicate, sortedKeys, sortedValues bool) {
	f := Fixture{
		t,
		m,
		allowDuplicate,
		sortedKeys,
		sortedValues,
	}
	f.checkEmpty("A")
	f.checkEach(nil)
	f.checkAssociations(nil)

	m.Put("A", 1)
	m.Put("B", 1)
	m.Put("B", 2)
	m.Put("C", 1)
	m.Put("C", 2)
	m.Put("C", 3)
	f.checkEmpty("D")
	f.checkEach([]entry{{"A", 1}, {"B", 1}, {"B", 2}, {"C", 1}, {"C", 2}, {"C", 3}})
	f.checkAssociations([]association{{"A", []int{1}}, {"B", []int{1, 2}}, {"C", []int{1, 2, 3}}})

	m.Put("C", 2)
	m.Put("C", 2)
	if allowDuplicate {
		f.checkEach([]entry{{"A", 1}, {"B", 1}, {"B", 2}, {"C", 1}, {"C", 2}, {"C", 2}, {"C", 2}, {"C", 3}})
		f.checkAssociations([]association{{"A", []int{1}}, {"B", []int{1, 2}}, {"C", []int{1, 2, 2, 2, 3}}})
		m.Remove("C", 2)
		m.Remove("C", 2)
	} else {
		f.checkEach([]entry{{"A", 1}, {"B", 1}, {"B", 2}, {"C", 1}, {"C", 2}, {"C", 3}})
		f.checkAssociations([]association{{"A", []int{1}}, {"B", []int{1, 2}}, {"C", []int{1, 2, 3}}})
	}

	m.Remove("D", 5)
	m.Remove("C", 4)
	m.Remove("C", 2)
	f.checkEach([]entry{{"A", 1}, {"B", 1}, {"B", 2}, {"C", 1}, {"C", 3}})
	f.checkAssociations([]association{{"A", []int{1}}, {"B", []int{1, 2}}, {"C", []int{1, 3}}})

	m.Remove("C", 3)
	m.Remove("C", 1)
	f.checkEmpty("C")
	f.checkEach([]entry{{"A", 1}, {"B", 1}, {"B", 2}})
	f.checkAssociations([]association{{"A", []int{1}}, {"B", []int{1, 2}}})

	m.RemoveAll("B")
	m.RemoveAll("D")
	f.checkEach([]entry{{"A", 1}})
	f.checkAssociations([]association{{"A", []int{1}}})

	m.Clear()
	f.checkEmpty("A")
	f.checkEach(nil)
	f.checkAssociations(nil)
}

func TestMapSlice(t *testing.T) {
	m := multimap.NewMapSlice[string, int]()
	testMultiMap(t, m, true, false, false)
}

func TestMapSet(t *testing.T) {
	m := multimap.NewMapSet[string](g.Less[int])
	testMultiMap(t, m, false, false, true)
}

func TestAvlSlice(t *testing.T) {
	m := multimap.NewAvlSlice[string, int](g.Less[string])
	testMultiMap(t, m, true, true, false)
}

func TestAvlSet(t *testing.T) {
	m := multimap.NewAvlSet(g.Less[string], g.Less[int])
	testMultiMap(t, m, false, true, true)
}
