package typedarena

import (
	"github.com/joetifa2003/mm-go/allocator"
	"github.com/joetifa2003/mm-go/vector"
)

type typedChunk[T any] struct {
	data  []T
	len   int
	alloc allocator.Allocator
}

func newChunk[T any](alloc allocator.Allocator, size int) *typedChunk[T] {
	chunk := allocator.Alloc[typedChunk[T]](alloc)
	chunk.data = allocator.AllocMany[T](alloc, size)
	chunk.alloc = alloc

	return chunk
}

func (c *typedChunk[T]) Alloc() *T {
	c.len++
	return &c.data[c.len-1]
}

func (c *typedChunk[T]) AllocMany(n int) []T {
	oldLen := c.len
	c.len += n
	return c.data[oldLen:c.len]
}

func (c *typedChunk[T]) Free() {
	allocator.FreeMany(c.alloc, c.data)
	allocator.Free(c.alloc, c)
}

// TypedArena is a growable typed arena
type TypedArena[T any] struct {
	chunks    *vector.Vector[*typedChunk[T]]
	chunkSize int
	alloc     allocator.Allocator
}

// New creates a typed arena with the specified chunk size.
// a chunk is the the unit of the arena, if T is int for example and the
// chunk size is 5, then each chunk is going to hold 5 ints. And if the
// chunk is filled it will allocate another chunk that can hold 5 ints.
// then you can call FreeArena and it will deallocate all chunks together
func New[T any](alloc allocator.Allocator, chunkSize int) *TypedArena[T] {
	tArena := allocator.Alloc[TypedArena[T]](alloc)
	tArena.chunkSize = chunkSize
	tArena.chunks = vector.New[*typedChunk[T]](alloc)

	firstChunk := newChunk[T](alloc, chunkSize)
	tArena.chunks.Push(firstChunk)

	return tArena
}

// Alloc allocates T from the arena
func (ta *TypedArena[T]) Alloc() *T {
	lastChunk := ta.chunks.Last()
	if lastChunk.len == ta.chunkSize {
		nc := newChunk[T](ta.alloc, ta.chunkSize)
		ta.chunks.Push(nc)
		return nc.Alloc()
	}
	return lastChunk.Alloc()
}

// AllocMany allocates n of T and returns a slice representing the heap.
// CAUTION: don't append to the slice, the purpose of it is to replace pointer
// arithmetic with slice indexing
// CAUTION: n cannot exceed chunk size
func (ta *TypedArena[T]) AllocMany(n int) []T {
	if n > ta.chunkSize {
		panic("cannot exceed chunk size")
	}

	lastChunk := ta.chunks.Last()
	if lastChunk.len+n > ta.chunkSize {
		nc := newChunk[T](ta.alloc, ta.chunkSize)
		ta.chunks.Push(nc)
		return nc.AllocMany(n)
	}

	return lastChunk.AllocMany(n)
}

// Free frees all allocated memory
func (ta *TypedArena[T]) Free() {
	for _, c := range ta.chunks.Slice() {
		c.Free()
	}
	ta.chunks.Free()
	allocator.Free(ta.alloc, ta)
}
