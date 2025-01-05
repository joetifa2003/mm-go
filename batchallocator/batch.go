// This allocator purpose is to reduce the overhead of calling CGO on every allocation/free, it also acts as an arena since it frees all the memory when `Destroy` is called.
// It allocats large chunks of memory at once and then divides them when you allocate, making it much faster.
// This allocator has to take another allocator for it to work, usually with the C allocator.
// You can optionally call `Free` on the pointers allocated by batchallocator manually, and it will free the memory as soon as it can.
// `Destroy` must be called to free internal resources and free all the memory allocated by the allocator.
package batchallocator

import (
	"os"
	"unsafe"

	"github.com/joetifa2003/mm-go/allocator"
	"github.com/joetifa2003/mm-go/minheap"
)

var pageSize = os.Getpagesize()

const (
	alignment     = unsafe.Alignof(uintptr(0))
	sizeOfPtrMeta = unsafe.Sizeof(ptrMeta{}) // Cache the size of ptrMeta type
)

type bucket struct {
	data   unsafe.Pointer // Base pointer of the bucket memory
	offset uintptr        // Number of used bytes in this bucket
	size   uintptr        // Total size of this bucket
	ptrs   int            // Number of pointers (allocations) inside the bucket
}

func (b *bucket) Free(a allocator.Allocator) {
	a.Free(b.data)
	a.Free(unsafe.Pointer(b))
}

// A metadata structure stored before each allocated memory block
type ptrMeta struct {
	bucket *bucket // Pointer back to the bucket
	size   int     // Size of the allocated memory block
}

// BatchAllocator manages a collection of memory buckets to optimize small allocations
type BatchAllocator struct {
	buckets    *minheap.MinHeap[*bucket] // Min-heap for managing buckets based on used space
	alloc      allocator.Allocator       // Underlying raw allocator (backed by malloc/free)
	bucketSize int                       // Configurable size for each new bucket
}

type BatchAllocatorOption func(alloc *BatchAllocator)

// WithBucketSize Option to specify bucket size when creating BatchAllocator
// You can allocate more memory than the bucketsize in one allocation, it will allocate a new bucket and put the data in it.
func WithBucketSize(size int) BatchAllocatorOption {
	return func(alloc *BatchAllocator) {
		alloc.bucketSize = size
	}
}

// New creates a new BatchAllocator and applies optional configuration using BatchAllocatorOption
func New(a allocator.Allocator, options ...BatchAllocatorOption) allocator.Allocator {
	balloc := allocator.Alloc[BatchAllocator](a)
	balloc.alloc = a

	// Apply configuration options to BatchAllocator
	for _, option := range options {
		option(balloc)
	}

	return allocator.NewAllocator(
		unsafe.Pointer(balloc),
		batchAllocatorAlloc,
		batchAllocatorFree,
		batchAllocatorRealloc,
		batchAllocatorDestroy,
	)
}

// Performs the allocation from the BatchAllocator
func batchAllocatorAlloc(allocator unsafe.Pointer, size int) unsafe.Pointer {
	balloc := (*BatchAllocator)(allocator)

	// Ensure we have a bucket to allocate from
	ensureBucketExists(balloc, size)

	// Check if the current top bucket can handle the allocation
	currentBucket := balloc.buckets.Peek()

	if currentBucket.offset+sizeOfPtrMeta+uintptr(size) <= currentBucket.size {
		currentBucket = balloc.buckets.Pop()

		freeStart := uintptr(currentBucket.data) + currentBucket.offset

		// Write the allocation metadata
		meta := (*ptrMeta)(unsafe.Pointer(freeStart))
		meta.bucket = currentBucket
		meta.size = size

		currentBucket.offset = align(currentBucket.offset+sizeOfPtrMeta+uintptr(size), alignment)
		currentBucket.ptrs++

		balloc.buckets.Push(meta.bucket)

		// Return the address of the memory after the metadata
		return unsafe.Pointer(freeStart + uintptr(sizeOfPtrMeta))
	}

	// If no bucket can accommodate the allocation, create a new one
	newBucket := allocateNewBucket(balloc, size)
	newBucket.offset = align(sizeOfPtrMeta+uintptr(size), alignment)
	newBucket.ptrs++
	balloc.buckets.Push(newBucket)

	// Write meta information at the base of the new bucket
	meta := (*ptrMeta)(unsafe.Pointer(newBucket.data))
	meta.bucket = newBucket
	meta.size = size

	return unsafe.Pointer(uintptr(unsafe.Pointer(meta)) + uintptr(sizeOfPtrMeta))
}

// Frees the allocated memory by decrementing reference count and freeing bucket if empty
func batchAllocatorFree(allocator unsafe.Pointer, ptr unsafe.Pointer) {
	balloc := (*BatchAllocator)(allocator)

	// Retrieve the metadata by moving back
	meta := (*ptrMeta)(unsafe.Pointer(uintptr(ptr) - sizeOfPtrMeta))
	meta.bucket.ptrs--

	// If no more pointers exist in the bucket, free the bucket
	if meta.bucket.ptrs == 0 {
		balloc.buckets.Remove(func(b *bucket) bool {
			return b == meta.bucket
		})
		meta.bucket.Free(balloc.alloc)
	}
}

// Reallocate a block of memory
func batchAllocatorRealloc(allocator unsafe.Pointer, ptr unsafe.Pointer, size int) unsafe.Pointer {
	newPtr := batchAllocatorAlloc(allocator, size)

	// Copy the data from the old location to the new one
	oldMeta := (*ptrMeta)(unsafe.Pointer(uintptr(ptr) - sizeOfPtrMeta))
	oldData := unsafe.Slice((*byte)(ptr), oldMeta.size)
	newData := unsafe.Slice((*byte)(newPtr), size)

	copy(newData, oldData)

	// Free the old memory
	batchAllocatorFree(allocator, ptr)

	return newPtr
}

// Destroys the batch allocator, freeing all buckets and underlying library resources
func batchAllocatorDestroy(a unsafe.Pointer) {
	balloc := (*BatchAllocator)(a)

	// Free all buckets in the heap
	for _, b := range balloc.buckets.Iter() {
		b.Free(balloc.alloc)
	}

	parentAllocator := balloc.alloc
	balloc.buckets.Free()
	allocator.Free(balloc.alloc, balloc)
	parentAllocator.Destroy()
}

// Helper function to handle memory alignment for a given pointer
func align(ptr uintptr, alignment uintptr) uintptr {
	mask := alignment - 1
	return (ptr + mask) &^ mask
}

// Allocates a new bucket with a given size, ensuring it's a multiple of the page size
func allocateNewBucket(balloc *BatchAllocator, size int) *bucket {
	size = max(balloc.bucketSize, size)

	nPages := size/pageSize + 1
	bucketSize := nPages * pageSize

	b := allocator.Alloc[bucket](balloc.alloc)
	b.data = balloc.alloc.Alloc(bucketSize)
	b.size = uintptr(bucketSize)
	b.offset = 0

	return b
}

// Pushes a new bucket into the bucket heap if needed
func ensureBucketExists(balloc *BatchAllocator, size int) {
	if balloc.buckets == nil {
		balloc.buckets = minheap.New(balloc.alloc, compareBucketFreeSpace)
	}

	if balloc.buckets.Len() == 0 {
		balloc.buckets.Push(allocateNewBucket(balloc, size))
	}
}

// Comparison function to prioritize buckets with more available space
func compareBucketFreeSpace(a, b *bucket) bool {
	return (a.size - a.offset) > (b.size - b.offset)
}
