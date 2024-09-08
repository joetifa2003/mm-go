package mm_test

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/joetifa2003/mm-go"
	"github.com/joetifa2003/mm-go/allocator"
	"github.com/joetifa2003/mm-go/batchallocator"
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

func linkedListFree[T any](alloc allocator.Allocator, list *LinkedList[T]) {
	currentNode := list.head
	for currentNode != nil {
		nextNode := currentNode.next
		allocator.Free(alloc, currentNode)
		currentNode = nextNode
	}
	allocator.Free(alloc, list)
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

func benchLinkedListManaged(b *testing.B, size int) {
	list := &LinkedList[int]{}
	for i := range size {
		linkedListPushManaged(list, i)
	}
	assertLinkedList(b, list)
}

func benchLinkedListCAlloc(b *testing.B, size int) {
	alloc := allocator.NewCallocator()
	defer alloc.Destroy()

	list := allocator.Alloc[LinkedList[int]](alloc)
	defer linkedListFree(alloc, list)

	for i := range size {
		linkedListPushAlloc(alloc, list, i)
	}

	assertLinkedList(b, list)
}

func benchLinkedListBatchAllocator(b *testing.B, size int, bucketSize int) {
	alloc := batchallocator.New(allocator.NewCallocator(),
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
