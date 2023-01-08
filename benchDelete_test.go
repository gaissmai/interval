package interval_test

import (
	"testing"
)

func BenchmarkDeleteFrom1(b *testing.B) {
	tree1 := mkTree(generateIvals(1))
	probe := generateIvals(1)[0]
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, _ = tree1.Delete(probe)
	}
}

func BenchmarkDeleteFrom10(b *testing.B) {
	tree10 := mkTree(generateIvals(10))
	probe := generateIvals(1)[0]
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, _ = tree10.Delete(probe)
	}
}

func BenchmarkDeleteFrom100(b *testing.B) {
	tree100 := mkTree(generateIvals(100))
	probe := generateIvals(1)[0]
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, _ = tree100.Delete(probe)
	}
}

func BenchmarkDeleteFrom1_000(b *testing.B) {
	tree1_000 := mkTree(generateIvals(1_000))
	probe := generateIvals(1)[0]
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, _ = tree1_000.Delete(probe)
	}
}

func BenchmarkDeleteFrom10_000(b *testing.B) {
	tree10_000 := mkTree(generateIvals(10_000))
	probe := generateIvals(1)[0]
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, _ = tree10_000.Delete(probe)
	}
}

func BenchmarkDeleteFrom100_000(b *testing.B) {
	tree100_000 := mkTree(generateIvals(100_000))
	probe := generateIvals(1)[0]
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, _ = tree100_000.Delete(probe)
	}
}

func BenchmarkDeleteFrom1_000_000(b *testing.B) {
	tree1_000_000 := mkTree(generateIvals(1_000_000))
	probe := generateIvals(1)[0]
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, _ = tree1_000_000.Delete(probe)
	}
}
