package batchallocator

import (
	"os"
	"unsafe"

	"github.com/joetifa2003/mm-go/allocator"
	"github.com/joetifa2003/mm-go/minheap"
)

var pageSize = os.Getpagesize()

const alignment = unsafe.Alignof(uintptr(0))

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

func align(ptr uintptr, align uintptr) uintptr {
	mask := align - 1
	return (ptr + mask) &^ mask
}

func batchAllocaator_alloc(allocator unsafe.Pointer, size int) unsafe.Pointer {
	b := (*BatchAllocator)(allocator)
	if b.buckets == nil {
		b.buckets = minheap.New(b.alloc, batch_less)
	}

	for b.buckets.Len() == 0 {
		b.buckets.Push(newBucket(b.alloc, size))
	}

	// Align size to be a multiple of the alignment
	alignedSize := uintptr(size)
	if alignedSize%alignment != 0 {
		alignedSize = (alignedSize + alignment - 1) &^ (alignment - 1)
	}

	bucket := b.buckets.Peek()
	usedEnd := bucket.used + int(sizeOfPtrMeta) + int(alignedSize)

	if usedEnd <= bucket.size {
		bucket = b.buckets.Pop()

		// Align the `bucket.data + bucket.used`
		addr := uintptr(bucket.data) + uintptr(bucket.used)
		alignedAddr := align(addr, alignment)

		// Create the metadata structure for ptr
		ptr := (*ptrMeta)(unsafe.Pointer(alignedAddr))
		bucket.used = int(alignedAddr-uintptr(bucket.data)) + int(sizeOfPtrMeta) + int(alignedSize)
		bucket.ptrs++
		ptr.bucket = bucket
		ptr.size = int(alignedSize)
		b.buckets.Push(bucket)

		return unsafe.Pointer(alignedAddr + uintptr(sizeOfPtrMeta))
	}

	// Bucket is too small, allocate a new bucket
	newBucket := newBucket(b.alloc, int(alignedSize))
	newBucket.used = int(sizeOfPtrMeta) + int(alignedSize)
	newBucket.ptrs++
	b.buckets.Push(newBucket)

	ptr := (*ptrMeta)(unsafe.Pointer(newBucket.data))
	ptr.bucket = newBucket
	ptr.size = int(alignedSize)

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
