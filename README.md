[![GoReportCard example](https://goreportcard.com/badge/github.com/joetifa2003/mm-go)](https://goreportcard.com/report/github.com/joetifa2003/mm-go)
[![GoDoc reference example](https://img.shields.io/badge/godoc-reference-blue.svg)](https://pkg.go.dev/github.com/joetifa2003/mm-go)

# mm-go Generic manual memory management for golang

Golang manages memory via GC and it's good for almost every use case but sometimes it can be a bottleneck.
and this is where mm-go comes in to play.

- [mm-go Generic manual memory management for golang](#mm-go-generic-manual-memory-management-for-golang)
  - [Before using mm-go](#before-using-mm-go)
  - [Installing](#installing)
  - [Alloc/Free](#allocfree)
  - [AllocMany/FreeMany](#allocmanyfreemany)
  - [ReAlloc](#realloc)
  - [Vector](#vector)
  - [Benchmarks](#benchmarks)

## Before using mm-go

-   Golang doesn't have any way to manually allocate/free memory, so how does mm-go allocate/free?
    It does so via cgo.
-   Before considering using this try to optimize your program to use less pointers, as golang GC most of the time performs worse when there is a lot of pointers, if you can't use this lib.
-   Manual memory management provides better performance (most of the time) but you are 100% responsible for managing it (bugs, segfaults, use after free, double free, ....)
-   Don't mix Manually and Managed memory (example if you put a slice in a manually managed struct it will get collected because go GC doesn't see the manually allocated struct, use Vector instead)
-   Lastly mm-go uses cgo for calloc/free and it's known that calling cgo has some overhead so try to minimize the calls to cgo (in hot loops for example)

## Installing

```
go get github.com/joetifa2003/mm-go
```

## Alloc/Free

Alloc is a generic function that allocates T and returns a pointer to it that u can free later using Free

```go
ptr := Alloc[int]() // allocates a single int and returns a ptr to it
defer Free(ptr)     // frees the int (defer recommended to prevent leaks)

assert.Equal(0, *ptr) // allocations are zeroed by default
*ptr = 15             // changes the value using the pointer
assert.Equal(15, *ptr)
```

```go
type Node struct {
    value int
}

ptr := Alloc[Node]() // allocates a single Node struct and returns a ptr to it
defer Free(ptr)     // frees the struct (defer recommended to prevent leaks)
```

## AllocMany/FreeMany

AllocMany is a generic function that allocates n of T and returns a slice that represents the heap (instead of pointer arithmetic => slice indexing) that u can free later using FreeMany

```go
allocated := AllocMany[int](2) // allocates 2 ints and returns it as a slice of ints with length 2
defer FreeMany(allocated)      // it's recommended to make sure the data gets deallocated (defer recommended to prevent leaks)
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
allocated := AllocMany[int](2) // allocates 2 int and returns it as a slice of ints with length 2
allocated[0] = 15
assert.Equal(2, len(allocated))
allocated = Reallocate(allocated, 3)
assert.Equal(3, len(allocated))
assert.Equal(15, allocated[0]) // data after reallocation stays the same
FreeMany(allocated)            // didn't use defer here because i'm doing a reallocation and changing the value of allocated variable (otherwise can segfault)
```

## Vector

You can think of the Vector as a manually managed slice that you can put in manually managed structs, if you put a slice in a manually managed struct it will get collected because go GC doesn't see the manually allocated struct, use Vector instead

```go
v := NewVector[int]()
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
v := NewVector[int](5)
defer v.Free()

assert.Equal(5, v.Len())
assert.Equal(5, v.Cap())
```

```go
v := NewVector[int](5, 6)
defer v.Free()

assert.Equal(5, v.Len())
assert.Equal(6, v.Cap())
```

```go
v := InitVector(1, 2, 3)
defer v.Free()

assert.Equal(3, v.Len())
assert.Equal(3, v.Cap())

assert.Equal(3, v.Pop())
assert.Equal(2, v.Pop())
assert.Equal(1, v.Pop())
```

## Benchmarks

Check the test files and github actions for the benchmarks (linux, macos, windows).
mm-go can sometimes be 5-10 times faster, if you are not careful it will be slower!

```
goos: linux
goarch: amd64
pkg: github.com/joetifa2003/mm-go
cpu: Intel(R) Xeon(R) CPU E5-2673 v4 @ 2.30GHz
BenchmarkHeapManaged-2   	1000000000	         0.09938 ns/op
BenchmarkHeapManaged-2   	1000000000	         0.09482 ns/op
BenchmarkHeapManaged-2   	1000000000	         0.09476 ns/op
BenchmarkHeapManaged-2   	1000000000	         0.09554 ns/op
BenchmarkHeapManaged-2   	1000000000	         0.08969 ns/op
BenchmarkArenaManual-2   	1000000000	         0.005030 ns/op
BenchmarkArenaManual-2   	1000000000	         0.008007 ns/op
BenchmarkArenaManual-2   	1000000000	         0.008574 ns/op
BenchmarkArenaManual-2   	1000000000	         0.008318 ns/op
BenchmarkArenaManual-2   	1000000000	         0.008397 ns/op
BenchmarkSlice-2         	   96015	     12658 ns/op
BenchmarkSlice-2         	  109916	     11724 ns/op
BenchmarkSlice-2         	  112123	     11606 ns/op
BenchmarkSlice-2         	  111018	     10988 ns/op
BenchmarkSlice-2         	   99549	     11017 ns/op
BenchmarkVector-2        	   50136	     23062 ns/op
BenchmarkVector-2        	   69448	     22905 ns/op
BenchmarkVector-2        	   75655	     16998 ns/op
BenchmarkVector-2        	   73281	     18884 ns/op
BenchmarkVector-2        	   77398	     17642 ns/op
```
