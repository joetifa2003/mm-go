package hashmap

import (
	"github.com/mitchellh/hashstructure/v2"

	"github.com/joetifa2003/mm-go"
	"github.com/joetifa2003/mm-go/linkedlist"
	"github.com/joetifa2003/mm-go/vector"
)

type pair[K comparable, V any] struct {
	key   K
	value V
}

// Hashmap Manually managed hashmap,
// keys can be hashmap.String, hashmap.Int or any type that
// implements the hashmap.Hashable interface
type Hashmap[K comparable, V any] struct {
	pairs      *vector.Vector[*linkedlist.LinkedList[pair[K, V]]]
	totalTaken int
}

// New creates a new Hashmap with key of type K and value of type V
func New[K comparable, V any]() *Hashmap[K, V] {
	hm := mm.Alloc[Hashmap[K, V]]()
	hm.pairs = vector.New[*linkedlist.LinkedList[pair[K, V]]](8)
	return hm
}

func (hm *Hashmap[K, V]) extend() {
	newPairs := vector.New[*linkedlist.LinkedList[pair[K, V]]](hm.pairs.Len() * 2)
	oldPairs := hm.pairs
	defer oldPairs.Free()

	hm.totalTaken = 0
	hm.pairs = newPairs

	for _, pairs := range oldPairs.Slice() {
		if pairs == nil {
			continue
		}

		pairs.ForEach(func(idx int, p pair[K, V]) {
			hm.Insert(p.key, p.value)
		})
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

	hash, err := hashstructure.Hash(key, hashstructure.FormatV2, nil)
	if err != nil {
		panic(err)
	}

	idx := int(hash % uint64(hm.pairs.Len()))
	pairs := hm.pairs.At(idx)
	if pairs == nil {
		newPairs := linkedlist.New[pair[K, V]]()
		hm.pairs.Set(idx, newPairs)
		pairs = newPairs
	}
	pairs.PushBack(pair[K, V]{key: key, value: value})
	hm.totalTaken++
}

// Get takes key K and return value V
func (hm *Hashmap[K, V]) Get(key K) (value V, exists bool) {
	hash, err := hashstructure.Hash(key, hashstructure.FormatV2, nil)
	if err != nil {
		panic(err)
	}

	idx := int(hash % uint64(hm.pairs.Len()))
	pairs := hm.pairs.At(idx)
	if pairs == nil {
		return *new(V), false
	}

	pairIdx, ok := pairs.FindIndex(func(value pair[K, V]) bool {
		return value.key == key
	})
	if !ok {
		return *new(V), false
	}

	return pairs.At(pairIdx).value, ok
}

// GetPtr takes key K and return a pointer to value V
func (hm *Hashmap[K, V]) GetPtr(key K) (value *V, exists bool) {
	hash, err := hashstructure.Hash(key, hashstructure.FormatV2, nil)
	if err != nil {
		panic(err)
	}

	idx := int(hash % uint64(hm.pairs.Len()))
	pairs := hm.pairs.At(idx)
	if pairs == nil {
		return nil, false
	}

	pairIdx, ok := pairs.FindIndex(func(value pair[K, V]) bool {
		return value.key == key
	})
	if !ok {
		return nil, false
	}

	return &pairs.AtPtr(pairIdx).value, ok
}

// ForEach iterates through all key/value pairs
func (hm *Hashmap[K, V]) ForEach(f func(key K, value V)) {
	for _, pairs := range hm.pairs.Slice() {
		if pairs == nil {
			continue
		}

		pairs.ForEach(func(idx int, p pair[K, V]) {
			f(p.key, p.value)
		})
	}
}

// Values returns all values as a slice
func (hm *Hashmap[K, V]) Values() []V {
	res := make([]V, 0)

	for _, pairs := range hm.pairs.Slice() {
		if pairs == nil {
			continue
		}

		pairs.ForEach(func(idx int, p pair[K, V]) {
			res = append(res, p.value)
		})
	}

	return res
}

// Keys returns all keys as a slice
func (hm *Hashmap[K, V]) Keys() []K {
	res := make([]K, 0)

	for _, pairs := range hm.pairs.Slice() {
		if pairs == nil {
			continue
		}

		pairs.ForEach(func(idx int, p pair[K, V]) {
			res = append(res, p.key)
		})
	}

	return res
}

// Delete delete value with key K
func (hm *Hashmap[K, V]) Delete(key K) {
	hash, err := hashstructure.Hash(key, hashstructure.FormatV2, nil)
	if err != nil {
		panic(err)
	}

	idx := int(hash % uint64(hm.pairs.Len()))
	pairs := hm.pairs.At(idx)
	if pairs == nil {
		return
	}

	pairs.Remove(func(idx int, p pair[K, V]) bool {
		return p.key == key
	})
}

// Free frees the Hashmap
func (hm *Hashmap[K, V]) Free() {
	for _, pairs := range hm.pairs.Slice() {
		if pairs != nil {
			pairs.Free()
		}
	}
	hm.pairs.Free()
	mm.Free(hm)
}
