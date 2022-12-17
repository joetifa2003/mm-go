package malloc

// #cgo CFLAGS: -I "./include"
// #include "./src/static.c"
import "C"
import "unsafe"

// CMalloc raw binding to mimalloc calloc(1, size)
func Malloc(size int) unsafe.Pointer {
	return C.mi_calloc(1, C.size_t(size))
}

// CMalloc raw binding to mimalloc free
func Free(ptr unsafe.Pointer) {
	C.mi_free(ptr)
}

// CMalloc raw binding to mimalloc realloc
func Realloc(ptr unsafe.Pointer, size int) unsafe.Pointer {
	return C.mi_realloc(ptr, C.size_t(size))
}
