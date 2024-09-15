package mmstring_test

import (
	"fmt"

	"github.com/joetifa2003/mm-go/allocator"
	"github.com/joetifa2003/mm-go/batchallocator"
	"github.com/joetifa2003/mm-go/hashmap"
	"github.com/joetifa2003/mm-go/mmstring"
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

func Example_complex() {
	alloc := batchallocator.New(allocator.NewC())
	defer alloc.Destroy() // all the memory allocated bellow will be freed, no need to free it manually.

	m := hashmap.New[*mmstring.MMString, *mmstring.MMString](alloc)
	m.Insert(mmstring.From(alloc, "hello"), mmstring.From(alloc, "world"))
	m.Insert(mmstring.From(alloc, "foo"), mmstring.From(alloc, "bar"))

	for k, v := range m.Iter() {
		fmt.Println(k.GetGoString(), v.GetGoString())
	}

	// Output:
	// hello world
	// foo bar
}
