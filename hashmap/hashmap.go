package hashmap

import (
	"iter"

	"github.com/dolthub/maphash"

	"github.com/joetifa2003/mm-go"
	"github.com/joetifa2003/mm-go/allocator"
	"github.com/joetifa2003/mm-go/linkedlist"
	"github.com/joetifa2003/mm-go/vector"
)

// Hashmap Manually managed hashmap,
type Hashmap[K comparable, V any] struct {
	pairs      *vector.Vector[*linkedlist.LinkedList[pair[K, V]]]
	totalTaken int
	mh         maphash.Hasher[K]
	alloc      allocator.Allocator
}

type pair[K comparable, V any] struct {
	key   K
	value V
}

// New creates a new Hashmap with key of type K and value of type V
func New[K comparable, V any](alloc allocator.Allocator) *Hashmap[K, V] {
	hm := allocator.Alloc[Hashmap[K, V]](alloc)
	hm.pairs = vector.New[*linkedlist.LinkedList[pair[K, V]]](alloc, 8)
	hm.mh = maphash.NewHasher[K]()
	hm.alloc = alloc
	return hm
}

func (hm *Hashmap[K, V]) extend() {
	newPairs := vector.New[*linkedlist.LinkedList[pair[K, V]]](hm.alloc, hm.pairs.Len()*2)
	oldPairs := hm.pairs
	defer oldPairs.Free()

	hm.totalTaken = 0
	hm.pairs = newPairs

	for _, pairs := range oldPairs.Iter() {
		if pairs == nil {
			continue
		}

		for _, p := range pairs.Iter() {
			hm.Insert(p.key, p.value)
		}
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

	hash := hm.mh.Hash(key)

	idx := int(hash % uint64(hm.pairs.Len()))
	pairs := hm.pairs.At(idx)
	if pairs == nil {
		newPairs := linkedlist.New[pair[K, V]](hm.alloc)
		hm.pairs.Set(idx, newPairs)
		pairs = newPairs
	}
	pairs.PushBack(pair[K, V]{key: key, value: value})
	hm.totalTaken++
}

// Get takes key K and return value V
func (hm *Hashmap[K, V]) Get(key K) (value V, exists bool) {
	hash := hm.mh.Hash(key)

	idx := int(hash % uint64(hm.pairs.Len()))
	pairs := hm.pairs.At(idx)
	if pairs == nil {
		return mm.Zero[V](), false
	}

	pairIdx, ok := pairs.FindIndex(func(value pair[K, V]) bool {
		return value.key == key
	})
	if !ok {
		return mm.Zero[V](), false
	}

	return pairs.At(pairIdx).value, ok
}

// GetPtr takes key K and return a pointer to value V
func (hm *Hashmap[K, V]) GetPtr(key K) (value *V, exists bool) {
	hash := hm.mh.Hash(key)

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

// Iter returns an iterator over all key/value pairs
func (hm *Hashmap[K, V]) Iter() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for _, pairs := range hm.pairs.Iter() {
			if pairs == nil {
				continue
			}

			for _, pair := range pairs.Iter() {
				if !yield(pair.key, pair.value) {
					return
				}
			}
		}
	}
}

// Values returns all values as a slice
func (hm *Hashmap[K, V]) Values() []V {
	res := make([]V, 0)

	for _, pairs := range hm.pairs.Iter() {
		if pairs == nil {
			continue
		}

		for _, p := range pairs.Iter() {
			res = append(res, p.value)
		}
	}

	return res
}

// Keys returns all keys as a slice
func (hm *Hashmap[K, V]) Keys() []K {
	res := make([]K, 0)

	for _, pairs := range hm.pairs.Iter() {
		if pairs == nil {
			continue
		}

		for _, p := range pairs.Iter() {
			res = append(res, p.key)
		}
	}

	return res
}

// Delete delete value with key K
func (hm *Hashmap[K, V]) Delete(key K) {
	hash := hm.mh.Hash(key)

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
	for _, pairs := range hm.pairs.Iter() {
		if pairs != nil {
			pairs.Free()
		}
	}
	hm.pairs.Free()
	allocator.Free(hm.alloc, hm)
}
