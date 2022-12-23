package hashmap_test

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/joetifa2003/mm-go/hashmap"
	"github.com/stretchr/testify/assert"
)

func TestHashmap(t *testing.T) {
	assert := assert.New(t)

	hm := hashmap.New[hashmap.String, string]()
	defer hm.Free()

	const TIMES = 1000

	for i := 0; i < TIMES; i++ {
		hm.Insert(hashmap.String(fmt.Sprint(i)), fmt.Sprint(i))
	}

	for i := 0; i < TIMES; i++ {
		value, exits := hm.Get(hashmap.String(fmt.Sprint(i)))
		assert.Equal(fmt.Sprint(i), value)
		assert.Equal(true, exits)
	}

	hm.Insert("name", "Foo")
	value, exists := hm.Get("name")
	assert.Equal(true, exists)
	assert.Equal("Foo", value)

	hm.Insert("name", "Bar")
	value, exists = hm.Get("name")
	assert.Equal(true, exists)
	assert.Equal("Bar", value)
}

const TIMES = 15000

func BenchmarkHashMap(b *testing.B) {
	for n := 0; n < b.N; n++ {
		hm := hashmap.New[hashmap.Int, string]()
		for i := 0; i < TIMES; i++ {
			hm.Insert(hashmap.Int(i), fmt.Sprint(i))
		}
		hm.Free()
		runtime.GC()
	}
}

func BenchmarkGoMap(b *testing.B) {
	for n := 0; n < b.N; n++ {
		hm := map[int]string{}
		for i := 0; i < TIMES; i++ {
			hm[i] = fmt.Sprint(i)
		}
		runtime.GC()
	}
}
