package batchallocator_test

import (
	"fmt"

	"github.com/joetifa2003/mm-go"
	"github.com/joetifa2003/mm-go/allocator"
	"github.com/joetifa2003/mm-go/batchallocator"
	"github.com/joetifa2003/mm-go/linkedlist"
	"github.com/joetifa2003/mm-go/mmstring"
	"github.com/joetifa2003/mm-go/vector"
)

func Example() {
	alloc := batchallocator.New(allocator.NewC()) // by default it allocates page, which is usually 4kb
	defer alloc.Destroy()                         // this frees all  memory allocated by the allocator automatically

	ptr := allocator.Alloc[int](alloc)
	// but you can still free the pointers manually if you want (will free buckets of memory if all pointers depending on it is freed)
	defer allocator.Free(alloc, ptr) // this can removed and the memory will be freed.
}

func Example_arena() {
	alloc := batchallocator.New(allocator.NewC())
	defer alloc.Destroy() // all the memory allocated bellow will be freed, no need to free it manually.

	v := vector.New[int](alloc)
	v.Push(15)
	v.Push(70)

	for _, i := range v.Iter() {
		fmt.Println(i)
	}

	l := linkedlist.New[*mmstring.MMString](alloc)
	l.PushBack(mmstring.From(alloc, "hello"))
	l.PushBack(mmstring.From(alloc, "world"))

	for _, i := range l.Iter() {
		fmt.Println(i.GetGoString())
	}

	// Output:
	// 15
	// 70
	// hello
	// world
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
