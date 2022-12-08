[![GoReportCard example](https://goreportcard.com/badge/github.com/joetifa2003/mm-go)](https://goreportcard.com/report/github.com/joetifa2003/mm-go)
[![GoDoc reference example](https://img.shields.io/badge/godoc-reference-blue.svg)](https://pkg.go.dev/github.com/joetifa2003/mm-go)

# mm-go Generic manual memory management for golang

Golang manages memory via GC and it's good for almost every use case but sometimes it can be a bottleneck.
and this is where mm-go comes in to play.

- [mm-go Generic manual memory management for golang](#mm-go-generic-manual-memory-management-for-golang)
  - [Before using mm-go](#before-using-mm-go)
  - [Installing](#installing)
  - [TypedArena (recommended)](#typedarena-recommended)
  - [Alloc/Free](#allocfree)
  - [AllocMany/FreeMany](#allocmanyfreemany)
  - [ReAlloc](#realloc)
  - [Vector](#vector)
    - [Methods](#methods)
      - [NewVector](#newvector)
      - [InitVector](#initvector)
      - [InitVector](#initvector-1)
      - [Push](#push)
      - [Pop](#pop)
      - [Len](#len)
      - [Cap](#cap)
      - [Slice](#slice)
      - [Last](#last)
      - [At](#at)
      - [Free](#free)
  - [Linked List](#linked-list)
    - [Methods](#methods-1)
      - [NewLinkedList](#newlinkedlist)
      - [PushBack](#pushback)
      - [PushFront](#pushfront)
      - [PopBack](#popback)
      - [PopFront](#popfront)
      - [ForEach](#foreach)
      - [At](#at-1)
      - [RemoveAt](#removeat)
      - [Remove](#remove)
      - [RemoveAll](#removeall)
      - [FindIndex](#findindex)
      - [FindIndexes](#findindexes)
      - [Len](#len-1)
      - [Free](#free-1)
  - [Benchmarks](#benchmarks)

## Before using mm-go

-   Golang doesn't have any way to manually allocate/free memory, so how does mm-go allocate/free?
    It does so via cgo.
-   Before considering using this try to optimize your program to use less pointers, as golang GC most of the time performs worse when there is a lot of pointers, if you can't use this lib.
-   Manual memory management provides better performance (most of the time) but you are 100% responsible for managing it (bugs, segfaults, use after free, double free, ....)
-   Don't mix Manually and Managed memory (example if you put a slice in a manually managed struct it will get collected because go GC doesn't see the manually allocated struct, use Vector instead)
-   All data structures provided by the package are manually managed and thus can be safely included in manually managed structs without the GC freeing them, but you have to free them yourself!
-   Try to minimize calls to cgo by preallocating (using Arena/AllocMany).
-   Check the docs, test files and read the README.

## Installing

```
go get github.com/joetifa2003/mm-go
```

## TypedArena (recommended)

NewTypedArena creates a typed arena with the specified chunk size.
a chunk is the the unit of the arena, if T is int for example and the
chunk size is 5, then each chunk is going to hold 5 ints. And if the
chunk is filled it will allocate another chunk that can hold 5 ints.
then you can call FreeArena and it will deallocate all chunks together.
Using this will simplify memory management.

```go
arena := mm.NewTypedArena[int](3) // 3 is the chunk size which gets preallocated, if you allocated more than 3 it will preallocate another chunk of 3 T
defer arena.Free() // freeing the arena using defer to prevent leaks

int1 := arena.Alloc()      // allocates 1 int from arena
*int1 = 1                  // changing it's value
ints := arena.AllocMany(2) // allocates 2 ints from the arena and returns a slice representing the heap (instead of pointer arithmetic)
ints[0] = 2                // changing the first value
ints[1] = 3                // changing the second value

// you can also take pointers from the slice
intPtr1 := &ints[0] // taking pointer from the manually managed heap
*intPtr1 = 15 // changing the value using pointers

assert.Equal(1, *int1)
assert.Equal(2, len(ints))
assert.Equal(15, ints[0])
assert.Equal(3, ints[1])
```

## Alloc/Free

Alloc is a generic function that allocates T and returns a pointer to it that you can free later using Free

```go
ptr := mm.Alloc[int]() // allocates a single int and returns a ptr to it
defer mm.Free(ptr)     // frees the int (defer recommended to prevent leaks)

assert.Equal(0, *ptr) // allocations are zeroed by default
*ptr = 15             // changes the value using the pointer
assert.Equal(15, *ptr)
```

```go
type Node struct {
    value int
}

ptr := mm.Alloc[Node]() // allocates a single Node struct and returns a ptr to it
defer mm.Free(ptr)     // frees the struct (defer recommended to prevent leaks)
```

## AllocMany/FreeMany

AllocMany is a generic function that allocates n of T and returns a slice that represents the heap (instead of pointer arithmetic => slice indexing) that you can free later using FreeMany

```go
allocated := mm.AllocMany[int](2) // allocates 2 ints and returns it as a slice of ints with length 2
defer mm.FreeMany(allocated)      // it's recommended to make sure the data gets deallocated (defer recommended to prevent leaks)
assert.Equal(2, len(allocated))
allocated[0] = 15    // changes the data in the slice (aka the heap)
ptr := &allocated[0] // takes a pointer to the first int in the heap
// Be careful if you do ptr := allocated[0] this will take a copy from the data on the heap
*ptr = 45            // changes the value from 15 to 45

assert.Equal(45, allocated[0])
```

## ReAlloc

Reallocate reallocates memory allocated with AllocMany and doesn't change underling data

```go
allocated := mm.AllocMany[int](2) // allocates 2 int and returns it as a slice of ints with length 2
allocated[0] = 15
assert.Equal(2, len(allocated))
allocated = mm.Reallocate(allocated, 3)
assert.Equal(3, len(allocated))
assert.Equal(15, allocated[0]) // data after reallocation stays the same
mm.FreeMany(allocated)            // didn't use defer here because i'm doing a reallocation and changing the value of allocated variable (otherwise can segfault)
```

## Vector

A contiguous growable array type.
You can think of the Vector as a manually managed slice that you can put in manually managed structs, if you put a slice in a manually managed struct it will get collected because go GC doesn't see the manually allocated struct.

```go
v := mm.NewVector[int]()
defer v.Free()

v.Push(1)
v.Push(2)
v.Push(3)

assert.Equal(3, v.Len())
assert.Equal(4, v.Cap())
assert.Equal([]int{1, 2, 3}, v.Slice())
assert.Equal(3, v.Pop())
assert.Equal(2, v.Pop())
assert.Equal(1, v.Pop())
```

```go
v := mm.NewVector[int](5)
defer v.Free()

assert.Equal(5, v.Len())
assert.Equal(5, v.Cap())
```

```go
v := mm.NewVector[int](5, 6)
defer v.Free()

assert.Equal(5, v.Len())
assert.Equal(6, v.Cap())
```

```go
v := mm.InitVector(1, 2, 3)
defer v.Free()

assert.Equal(3, v.Len())
assert.Equal(3, v.Cap())

assert.Equal(3, v.Pop())
assert.Equal(2, v.Pop())
assert.Equal(1, v.Pop())
```

### Methods

#### NewVector

```go
// NewVector creates a new empty vector, if args not provided
// it will create an empty vector, if only one arg is provided
// it will init a vector with len and cap equal to the provided arg,
// if two args are provided it will init a vector with len = args[0] cap = args[1]
func NewVector[T any](args ...int) *Vector[T]
```

#### InitVector

```go
// InitVector initializes a new vector with the T elements provided and sets
// it's len and cap to len(values)
func InitVector[T any](values ...T) *Vector[T]
```

#### InitVector

```go
// InitVector initializes a new vector with the T elements provided and sets
// it's len and cap to len(values)
func InitVector[T any](values ...T) *Vector[T]
```

#### Push

```go
// Push pushes value T to the vector, grows if needed.
func (v *Vector[T]) Push(value T)
```

#### Pop

```go
// Pop pops value T from the vector and returns it
func (v *Vector[T]) Pop() T
```

#### Len

```go
// Len gets vector length
func (v *Vector[T]) Len() int
```

#### Cap

```go
// Cap gets vector capacity (underling memory length).
func (v *Vector[T]) Cap() int
```

#### Slice

```go
// Slice gets a slice representing the vector
// CAUTION: don't append to this slice, this is only used
// if you want to loop on the vec elements
func (v *Vector[T]) Slice() []T
```

#### Last

```go
// Last gets the last element from a vector
func (v *Vector[T]) Last() T
```

#### At

```go
// At gets element T at specified index
func (v *Vector[T]) At(idx int) T
```

#### Free

```go
// Free deallocats the vector
func (v *Vector[T]) Free()
```

## Linked List

LinkedList a doubly-linked list.
Note: can be a lot slower than Vector but sometimes faster in specific use cases

### Methods

#### NewLinkedList

```go
// NewLinkedList creates a new linked list.
func NewLinkedList[T any]() *LinkedList[T]
```

#### PushBack

```go
// PushBack pushes value T to the back of the linked list.
func (ll *LinkedList[T]) PushBack(value T)
```

#### PushFront

```go
// PushFront pushes value T to the back of the linked list.
func (ll *LinkedList[T]) PushFront(value T)
```

#### PopBack

```go
// PopBack pops and returns value T from the back of the linked list.
func (ll *LinkedList[T]) PopBack() T
```

#### PopFront

```go
// PopFront pops and returns value T from the front of the linked list.
func (ll *LinkedList[T]) PopFront() T
```

#### ForEach

```go
// ForEach iterates through the linked list.
func (ll *LinkedList[T]) ForEach(f func(idx int, value T))
```

#### At

```go
// At gets value T at idx.
func (ll *LinkedList[T]) At(idx int) T
```

#### RemoveAt

```go
// RemoveAt removes value T at specified index and returns it.
func (ll *LinkedList[T]) RemoveAt(idx int) T
```

#### Remove

```go
// Remove removes the first value T that pass the test implemented by the provided function.
// if the test function succeeded it will return the value and true
func (ll *LinkedList[T]) Remove(f func(idx int, value T) bool) (value T, ok bool)
```

#### RemoveAll

```go
// RemoveAll removes all values of T that pass the test implemented by the provided function.
func (ll *LinkedList[T]) RemoveAll(f func(idx int, value T) bool) []T
```

#### FindIndex

```go
// FindIndex returns the first index of value T that pass the test implemented by the provided function.
func (ll *LinkedList[T]) FindIndex(f func(value T) bool) (idx int, ok bool)
```

#### FindIndexes

```go
// FindIndex returns all indexes of value T that pass the test implemented by the provided function.
func (ll *LinkedList[T]) FindIndexes(f func(value T) bool) []int
```

#### Len

```go
// Len gets linked list length.
func (ll *LinkedList[T]) Len() int
```

#### Free

```go
// Free frees the linked list.
func (ll *LinkedList[T]) Free()
```

## Benchmarks

Check the test files and github actions for the benchmarks (linux, macos, windows).
mm-go can sometimes be 5-10 times faster.

```
name                                time/op

Slice-2                             4.17µs ± 3%
Vector-2                            4.70µs ± 1%
LinkedList-2                         408µs ± 2%
HeapManaged/node_count_10000-2       863µs ± 3%
HeapManaged/node_count_100000-2     4.85ms ± 1%
HeapManaged/node_count_10000000-2    881ms ±14%
HeapManaged/node_count_100000000-2   8.54s ±12%
Manual/node_count_10000-2            693µs ± 6%
Manual/node_count_100000-2          2.49ms ± 4%
Manual/node_count_10000000-2         196ms ± 4%
Manual/node_count_100000000-2        2.05s ± 7%
ArenaManual/node_count_10000-2       682µs ± 4%
ArenaManual/node_count_100000-2     2.55ms ± 6%
ArenaManual/node_count_10000000-2    197ms ± 3%
ArenaManual/node_count_100000000-2   2.06s ± 9%
```
