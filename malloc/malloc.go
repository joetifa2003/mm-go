package malloc

import (
	"unsafe"
)

// #include <stdlib.h>
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

// var calloc func(n int, size int) unsafe.Pointer
// var realloc func(ptr unsafe.Pointer, size int) unsafe.Pointer
// var free func(ptr unsafe.Pointer)

// func getSystemLibrary() string {
// 	switch runtime.GOOS {
// 	case "darwin":
// 		return "/usr/lib/libSystem.B.dylib"
// 	case "linux":
// 		return "libc.so.6"
// 	case "windows":
// 		return "ucrtbase.dll"
// 	default:
// 		panic(fmt.Errorf("GOOS=%s is not supported", runtime.GOOS))
// 	}
// }
//
// func init() {
// 	libc, err := openLibrary(getSystemLibrary())
// 	if err != nil {
// 		panic(err)
// 	}
// 	purego.RegisterLibFunc(&calloc, libc, "calloc")
// 	purego.RegisterLibFunc(&realloc, libc, "realloc")
// 	purego.RegisterLibFunc(&free, libc, "free")
// }
