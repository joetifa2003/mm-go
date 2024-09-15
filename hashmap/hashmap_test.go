package hashmap_test

import (
	"runtime"
	"testing"

	"github.com/joetifa2003/mm-go/allocator"
	"github.com/joetifa2003/mm-go/batchallocator"
	"github.com/joetifa2003/mm-go/hashmap"
)

const TIMES = 500

func BenchmarkHashmapGo(b *testing.B) {
	for i := 0; i < b.N; i++ {
		h := newMap()

		for i := 0; i < TIMES; i++ {
			h[i] = i
		}

		runtime.GC()
	}
}

func BenchmarkHashmapCAlloc(b *testing.B) {
	for i := 0; i < b.N; i++ {
		alloc := allocator.NewC()
		h := hashmap.New[int, int](alloc)

		for i := 0; i < TIMES; i++ {
			h.Set(i, i)
		}

		h.Free()
		alloc.Destroy()
	}
}

func BenchmarkHashmapBatchAlloc(b *testing.B) {
	for i := 0; i < b.N; i++ {
		alloc := batchallocator.New(allocator.NewC())
		h := hashmap.New[int, int](alloc)

		for i := 0; i < TIMES; i++ {
			h.Set(i, i)
		}

		h.Free()
		alloc.Destroy()
	}
}

func newMap() map[int]int {
	return make(map[int]int)
}
