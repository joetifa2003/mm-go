package mm

import "unsafe"

// SizeOf returns the size of T in bytes
func SizeOf[T any]() int {
	var zeroV T
	return int(unsafe.Sizeof(zeroV))
}
