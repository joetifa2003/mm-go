package mm

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
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
	allocatedNodes := AllocMany[Node](nodes)

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

	FreeMany(allocatedNodes)
	runtime.GC()
}

func BenchmarkManual(b *testing.B) {
	benchMarkSuit(b, manual)
}

func arenaManual(nodes int) {
	arena := NewTypedArena[Node](nodes)
	allocatedNodes := arena.AllocMany(nodes)

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

	arena.Free()
	runtime.GC()
}

func BenchmarkArenaManual(b *testing.B) {
	benchMarkSuit(b, arenaManual)
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

type TreeNode struct {
	value       int
	left, right *TreeNode
}

func createTreeManaged(depth int) *TreeNode {
	if depth != 0 {
		return &TreeNode{
			value: depth,
			left:  createTreeManaged(depth - 1),
			right: createTreeManaged(depth - 1),
		}
	}

	return nil
}

func createTreeManual(depth int, arena *TypedArena[TreeNode]) *TreeNode {
	if depth != 0 {
		node := arena.Alloc()
		node.left = createTreeManual(depth-1, arena)
		node.right = createTreeManual(depth-1, arena)
		return node
	}

	return nil
}

func sumBinaryTree(tree *TreeNode) int {
	if tree.left == nil && tree.right == nil {
		return tree.value
	}

	return sumBinaryTree(tree.left) + sumBinaryTree(tree.right)
}

func BenchmarkBinaryTree(b *testing.B) {
	b.Run("managed", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			tree := createTreeManaged(25)
			sumBinaryTree(tree)
			runtime.GC()
		}
	})

	b.Run("arena manual", func(b *testing.B) {
		for _, chunkSize := range []int{50, 100, 150, 250} {
			b.Run(fmt.Sprintf("chunk size %d", chunkSize), func(b *testing.B) {
				for n := 0; n < b.N; n++ {
					arena := NewTypedArena[TreeNode](chunkSize)
					tree := createTreeManual(25, arena)
					sumBinaryTree(tree)
					arena.Free()
					runtime.GC()
				}
			})
		}
	})
}

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

	*ptr = 15
	assert.Equal(15, *ptr)

	ptr2 := Alloc[[1e3]int]() // creates large array to make malloc mmap new chunk
	defer Free(ptr2)
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

func TestUnmapChunk(t *testing.T) {
	data1 := AllocMany[int](1e6)
	data2 := AllocMany[int](1e6)
	data3 := AllocMany[int](1e6)
	FreeMany(data2)
	FreeMany(data1)
	FreeMany(data3)
}
