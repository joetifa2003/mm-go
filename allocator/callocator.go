package allocator

import (
	"unsafe"

	"github.com/joetifa2003/mm-go/malloc"
)

type CAllocator struct{}

func (c CAllocator) Alloc(size int) unsafe.Pointer {
	return malloc.Malloc(size)
}

func (c CAllocator) Free(ptr unsafe.Pointer) {
	malloc.Free(ptr)
}

func (c CAllocator) Realloc(ptr unsafe.Pointer, size int) unsafe.Pointer {
	return malloc.Realloc(ptr, int(size))
}

func (c CAllocator) Destroy() {}

func NewC() Allocator {
	return CAllocator{}
}
