package interval_test

import (
	"testing"
)

func BenchmarkDeleteFrom1(b *testing.B) {
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, _ = tree1.Delete(probe)
	}
}

func BenchmarkDeleteFrom10(b *testing.B) {
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, _ = tree10.Delete(probe)
	}
}

func BenchmarkDeleteFrom100(b *testing.B) {
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, _ = tree100.Delete(probe)
	}
}

func BenchmarkDeleteFrom1_000(b *testing.B) {
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, _ = tree1_000.Delete(probe)
	}
}

func BenchmarkDeleteFrom10_000(b *testing.B) {
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, _ = tree10_000.Delete(probe)
	}
}

func BenchmarkDeleteFrom100_000(b *testing.B) {
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, _ = tree100_000.Delete(probe)
	}
}

func BenchmarkDeleteFrom1_000_000(b *testing.B) {
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, _ = tree1_000_000.Delete(probe)
	}
}
