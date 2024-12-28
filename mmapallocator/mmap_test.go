package mmapallocator_test

import (
	"testing"

	"github.com/joetifa2003/mm-go/allocator"
	"github.com/joetifa2003/mm-go/batchallocator"
	"github.com/joetifa2003/mm-go/mmapallocator"
)

func TestMMapAllocator(t *testing.T) {
	alloc := batchallocator.New(mmapallocator.NewMMapAllocator())

	ptr := allocator.Alloc[int](alloc)
	*ptr = 1
	ptr2 := allocator.Alloc[int](alloc)
	*ptr2 = 2
	allocator.Free(alloc, ptr2)
	allocator.Free(alloc, ptr)
	alloc.Destroy()

	alloc2 := batchallocator.New(mmapallocator.NewMMapAllocator())

	ptr = allocator.Alloc[int](alloc2)
	*ptr = 5

	alloc2.Destroy()
}
