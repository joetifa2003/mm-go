package hashmap_test

import (
	"runtime"
	"testing"

	"github.com/joetifa2003/mm-go/hashmap"
	"github.com/joetifa2003/mm-go/mmstring"
)

const TIMES = 5000

func BenchmarkHashmap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		h := hashmap.New[int, *mmstring.MMString]()

		for i := 0; i < TIMES; i++ {
			h.Insert(i, mmstring.From("foo bar"))
		}

		h.Free()
		runtime.GC()
	}
}

func BenchmarkHashmapGo(b *testing.B) {
	for i := 0; i < b.N; i++ {
		h := map[string]string{}

		for i := 0; i < TIMES; i++ {
			h["foo"] = "foo bar"
		}

		_ = h
		runtime.GC()
	}
}
