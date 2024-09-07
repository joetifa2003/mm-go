package allocator

import "unsafe"

type Allocator struct {
	allocator unsafe.Pointer
	alloc     func(allocator unsafe.Pointer, size int) unsafe.Pointer
	free      func(allocator unsafe.Pointer, ptr unsafe.Pointer)
	realloc   func(allocator unsafe.Pointer, ptr unsafe.Pointer, size int) unsafe.Pointer
	destroy   func(allocator unsafe.Pointer)
}

func NewAllocator(
	allocator unsafe.Pointer,
	alloc func(allocator unsafe.Pointer, size int) unsafe.Pointer,
	free func(allocator unsafe.Pointer, ptr unsafe.Pointer),
	realloc func(allocator unsafe.Pointer, ptr unsafe.Pointer, size int) unsafe.Pointer,
	destroy func(allocator unsafe.Pointer),
) Allocator {
	return Allocator{
		allocator: allocator,
		alloc:     alloc,
		free:      free,
		realloc:   realloc,
		destroy:   destroy,
	}
}

func (a Allocator) Alloc(size int) unsafe.Pointer {
	return a.alloc(a.allocator, size)
}

func (a Allocator) Free(ptr unsafe.Pointer) {
	a.free(a.allocator, ptr)
}

func (a Allocator) Realloc(ptr unsafe.Pointer, size int) unsafe.Pointer {
	return a.realloc(a.allocator, ptr, size)
}

func (a Allocator) Destroy() {
	a.destroy(a.allocator)
}

func getSize[T any]() int {
	var zeroV T
	return int(unsafe.Sizeof(zeroV))
}

// Alloc allocates T and returns a pointer to it.
func Alloc[T any](a Allocator) *T {
	ptr := a.Alloc(getSize[T]())
	return (*T)(unsafe.Pointer(ptr))
}

// FreeMany frees memory allocated by Alloc takes a ptr
func Free[T any](a Allocator, ptr *T) {
	a.Free(unsafe.Pointer(ptr))
}

func AllocMany[T any](a Allocator, n int) []T {
	ptr := a.Alloc(getSize[T]() * n)
	return unsafe.Slice(
		(*T)(ptr),
		n,
	)
}

func FreeMany[T any](a Allocator, slice []T) {
	a.Free(unsafe.Pointer(&slice[0]))
}
