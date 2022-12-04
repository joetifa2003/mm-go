[![GoReportCard example](https://goreportcard.com/badge/github.com/joetifa2003/mm-go)](https://goreportcard.com/report/github.com/joetifa2003/mm-go)
[![GoDoc reference example](https://img.shields.io/badge/godoc-reference-blue.svg)](https://pkg.go.dev/github.com/joetifa2003/mm-go)

# mm-go Generic manual memory management for golang

Golang manages memory via GC and it's good for almost every use case but sometimes it can be a bottleneck.
and this is where mm-go comes in to play.

- [mm-go Generic manual memory management for golang](#mm-go-generic-manual-memory-management-for-golang)
  - [Before using mm-go](#before-using-mm-go)
  - [Installing](#installing)
  - [TypedArena (recommended)](#typedarena-recommended)
  - [Vector](#vector)
  - [Alloc/Free](#allocfree)
  - [AllocMany/FreeMany](#allocmanyfreemany)
  - [ReAlloc](#realloc)
  - [Benchmarks](#benchmarks)

## Before using mm-go

-   Golang doesn't have any way to manually allocate/free memory, so how does mm-go allocate/free?
    It does so via a custom allocator (see malloc.go) using direct system calls.
-   Before considering using this try to optimize your program to use less pointers, as golang GC most of the time performs worse when there is a lot of pointers, if you can't use this lib.
-   Manual memory management provides better performance (most of the time) but you are 100% responsible for managing it (bugs, segfaults, use after free, double free, ....)
-   Don't mix Manually and Managed memory (example if you put a slice in a manually managed struct it will get collected because go GC doesn't see the manually allocated struct, use Vector instead)

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
arena := mm.NewTypedArena[int](3)
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

## Vector

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

## Benchmarks

Check the test files and github actions for the benchmarks (linux, macos, windows).
mm-go can sometimes be 5-10 times faster.

```
Run benchstat out.txt
name                                time/op
HeapManaged/node_count_10000-2       510µs ± 1%
HeapManaged/node_count_100000-2     3.64ms ± 1%
HeapManaged/node_count_10000000-2    645ms ± 1%
HeapManaged/node_count_100000000-2   6.33s ± 3%
Manual/node_count_10000-2            223µs ± 1%
Manual/node_count_100000-2           580µs ± 1%
Manual/node_count_10000000-2        69.0ms ± 1%
Manual/node_count_100000000-2        684ms ± 0%
ArenaManual/node_count_10000-2       223µs ± 1%
ArenaManual/node_count_100000-2      560µs ± 0%
ArenaManual/node_count_10000000-2   67.6ms ± 1%
ArenaManual/node_count_100000000-2   666ms ± 0%
Slice-2                             9.26µs ± 1%
Vector-2                            10.2µs ± 3%
```
