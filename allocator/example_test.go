package allocator_test

import (
	"fmt"
	"unsafe"

	"github.com/joetifa2003/mm-go"
	"github.com/joetifa2003/mm-go/allocator"
)

func ExampleNewC() {
	alloc := allocator.NewC()
	defer alloc.Destroy()

	ptr := allocator.Alloc[int](alloc)
	defer allocator.Free(alloc, ptr)

	*ptr = 15
	fmt.Println(*ptr)

	// Output: 15
}

func ExampleAlloc() {
	alloc := allocator.NewC()
	defer alloc.Destroy()

	// So you can do this:
	ptr := allocator.Alloc[int](alloc) // allocates a single int and returns a ptr to it
	defer allocator.Free(alloc, ptr)   // frees the int (defer recommended to prevent leaks)
	*ptr = 15
	fmt.Println(*ptr)

	// instead of doing this:
	ptr2 := (*int)(alloc.Alloc(mm.SizeOf[int]()))
	defer alloc.Free(unsafe.Pointer(ptr2))
	*ptr2 = 15

	fmt.Println(*ptr2)

	// Output:
	// 15
	// 15
}

func ExampleAllocMany() {
	alloc := allocator.NewC()
	defer alloc.Destroy()

	heap := allocator.AllocMany[int](alloc, 2) // allocates 2 ints and returns it as a slice of ints with length 2
	defer allocator.FreeMany(alloc, heap)      // it's recommended to make sure the data gets deallocated (defer recommended to prevent leaks)

	heap[0] = 15    // changes the data in the slice (aka the heap)
	ptr := &heap[0] // takes a pointer to the first int in the heap
	// Be careful if you do ptr := heap[0] this will take a copy from the data on the heap
	*ptr = 45 // changes the value from 15 to 45
	heap[1] = 70

	fmt.Println(heap[0])
	fmt.Println(heap[1])

	// Output:
	// 45
	// 70
}

func ExampleRealloc() {
	alloc := allocator.NewC()
	defer alloc.Destroy()

	heap := allocator.AllocMany[int](alloc, 2) // allocates 2 int and returns it as a slice of ints with length 2

	heap[0] = 15
	heap[1] = 70

	heap = allocator.Realloc(alloc, heap, 3)
	heap[2] = 100

	fmt.Println(heap[0])
	fmt.Println(heap[1])
	fmt.Println(heap[2])

	allocator.FreeMany(alloc, heap)

	// Output:
	// 15
	// 70
	// 100
}

func ExampleNewAllocator() {
	// Create a custom allocator
	alloc := allocator.NewAllocator(
		nil,
		myallocator_alloc,
		myallocator_free,
		myallocator_realloc,
		myallocator_destroy,
	)

	// Check how C allocator is implemented
	// or batchallocator soruce for a reference

	_ = alloc
}

func myallocator_alloc(allocator unsafe.Pointer, size int) unsafe.Pointer {
	return nil
}

func myallocator_free(allocator unsafe.Pointer, ptr unsafe.Pointer) {
}

func myallocator_realloc(allocator unsafe.Pointer, ptr unsafe.Pointer, size int) unsafe.Pointer {
	return nil
}

func myallocator_destroy(allocator unsafe.Pointer) {
}
