package mm

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

const NNodes = 1000000

type Node struct {
	Value int
	Prev  *Node
	Next  *Node
}

// DoNotOptimiseNodePointer prevent the compiler removing the function bodies
var DoNotOptimiseNodePointer *Node

func BenchmarkHeapManaged(b *testing.B) {
	// Start with a clean slate
	DoNotOptimiseNodePointer = nil
	runtime.GC()
	b.ResetTimer()

	for i := b.N; i != 0; i-- {
		allocatedNodes := make([]Node, NNodes)

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

		DoNotOptimiseNodePointer = &allocatedNodes[len(allocatedNodes)-1]
	}

	runtime.GC()

	DoNotOptimiseNodePointer = nil
}

func BenchmarkArenaManual(b *testing.B) {
	// Start with a clean slate
	DoNotOptimiseNodePointer = nil
	runtime.GC()
	b.ResetTimer()

	for i := b.N; i != 0; i-- {
		allocatedNodes := AllocMany[Node](NNodes)

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

		DoNotOptimiseNodePointer = &allocatedNodes[len(allocatedNodes)-1]
		FreeMany(allocatedNodes)
	}

	DoNotOptimiseNodePointer = nil
}

const LOOP_TIMES = 1500

func TestAllocMany(t *testing.T) {
	assert := assert.New(t)

	allocated := AllocMany[int](2) // allocates 2 ints and returns it as a slice of ints with length 2
	defer FreeMany(allocated)      // it's recommended to make sure the data gets deallocated (defer recommended to prevent leaks)
	assert.Equal(2, len(allocated))
	allocated[0] = 15    // changes the data in the slice (aka the heap)
	ptr := &allocated[0] // takes a pointer to the data in the heap
	*ptr = 45            // changes the value from 15 to 45

	assert.Equal(45, allocated[0])
}

func TestAlloc(t *testing.T) {
	assert := assert.New(t)

	ptr := Alloc[int]() // allocates a single int and returns a ptr to it
	defer Free(ptr)     // frees the int (defer recommended to prevent leaks)

	assert.Equal(0, *ptr) // allocations are zeroed by default
	*ptr = 15
	assert.Equal(15, *ptr)
}

func TestReallocate(t *testing.T) {
	assert := assert.New(t)

	allocated := AllocMany[int](2) // allocates 2 int and returns it as a slice of ints with length 2
	allocated[0] = 15
	assert.Equal(2, len(allocated))
	allocated = Reallocate(allocated, 3)
	assert.Equal(3, len(allocated))
	assert.Equal(15, allocated[0]) // data after reallocation stays the same
	FreeMany(allocated)            // didn't use defer here because i'm doing a reallocation and changing the value of allocated variable (otherwise can segfault)
}
