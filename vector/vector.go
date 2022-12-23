package vector

import (
	"fmt"

	"github.com/joetifa2003/mm-go"
)

// Vector a contiguous growable array type
type Vector[T any] struct {
	data []T
	cap  int
	len  int
}

func createVector[T any](len int, cap int) *Vector[T] {
	vector := mm.Alloc[Vector[T]]()
	vector.cap = cap
	vector.len = len
	vector.data = mm.AllocMany[T](vector.cap)

	return vector
}

// New creates a new empty vector, if args not provided
// it will create an empty vector, if only one arg is provided
// it will init a vector with len and cap equal to the provided arg,
// if two args are provided it will init a vector with len = args[0] cap = args[1]
func New[T any](args ...int) *Vector[T] {
	switch len(args) {
	case 0:
		return createVector[T](0, 1)
	case 1:
		return createVector[T](args[0], args[0])
	default:
		return createVector[T](args[0], args[1])
	}
}

// Init initializes a new vector with the T elements provided and sets
// it's len and cap to len(values)
func Init[T any](values ...T) *Vector[T] {
	vector := createVector[T](len(values), len(values))
	copy(vector.data, values)
	return vector
}

// Push pushes value T to the vector, grows if needed.
func (v *Vector[T]) Push(value T) {
	if v.len == v.cap {
		v.data = mm.Reallocate(v.data, v.cap*2)
		v.cap *= 2
	}

	v.data[v.len] = value
	v.len++
}

// Pop pops value T from the vector and returns it
func (v *Vector[T]) Pop() T {
	v.len--
	return v.data[v.len]
}

// Len gets vector length
func (v *Vector[T]) Len() int {
	return v.len
}

// Cap gets vector capacity (underling memory length).
func (v *Vector[T]) Cap() int {
	return v.cap
}

// Slice gets a slice representing the vector
// CAUTION: don't append to this slice, this is only used
// if you want to loop on the vec elements
func (v *Vector[T]) Slice() []T {
	return v.data[:v.len]
}

// Last gets the last element from a vector
func (v *Vector[T]) Last() T {
	return v.data[v.len-1]
}

// At gets element T at specified index
func (v *Vector[T]) At(idx int) T {
	if idx >= v.len {
		panic(fmt.Sprintf("cannot index %d in a vector with length %d", idx, v.len))
	}

	return v.data[idx]
}

// AtPtr gets element a pointer of T at specified index
func (v *Vector[T]) AtPtr(idx int) *T {
	if idx >= v.len {
		panic(fmt.Sprintf("cannot index %d in a vector with length %d", idx, v.len))
	}

	return &v.data[idx]
}

// Set sets element T at specified index
func (v *Vector[T]) Set(idx int, value T) {
	if idx >= v.len {
		panic(fmt.Sprintf("cannot set %d in a vector with length %d", idx, v.len))
	}

	v.data[idx] = value
}

// Free deallocats the vector
func (v *Vector[T]) Free() {
	mm.FreeMany(v.data)
	mm.Free(v)
}
