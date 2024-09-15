package mm

import "unsafe"

// SizeOf returns the size of T in bytes
func SizeOf[T any]() int {
	var zeroV T
	return int(unsafe.Sizeof(zeroV))
}

// Zero returns a zero value of T
func Zero[T any]() T {
	var zeroV T
	return zeroV
}
