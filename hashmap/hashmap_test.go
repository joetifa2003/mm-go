package hashmap_test

import (
	"fmt"
	"testing"

	"github.com/joetifa2003/mm-go/hashmap"
	"github.com/stretchr/testify/assert"
)

func TestHashmap(t *testing.T) {
	t.Run("insert", testInsert)
	t.Run("keys and values", testHashmapKeysValues)

	t.Run("delete", testDelete)
}

func testInsert(t *testing.T) {
	assert := assert.New(t)

	hm := hashmap.New[hashmap.String, string]()
	defer hm.Free()

	const TIMES = 1000

	for i := 0; i < TIMES; i++ {
		hm.Insert(hashmap.String(fmt.Sprint(i)), fmt.Sprint(i))
	}

	for i := 0; i < TIMES; i++ {
		if i%2 == 0 {
			hm.Delete(hashmap.String(fmt.Sprint(i)))
		}
	}

	for i := 0; i < TIMES; i++ {
		if i%2 == 0 {
			continue
		}

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

func testHashmapKeysValues(t *testing.T) {
	assert := assert.New(t)

	hm := hashmap.New[hashmap.String, int]()
	defer hm.Free()

	hm.Insert("foo", 0)
	hm.Insert("bar", 1)

	values := hm.Values()
	assert.Equal(2, len(values))
	assert.Contains(values, 0)
	assert.Contains(values, 1)

	keys := hm.Keys()
	assert.Equal(2, len(keys))
	assert.Contains(keys, hashmap.String("foo"))
	assert.Contains(keys, hashmap.String("bar"))

	type pair struct {
		key   hashmap.String
		value int
	}

	expectedPairs := []pair{
		{key: "foo", value: 0},
		{key: "bar", value: 1},
	}

	i := 0
	hm.ForEach(func(key hashmap.String, value int) {
		assert.Contains(expectedPairs, pair{key: key, value: value})

		i++
	})

	assert.Equal(i, 2)
}

func testDelete(t *testing.T) {
	assert := assert.New(t)

	hm := hashmap.New[hashmap.Int, int]()
	hm.Insert(1, 1)
	assert.Equal([]int{1}, hm.Values())
	hm.Delete(1)
	assert.Equal([]int{}, hm.Values())
}

const TIMES = 15000

func BenchmarkHashMap(b *testing.B) {
	for n := 0; n < b.N; n++ {
		hm := hashmap.New[hashmap.Int, int]()
		for i := 0; i < TIMES; i++ {
			hm.Insert(hashmap.Int(i), i)
		}

		sum := 0
		for i := 0; i < TIMES; i++ {
			v, _ := hm.Get(hashmap.Int(i))
			sum += v
		}

		_ = sum

		hm.Free()
	}
}

func BenchmarkGoMap(b *testing.B) {
	for n := 0; n < b.N; n++ {
		hm := map[int]int{}
		for i := 0; i < TIMES; i++ {
			hm[i] = i
		}

		sum := 0
		for i := 0; i < TIMES; i++ {
			sum += hm[i]
		}

		_ = sum
	}
}
