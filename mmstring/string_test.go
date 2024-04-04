package mmstring_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/joetifa2003/mm-go/mmstring"
)

func TestString(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()

	mmString := mmstring.From(ctx, "hi")
	defer mmString.Free()

	assert.Equal('h', mmString.At(0))
	assert.Equal('i', mmString.At(1))
	assert.Panics(func() {
		mmString.At(3)
	})

	assert.Equal("hi", mmString.GetGoString())
}
