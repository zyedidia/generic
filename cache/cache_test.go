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

	c.Each(func(key, val int) {
		fmt.Println("each", key)
	})

	c.Resize(3)
	fmt.Println("size", c.Size())
	fmt.Println("capacity", c.Capacity())

	c.SetEvictCallback(func(key, val int) {
		fmt.Println("evict", key)
	})
	c.Put(1, 1)
	c.Put(2, 2) // evicts 42
	c.Remove(3) // no effect
	c.Resize(1) // evicts 0 and 1

	c.Each(func(key, val int) {
		fmt.Println("each", key)
	})

	// Output:
	// each 0
	// each 42
	// size 2
	// capacity 3
	// evict 42
	// evict 0
	// evict 1
	// each 2
}
