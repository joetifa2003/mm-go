package allocator

import "unsafe"

// Allocator is an interface that defines some methods needed for most allocators.
// It's not a golang interface, so it's safe to use in manually managed structs (will not get garbage collected).
type Allocator struct {
	allocator unsafe.Pointer
	alloc     func(allocator unsafe.Pointer, size int) unsafe.Pointer
	free      func(allocator unsafe.Pointer, ptr unsafe.Pointer)
	realloc   func(allocator unsafe.Pointer, ptr unsafe.Pointer, size int) unsafe.Pointer
	destroy   func(allocator unsafe.Pointer)
}

// NewAllocator creates a new Allocator
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

// Alloc allocates size bytes and returns an unsafe pointer to it.
func (a Allocator) Alloc(size int) unsafe.Pointer {
	return a.alloc(a.allocator, size)
}

// Free frees the memory pointed by ptr
func (a Allocator) Free(ptr unsafe.Pointer) {
	a.free(a.allocator, ptr)
}

// Realloc reallocates the memory pointed by ptr with a new size and returns a new pointer to it.
func (a Allocator) Realloc(ptr unsafe.Pointer, size int) unsafe.Pointer {
	return a.realloc(a.allocator, ptr, size)
}

// Destroy destroys the allocator.
// After calling this, the allocator is no longer usable.
// This is useful for cleanup, freeing allocator internal resources, etc.
func (a Allocator) Destroy() {
	a.destroy(a.allocator)
}

func getSize[T any]() int {
	var zeroV T
	return int(unsafe.Sizeof(zeroV))
}

// Alloc allocates T and returns a pointer to it.
func Alloc[T any](a Allocator) *T {
	ptr := a.alloc(a.allocator, getSize[T]())
	return (*T)(unsafe.Pointer(ptr))
}

// FreeMany frees memory allocated by Alloc takes a ptr
// CAUTION: be careful not to double free, and prefer using defer to deallocate
func Free[T any](a Allocator, ptr *T) {
	a.free(a.allocator, unsafe.Pointer(ptr))
}

// AllocMany allocates n of T and returns a slice representing the heap.
// CAUTION: don't append to the slice, the purpose of it is to replace pointer
// arithmetic with slice indexing
func AllocMany[T any](a Allocator, n int) []T {
	ptr := a.alloc(a.allocator, getSize[T]()*n)
	return unsafe.Slice(
		(*T)(ptr),
		n,
	)
}

// FreeMany frees memory allocated by AllocMany takes in the slice (aka the heap)
// CAUTION: be careful not to double free, and prefer using defer to deallocate
func FreeMany[T any](a Allocator, slice []T) {
	a.free(a.allocator, unsafe.Pointer(&slice[0]))
}

// Realloc reallocates memory allocated with AllocMany and doesn't change underling data
func Realloc[T any](a Allocator, slice []T, newN int) []T {
	ptr := a.realloc(a.allocator, unsafe.Pointer(&slice[0]), getSize[T]()*newN)
	return unsafe.Slice(
		(*T)(ptr),
		newN,
	)
}
