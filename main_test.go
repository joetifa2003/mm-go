package main

import (
	"testing"
)

func BenchmarkHeapManaged(b *testing.B) {
	for i := b.N; i <= b.N; i++ {
		managed()
	}
}

func BenchmarkManual(b *testing.B) {
	for i := b.N; i <= b.N; i++ {
		unManaged()
	}
}
