package mm_test

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/joetifa2003/mm-go"
	"github.com/joetifa2003/mm-go/allocator"
	"github.com/joetifa2003/mm-go/batchallocator"
	"github.com/joetifa2003/mm-go/typedarena"
)

type Node[T any] struct {
	value T
	next  *Node[T]
	prev  *Node[T]
}

type LinkedList[T any] struct {
	head *Node[T]
	tail *Node[T]
}

func linkedListPushManaged[T any](list *LinkedList[T], value T) {
	node := &Node[T]{value: value}
	if list.head == nil {
		list.head = node
		list.tail = node
	} else {
		list.tail.next = node
		node.prev = list.tail
		list.tail = node
	}
}

func linkedListPushAlloc[T any](alloc allocator.Allocator, list *LinkedList[T], value T) {
	node := allocator.Alloc[Node[T]](alloc)
	node.value = value

	if list.head == nil {
		list.head = node
		list.tail = node
	} else {
		list.tail.next = node
		node.prev = list.tail
		list.tail = node
	}
}

func linkedListPushArena[T any](arena *typedarena.TypedArena[Node[T]], list *LinkedList[T], value T) {
	node := arena.Alloc()
	node.value = value

	if list.head == nil {
		list.head = node
		list.tail = node
	} else {
		list.tail.next = node
		node.prev = list.tail
		list.tail = node
	}
}

func linkedListFree[T any](alloc allocator.Allocator, list *LinkedList[T]) {
	currentNode := list.head
	for currentNode != nil {
		nextNode := currentNode.next
		allocator.Free(alloc, currentNode)
		currentNode = nextNode
	}
}

const LINKED_LIST_SIZE = 10000

func BenchmarkLinkedListManaged(b *testing.B) {
	for range b.N {
		benchLinkedListManaged(b, LINKED_LIST_SIZE)
		runtime.GC()
	}
}

func BenchmarkLinkedListCAlloc(b *testing.B) {
	for range b.N {
		benchLinkedListCAlloc(b, LINKED_LIST_SIZE)
	}
}

func BenchmarkLinkedListBatchAllocator(b *testing.B) {
	for _, bucketSize := range []int{100, 200, 500, LINKED_LIST_SIZE} {
		b.Run(fmt.Sprintf("bucket size %d", bucketSize), func(b *testing.B) {
			for range b.N {
				benchLinkedListBatchAllocator(b, LINKED_LIST_SIZE, bucketSize)
			}
		})
	}
}

func BenchmarkLinkedListTypedArena(b *testing.B) {
	for _, chunkSize := range []int{100, 200, 500, LINKED_LIST_SIZE} {
		b.Run(fmt.Sprintf("chunk size %d", chunkSize), func(b *testing.B) {
			for range b.N {
				benchLinkedListTypedArena(b, LINKED_LIST_SIZE, chunkSize)
			}
		})
	}
}

func benchLinkedListTypedArena(b *testing.B, size int, chunkSize int) {
	alloc := allocator.NewC()
	defer alloc.Destroy()

	arena := typedarena.New[Node[int]](alloc, chunkSize)
	defer arena.Free()

	list := allocator.Alloc[LinkedList[int]](alloc)
	defer allocator.Free(alloc, list)

	for i := range size {
		linkedListPushArena(arena, list, i)
	}

	assertLinkedList(b, list)
}

func benchLinkedListManaged(b *testing.B, size int) {
	list := &LinkedList[int]{}
	for i := range size {
		linkedListPushManaged(list, i)
	}
	assertLinkedList(b, list)
}

func benchLinkedListCAlloc(b *testing.B, size int) {
	alloc := allocator.NewC()
	defer alloc.Destroy()

	list := allocator.Alloc[LinkedList[int]](alloc)
	defer linkedListFree(alloc, list)

	for i := range size {
		linkedListPushAlloc(alloc, list, i)
	}

	assertLinkedList(b, list)
}

func benchLinkedListBatchAllocator(b *testing.B, size int, bucketSize int) {
	alloc := batchallocator.New(allocator.NewC(),
		batchallocator.WithBucketSize(mm.SizeOf[Node[int]]()*bucketSize),
	)
	defer alloc.Destroy()

	list := allocator.Alloc[LinkedList[int]](alloc)
	for i := range size {
		linkedListPushAlloc(alloc, list, i)
	}
	assertLinkedList(b, list)
}

func assertLinkedList(t *testing.B, list *LinkedList[int]) {
	if list.head == nil {
		t.Fatal("list head is nil")
	}
	if list.tail == nil {
		t.Fatal("list tail is nil")
	}

	currentNode := list.head
	i := 0
	for currentNode != nil {
		if currentNode.value != i {
			t.Fatalf("list value at index %d is %d, expected %d", i, currentNode.value, i)
		}
		i++
		currentNode = currentNode.next
	}
}
