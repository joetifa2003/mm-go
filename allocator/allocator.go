package allocator

import (
	"context"
	"unsafe"
)

type ctxKey int

const (
	AllocatorKey ctxKey = iota
)

type Allocator interface {
	Alloc(size int) unsafe.Pointer
	Realloc(ptr unsafe.Pointer, size int) unsafe.Pointer
	Free(ptr unsafe.Pointer)
	Destroy()
}

func WithAllocator(ctx context.Context, allocator Allocator) context.Context {
	return context.WithValue(ctx, AllocatorKey, allocator)
}

func GetAllocator(ctx context.Context) Allocator {
	allocator, ok := ctx.Value(AllocatorKey).(Allocator)
	if !ok {
		return CAllocator{}
	}

	return allocator
}
