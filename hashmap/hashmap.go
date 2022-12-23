package hashmap

import (
	"hash/fnv"

	"github.com/joetifa2003/mm-go"
	"github.com/joetifa2003/mm-go/vector"
)

// Hashable keys must implement this interface
// or use type hashmap.String and hashmap.Int
// which implements the interface
type Hashable interface {
	comparable
	Hash() uint32
}

// String a string type that implements Hashable Interface
type String string

func (s String) Hash() uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

// Int an int type that implements Hashable Interface
type Int int

func (i Int) Hash() uint32 {
	return uint32(i)
}

type pair[K Hashable, V any] struct {
	key   K
	value V
	taken bool
}

// Hashmap Manually managed hashmap,
// keys can be hashmap.String, hashmap.Int or any type that
// implements the hashmap.Hashable interface
type Hashmap[K Hashable, V any] struct {
	pairs      *vector.Vector[pair[K, V]]
	totalTaken int
}

// New creates a new Hashmap with key of type K and value of type V
func New[K Hashable, V any]() *Hashmap[K, V] {
	hm := mm.Alloc[Hashmap[K, V]]()
	hm.pairs = vector.New[pair[K, V]](1)
	return hm
}

func (hm *Hashmap[K, V]) extend() {
	newPairs := vector.New[pair[K, V]](hm.pairs.Len() * 2)
	oldPairs := hm.pairs
	defer oldPairs.Free()

	hm.totalTaken = 0
	hm.pairs = newPairs

	for _, c := range oldPairs.Slice() {
		hm.Insert(c.key, c.value)
	}
}

// Insert inserts a new value V if key K doesn't exist,
// Otherwise update the key K with value V
func (hm *Hashmap[K, V]) Insert(key K, value V) {
	if ptr, exists := hm.GetPtr(key); exists {
		*ptr = value
		return
	}

	if hm.totalTaken == hm.pairs.Len() {
		hm.extend()
	}

	idx := int(key.Hash() % uint32(hm.pairs.Len()))
	for hm.pairs.At(idx).taken {
		idx = (idx + 1) % hm.pairs.Len()
	}

	hm.pairs.Set(idx, pair[K, V]{key: key, value: value, taken: true})
	hm.totalTaken++
}

func (hm *Hashmap[K, V]) getIdx(key K) (index int, exists bool) {
	idx := int(key.Hash() % uint32(hm.pairs.Len()))
	startingIdx := idx

	for {
		pair := hm.pairs.At(idx)
		if pair.taken && pair.key == key {
			return idx, true
		}

		idx = (idx + 1) % hm.pairs.Len()

		if idx == startingIdx {
			return 0, false
		}
	}
}

// Get takes key K and return value V
func (hm *Hashmap[K, V]) Get(key K) (value V, exists bool) {
	idx, exists := hm.getIdx(key)
	if exists {
		value = hm.pairs.At(idx).value
	}
	return value, exists
}

// GetPtr takes key K and return a pointer to value V
func (hm *Hashmap[K, V]) GetPtr(key K) (value *V, exists bool) {
	idx, exists := hm.getIdx(key)
	if exists {
		value = &hm.pairs.AtPtr(idx).value
	}
	return value, exists
}

// Free frees the Hashmap
func (hm *Hashmap[K, V]) Free() {
	hm.pairs.Free()
	mm.Free(hm)
}
