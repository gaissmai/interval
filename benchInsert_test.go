package interval_test

import (
	"testing"
)

func BenchmarkInsertInto1(b *testing.B) {
	tree1 := mkTree(generateIvals(1))
	probe := generateIvals(1)[0]
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = tree1.Insert(probe)
	}
}

func BenchmarkInsertInto10(b *testing.B) {
	tree10 := mkTree(generateIvals(10))
	probe := generateIvals(1)[0]
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = tree10.Insert(probe)
	}
}

func BenchmarkInsertInto100(b *testing.B) {
	tree100 := mkTree(generateIvals(100))
	probe := generateIvals(1)[0]
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = tree100.Insert(probe)
	}
}

func BenchmarkInsertInto1_000(b *testing.B) {
	tree1_000 := mkTree(generateIvals(1_000))
	probe := generateIvals(1)[0]
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = tree1_000.Insert(probe)
	}
}

func BenchmarkInsertInto10_000(b *testing.B) {
	tree10_000 := mkTree(generateIvals(10_000))
	probe := generateIvals(1)[0]
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = tree10_000.Insert(probe)
	}
}

func BenchmarkInsertInto100_000(b *testing.B) {
	tree100_000 := mkTree(generateIvals(100_000))
	probe := generateIvals(1)[0]
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = tree100_000.Insert(probe)
	}
}

func BenchmarkInsertInto1_000_000(b *testing.B) {
	tree1_000_000 := mkTree(generateIvals(1_000_000))
	probe := generateIvals(1)[0]
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = tree1_000_000.Insert(probe)
	}
}
