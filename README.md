[![GoReportCard example](https://goreportcard.com/badge/github.com/joetifa2003/mm-go)](https://goreportcard.com/report/github.com/joetifa2003/mm-go)
[![GoDoc reference example](https://img.shields.io/badge/godoc-reference-blue.svg)](https://pkg.go.dev/github.com/joetifa2003/mm-go)

# mm-go Generic manual memory management for golang

<!--toc:start-->
- [mm-go Generic manual memory management for golang](#mm-go-generic-manual-memory-management-for-golang)
  - [Before using mm-go](#before-using-mm-go)
  - [Installing](#installing)
  - [Packages](#packages)
  - [Allocators](#allocators)
    - [C allocator](#c-allocator)
    - [BatchAllocator](#batchallocator)
    - [Generic Helpers](#generic-helpers)
      - [Alloc/Free](#allocfree)
      - [AllocMany/FreeMany](#allocmanyfreemany)
    - [ReAlloc](#realloc)
  - [typedarena](#typedarena)
    - [Why does this exists while there is BatchAllocator?](#why-does-this-exists-while-there-is-batchallocator)
  - [vector](#vector)
  - [Benchmarks](#benchmarks)
<!--toc:end-->

Golang manages memory via GC and it's good for almost every use case but sometimes it can be a bottleneck.
and this is where mm-go comes in to play.

## Before using mm-go

-   Golang doesn't have any way to manually allocate/free memory, so how does mm-go allocate/free?
    It does so via **cgo**.
-   Before considering using this try to optimize your program to use less pointers, as golang GC most of the time performs worse when there is a lot of pointers, if you can't use this lib.
-   Manual memory management provides better performance (most of the time) but you are **100% responsible** for managing it (bugs, segfaults, use after free, double free, ....)
-   **Don't mix** Manually and Managed memory (example if you put a slice in a manually managed struct it will get collected because go GC doesn't see the manually allocated struct, use Vector instead)
-   All data structures provided by the package are manually managed and thus can be safely included in manually managed structs without the GC freeing them, but **you have to free them yourself!**
-   Try to minimize calls to cgo by preallocating (using batchallocator/Arena/AllocMany).
-   Check the docs, test files and read the README.

## Installing

```
go get -u github.com/joetifa2003/mm-go
```

## Packages

`mm` - basic generic memory management functions.

`typedarena` - contains TypedArena which allocates many objects and free them all at once.

`vector` - contains a manually managed Vector implementation.

`linkedlist` - contains a manually managed Linkedlist implementation.

`mmstring` - contains a manually managed string implementation.

`malloc` - contains wrappers to raw C malloc and free.

`allocator` - contains the Allocator interface and the C allocator implementation.

`batchallocator` - contains implementation an allocator, this can be used as an arena, and a way to reduce CGO overhead.

## Allocators

Allocator is an interface that defines some methods needed for most allocators.

```go
Alloc(size int) unsafe.Pointer // returns a pointer to the allocated memory
Free(ptr unsafe.Pointer)      // frees the memory pointed by ptr
Realloc(ptr unsafe.Pointer, size int) unsafe.Pointer // reallocates the memory pointed by ptr
Destroy() // any cleanup that the allocator needs to do
```

Currently there is two allocators implemented.

### C allocator

The C allocator is using CGO underthehood to call calloc, realloc and free.

```go
alloc := allocator.NewC()
defer alloc.Destroy()

ptr := allocator.Alloc[int](alloc)
defer allocator.Free(alloc, ptr)

*ptr = 15
```

### BatchAllocator

This allocator purpose is to reduce the overhead of calling CGO on every allocation/free, it also acts as an arena since it frees all the memory when `Destroy` is called.

Instead, it allocats large chunks of memory at once and then divides them when you allocate, making it much faster.

This allocator has to take another allocator for it to work, usually with the C allocator.

```go
alloc := batchallocator.New(allocator.NewC())
defer alloc.Destroy()

ptr := allocator.Alloc[int](alloc)
defer allocator.Free(alloc, ptr)

*ptr = 15
```

With this allocator, calling `Free/FreeMany` on pointers allocated with `Alloc/AllocMany` is optional, since when you call `Destroy` all memory is freed by default.

But if you call `Free` the memory will be freed, so it acts as both Slab allocator and an Arena.

You can specify the size of chunks that are allocated by using options.

```go
alloc := batchallocator.New(allocator.NewC(),
    batchallocator.WithBucketSize(mm.SizeOf[int]()*15),
)
```

For example this configures the batch allocator to allocate at minimum 15 ints at a time (by default it allocates ` page, which is usually 4kb).

You can also allocate more than this configured amount in one big allocation, and it will work fine, unlike `typedarena`, more on that later.

### Generic Helpers

As you saw in the examples above, there are some helper functions that automatically detrimine the size of the type you want to allocate, and it also automatically does type casting from `unsafe.Pointer`.

So instead of doing this:

```go
alloc := batchallocator.New(allocator.NewC())
defer alloc.Destroy()

ptr := (*int)(alloc.Alloc(int(unsafe.Sizeof(int))))
defer alloc.Free(unsafe.Pointer(ptr))
```

You can do this:

```go
alloc := batchallocator.New(allocator.NewC())
defer alloc.Destroy()

ptr := allocator.Alloc[int](alloc)
defer allocator.Free(alloc, ptr)
```

Yes, go doesn't have generic on pointer receivers, so these had to be implemented as functions.

#### Alloc/Free

Alloc is a generic function that allocates T and returns a pointer to it that you can free later using Free

```go
alloc := batchallocator.New(allocator.NewC())
defer alloc.Destroy()

ptr := allocator.Alloc[int](alloc) // allocates a single int and returns a ptr to it
defer allocator.Free(alloc, ptr)     // frees the int (defer recommended to prevent leaks)

assert.Equal(0, *ptr) // allocations are zeroed by default
*ptr = 15             // changes the value using the pointer
assert.Equal(15, *ptr)
```

```go
type Node struct {
    value int
}

alloc := batchallocator.New(allocator.NewC())
ptr := allocator.Alloc[Node](alloc) // allocates a single Node struct and returns a ptr to it
defer allocator.Free(alloc, ptr)     // frees the struct (defer recommended to prevent leaks)
```

#### AllocMany/FreeMany

AllocMany is a generic function that allocates n of T and returns a slice that represents the heap (instead of pointer arithmetic => slice indexing) that you can free later using FreeMany

```go
alloc := allocator.NewC()
defer allocator.Destroy()

heap := allocator.AllocMany[int](alloc, 2) // allocates 2 ints and returns it as a slice of ints with length 2
defer allocator.FreeMany(heap)      // it's recommended to make sure the data gets deallocated (defer recommended to prevent leaks)

assert.Equal(2, len(heap))
heap[0] = 15    // changes the data in the slice (aka the heap)
ptr := &heap[0] // takes a pointer to the first int in the heap
// Be careful if you do ptr := heap[0] this will take a copy from the data on the heap
*ptr = 45            // changes the value from 15 to 45

assert.Equal(45, heap[0])
assert.Equal(0, heap[1])
```

WARNING: Do not append to the slice, this is only used to avoid pointer arithmetic and unsafe code.

### ReAlloc

Reallocate reallocates memory allocated with AllocMany and doesn't change underling data

```go
alloc := allocator.NewC()
defer alloc.Destroy()

heap := allocator.AllocMany[int](alloc, 2) // allocates 2 int and returns it as a slice of ints with length 2
heap[0] = 15
assert.Equal(2, len(heap))

heap = allocator.Realloc(allocated, 3)

assert.Equal(3, len(heap))
assert.Equal(15, heap[0]) // data after reallocation stays the same

allocator.FreeMany(heap)            // didn't use defer here because i'm doing a reallocation and changing the value of allocated variable (otherwise can segfault)
```


## typedarena

New creates a typed arena with the specified chunk size.
a chunk is the the unit of the arena, if T is int for example and the
chunk size is 5, then each chunk is going to hold 5 ints. And if the
chunk is filled it will allocate another chunk that can hold 5 ints.
then you can call FreeArena and it will deallocate all chunks together.
Using this will simplify memory management.

```go
alloc := allocator.NewC()
defer alloc.Destroy()

arena := typedarena.New[int](alloc, 3) // 3 is the chunk size which gets preallocated, if you allocated more than 3 it will preallocate another chunk of 3 T
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

### Why does this exists while there is BatchAllocator?

- `typedarena` is much faster because it only works with one single type. 
- `batchallocator` is more generic and works for any type, even multiple types all at once.

```go
// You cannot do this with `typedarena` because it only works with one single type.

alloc := batchallocator.New(allocator.NewC())
defer alloc.Destroy()

i := allocator.Alloc[int](alloc)
s := allocator.Alloc[string](alloc)
x := allocator.Alloc[float64](alloc)

// they are all freed automatically because of Destroy above
```

- `batchallocator` can be passed to multiple data structures, like `vector` and `hashmap` and it will be automatically Freed when `Destroy` is called.

- Also `typedarena.AllocMany` cannot exceed chunk size, but with `batchallocator` you can request any amount of memory.

```go
alloc := allocator.NewC()
defer alloc.Destroy()

arena := typedarena.New[int](alloc, 3)

heap := arena.AllocMany(2) // fine
heap2 := arena.AllocMany(2) // also fine
heap2 := arena.AllocMany(5) // panics 
```


## vector

A contiguous growable array type.
You can think of the Vector as a manually managed slice that you can put in manually managed structs, if you put a slice in a manually managed struct it will get collected because go GC doesn't see the manually allocated struct.

```go
v := vector.New[int]()
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
v := vector.New[int](5)
defer v.Free()

assert.Equal(5, v.Len())
assert.Equal(5, v.Cap())
```

```go
v := vector.New[int](5, 6)
defer v.Free()

assert.Equal(5, v.Len())
assert.Equal(6, v.Cap())
```

```go
v := vector.Init(1, 2, 3)
defer v.Free()

assert.Equal(3, v.Len())
assert.Equal(3, v.Cap())

assert.Equal(3, v.Pop())
assert.Equal(2, v.Pop())
assert.Equal(1, v.Pop())
```


## Benchmarks

Check the test files and github actions for the benchmarks (linux, macos, windows).
mm-go can sometimes be 5-10 times faster.

```
Run go test ./... -bench=. -count 5 > out.txt && benchstat out.txt

name                                time/op
pkg:github.com/joetifa2003/mm-go goos:linux goarch:amd64
HeapManaged/node_count_10000-2       504µs ± 1%
HeapManaged/node_count_100000-2     3.73ms ± 6%
HeapManaged/node_count_10000000-2    664ms ± 8%
HeapManaged/node_count_100000000-2   6.30s ± 4%
Manual/node_count_10000-2            226µs ± 1%
Manual/node_count_100000-2           576µs ± 1%
Manual/node_count_10000000-2        70.6ms ± 1%
Manual/node_count_100000000-2        702ms ± 1%
ArenaManual/node_count_10000-2       226µs ± 1%
ArenaManual/node_count_100000-2      553µs ± 0%
ArenaManual/node_count_10000000-2   69.1ms ± 0%
ArenaManual/node_count_100000000-2   681ms ± 1%
BinaryTreeManaged-2                  6.07s ±10%
BinaryTreeArena/chunk_size_50-2      2.30s ±21%
BinaryTreeArena/chunk_size_100-2     1.47s ± 5%
BinaryTreeArena/chunk_size_150-2     1.42s ±36%
BinaryTreeArena/chunk_size_250-2     1.11s ± 0%
BinaryTreeArena/chunk_size_500-2     1.00s ± 0%
```
