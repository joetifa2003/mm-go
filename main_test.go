package main

import (
	"testing"
)

const NNodes = 1000000

type Node struct {
	Value int
	Prev  *Node
	Next  *Node
}

func unManaged() {
	allocatedNodes := AllocMany[Node](NNodes)
	defer FreeMany(allocatedNodes)

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
}

func managed() {
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
}

func BenchmarkHeapManaged(b *testing.B) {
	for i := b.N; i <= b.N; i++ {
		managed()
	}
}

func BenchmarkManual(b *testing.B) {
	for i := b.N; i <= b.N; i++ {
		unManaged()
	}
}
