package mm

import (
	"context"
	"unsafe"

	"github.com/joetifa2003/mm-go/allocator"
)

func getSize[T any]() int {
	var zeroV T
	return int(unsafe.Sizeof(zeroV))
}

// Alloc allocates T and returns a pointer to it.
func Alloc[T any](ctx context.Context) *T {
	allocator := allocator.GetAllocator(ctx)
	ptr := allocator.Alloc(getSize[T]())
	return (*T)(unsafe.Pointer(ptr))
}

// FreeMany frees memory allocated by Alloc takes a ptr
// CAUTION: be careful not to double free, and prefer using defer to deallocate
func Free[T any](ctx context.Context, ptr *T) {
	allocator := allocator.GetAllocator(ctx)
	allocator.Free(unsafe.Pointer(ptr))
}

// AllocMany allocates n of T and returns a slice representing the heap.
// CAUTION: don't append to the slice, the purpose of it is to replace pointer
// arithmetic with slice indexing
func AllocMany[T any](ctx context.Context, n int) []T {
	allocator := allocator.GetAllocator(ctx)
	ptr := allocator.Alloc(getSize[T]() * n)
	return unsafe.Slice(
		(*T)(ptr),
		n,
	)
}

// FreeMany frees memory allocated by AllocMany takes in the slice (aka the heap)
// CAUTION: be careful not to double free, and prefer using defer to deallocate
func FreeMany[T any](ctx context.Context, slice []T) {
	allocator := allocator.GetAllocator(ctx)
	allocator.Free(unsafe.Pointer(&slice[0]))
}

// Reallocate reallocates memory allocated with AllocMany and doesn't change underling data
func Reallocate[T any](ctx context.Context, slice []T, newN int) []T {
	allocator := allocator.GetAllocator(ctx)
	ptr := allocator.Realloc(unsafe.Pointer(&slice[0]), getSize[T]()*newN)
	return unsafe.Slice(
		(*T)(ptr),
		newN,
	)
}
