package mmstring

import (
	"context"

	"github.com/joetifa2003/mm-go/vector"
)

// MMString is a manually manged string that is basically a *Vector[rune]
// and contains all the methods of a vector plus additional helper functions
type MMString struct {
	*vector.Vector[rune]
}

// New create a new manually managed string
func New(ctx context.Context) MMString {
	mmString := MMString{}
	mmString.Vector = vector.New[rune](ctx)
	return mmString
}

// From creates a new manually managed string,
// And initialize it with a go string
func From(ctx context.Context, input string) MMString {
	mmString := New(ctx)

	for _, r := range input {
		mmString.Push(r)
	}

	return mmString
}

// GetGoString returns go string from manually managed string.
// CAUTION: You also have to free the MMString
func (s *MMString) GetGoString() string {
	runes := s.Slice()
	return string(runes)
}

// AppendGoString appends go string to manually managed string
func (s *MMString) AppendGoString(input string) {
	for _, r := range input {
		s.Push(r)
	}
}

// Free frees MMString
func (s *MMString) Free() {
	s.Vector.Free()
}
