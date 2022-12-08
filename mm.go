package mm

import (
	"unsafe"
)

// #include <stdlib.h>
import "C"

func c_malloc(size int) unsafe.Pointer {
	return C.calloc(1, C.size_t(size))
}

func c_free(ptr unsafe.Pointer) {
	C.free(ptr)
}

func c_realloc(ptr unsafe.Pointer, size int) unsafe.Pointer {
	return C.realloc(ptr, C.size_t(size))
}

func getSize[T any]() int {
	var zeroV T
	return int(unsafe.Sizeof(zeroV))
}

// Alloc allocates T and returns a pointer to it.
func Alloc[T any]() *T {
	ptr := c_malloc(getSize[T]())
	return (*T)(unsafe.Pointer(ptr))
}

// FreeMany frees memory allocated by Alloc takes a ptr
// CAUTION: be careful not to double free, and prefer using defer to deallocate
func Free[T any](ptr *T) {
	c_free(unsafe.Pointer(ptr))
}

// AllocMany allocates n of T and returns a slice representing the heap.
// CAUTION: don't append to the slice, the purpose of it is to replace pointer
// arithmetic with slice indexing
func AllocMany[T any](n int) []T {
	ptr := c_malloc(getSize[T]() * n)
	return unsafe.Slice(
		(*T)(ptr),
		n,
	)
}

// FreeMany frees memory allocated by AllocMany takes in the slice (aka the heap)
// CAUTION: be careful not to double free, and prefer using defer to deallocate
func FreeMany[T any](slice []T) {
	c_free(unsafe.Pointer(&slice[0]))
}

// Reallocate reallocates memory allocated with AllocMany and doesn't change underling data
func Reallocate[T any](slice []T, newN int) []T {
	ptr := c_realloc(unsafe.Pointer(&slice[0]), getSize[T]()*newN)
	return unsafe.Slice(
		(*T)(ptr),
		newN,
	)
}
