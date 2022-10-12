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
	// Output:
	// true
	// false
	// false
	// true
	// false
	// false
}
