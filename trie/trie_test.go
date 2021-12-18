package trie_test

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/zyedidia/generic/trie"
)

var letters = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randstring(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func checkeq(tr *trie.Trie[int], m map[string]int, t *testing.T) {
	for k, v := range m {
		tv, ok := tr.Get(k)
		if !ok {
			t.Fatalf("%v should exist", k)
		}
		if tv != v {
			t.Fatalf("%v != %v, key: %v", tv, v, k)
		}
	}
}

func TestCrossCheck(t *testing.T) {
	stdm := make(map[string]int)
	tree := trie.New[int]()

	const nops = 1000

	for i := 0; i < nops; i++ {
		key := randstring(rand.Intn(20) + 1)
		val := rand.Int()
		op := rand.Intn(2)

		switch op {
		case 0:
			stdm[key] = val
			tree.Put(key, val)
		case 1:
			var del string
			for k := range stdm {
				del = k
				break
			}
			delete(stdm, del)
			tree.Remove(del)
		}

		checkeq(tree, stdm, t)
	}
}

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
