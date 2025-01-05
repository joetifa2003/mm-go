package mmapallocator_test

import (
	"testing"

	"github.com/joetifa2003/mm-go/allocator"
	"github.com/joetifa2003/mm-go/mmapallocator"
)

func TestMMapAllocator(t *testing.T) {
	alloc := mmapallocator.NewMMapAllocator()

	ptr := allocator.Alloc[int](alloc)
	*ptr = 1
	ptr2 := allocator.Alloc[int](alloc)
	*ptr2 = 2
	allocator.Free(alloc, ptr2)
	ptr3 := allocator.Alloc[int](alloc)
	_ = ptr3
	allocator.Free(alloc, ptr)
	allocator.Free(alloc, ptr3)
	alloc.Destroy()
}
