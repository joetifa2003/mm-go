package minheap

import (
	"testing"

	"github.com/joetifa2003/mm-go/allocator"
)

func TestMinHeap(t *testing.T) {
	alloc := allocator.NewC()
	heap := New[int](alloc, func(a, b int) bool { return a < b })

	heap.Push(3)
	heap.Push(4)
	heap.Push(1)
	heap.Push(0)
}
