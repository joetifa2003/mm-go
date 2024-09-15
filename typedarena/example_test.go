package typedarena_test

import (
	"fmt"

	"github.com/joetifa2003/mm-go/allocator"
	"github.com/joetifa2003/mm-go/typedarena"
)

type Entity struct {
	VelocityX float32
	VelocityY float32
	PositionX float32
	PositionY float32
}

func Example() {
	alloc := allocator.NewC()
	defer alloc.Destroy()

	arena := typedarena.New[Entity](
		alloc,
		10,
	)
	defer arena.Free() // frees all memory

	for i := 0; i < 10; i++ {
		e := arena.Alloc() // *Entity
		e.VelocityX = float32(i)
		e.VelocityY = float32(i)
		e.PositionX = float32(i)
		e.PositionY = float32(i)
		fmt.Println(e.VelocityX, e.VelocityY, e.PositionX, e.PositionY)
	}

	entities := arena.AllocMany(10) // allocate slice of 10 entities (cannot exceed 10 here because chunk size is 10 above, this limitation doesn't exist in batchallocator)

	_ = entities

	// Output:
	// 0 0 0 0
	// 1 1 1 1
	// 2 2 2 2
	// 3 3 3 3
	// 4 4 4 4
	// 5 5 5 5
	// 6 6 6 6
	// 7 7 7 7
	// 8 8 8 8
	// 9 9 9 9
}
