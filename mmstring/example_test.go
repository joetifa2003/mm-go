package mmstring_test

import (
	"fmt"

	"github.com/joetifa2003/mm-go/allocator"
	"github.com/joetifa2003/mm-go/batchallocator"
	"github.com/joetifa2003/mm-go/mmstring"
	"github.com/joetifa2003/mm-go/vector"
)

func Example() {
	alloc := allocator.NewC()
	defer alloc.Destroy()

	s := mmstring.New(alloc)
	defer s.Free()

	s.AppendGoString("Hello ")
	s.AppendGoString("World")

	s2 := mmstring.From(alloc, "Foo Bar")
	defer s2.Free()

	fmt.Println(s.GetGoString())
	fmt.Println(s2.GetGoString())

	// Output:
	// Hello World
	// Foo Bar
}

func Example_datastructures() {
	alloc := batchallocator.New(allocator.NewC())
	defer alloc.Destroy() // all the memory allocated bellow will be freed, no need to free it manually.

	m := vector.New[*mmstring.MMString](alloc)
	m.Push(mmstring.From(alloc, "hello"))
	m.Push(mmstring.From(alloc, "world"))

	for k, v := range m.Iter() {
		fmt.Println(k, v.GetGoString())
	}

	// Output:
	// 0 hello
	// 1 world
}
