package malloc

// #cgo CFLAGS: -I "./mimalloc/include"
// #include "./mimalloc/src/static.c"
import "C"
import "unsafe"

func CMalloc(size int) unsafe.Pointer {
	return C.mi_calloc(1, C.size_t(size))
}

func CFree(ptr unsafe.Pointer) {
	C.mi_free(ptr)
}

func CRealloc(ptr unsafe.Pointer, size int) unsafe.Pointer {
	return C.mi_realloc(ptr, C.size_t(size))
}
