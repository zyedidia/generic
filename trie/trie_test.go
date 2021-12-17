package trie_test

import (
	"fmt"

	"github.com/zyedidia/generic/trie"
)

func Example() {
	tr := trie.New[int]()
	tr.Put([]byte("foo"), 1)
	tr.Put([]byte("fo"), 2)
	tr.Put([]byte("bar"), 3)

	fmt.Println(tr.Contains([]byte("f")))
	fmt.Println(tr.KeysWithPrefix([]byte("")))
}
