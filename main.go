package main

import (
	"unsafe"
)

// #include <stdlib.h>
import "C"

// Alloc allocates T on the heap using mmap
func Alloc[T any]() *T {
	var zeroV T
	size := int(unsafe.Sizeof(zeroV))
	ptr := C.calloc(1, C.size_t(size))
	return (*T)(unsafe.Pointer(ptr))
}

func Free[T any](ptr *T) {
	C.free(unsafe.Pointer(ptr))
}

// AllocMany allocates n of T using mmap and returns a slice representing
// the heap.
func AllocMany[T any](n int) []T {
	var zeroV T
	size := int(unsafe.Sizeof(zeroV))
	ptr := C.calloc(C.size_t(n), C.size_t(size))
	return unsafe.Slice(
		(*T)(ptr),
		n,
	)
}

func FreeMany[T any](ptr []T) {
	C.free(unsafe.Pointer(&ptr[0]))
}

type Value struct {
	x, y, z int
}
