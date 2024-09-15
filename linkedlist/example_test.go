package linkedlist

import (
	"fmt"

	"github.com/joetifa2003/mm-go/allocator"
)

func Example() {
	alloc := allocator.NewC()
	defer alloc.Destroy()

	ll := New[int](alloc)
	defer ll.Free()

	ll.PushBack(1)
	ll.PushBack(2)
	ll.PushBack(3)
	ll.PushBack(4)

	fmt.Println("PopBack:", ll.PopBack())
	fmt.Println("PopFront:", ll.PopFront())

	for _, i := range ll.Iter() {
		fmt.Println(i)
	}

	// Output:
	// PopBack: 4
	// PopFront: 1
	// 2
	// 3
}
