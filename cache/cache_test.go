package cache_test

import (
	"testing"

	g "github.com/zyedidia/generic"
	"github.com/zyedidia/generic/cache"
)

func TestSimple(t *testing.T) {
	c := cache.New[g.Int, g.Int](2)

	c.Put(42, 42)
	c.Put(10, 10)
	c.Get(42)
	c.Put(0, 0) // evicts 10

	contents := make([]g.Int, 0, 2)
	c.Iter().For(func(kv cache.KV[g.Int, g.Int]) {
		contents = append(contents, kv.Key)
	})

	if contents[0] != 0 || contents[1] != 42 {
		t.Fatal("incorrect")
	}
}
