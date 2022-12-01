package mm

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

// DoNotOptimise prevent the compiler removing the function bodies
var DoNotOptimiseInts []int
var DoNotOptimiseVector *Vector[int]

func BenchmarkSlice(b *testing.B) {
	// Start with a clean slate
	DoNotOptimiseInts = nil
	runtime.GC()
	b.ResetTimer()

	for i := b.N; i != 0; i-- {
		numbers := []int{}

		for j := 0; j < LOOP_TIMES; j++ {
			numbers = append(numbers, j)
		}

		for j := 0; j < LOOP_TIMES; j++ {
			// Pop
			numbers = numbers[:len(numbers)-1]
		}

		DoNotOptimiseInts = numbers
	}

	runtime.GC()

	DoNotOptimiseInts = nil
}

func BenchmarkVector(b *testing.B) {
	// Start with a clean slate
	DoNotOptimiseVector = nil
	runtime.GC()
	b.ResetTimer()

	for i := b.N; i != 0; i-- {
		numbersVec := NewVector[int]()

		for j := 0; j < LOOP_TIMES; j++ {
			numbersVec.Push(j)
		}

		for j := 0; j < LOOP_TIMES; j++ {
			numbersVec.Pop()
		}

		DoNotOptimiseVector = numbersVec
		numbersVec.Free()
	}

	DoNotOptimiseVector = nil
}

func TestVector(t *testing.T) {
	assert := assert.New(t)

	v := NewVector[int]()
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
}

func TestVectorInit(t *testing.T) {
	t.Run("Init with no args", func(t *testing.T) {
		assert := assert.New(t)

		v := NewVector[int]()
		defer v.Free()

		assert.Equal(0, v.Len())
		assert.Equal(1, v.Cap())
	})

	t.Run("Init with one arg", func(t *testing.T) {
		assert := assert.New(t)

		v := NewVector[int](5)
		defer v.Free()

		assert.Equal(5, v.Len())
		assert.Equal(5, v.Cap())
	})

	t.Run("Init with two args", func(t *testing.T) {
		assert := assert.New(t)

		v := NewVector[int](5, 6)
		defer v.Free()

		assert.Equal(5, v.Len())
		assert.Equal(6, v.Cap())
	})

	t.Run("Init vector with slice", func(t *testing.T) {
		assert := assert.New(t)

		v := InitVector(1, 2, 3)
		defer v.Free()

		assert.Equal(3, v.Len())
		assert.Equal(3, v.Cap())

		assert.Equal(3, v.Pop())
		assert.Equal(2, v.Pop())
		assert.Equal(1, v.Pop())
	})
}
