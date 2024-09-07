package minheap

import (
	"iter"

	"github.com/joetifa2003/mm-go/allocator"
	"github.com/joetifa2003/mm-go/vector"
)

type MinHeap[T any] struct {
	alloc allocator.Allocator
	data  *vector.Vector[T]
	less  func(a, b T) bool // check if a < b
}

// New creates a new MinHeap.
func New[T any](alloc allocator.Allocator, less func(a, b T) bool) *MinHeap[T] {
	minHeap := allocator.Alloc[MinHeap[T]](alloc)
	minHeap.alloc = alloc
	minHeap.data = vector.New[T](alloc) // Start with an initial capacity of 16
	minHeap.less = less
	return minHeap
}

// Push adds a value to the heap.
func (h *MinHeap[T]) Push(value T) {
	h.data.Push(value)
	h.heapifyUp(h.data.Len() - 1)
}

// Pop removes and returns the minimum value from the heap.
func (h *MinHeap[T]) Pop() T {
	if h.data.Len() == 0 {
		panic("cannot pop from empty heap")
	}

	minValue := h.data.At(0)
	h.data.RemoveAt(0)
	h.heapifyDown(0)

	return minValue
}

// Peek returns the minimum value from the heap without removing it.
func (h *MinHeap[T]) Peek() T {
	if h.data.Len() == 0 {
		panic("cannot peek into empty heap")
	}
	return h.data.At(0)
}

// Len returns the number of elements in the heap.
func (h *MinHeap[T]) Len() int {
	return h.data.Len()
}

// Free frees the heap.
func (h *MinHeap[T]) Free() {
	h.data.Free()
	allocator.Free(h.alloc, h)
}

func (h *MinHeap[T]) Remove(f func(T) bool) {
	for i := 0; i < h.data.Len(); i++ {
		if f(h.data.At(i)) {
			h.removeAt(i)
			return
		}
	}
}

func (h *MinHeap[T]) heapifyUp(index int) {
	for index > 0 {
		parentIndex := (index - 1) / 2
		if h.less(h.data.At(parentIndex), h.data.At(index)) {
			break
		}
		h.swap(parentIndex, index)
		index = parentIndex
	}
}

func (h *MinHeap[T]) heapifyDown(index int) {
	for {
		leftChildIndex := 2*index + 1
		rightChildIndex := 2*index + 2
		smallestIndex := index

		if leftChildIndex < h.data.Len() && h.less(h.data.At(leftChildIndex), h.data.At(smallestIndex)) {
			smallestIndex = leftChildIndex
		}

		if rightChildIndex < h.data.Len() && h.less(h.data.At(rightChildIndex), h.data.At(smallestIndex)) {
			smallestIndex = rightChildIndex
		}

		if smallestIndex == index {
			break
		}

		h.swap(index, smallestIndex)
		index = smallestIndex
	}
}

func (h *MinHeap[T]) swap(i, j int) {
	temp := h.data.At(i)
	h.data.Set(i, h.data.At(j))
	h.data.Set(j, temp)
}

// removeAt removes the element at the specified index from the heap.
func (h *MinHeap[T]) removeAt(index int) {
	if index == h.data.Len()-1 {
		h.data.Pop()
	} else {
		h.swap(index, h.data.Len()-1)
		h.data.Pop()
		h.heapifyDown(index)
		h.heapifyUp(index)
	}
}

func (h *MinHeap[T]) Iter() iter.Seq[T] {
	return h.data.Iter()
}
