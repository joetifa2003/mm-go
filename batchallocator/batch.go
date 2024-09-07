package batchallocator

import (
	"os"
	"unsafe"

	"github.com/joetifa2003/mm-go/allocator"
	"github.com/joetifa2003/mm-go/minheap"
)

var pageSize = os.Getpagesize() * 15

type bucket struct {
	data unsafe.Pointer
	used int
	size int
	ptrs int
}

func newBucket(a allocator.Allocator, size int) *bucket {
	nPages := size/pageSize + 1
	size = nPages * pageSize

	b := allocator.Alloc[bucket](a)
	b.data = a.Alloc(size)
	b.used = 0
	b.size = size

	return b
}

func (b *bucket) Free(a allocator.Allocator) {
	a.Free(b.data)
	a.Free(unsafe.Pointer(b))
}

func batch_less(a, b *bucket) bool {
	return (a.size - a.used) > (b.size - b.used)
}

type ptrMeta struct {
	bucket *bucket
	size   int
}

var sizeOfPtrMeta = unsafe.Sizeof(ptrMeta{})

type BatchAllocator struct {
	buckets *minheap.MinHeap[*bucket]
	alloc   allocator.Allocator
}

func New(a allocator.Allocator) allocator.Allocator {
	ba := allocator.Alloc[BatchAllocator](a)
	ba.alloc = a
	return allocator.NewAllocator(unsafe.Pointer(ba), batchAllocaator_alloc, batchAllocaator_free, batchAllocaator_realloc, batchAllocaator_destroy)
}

func batchAllocaator_alloc(allocator unsafe.Pointer, size int) unsafe.Pointer {
	b := (*BatchAllocator)(allocator)
	if b.buckets == nil {
		b.buckets = minheap.New(b.alloc, batch_less)
	}

	for b.buckets.Len() == 0 {
		b.buckets.Push(newBucket(b.alloc, size))
	}

	if b.buckets.Peek().used+int(sizeOfPtrMeta)+size <= b.buckets.Peek().size {
		bucket := b.buckets.Pop()
		ptr := (*ptrMeta)(unsafe.Pointer(uintptr(bucket.data) + uintptr(bucket.used)))
		bucket.used += int(sizeOfPtrMeta) + size
		bucket.ptrs++
		ptr.bucket = bucket
		ptr.size = size
		b.buckets.Push(bucket)

		return unsafe.Pointer(uintptr(unsafe.Pointer(ptr)) + uintptr(sizeOfPtrMeta))
	}

	newBucket := newBucket(b.alloc, size)
	newBucket.used = int(sizeOfPtrMeta) + size
	newBucket.ptrs++
	b.buckets.Push(newBucket)

	ptr := (*ptrMeta)(unsafe.Pointer(uintptr(newBucket.data)))
	ptr.bucket = newBucket
	ptr.size = size

	return unsafe.Pointer(uintptr(unsafe.Pointer(ptr)) + uintptr(sizeOfPtrMeta))
}

func batchAllocaator_free(allocator unsafe.Pointer, ptr unsafe.Pointer) {
	ba := (*BatchAllocator)(allocator)

	meta := (*ptrMeta)(unsafe.Pointer(uintptr(ptr) - sizeOfPtrMeta))
	meta.bucket.ptrs--

	if meta.bucket.ptrs == 0 {
		ba.buckets.Remove(func(b *bucket) bool {
			return b == meta.bucket
		})
		meta.bucket.Free(ba.alloc)
	}
}

func batchAllocaator_realloc(allocator unsafe.Pointer, ptr unsafe.Pointer, size int) unsafe.Pointer {
	newPtr := batchAllocaator_alloc(allocator, size)

	oldPtrMeta := (*ptrMeta)(unsafe.Pointer(uintptr(ptr) - sizeOfPtrMeta))
	oldPtrData := unsafe.Slice((*byte)(ptr), oldPtrMeta.size)
	newPtrData := unsafe.Slice((*byte)(newPtr), size)

	copy(newPtrData, oldPtrData)

	batchAllocaator_free(allocator, ptr)

	return newPtr
}

func batchAllocaator_destroy(alloc unsafe.Pointer) {
	ba := (*BatchAllocator)(alloc)
	for b := range ba.buckets.Iter() {
		b.Free(ba.alloc)
	}
	ba.buckets.Free()
	allocator.Free(ba.alloc, ba)
}
