package mmstring_test

import (
	"testing"

	"github.com/joetifa2003/mm-go/mmstring"
	"github.com/stretchr/testify/assert"
)

func TestString(t *testing.T) {
	assert := assert.New(t)
	mmString := mmstring.From("hi")
	defer mmString.Free()

	assert.Equal('h', mmString.At(0))
	assert.Equal('i', mmString.At(1))
	assert.Panics(func() {
		mmString.At(3)
	})

	assert.Equal("hi", mmString.GetGoString())
}
