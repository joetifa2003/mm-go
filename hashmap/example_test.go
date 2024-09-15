package hashmap

import (
	"fmt"

	"github.com/joetifa2003/mm-go/allocator"
	"github.com/joetifa2003/mm-go/batchallocator"
)

func Example() {
	alloc := batchallocator.New(allocator.NewC())
	defer alloc.Destroy()

	hm := New[int, int](alloc)
	defer hm.Free() // can be removed

	hm.Set(1, 10)
	hm.Set(2, 20)
	hm.Set(3, 30)

	sumKeys := 0
	sumValues := 0
	for k, v := range hm.Iter() {
		sumKeys += k
		sumValues += v
	}

	fmt.Println(sumKeys)
	fmt.Println(sumValues)

	// Output:
	// 6
	// 60
}
