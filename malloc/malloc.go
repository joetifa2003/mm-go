package malloc

// #cgo CFLAGS: -I "./mimalloc/include"
// #include "./mimalloc/src/static.c"
import "C"
import "unsafe"

// CMalloc raw binding to mimalloc calloc(1, size)
func CMalloc(size int) unsafe.Pointer {
	return C.mi_calloc(1, C.size_t(size))
}

// CMalloc raw binding to mimalloc free
func CFree(ptr unsafe.Pointer) {
	C.mi_free(ptr)
}

// CMalloc raw binding to mimalloc realloc
func CRealloc(ptr unsafe.Pointer, size int) unsafe.Pointer {
	return C.mi_realloc(ptr, C.size_t(size))
}
