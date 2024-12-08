package run_8

import (
	"testing"
	"unsafe"
)

func BenchmarkString(b *testing.B) {
	data := []byte("Hallo, Welt!")

	for i := 0; i < b.N; i++ {
		_ = string(data)
	}
}

func BenchmarkUnsafeString(b *testing.B) {
	data := []byte("Hallo, Welt!")

	for i := 0; i < b.N; i++ {
		_ = unsafe.String(unsafe.SliceData(data), len(data))
	}
}
