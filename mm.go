package mm

import (
	"unsafe"
)

// #include <stdlib.h>
import "C"

// Alloc allocates T and returns a pointer to it.
func Alloc[T any]() *T {
	var zeroV T
	size := int(unsafe.Sizeof(zeroV))
	ptr := C.calloc(1, C.size_t(size))
	return (*T)(unsafe.Pointer(ptr))
}

// FreeMany frees memory allocated by Alloc takes a ptr
// CAUTION: be careful not to double free, and prefer using defer to deallocate
func Free[T any](ptr *T) {
	C.free(unsafe.Pointer(ptr))
}

// AllocMany allocates n of T and returns a slice representing the heap.
func AllocMany[T any](n int) []T {
	var zeroV T
	size := int(unsafe.Sizeof(zeroV))
	ptr := C.calloc(C.size_t(n), C.size_t(size))
	return unsafe.Slice(
		(*T)(ptr),
		n,
	)
}

// FreeMany frees memory allocated by AllocMany takes in the slice (aka the heap)
// CAUTION: be careful not to double free, and prefer using defer to deallocate
func FreeMany[T any](slice []T) {
	C.free(unsafe.Pointer(&slice[0]))
}

// Reallocate reallocates memory allocated with AllocMany and doesn't change underling data
func Reallocate[T any](slice []T, newN int) []T {
	var zeroV T
	size := int(unsafe.Sizeof(zeroV))
	ptr := C.realloc(unsafe.Pointer(&slice[0]), C.size_t(size*newN))
	return unsafe.Slice(
		(*T)(ptr),
		newN,
	)
}
