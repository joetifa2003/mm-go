package malloc

import (
	"unsafe"
)

//#include <stdlib.h>
import "C"

// CMalloc raw binding to c calloc(1, size)
func Malloc(size int) unsafe.Pointer {
	return C.calloc(1, C.size_t(size))
}

// CMalloc raw binding to c free
func Free(ptr unsafe.Pointer) {
	C.free(ptr)
}

// CMalloc raw binding to c realloc
func Realloc(ptr unsafe.Pointer, size int) unsafe.Pointer {
	return C.realloc(ptr, C.size_t(size))
}
