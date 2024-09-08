package batchallocator

import (
	"testing"
	"unsafe"

	"github.com/stretchr/testify/require"

	"github.com/joetifa2003/mm-go/allocator"
)

func TestBatchAllocator(t *testing.T) {
	assert := require.New(t)

	alloc := New(allocator.NewCallocator())

	i := allocator.Alloc[int](alloc)
	*i = 1
	j := allocator.Alloc[int](alloc)
	*j = 2

	assert.Equal(1, *i)
	assert.Equal(2, *j)

	arr := allocator.Alloc[[3000]int](alloc)
	for i := 0; i < 32; i++ {
		arr[i] = i
	}

	allocator.Free(alloc, i)

	assert.Equal(2, *j)

	allocator.Free(alloc, j)

	allocator.Free(alloc, arr)

	alloc.Destroy()
}

func TestBatchAllocatorAligned(t *testing.T) {
	assert := require.New(t)

	alloc := New(allocator.NewCallocator())

	alloc.Alloc(13)
	alloc.Alloc(11)
	y := allocator.Alloc[int](alloc)
	*y = 2

	assert.Equal(0, int(uintptr(unsafe.Pointer(y))%8))
}
