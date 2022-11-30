package main

import (
	"unsafe"

	"github.com/edsrzf/mmap-go"
)

// allocate a []byte with size using mmap
func allocBytes(size int) []byte {
	bytes, err := mmap.MapRegion(nil, size, mmap.RDWR, mmap.ANON, 0)
	if err != nil {
		panic(err)
	}

	return bytes
}

// converts a pointer to a MMap ([]byte) with specified size
func convertToMMap[T any](ptr *T, size int) mmap.MMap {
	return mmap.MMap(unsafe.Slice((*byte)(unsafe.Pointer(ptr)), size))
}

// Alloc allocates T on the heap using mmap
func Alloc[T any]() *T {
	var zeroV T
	size := int(unsafe.Sizeof(zeroV))
	bytes := allocBytes(size)
	return (*T)(unsafe.Pointer(&bytes[0]))
}

// AllocMany allocates n of T using mmap and returns a slice representing
// the heap.
func AllocMany[T any](n int) []T {
	var zeroV T
	size := int(unsafe.Sizeof(zeroV))
	bytes := allocBytes(size * n)
	return unsafe.Slice(
		(*T)(unsafe.Pointer(&bytes[0])),
		n,
	)
}

func FreeMany[T any](slice []T) {
	var zeroV T
	size := int(unsafe.Sizeof(zeroV))
	mmap := convertToMMap(&slice[0], len(slice)*size)
	err := mmap.Unmap()
	if err != nil {
		panic(err)
	}
}

func Free[T any](ptr *T) {
	var zeroV T
	size := int(unsafe.Sizeof(zeroV))
	mmap := convertToMMap(ptr, size)
	err := mmap.Unmap()
	if err != nil {
		panic(err)
	}
}

const NNodes = 100000

type Node struct {
	Value int
	Prev  *Node
	Next  *Node
}

func unManaged() {
	allocatedNodes := AllocMany[Node](NNodes)
	defer FreeMany(allocatedNodes)

	for j := 0; j < NNodes; j++ {
		var prev *Node
		var next *Node
		if j != 0 {
			prev = &allocatedNodes[j-1]
		}
		if j != NNodes-1 {
			next = &allocatedNodes[j+1]
		}

		allocatedNodes[j] = Node{
			Value: j,
			Prev:  prev,
			Next:  next,
		}
	}
}

func managed() {
	nodes := make([]*Node, NNodes)

	for j := 0; j < NNodes; j++ {
		var prev *Node
		var next *Node
		if j != 0 {
			prev = nodes[j-1]
		}
		if j != NNodes-1 {
			next = nodes[j+1]
		}

		nodes[j] = &Node{
			Value: j,
			Prev:  prev,
			Next:  next,
		}
	}
}

func main() {
	unManaged()
}
