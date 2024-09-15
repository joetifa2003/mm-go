package mm_test

import (
	"fmt"

	"github.com/joetifa2003/mm-go"
)

func ExampleSizeOf() {
	fmt.Println(mm.SizeOf[int32]())
	fmt.Println(mm.SizeOf[int64]())
	// Output:
	// 4
	// 8
}
