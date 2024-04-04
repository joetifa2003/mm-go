package mm_test

import (
	"context"
	"fmt"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/joetifa2003/mm-go"
)

type Node struct {
	Value int
	Prev  *Node
	Next  *Node
}

func heapManaged(nodes int) {
	allocated := make([]*Node, nodes)

	for j := 0; j < nodes; j++ {
		var prev *Node
		var next *Node
		if j != 0 {
			prev = allocated[j-1]
		}
		if j != nodes-1 {
			next = allocated[j+1]
		}

		allocated[j] = &Node{
			Value: j,
			Prev:  prev,
			Next:  next,
		}
	}

	runtime.GC()
}

func BenchmarkHeapManaged(b *testing.B) {
	benchMarkSuit(b, heapManaged)
}

func manual(nodes int) {
	ctx := context.Background()
	allocatedNodes := mm.AllocMany[Node](ctx, nodes)

	for j := 0; j < nodes; j++ {
		var prev *Node
		var next *Node
		if j != 0 {
			prev = &allocatedNodes[j-1]
		}
		if j != nodes-1 {
			next = &allocatedNodes[j+1]
		}

		allocatedNodes[j] = Node{
			Value: j,
			Prev:  prev,
			Next:  next,
		}
	}

	mm.FreeMany(ctx, allocatedNodes)
	runtime.GC()
}

func BenchmarkManual(b *testing.B) {
	benchMarkSuit(b, manual)
}

func benchMarkSuit(b *testing.B, f func(int)) {
	nodeCounts := []int{10000, 100000, 10000000, 100000000}
	for _, nc := range nodeCounts {
		b.Run(fmt.Sprintf("node count %d", nc), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				f(nc)
			}
		})
	}
}

func TestAllocMany(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()

	allocated := mm.AllocMany[int](ctx, 2) // allocates 2 ints and returns it as a slice of ints with length 2
	defer mm.FreeMany(ctx, allocated)      // it's recommended to make sure the data gets deallocated (defer recommended to prevent leaks)
	assert.Equal(2, len(allocated))
	allocated[0] = 15    // changes the data in the slice (aka the heap)
	ptr := &allocated[0] // takes a pointer to the data in the heap
	*ptr = 45            // changes the value from 15 to 45

	assert.Equal(45, allocated[0])
}

func TestAlloc(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()

	ptr := mm.Alloc[int](ctx) // allocates a single int and returns a ptr to it
	defer mm.Free(ctx, ptr)   // frees the int (defer recommended to prevent leaks)

	*ptr = 15
	assert.Equal(15, *ptr)

	ptr2 := mm.Alloc[[1e3]int](ctx) // creates large array to make malloc mmap new chunk
	defer mm.Free(ctx, ptr2)
}

func TestReallocate(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()

	allocated := mm.AllocMany[int](ctx, 2) // allocates 2 int and returns it as a slice of ints with length 2
	allocated[0] = 15
	assert.Equal(2, len(allocated))
	allocated = mm.Reallocate(ctx, allocated, 3)
	assert.Equal(3, len(allocated))
	assert.Equal(15, allocated[0]) // data after reallocation stays the same
	mm.FreeMany(ctx, allocated)    // didn't use defer here because i'm doing a reallocation and changing the value of allocated variable (otherwise can segfault)
}

func TestUnmapChunk(t *testing.T) {
	ctx := context.Background()

	data1 := mm.AllocMany[int](ctx, 1e6)
	data2 := mm.AllocMany[int](ctx, 1e6)
	data3 := mm.AllocMany[int](ctx, 1e6)
	mm.FreeMany(ctx, data2)
	mm.FreeMany(ctx, data1)
	mm.FreeMany(ctx, data3)
}
