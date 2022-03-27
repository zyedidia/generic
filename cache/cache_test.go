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

	c.Each(func(key int, val int) {
		fmt.Println(key)
	})
	// Output:
	// 0
	// 42
}
