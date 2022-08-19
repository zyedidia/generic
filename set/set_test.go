package set

import (
	"fmt"
	"sort"
	"testing"

	"github.com/zyedidia/generic"
)

func ExampleSet_ConstUnion() {
	fmt.Print(NewMapset(1, 4, 7).ConstUnion(2, 3, 5, 6))
	// Output: [1 2 3 4 5 6 7]
}

func ExampleSet_ConstDifference() {
	fmt.Print(NewMapset(1.2, 1.8, 2.6, 3.5).ConstDifference(1.2, 2.6))
	// Output: [1.8 3.5]
}

func ExampleSet_ConstIntersection() {
	fmt.Print(NewMapset("a", "b", "c").ConstIntersection("b", "c", "e"))
	// Output: [b c]
}

func ExampleSet_ConstSymmetricDifference() {
	fmt.Print(NewMapset(1, 2, 3).ConstSymmetricDifference(2, 3, 4))
	// Output: [1 4]
}

func ExampleSet_SymmetricDifference() {
	one := NewMapset(2, 3, 4)
	two := NewMapset(4, 5, 6)
	fmt.Print(NewMapset(1, 2, 3).SymmetricDifference(one, two))
	// Output: [1 5 6]
}

func ExampleSet_Union() {
	one := NewMapset(2, 3, 4)
	two := NewMapset(4, 5, 6)
	fmt.Print(NewMapset(1, 2, 3).Union(one, two))
	// Output: [1 2 3 4 5 6]
}

func ExampleSet_InPlaceIntersection() {
	one := NewMapset(2, 3, 4)
	two := NewMapset(4, 5, 6)
	one.InPlaceIntersection(two)
	fmt.Print(one)
	// Output: [4]
}

func ExampleSet_Keys() {
	keys := NewMapset("one", "two").Keys()
	sort.Strings(keys)
	fmt.Println(keys)
	// Output:
	// [one two]
}

func ExampleSet_Intersection() {
	s := NewMapset(1, 4, 7)
	o := NewHashset(1, generic.Equals[int], generic.HashInt, 1, 7)
	inter := s.Intersection(o)

	fmt.Println(inter)
	fmt.Printf("%T", inter.SetOf) // the same type as the receiver
	// Output:
	// [1 7]
	// mapset.Set[int]
}

func ExampleSet_Difference() {
	s := NewHashset(1, generic.Equals[int], generic.HashInt, 3, 4, 5, 6, 7)
	o := NewMapset(5)
	diff := s.Difference(o)

	fmt.Println(diff)
	fmt.Printf("%T", diff.SetOf) // the same type as the receiver
	// Output:
	// [3 4 6 7]
	// *hashset.Set[int]
}

func TestSetTypes(t *testing.T) {
	type something struct{ string }
	set := NewMapset(something{"hello"}, something{"world"})
	if !set.Has(something{"hello"}) {
		t.Errorf(`set is missing a value`)
	}
	if set.Has(something{"mystery"}) {
		t.Errorf(`set has an unexpected value`)
	}
}

func FuzzDifference(f *testing.F) {
	f.Fuzz(func(t *testing.T, needle, hay1, hay2 int) {
		found := needle == hay1 || needle == hay2
		search := NewMapset(hay1, hay2)
		initialize := search.Size()

		diff := search.Difference(NewMapset(needle))
		t.Logf("diff %d,%d - %d: %v", hay1, hay2, needle, diff.Keys())
		if found != (diff.Size() < initialize) {
			t.Error("unexpected result from diff")
		}

		inter := search.Intersection(NewMapset(needle))
		t.Logf("inter %d,%d + %d: %v", hay1, hay2, needle, diff.Keys())
		if found == (inter.Size() == 0) {
			t.Error("unexpected result from inter")
		}
	})
}
