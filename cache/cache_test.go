package cache_test

import (
	"fmt"

	"github.com/zyedidia/generic/cache"
)

func Example() {
	c := cache.New[int, int](2)

	c.Put(42, 42)
	c.Put(10, 10)
	c.Get(42)
	c.Put(0, 0) // evicts 10

	c.Iter().For(func(kv cache.KV[int, int]) {
		fmt.Println(kv.Key)
	})
	// Output:
	// 0
	// 42
}
