package hashset_test

import (
	"fmt"
	"math/rand"
	"testing"

	g "github.com/zyedidia/generic"
	"github.com/zyedidia/generic/hashset"
)

func checkeq[K any](set *hashset.Set[K], get func(k K) bool, t *testing.T) {
	set.Each(func(key K) {
		if !get(key) {
			t.Fatalf("value %v should be in the set", key)
		}
	})
}

func TestCrossCheck(t *testing.T) {
	stdm := make(map[int]bool)
	set := hashset.New[int](1, g.Equals[int], g.HashInt)

	const nops = 1000
	for i := 0; i < nops; i++ {
		op := rand.Intn(2)
		switch op {
		case 0:
			key := rand.Int()
			stdm[key] = true
			set.Put(key)
		case 1:
			var del int
			for k := range stdm {
				del = k
				break
			}
			delete(stdm, del)
			set.Remove(del)
		}

		checkeq(set, func(k int) bool {
			_, ok := stdm[int(k)]
			return ok
		}, t)
	}
}

func TestOf(t *testing.T) {
	testcases := []struct {
		name  string
		input []string
	}{
		{"init with several items", []string{"foo", "bar", "baz"}},
		{"init without values", []string{}},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			set := hashset.Of[string](10, g.Equals[string], g.HashString, tc.input...)

			if len(tc.input) != set.Size() {
				t.Fatalf("expected %d elements in set, got %d", len(tc.input), set.Size())
			}
			for _, val := range tc.input {
				if !set.Has(val) {
					t.Fatalf("expected to find val '%s' in set but did not", val)
				}
			}
		})
	}
}

func TestUnion(t *testing.T) {
	testcases := []struct {
		name     string
		a        []int
		b        []int
		expected []int
	}{
		{
			name:     "overlapping sets",
			a:        []int{1, 2, 3},
			b:        []int{3, 4, 5},
			expected: []int{1, 2, 3, 4, 5},
		},
		{
			name:     "disjoint sets",
			a:        []int{1, 2},
			b:        []int{3, 4},
			expected: []int{1, 2, 3, 4},
		},
		{
			name:     "first set empty",
			a:        []int{},
			b:        []int{1, 2, 3},
			expected: []int{1, 2, 3},
		},
		{
			name:     "second set empty",
			a:        []int{1, 2, 3},
			b:        []int{},
			expected: []int{1, 2, 3},
		},
		{
			name:     "both empty",
			a:        []int{},
			b:        []int{},
			expected: []int{},
		},
		{
			name:     "identical sets",
			a:        []int{1, 2, 3},
			b:        []int{1, 2, 3},
			expected: []int{1, 2, 3},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			s1 := hashset.Of[int](uint64(len(tc.a)), g.Equals[int], g.HashInt, tc.a...)
			s2 := hashset.Of[int](uint64(len(tc.b)), g.Equals[int], g.HashInt, tc.b...)

			union := s1.Union(s2)

			if union.Size() != len(tc.expected) {
				t.Fatalf("expected size %d, got %d", len(tc.expected), union.Size())
			}

			for _, v := range tc.expected {
				if !union.Has(v) {
					t.Fatalf("expected value %d in union, but was missing", v)
				}
			}

			if s1.Size() != len(tc.a) {
				t.Fatalf("original set s1 was mutated: expected size %d, got %d", len(tc.a), s1.Size())
			}
			if s2.Size() != len(tc.b) {
				t.Fatalf("original set s2 was mutated: expected size %d, got %d", len(tc.b), s2.Size())
			}
		})
	}
}

func ExampleSet_Union() {
	s1 := hashset.Of[string](3, g.Equals[string], g.HashString, "apple", "banana")
	s2 := hashset.Of[string](3, g.Equals[string], g.HashString, "banana", "cherry")

	u := s1.Union(s2)

	fmt.Println(u.Has("apple"))
	fmt.Println(u.Has("banana"))
	fmt.Println(u.Has("cherry"))
	fmt.Println(u.Has("grape"))
	// Output:
	// true
	// true
	// true
	// false
}

func TestIntersection(t *testing.T) {
	testcases := []struct {
		name     string
		a        []int
		b        []int
		expected []int
	}{
		{
			name:     "overlapping sets",
			a:        []int{1, 2, 3},
			b:        []int{3, 4, 5},
			expected: []int{3},
		},
		{
			name:     "disjoint sets",
			a:        []int{1, 2},
			b:        []int{3, 4},
			expected: []int{},
		},
		{
			name:     "first set empty",
			a:        []int{},
			b:        []int{1, 2, 3},
			expected: []int{},
		},
		{
			name:     "second set empty",
			a:        []int{1, 2, 3},
			b:        []int{},
			expected: []int{},
		},
		{
			name:     "both empty",
			a:        []int{},
			b:        []int{},
			expected: []int{},
		},
		{
			name:     "identical sets",
			a:        []int{1, 2, 3},
			b:        []int{1, 2, 3},
			expected: []int{1, 2, 3},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			s1 := hashset.Of[int](uint64(len(tc.a)), g.Equals[int], g.HashInt, tc.a...)
			s2 := hashset.Of[int](uint64(len(tc.b)), g.Equals[int], g.HashInt, tc.b...)

			intersection := s1.Intersection(s2)

			if intersection.Size() != len(tc.expected) {
				t.Fatalf("expected size %d, got %d", len(tc.expected), intersection.Size())
			}

			for _, v := range tc.expected {
				if !intersection.Has(v) {
					t.Fatalf("expected value %d in intersection, but was missing", v)
				}
			}

			if s1.Size() != len(tc.a) {
				t.Fatalf("original set s1 was mutated: expected size %d, got %d", len(tc.a), s1.Size())
			}
			if s2.Size() != len(tc.b) {
				t.Fatalf("original set s2 was mutated: expected size %d, got %d", len(tc.b), s2.Size())
			}
		})
	}
}

func Example() {
	set := hashset.New[string](3, g.Equals[string], g.HashString)
	set.Put("foo")
	set.Put("bar")
	set.Put("baz")

	fmt.Println(set.Has("foo"))
	fmt.Println(set.Has("quux"))

	set.Remove("foo")

	fmt.Println(set.Has("foo"))
	fmt.Println(set.Has("bar"))

	set.Clear()

	fmt.Println(set.Has("foo"))
	fmt.Println(set.Has("bar"))

	s1 := hashset.Of[string](3, g.Equals[string], g.HashString, "apple", "banana")
	s2 := hashset.Of[string](3, g.Equals[string], g.HashString, "banana", "cherry")

	i := s1.Intersection(s2)

	fmt.Println(i.Has("apple"))
	fmt.Println(i.Has("banana"))
	fmt.Println(i.Has("cherry"))

	// Output:
	// true
	// false
	// false
	// true
	// false
	// false
	// false
	// true
	// false
}
