package trie_test

import (
	"fmt"
	"math/rand"
	"sort"
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
	keys := make([]string, 0, len(m))
	for k, v := range m {
		keys = append(keys, k)
		tv, ok := tr.Get(k)
		if !ok {
			t.Fatalf("%v should exist", k)
		}
		if tv != v {
			t.Fatalf("%v != %v, key: %v", tv, v, k)
		}
	}
	trieKeys := tr.Keys()
	sort.Strings(keys)
	for i := range keys {
		if keys[i] != trieKeys[i] {
			t.Fatalf("%d, %s != %s", i, keys, trieKeys)
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

func TestKeys(t *testing.T) {
	tr := trie.New[int]()
	tr.Put("topic1", 1)
	tr.Put("topic2", 2)

	keys := tr.Keys()
	if len(keys) != 2 || keys[0] != "topic1" || keys[1] != "topic2" {
		t.Fatal(keys)
	}
}

func Example() {
	tr := trie.New[int]()
	tr.Put("f§oo", 1)
	tr.Put("f§o", 2)
	tr.Put("bar", 3)

	fmt.Println(tr.Contains("§"))
	fmt.Println(tr.KeysWithPrefix(""))
	fmt.Println(tr.KeysWithPrefix("f§"))
	// Output:
	// false
	// [bar f§o f§oo]
	// [f§o f§oo]
}

func TestRemoveSize(t *testing.T) {
	tr := trie.New[int]()
	tr.Remove("missing")
	if tr.Size() != 0 {
		t.Fatalf("removing from an empty trie changed its size to %d", tr.Size())
	}
	tr.Put("cat", 1)
	tr.Put("car", 2)
	for _, key := range []string{"", "ca", "cab", "cats", "other"} {
		tr.Remove(key)
		if tr.Size() != 2 || !tr.Contains("cat") || !tr.Contains("car") {
			t.Fatalf("removing absent key %q changed the trie: size=%d keys=%v", key, tr.Size(), tr.Keys())
		}
	}
	tr.Remove("cat")
	tr.Remove("cat")
	if tr.Size() != 1 || tr.Contains("cat") || !tr.Contains("car") {
		t.Fatalf("repeated removal corrupted size or keys: size=%d keys=%v", tr.Size(), tr.Keys())
	}
	tr.Remove("car")
	if tr.Size() != 0 || len(tr.Keys()) != 0 {
		t.Fatalf("trie should be empty: size=%d keys=%v", tr.Size(), tr.Keys())
	}
}
