package hashmap_test

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/joetifa2003/mm-go/hashmap"
)

const TIMES = 5000

func BenchmarkHashmap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		h := hashmap.New[hashmap.Int, int]()

		for i := 0; i < TIMES; i++ {
			h.Insert(hashmap.Int(i), i)
		}
		h.Free()
		runtime.GC()
	}
}

func BenchmarkHashmapGo(b *testing.B) {
	for i := 0; i < b.N; i++ {
		h := map[string]int{}
		for i := 0; i < TIMES; i++ {
			h[fmt.Sprint(i)] = i
		}
		_ = h
		runtime.GC()
	}
}
