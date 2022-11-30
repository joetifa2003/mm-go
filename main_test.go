package main

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

func BenchmarkHeapManaged(b *testing.B) {
	for i := b.N; i <= b.N; i++ {
		nodes := make([]*Node, NNodes)

		for j := 0; j < NNodes; j++ {
			var prev *Node
			var next *Node
			if j != 0 {
				prev = nodes[j-1]
			}
			if j != NNodes-1 {
				next = nodes[j+1]
			}

			nodes[j] = &Node{
				Value: j,
				Prev:  prev,
				Next:  next,
			}
		}

		runtime.GC()
	}
}

func BenchmarkArenaManual(b *testing.B) {
	for i := b.N; i <= b.N; i++ {
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

		FreeMany(allocatedNodes)
		runtime.GC()
	}
}

func TestAllocMany(t *testing.T) {
	assert := assert.New(t)

	allocated := AllocMany[int](2) // allocates 1 int and returns it as a slice of ints with length 2
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
