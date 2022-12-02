package mm

import "testing"

func TestTypedArena(t *testing.T) {
	ta := NewTypedArena[int](1)
	int1 := ta.Alloc()
	*int1 = 15
	int2 := ta.Alloc()
	*int2 = 20
	ta.Free()
}
