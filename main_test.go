package main

import (
	"runtime"
	"testing"
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

// Example of cgo overhead makes it actually slower than go
// don't alloc and free in a hot loop and try to allocate in chunks
// to avoid cgo overhead (look at BenchmarkArenaManual)
func BenchmarkManual(b *testing.B) {
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

			node := Alloc[Node]()
			*node = Node{
				Value: j,
				Prev:  prev,
				Next:  next,
			}
			nodes[j] = node
		}

		for _, n := range nodes {
			Free(n)
		}

		runtime.GC()
	}
}
