package trie_test

import (
	"fmt"

	"github.com/zyedidia/generic/trie"
)

func Example() {
	tr := trie.New[int]()
	tr.Put("foo", 1)
	tr.Put("fo", 2)
	tr.Put("bar", 3)

	fmt.Println(tr.Contains("f"))
	fmt.Println(tr.KeysWithPrefix(""))
	fmt.Println(tr.KeysWithPrefix("f"))
	// Output:
	// false
	// [bar fo foo]
	// [fo foo]
}
