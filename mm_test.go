package mm_test

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/joetifa2003/mm-go/allocator"
	"github.com/joetifa2003/mm-go/batchallocator"
	"github.com/joetifa2003/mm-go/typedarena"
)

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

func createTreeManual(depth int, arena *typedarena.TypedArena[TreeNode]) *TreeNode {
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

const TREE_DEPTH = 13

func BenchmarkBinaryTreeManaged(b *testing.B) {
	for n := 0; n < b.N; n++ {
		tree := createTreeManaged(TREE_DEPTH)
		runtime.GC()
		sumBinaryTree(tree)
	}
}

func BenchmarkBinaryTreeArena(b *testing.B) {
	for _, chunkSize := range []int{50, 100, 150, 250, 500} {
		b.Run(fmt.Sprintf("callocator chunk size %d", chunkSize), func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				alloc := allocator.NewCallocator()
				arena := typedarena.New[TreeNode](alloc, chunkSize)
				tree := createTreeManual(TREE_DEPTH, arena)
				runtime.GC()
				sumBinaryTree(tree)
				arena.Free()
				alloc.Destroy()
			}
		})

		b.Run(fmt.Sprintf("batchallocator chunk size %d", chunkSize), func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				alloc := batchallocator.New(allocator.NewCallocator())
				arena := typedarena.New[TreeNode](alloc, chunkSize)
				tree := createTreeManual(TREE_DEPTH, arena)
				runtime.GC()
				sumBinaryTree(tree)
				arena.Free()
				alloc.Destroy()
			}
		})
	}
}
