package list_test

import (
	"fmt"

	"github.com/zyedidia/generic/list"
)

func Example() {
	l := list.New[int]()
	l.PushBack(0)
	l.PushBack(1)
	l.PushBack(2)
	l.PushBack(3)

	l.Front.Iter().For(func(i int) {
		fmt.Println(i)
	})
	// Output:
	// 0
	// 1
	// 2
	// 3
}
