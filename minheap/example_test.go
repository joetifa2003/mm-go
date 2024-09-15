package minheap_test

import (
	"fmt"

	"github.com/joetifa2003/mm-go/allocator"
	"github.com/joetifa2003/mm-go/minheap"
)

func int_less(a, b int) bool { return a < b }

func Example() {
	alloc := allocator.NewC()
	defer alloc.Destroy()

	h := minheap.New[int](alloc, int_less)

	// Push some values onto the heap
	h.Push(2)
	h.Push(1)
	h.Push(4)
	h.Push(3)
	h.Push(5)

	// Pop the minimum value from the heap
	fmt.Println(h.Pop())
	fmt.Println(h.Pop())

	// Output:
	// 1
	// 2
}

func int_greater(a, b int) bool { return a > b }

func Example_MaxHeap() {
	alloc := allocator.NewC()
	defer alloc.Destroy()

	h := minheap.New[int](alloc, int_greater)

	// Push some values onto the heap
	h.Push(2)
	h.Push(1)
	h.Push(4)
	h.Push(3)
	h.Push(5)

	// Pop the max value from the heap
	fmt.Println(h.Pop())
	fmt.Println(h.Pop())

	// Output:
	// 5
	// 4
}
