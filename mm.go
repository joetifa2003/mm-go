package mm

import "unsafe"

func SizeOf[T any]() int {
	var zeroV T
	return int(unsafe.Sizeof(zeroV))
}
