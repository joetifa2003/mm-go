package typedarena_test

import (
	"testing"

	"github.com/joetifa2003/mm-go/typedarena"
	"github.com/stretchr/testify/assert"
)

func TestTypedArena(t *testing.T) {
	assert := assert.New(t)

	arena := typedarena.New[int](4)
	defer arena.Free()

	int1 := arena.Alloc()      // allocates 1 int from arena
	*int1 = 1                  // changing it's value
	ints := arena.AllocMany(2) // allocates 2 ints from the arena and returns a slice representing the heap (instead of pointer arithmetic)
	ints[0] = 2                // changing the first value
	ints[1] = 3                // changing the second value

	// you can also take pointers from the slice
	intPtr1 := &ints[0]
	*intPtr1 = 15

	assert.Panics(func() {
		arena.AllocMany(10) // bigger than chunk size
	})

	arena.AllocMany(4) // creates a new chunk to fit in new values

	arena.Alloc() // also creates a new chunk to fit in new values

	assert.Equal(1, *int1)
	assert.Equal(2, len(ints))
	assert.Equal(15, ints[0])
	assert.Equal(3, ints[1])
}
