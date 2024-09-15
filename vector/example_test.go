package vector_test

import (
	"fmt"

	"github.com/joetifa2003/mm-go/allocator"
	"github.com/joetifa2003/mm-go/vector"
)

func Example() {
	alloc := allocator.NewC()
	v := vector.New[int](alloc)
	v.Push(1)
	v.Push(2)
	v.Push(3)

	fmt.Println("Length:", v.Len())
	for i := 0; i < v.Len(); i++ {
		fmt.Println(v.At(i))
	}

	for _, k := range v.Iter() {
		fmt.Println(k)
	}

	// Output:
	// Length: 3
	// 1
	// 2
	// 3
	// 1
	// 2
	// 3
}
