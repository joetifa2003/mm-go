package allocator

import (
	"os"
	"unsafe"
)

var pageSize = os.Getpagesize() * 1

var chunkSize = int(unsafe.Sizeof(Chunk{}))

type Chunk struct {
	data      unsafe.Pointer
	allocated int
	size      int
	ptrCount  int
	next      *Chunk
	prev      *Chunk
}

type ArenaAllocator struct {
	tail      *Chunk
	allocator Allocator
}

var arenaAllocatorSize = int(unsafe.Sizeof(ArenaAllocator{}))

func NewArena(allocator Allocator) Allocator {
	a := (*ArenaAllocator)(allocator.Alloc(arenaAllocatorSize))
	a.allocator = allocator

	return a
}

func (a *ArenaAllocator) newChunk(n int) *Chunk {
	ch := (*Chunk)(a.allocator.Alloc(chunkSize))
	nPages := (n / pageSize) + 1
	ch.size = nPages * pageSize
	ch.data = a.allocator.Alloc(nPages * pageSize)

	return ch
}

type Meta struct {
	chunk *Chunk
	size  int
}

var metaSize = int(unsafe.Sizeof(Meta{}))

func (a *ArenaAllocator) Alloc(n int) unsafe.Pointer {
	if a.tail == nil {
		a.tail = a.newChunk(n)
	}

	if a.tail.allocated+n+metaSize > a.tail.size {
		ch := a.newChunk(n + metaSize)
		ch.prev = a.tail
		a.tail.next = ch
		a.tail = ch
	}
	meta := (*Meta)(unsafe.Pointer(uintptr(a.tail.data) + uintptr(a.tail.allocated)))
	meta.chunk = a.tail
	meta.size = n
	ptr := unsafe.Pointer(uintptr(unsafe.Pointer(meta)) + uintptr(metaSize))
	a.tail.allocated += n + metaSize
	a.tail.ptrCount++

	return ptr
}

func (a *ArenaAllocator) freeChunk(ch *Chunk) {
	if ch.prev != nil {
		ch.prev.next = ch.next
	}
	if ch.next != nil {
		ch.next.prev = ch.prev
	}

	if ch == a.tail {
		a.tail = ch.prev
	}

	a.allocator.Free(ch.data)
	a.allocator.Free(unsafe.Pointer(ch))
}

func (a *ArenaAllocator) Free(ptr unsafe.Pointer) {
	meta := (*Meta)(unsafe.Pointer(uintptr(ptr) - uintptr(metaSize)))
	meta.chunk.ptrCount--
}

func (a *ArenaAllocator) Destroy() {
	for ch := a.tail; ch != nil; ch = ch.next {
		a.freeChunk(ch)
	}
}

func (a *ArenaAllocator) Realloc(ptr unsafe.Pointer, size int) unsafe.Pointer {
	newPtr := a.Alloc(size)

	oldMeta := (*Meta)(unsafe.Pointer(uintptr(ptr) - uintptr(metaSize)))
	newPtrData := unsafe.Slice((*byte)(newPtr), size)
	oldPtrData := unsafe.Slice((*byte)(ptr), oldMeta.size)
	copy(newPtrData, oldPtrData)

	return newPtr
}
