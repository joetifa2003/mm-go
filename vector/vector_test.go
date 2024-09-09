package vector_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/joetifa2003/mm-go/allocator"
	"github.com/joetifa2003/mm-go/vector"
)

func TestVector(t *testing.T) {
	alloc := allocator.NewC()

	assert := assert.New(t)

	v := vector.New[int](alloc)
	defer v.Free()

	v.Push(1)
	v.Push(2)
	v.Push(3)

	assert.Equal(3, v.Len())
	assert.Equal(4, v.Cap())
	assert.Equal([]int{1, 2, 3}, v.Slice())
	assert.Equal(3, v.Pop())
	assert.Equal(2, v.Pop())
	assert.Equal(1, v.Pop())

	v.Push(1)
	v.Push(2)

	assert.Equal(1, v.At(0))
	assert.Equal(2, v.At(1))
	assert.Panics(func() {
		v.At(3)
	})
}

func TestVectorInit(t *testing.T) {
	t.Run("Init with no args", func(t *testing.T) {
		alloc := allocator.NewC()
		assert := assert.New(t)

		v := vector.New[int](alloc)
		defer v.Free()

		assert.Equal(0, v.Len())
		assert.Equal(1, v.Cap())
	})

	t.Run("Init with one arg", func(t *testing.T) {
		alloc := allocator.NewC()
		assert := assert.New(t)

		v := vector.New[int](alloc, 5)
		defer v.Free()

		assert.Equal(5, v.Len())
		assert.Equal(5, v.Cap())
	})

	t.Run("Init with two args", func(t *testing.T) {
		alloc := allocator.NewC()
		assert := assert.New(t)

		v := vector.New[int](alloc, 5, 6)
		defer v.Free()

		assert.Equal(5, v.Len())
		assert.Equal(6, v.Cap())
	})

	t.Run("Init vector with slice", func(t *testing.T) {
		alloc := allocator.NewC()
		assert := assert.New(t)

		v := vector.Init(alloc, 1, 2, 3)
		defer v.Free()

		assert.Equal(3, v.Len())
		assert.Equal(3, v.Cap())

		assert.Equal(3, v.Pop())
		assert.Equal(2, v.Pop())
		assert.Equal(1, v.Pop())
	})
}
