package cache_test

import (
	"fmt"

	g "github.com/zyedidia/generic"
	"github.com/zyedidia/generic/cache"
)

func Example() {
	c := cache.New[g.Int, g.Int](2)

	c.Put(42, 42)
	c.Put(10, 10)
	c.Get(42)
	c.Put(0, 0) // evicts 10

	c.Iter().For(func(kv cache.KV[g.Int, g.Int]) {
		fmt.Println(kv.Key)
	})
	// Output:
	// 0
	// 42
}
