package mm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTypedArena(t *testing.T) {
	assert := assert.New(t)

	ta := NewTypedArena[int](3)
	int1 := ta.Alloc()
	*int1 = 1
	ints := ta.AllocMany(2)
	ints[0] = 2
	ints[1] = 3

	assert.Equal(1, ta.chunks.len)
	assert.Equal(1, *int1)
	assert.Equal(2, len(ints))
	assert.Equal(2, ints[0])
	assert.Equal(3, ints[1])
	ta.Free()
}
