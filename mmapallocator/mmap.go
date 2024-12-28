package mmapallocator

import (
	"os"
	"syscall"
	"unsafe"

	"github.com/edsrzf/mmap-go"

	"github.com/joetifa2003/mm-go"
	"github.com/joetifa2003/mm-go/allocator"
)

type mmapMeta struct {
	size  int
	chunk *chunk
}

type chunk struct {
	data   unsafe.Pointer
	offset uintptr
	size   uintptr
	ptrs   uintptr
	next   *chunk
	prev   *chunk
}

var pageSize = os.Getpagesize() * 15

var chunkHead *chunk

const (
	sizeOfPtrMeta = unsafe.Sizeof(mmapMeta{})
	sizeOfChunk   = unsafe.Sizeof(chunk{})
	alignment     = unsafe.Alignof(uintptr(0))
)

func NewMMapAllocator() allocator.Allocator {
	return allocator.NewAllocator(nil, mmap_alloc, mmap_free, mmap_realloc, mmap_destroy)
}

func mmap_alloc(allocator unsafe.Pointer, size int) unsafe.Pointer {
	var viableChunk *chunk
	currentChunk := chunkHead

	for currentChunk != nil {
		if currentChunk.offset+uintptr(size)+sizeOfPtrMeta <= currentChunk.size {
			viableChunk = currentChunk
			break
		}
		currentChunk = currentChunk.next
	}

	if viableChunk == nil {
		totalSize := int(mm.Align(uintptr(size)+sizeOfChunk+sizeOfPtrMeta, uintptr(pageSize)))
		res, err := mmap.MapRegion(nil, totalSize, syscall.PROT_READ|syscall.PROT_WRITE, mmap.ANON|mmap.RDWR, 0)
		if err != nil {
			panic(err)
		}

		viableChunk = (*chunk)(unsafe.Pointer(&res[0]))
		viableChunk.data = unsafe.Pointer(uintptr(unsafe.Pointer(&res[0])) + sizeOfChunk)
		viableChunk.size = uintptr(totalSize - int(sizeOfChunk))
		viableChunk.offset = 0
		viableChunk.ptrs = 0

		if chunkHead != nil {
			viableChunk.next = chunkHead
			chunkHead.prev = viableChunk
			chunkHead = viableChunk
		} else {
			chunkHead = viableChunk
		}
	}

	ptrMeta := (*mmapMeta)(unsafe.Pointer(uintptr(viableChunk.data) + viableChunk.offset))
	ptrMeta.size = size
	ptrMeta.chunk = viableChunk
	viableChunk.offset = mm.Align(viableChunk.offset+uintptr(size)+sizeOfPtrMeta, alignment)
	viableChunk.ptrs++

	ptr := unsafe.Pointer(uintptr(unsafe.Pointer(ptrMeta)) + uintptr(sizeOfPtrMeta))
	return ptr
}

func mmap_free(allocator unsafe.Pointer, ptr unsafe.Pointer) {
	ptrMeta := (*mmapMeta)(unsafe.Pointer(uintptr(ptr) - sizeOfPtrMeta))
	chunk := ptrMeta.chunk
	chunk.ptrs--

	if chunk.ptrs == 0 {
		if chunk == chunkHead {
			chunkHead = chunkHead.next
		} else {
			if chunk.prev != nil {
				chunk.prev.next = chunk.next
			}

			if chunk.next != nil {
				chunk.next.prev = chunk.prev
			}
		}

		m := mmap.MMap(
			unsafe.Slice(
				(*byte)(
					unsafe.Pointer(uintptr(chunk.data)-sizeOfChunk),
				), chunk.size+sizeOfChunk),
		)
		err := m.Unmap()
		if err != nil {
			panic(err)
		}
	}
}

func mmap_realloc(allocator unsafe.Pointer, ptr unsafe.Pointer, size int) unsafe.Pointer {
	oldMeta := (*mmapMeta)(unsafe.Pointer(uintptr(ptr) - sizeOfPtrMeta))
	oldData := unsafe.Slice((*byte)(ptr), oldMeta.size)
	newPtr := mmap_alloc(allocator, size)
	newData := unsafe.Slice((*byte)(newPtr), size)
	copy(newData, oldData)
	mmap_free(allocator, ptr)

	return newPtr
}

func mmap_destroy(allocator unsafe.Pointer) {
	currentChunk := chunkHead
	for currentChunk != nil {
		nextChunk := currentChunk.next
		m := mmap.MMap(
			unsafe.Slice(
				(*byte)(
					unsafe.Pointer(uintptr(currentChunk.data)-sizeOfChunk),
				), currentChunk.size+sizeOfChunk),
		)
		err := m.Unmap()
		if err != nil {
			panic(err)
		}
		currentChunk = nextChunk
	}

	chunkHead = nil
}
