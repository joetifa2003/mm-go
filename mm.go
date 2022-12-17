package mm

import (
	"unsafe"

	"github.com/joetifa2003/mm-go/malloc"
)

func getSize[T any]() int {
	var zeroV T
	return int(unsafe.Sizeof(zeroV))
}

// Alloc allocates T and returns a pointer to it.
func Alloc[T any]() *T {
	ptr := malloc.CMalloc(getSize[T]())
	return (*T)(unsafe.Pointer(ptr))
}

// FreeMany frees memory allocated by Alloc takes a ptr
// CAUTION: be careful not to double free, and prefer using defer to deallocate
func Free[T any](ptr *T) {
	malloc.CFree(unsafe.Pointer(ptr))
}

// AllocMany allocates n of T and returns a slice representing the heap.
// CAUTION: don't append to the slice, the purpose of it is to replace pointer
// arithmetic with slice indexing
func AllocMany[T any](n int) []T {
	ptr := malloc.CMalloc(getSize[T]() * n)
	return unsafe.Slice(
		(*T)(ptr),
		n,
	)
}

// FreeMany frees memory allocated by AllocMany takes in the slice (aka the heap)
// CAUTION: be careful not to double free, and prefer using defer to deallocate
func FreeMany[T any](slice []T) {
	malloc.CFree(unsafe.Pointer(&slice[0]))
}

// Reallocate reallocates memory allocated with AllocMany and doesn't change underling data
func Reallocate[T any](slice []T, newN int) []T {
	ptr := malloc.CRealloc(unsafe.Pointer(&slice[0]), getSize[T]()*newN)
	return unsafe.Slice(
		(*T)(ptr),
		newN,
	)
}
