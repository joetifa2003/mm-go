package allocator

// #include <stdlib.h>
import "C"

import "unsafe"

func NewCallocator() Allocator {
	return NewAllocator(nil, callocator_alloc, callocator_free, callocator_realloc, callocator_destroy)
}

func callocator_alloc(allocator unsafe.Pointer, size int) unsafe.Pointer {
	return C.calloc(1, C.size_t(size))
}

func callocator_free(allocator unsafe.Pointer, ptr unsafe.Pointer) {
	C.free(ptr)
}

func callocator_realloc(allocator unsafe.Pointer, ptr unsafe.Pointer, size int) unsafe.Pointer {
	return C.realloc(ptr, C.size_t(size))
}

func callocator_destroy(allocator unsafe.Pointer) {}
