package batchallocator_test

import (
	"github.com/joetifa2003/mm-go"
	"github.com/joetifa2003/mm-go/allocator"
	"github.com/joetifa2003/mm-go/batchallocator"
)

func Example() {
	alloc := batchallocator.New(allocator.NewC()) // by default it allocates page, which is usually 4kb
	defer alloc.Destroy()                         // this frees all  memory allocated by the allocator automatically

	ptr := allocator.Alloc[int](alloc)
	// but you can still free the pointers manually if you want (will free buckets of memory if all pointers depending on it is freed)
	defer allocator.Free(alloc, ptr) // this can removed and the memory will be freed.
}

func ExampleWithBucketSize() {
	alloc := batchallocator.New(
		allocator.NewC(),
		batchallocator.WithBucketSize(mm.SizeOf[int]()*15), // configure the allocator to allocate size of 15 ints in one go.
	)
	defer alloc.Destroy()

	ptr := allocator.Alloc[int](alloc)
	defer allocator.Free(alloc, ptr) // this can be removed and the memory will still be freed on Destroy.

	ptr2 := allocator.Alloc[int](alloc) // will not call CGO because there is still enough memory in the Bucket.
	defer allocator.Free(alloc, ptr2)   // this can be removed and the memory will still be freed on Destroy.

}
