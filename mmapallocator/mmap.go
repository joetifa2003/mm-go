package mmapallocator

import (
	"os"
	"unsafe"

	"github.com/edsrzf/mmap-go"

	"github.com/joetifa2003/mm-go"
	"github.com/joetifa2003/mm-go/allocator"
)

const (
	sizeClassCount = 9
	bigSizeClass   = -1
)

var sizeClasses = [sizeClassCount]int{32, 64, 128, 256, 512, 1024, 2048, 4096, bigSizeClass}

func getSizeClass(size int) (int, int) {
	for i, sizeClass := range sizeClasses {
		if size <= sizeClass {
			return i, sizeClass
		}
	}

	return len(sizeClasses) - 1, bigSizeClass
}

var freeLists [sizeClassCount]*block

type block struct {
	data  unsafe.Pointer
	chunk *chunk
	next  *block
	size  int
}

type ptrMeta struct {
	block *block
}

type chunk struct {
	data   unsafe.Pointer
	offset uintptr
	size   uintptr
	ptrs   uintptr
}

var pageSize = os.Getpagesize() * 15

const (
	sizeOfPtrMeta = unsafe.Sizeof(ptrMeta{})
	sizeOfBlock   = unsafe.Sizeof(block{})
	sizeOfChunk   = unsafe.Sizeof(chunk{})
	alignment     = unsafe.Alignof(uintptr(0))
)

func NewMMapAllocator() allocator.Allocator {
	return allocator.NewAllocator(nil, mmap_alloc, mmap_free, mmap_realloc, mmap_destroy)
}

func mmap_alloc(allocator unsafe.Pointer, size int) unsafe.Pointer {
	sizeClassIdx, sizeClass := getSizeClass(int(sizeOfPtrMeta) + size)

	b := freeLists[sizeClassIdx]
	if b != nil {
		ptrMeta := (*ptrMeta)(unsafe.Pointer(uintptr(b.data)))
		ptrMeta.block = b
		b.chunk.ptrs++
		freeLists[sizeClassIdx] = b.next
		return unsafe.Pointer(uintptr(unsafe.Pointer(ptrMeta)) + sizeOfPtrMeta)
	}

	// init size class
	totalSize := mm.Align(sizeOfChunk+uintptr(sizeClass)+sizeOfBlock, uintptr(pageSize))
	m, err := mmap.MapRegion(nil, int(totalSize), mmap.RDWR, mmap.ANON, 0)
	if err != nil {
		panic(err)
	}

	chunk := (*chunk)(unsafe.Pointer(&m[0]))
	chunk.data = unsafe.Pointer(unsafe.Pointer(&m[sizeOfChunk]))
	chunk.offset = 0
	chunk.size = totalSize - sizeOfChunk

	nBlocks := chunk.size / (uintptr(sizeClass) + sizeOfBlock)
	for i := uintptr(0); i < nBlocks; i++ {
		b := (*block)(unsafe.Pointer(uintptr(chunk.data) + uintptr(i*sizeOfBlock)))
		b.data = unsafe.Pointer(uintptr(chunk.data) + uintptr(i*(uintptr(sizeClass)+sizeOfBlock)))
		b.next = freeLists[sizeClassIdx]
		b.size = sizeClass
		b.chunk = chunk
		freeLists[sizeClassIdx] = b
	}

	return mmap_alloc(allocator, size)
}

func mmap_free(allocator unsafe.Pointer, ptr unsafe.Pointer) {
	ptrMeta := (*ptrMeta)(unsafe.Pointer(uintptr(ptr) - sizeOfPtrMeta))
	block := ptrMeta.block
	sizeClassIdx, _ := getSizeClass(block.size)
	block.chunk.ptrs--
	head := freeLists[sizeClassIdx]
	block.next = head
	freeLists[sizeClassIdx] = block
}

func mmap_realloc(allocator unsafe.Pointer, ptr unsafe.Pointer, size int) unsafe.Pointer {
	return nil
}

func mmap_destroy(allocator unsafe.Pointer) {}
