package vector

import (
	"fmt"
	"iter"

	"github.com/joetifa2003/mm-go/allocator"
)

// Vector a contiguous growable array type
type Vector[T any] struct {
	data  []T
	len   int
	alloc allocator.Allocator
}

func createVector[T any](alloc allocator.Allocator, len int, cap int) *Vector[T] {
	vector := allocator.Alloc[Vector[T]](alloc)
	vector.len = len
	vector.data = allocator.AllocMany[T](alloc, cap)
	vector.alloc = alloc

	return vector
}

// New creates a new empty vector, if args not provided
// it will create an empty vector, if only one arg is provided
// it will init a vector with len and cap equal to the provided arg,
// if two args are provided it will init a vector with len = args[0] cap = args[1]
func New[T any](aloc allocator.Allocator, args ...int) *Vector[T] {
	switch len(args) {
	case 0:
		return createVector[T](aloc, 0, 1)
	case 1:
		return createVector[T](aloc, args[0], args[0])
	default:
		return createVector[T](aloc, args[0], args[1])
	}
}

// Init initializes a new vector with the T elements provided and sets
// it's len and cap to len(values)
func Init[T any](alloc allocator.Allocator, values ...T) *Vector[T] {
	vector := createVector[T](alloc, len(values), len(values))
	copy(vector.data, values)
	return vector
}

// Push pushes value T to the vector, grows if needed.
func (v *Vector[T]) Push(value T) {
	if v.len == v.Cap() {
		v.data = allocator.Realloc(v.alloc, v.data, v.Cap()*2)
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
	return cap(v.data)
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

// UnsafeAT gets element T at specified index without bounds checking
func (v *Vector[T]) UnsafeAt(idx int) T {
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
	allocator.FreeMany[T](v.alloc, v.data)
	allocator.Free(v.alloc, v)
}

func (v *Vector[T]) RemoveAt(idx int) T {
	if idx >= v.len {
		panic(fmt.Sprintf("cannot remove %d in a vector with length %d", idx, v.len))
	}

	tmp := v.data[idx]
	v.data[idx] = v.data[v.len-1]
	v.len--

	return tmp
}

func (v *Vector[T]) Iter() iter.Seq[T] {
	return func(yield func(T) bool) {
		for i := 0; i < v.len; i++ {
			if !yield(v.data[i]) {
				return
			}
		}
	}
}
