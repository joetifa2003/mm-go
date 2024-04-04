package linkedlist

import (
	"context"
	"fmt"

	"github.com/joetifa2003/mm-go"
)

var popEmptyMsg = "cannot pop empty linked list"

type linkedListNode[T any] struct {
	value T
	next  *linkedListNode[T]
	prev  *linkedListNode[T]
}

// LinkedList a doubly-linked list.
// Note: can be a lot slower than Vector but sometimes faster in specific use cases
type LinkedList[T any] struct {
	head   *linkedListNode[T]
	tail   *linkedListNode[T]
	length int
	ctx    context.Context
}

// New creates a new linked list.
func New[T any](ctx context.Context) *LinkedList[T] {
	linkedList := mm.Alloc[LinkedList[T]](ctx)
	linkedList.ctx = ctx

	return linkedList
}

func (ll *LinkedList[T]) init(value T) {
	ll.head = mm.Alloc[linkedListNode[T]](ll.ctx)
	ll.head.value = value
	ll.tail = ll.head
	ll.length++
}

func (ll *LinkedList[T]) popLast() T {
	value := ll.tail.value
	mm.Free(ll.ctx, ll.tail)
	ll.tail = nil
	ll.head = nil
	ll.length--
	return value
}

// PushBack pushes value T to the back of the linked list.
func (ll *LinkedList[T]) PushBack(value T) {
	// initialize the linked list
	if ll.head == nil && ll.tail == nil {
		ll.init(value)
		return
	}

	newNode := mm.Alloc[linkedListNode[T]](ll.ctx)
	newNode.value = value
	newNode.prev = ll.tail
	ll.tail.next = newNode
	ll.tail = newNode
	ll.length++
}

// PushFront pushes value T to the back of the linked list.
func (ll *LinkedList[T]) PushFront(value T) {
	// initialize the linked list
	if ll.head == nil && ll.tail == nil {
		ll.init(value)
		return
	}

	newNode := mm.Alloc[linkedListNode[T]](ll.ctx)
	newNode.value = value
	newNode.next = ll.head
	ll.head.prev = newNode
	ll.head = newNode
	ll.length++
}

// PopBack pops and returns value T from the back of the linked list.
func (ll *LinkedList[T]) PopBack() T {
	if ll.length == 0 {
		panic(popEmptyMsg)
	}

	if ll.length == 1 {
		return ll.popLast()
	}

	value := ll.tail.value
	newTail := ll.tail.prev
	newTail.next = nil
	mm.Free(ll.ctx, ll.tail)
	ll.tail = newTail
	ll.length--

	return value
}

// PopFront pops and returns value T from the front of the linked list.
func (ll *LinkedList[T]) PopFront() T {
	if ll.length == 0 {
		panic(popEmptyMsg)
	}

	if ll.length == 1 {
		return ll.popLast()
	}

	value := ll.head.value
	newHead := ll.head.next
	newHead.prev = nil
	mm.Free(ll.ctx, ll.head)
	ll.head = newHead
	ll.length--

	return value
}

// ForEach iterates through the linked list.
func (ll *LinkedList[T]) ForEach(f func(idx int, value T)) {
	idx := 0

	currentNode := ll.head
	for currentNode != nil {
		f(idx, currentNode.value)
		currentNode = currentNode.next
		idx++
	}
}

func (ll *LinkedList[T]) nodeAt(idx int) *linkedListNode[T] {
	if idx >= ll.length {
		panic(fmt.Sprintf("cannot index %d in a linked list with length %d", idx, ll.length))
	}

	i := 0
	currentNode := ll.head
	for i != idx {
		currentNode = currentNode.next
		i++
	}

	return currentNode
}

// At gets value T at idx.
func (ll *LinkedList[T]) At(idx int) T {
	return ll.nodeAt(idx).value
}

// AtPtr gets a pointer to value T at idx.
func (ll *LinkedList[T]) AtPtr(idx int) *T {
	return &ll.nodeAt(idx).value
}

// RemoveAt removes value T at specified index and returns it.
func (ll *LinkedList[T]) RemoveAt(idx int) T {
	node := ll.nodeAt(idx)
	if node.prev == nil {
		return ll.PopFront()
	}
	if node.next == nil {
		return ll.PopBack()
	}

	value := node.value

	nextNode := node.next
	prevNode := node.prev
	nextNode.prev = prevNode
	prevNode.next = nextNode
	ll.length--

	mm.Free(ll.ctx, node)

	return value
}

// Remove removes the first value T that pass the test implemented by the provided function.
// if the test succeeded it will return the value and true
func (ll *LinkedList[T]) Remove(f func(idx int, value T) bool) (value T, ok bool) {
	i := 0
	currentNode := ll.head
	for currentNode != nil {
		nextNode := currentNode.next

		if f(i, currentNode.value) {
			return ll.RemoveAt(i), true
		}

		currentNode = nextNode
		i++
	}

	ok = false
	return
}

// RemoveAll removes all values of T that pass the test implemented by the provided function.
func (ll *LinkedList[T]) RemoveAll(f func(idx int, value T) bool) []T {
	res := []T{}

	i := 0
	currentNode := ll.head
	for currentNode != nil {
		nextNode := currentNode.next

		if f(i, currentNode.value) {
			res = append(res, ll.RemoveAt(i))
			i--
		}

		currentNode = nextNode
		i++
	}

	return res
}

// FindIndex returns the first index of value T that pass the test implemented by the provided function.
func (ll *LinkedList[T]) FindIndex(f func(value T) bool) (idx int, ok bool) {
	i := 0
	currentNode := ll.head
	for currentNode != nil {
		nextNode := currentNode.next

		if f(currentNode.value) {
			return i, true
		}

		currentNode = nextNode
		i++
	}

	return 0, false
}

// FindIndex returns all indexes of value T that pass the test implemented by the provided function.
func (ll *LinkedList[T]) FindIndexes(f func(value T) bool) []int {
	res := []int{}

	i := 0
	currentNode := ll.head
	for currentNode != nil {
		nextNode := currentNode.next

		if f(currentNode.value) {
			res = append(res, i)
		}

		currentNode = nextNode
		i++
	}

	return res
}

// Len gets linked list length.
func (ll *LinkedList[T]) Len() int {
	return ll.length
}

// Free frees the linked list.
func (ll *LinkedList[T]) Free() {
	currentNode := ll.head

	for currentNode != nil {
		nextNode := currentNode.next
		mm.Free(ll.ctx, currentNode)
		currentNode = nextNode
	}

	mm.Free(ll.ctx, ll)
}
