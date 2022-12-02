package mm

type chunk[T any] struct {
	data []T
	len  int
}

func newChunk[T any](size int) *chunk[T] {
	chunk := Alloc[chunk[T]]()
	chunk.data = AllocMany[T](size)

	return chunk
}

func (c *chunk[T]) Alloc() *T {
	c.len++
	return &c.data[c.len-1]
}

func (c *chunk[T]) Free() {
	FreeMany(c.data)
	Free(c)
}

// TypedArena is a growable typed arena
type TypedArena[T any] struct {
	chunks    *Vector[*chunk[T]]
	chunkSize int
}

// NewTypedArena creates a typed arena with the specified chunk size.
// a chunk is the the unit of the arena, if T is int for example and the
// chunk size is 5, then each chunk is going to hold 5 ints. And if the
// chunk is filled it will allocate another chunk that can hold 5 ints.
// then you can call FreeArena and it will deallocate all chunks together
func NewTypedArena[T any](chunkSize int) *TypedArena[T] {
	tArena := Alloc[TypedArena[T]]()
	tArena.chunkSize = chunkSize
	tArena.chunks = NewVector[*chunk[T]]()

	firstChunk := newChunk[T](chunkSize)
	tArena.chunks.Push(firstChunk)

	return tArena
}

// Alloc allocates T from the arena
func (ta *TypedArena[T]) Alloc() *T {
	lastChunk := ta.chunks.Last()
	if lastChunk.len == len(lastChunk.data) {
		nc := newChunk[T](ta.chunkSize)
		ta.chunks.Push(nc)
		return nc.Alloc()
	}
	return lastChunk.Alloc()
}

// Free frees all allocated memory
func (ta *TypedArena[T]) Free() {
	for _, c := range ta.chunks.Slice() {
		c.Free()
	}
	ta.chunks.Free()
	Free(ta)
}
