package mm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func testPushAndPop(t *testing.T) {
	assert := assert.New(t)

	ll := NewLinkedList[int]()
	defer ll.Free()

	ll.PushBack(15)
	ll.PushBack(16)
	ll.PushFront(14)

	assert.Equal(16, ll.PopBack())
	assert.Equal(2, ll.Len())
	assert.Equal(15, ll.PopBack())
	assert.Equal(1, ll.Len())
	assert.Equal(14, ll.PopBack())
	assert.Equal(0, ll.Len())

	ll.PushFront(10)

	assert.Equal(10, ll.PopBack())

	// pop on empty linked list
	assert.Panics(func() {
		ll.PopBack()
	})
	assert.Panics(func() {
		ll.PopFront()
	})

	ll.PushBack(10)
	ll.PushBack(15)

	assert.Equal(10, ll.PopFront())
	assert.Equal(15, ll.PopFront())

	ll.PushFront(15)
}

func testForEach(t *testing.T) {
	assert := assert.New(t)

	ll := NewLinkedList[int]()
	defer ll.Free()

	ll.PushBack(2)
	ll.PushBack(3)
	ll.PushBack(4)

	even := 0
	idxSum := 0
	ll.ForEach(func(idx, value int) {
		if value%2 == 0 {
			even++
		}

		idxSum += idx
	})

	assert.Equal(2, even)
	assert.Equal(3, idxSum)
}

func testIndexing(t *testing.T) {
	assert := assert.New(t)

	ll := NewLinkedList[int]()
	defer ll.Free()

	ll.PushBack(1)
	ll.PushBack(2)

	assert.Equal(1, ll.At(0))
	assert.Equal(2, ll.At(1))

	assert.Panics(func() {
		ll.At(999)
	})
}

func testRemove(t *testing.T) {
	assert := assert.New(t)

	ll := NewLinkedList[int]()
	defer ll.Free()

	ll.PushBack(1)
	ll.PushBack(2)
	ll.PushBack(3)

	assert.Equal(2, ll.RemoveAt(1))
	assert.Equal(2, ll.Len())
	assert.Equal(1, ll.RemoveAt(0))
	assert.Equal(1, ll.Len())
	assert.Equal(3, ll.PopFront())

	ll.PushBack(1)
	ll.PushBack(2)
	ll.PushBack(3)
	ll.PushBack(4)

	// remove the first even element
	ll.Remove(func(idx, value int) bool {
		return value%2 == 0
	})

	assert.Equal(4, ll.PopBack())
	assert.Equal(3, ll.PopBack())
	assert.Equal(1, ll.PopBack())

	ll.PushBack(1)
	ll.PushBack(2)
	ll.PushBack(6)

	ll.RemoveAll(func(idx, value int) bool {
		return value%2 == 0
	})

	assert.Equal(1, ll.PopBack())
}

func TestLinkedList(t *testing.T) {
	t.Run("push and pop", testPushAndPop)
	t.Run("for each", testForEach)
	t.Run("indexing", testIndexing)
	t.Run("remove", testRemove)
}
