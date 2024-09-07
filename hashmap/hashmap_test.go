package hashmap_test

import (
	"testing"

	"github.com/joetifa2003/mm-go/hashmap"
)

const TIMES = 5000

func BenchmarkHashmapGo(b *testing.B) {
	for i := 0; i < b.N; i++ {
		h := newMap()

		for i := 0; i < TIMES; i++ {
			h[i] = i
		}
	}
}

func BenchmarkHashmap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		h := hashmap.New[int, int]()

		for i := 0; i < TIMES; i++ {
			h.Insert(i, i)
		}

		h.Free()
	}
}

func newMap() map[int]int {
	return make(map[int]int)
}
