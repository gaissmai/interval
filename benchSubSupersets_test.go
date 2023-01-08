package interval_test

import (
	"testing"
)

func BenchmarkSubsetsIn1(b *testing.B) {
	tree1 := mkTree(generateIvals(1))
	probe := generateIvals(1)[0]
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = tree1.Subsets(probe)
	}
}

func BenchmarkSubsetsIn10(b *testing.B) {
	tree10 := mkTree(generateIvals(10))
	probe := generateIvals(1)[0]
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = tree10.Subsets(probe)
	}
}

func BenchmarkSubsetsIn100(b *testing.B) {
	tree100 := mkTree(generateIvals(100))
	probe := generateIvals(1)[0]
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = tree100.Subsets(probe)
	}
}

func BenchmarkSubsetsIn1_000(b *testing.B) {
	tree1_000 := mkTree(generateIvals(1_000))
	probe := generateIvals(1)[0]
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = tree1_000.Subsets(probe)
	}
}

func BenchmarkSubsetsIn10_000(b *testing.B) {
	tree10_000 := mkTree(generateIvals(10_000))
	probe := generateIvals(1)[0]
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = tree10_000.Subsets(probe)
	}
}

func BenchmarkSubsetsIn100_000(b *testing.B) {
	tree100_000 := mkTree(generateIvals(100_000))
	probe := generateIvals(1)[0]
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = tree100_000.Subsets(probe)
	}
}

func BenchmarkSubsetsIn1_000_000(b *testing.B) {
	tree1_000_000 := mkTree(generateIvals(1_000_000))
	probe := generateIvals(1)[0]
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = tree1_000_000.Subsets(probe)
	}
}

// #################################################

func BenchmarkSupersetsIn1(b *testing.B) {
	tree1 := mkTree(generateIvals(1))
	probe := generateIvals(1)[0]
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = tree1.Supersets(probe)
	}
}

func BenchmarkSupersetsIn10(b *testing.B) {
	tree10 := mkTree(generateIvals(10))
	probe := generateIvals(1)[0]
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = tree10.Supersets(probe)
	}
}

func BenchmarkSupersetsIn100(b *testing.B) {
	tree100 := mkTree(generateIvals(100))
	probe := generateIvals(1)[0]
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = tree100.Supersets(probe)
	}
}

func BenchmarkSupersetsIn1_000(b *testing.B) {
	tree1_000 := mkTree(generateIvals(1_000))
	probe := generateIvals(1)[0]
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = tree1_000.Supersets(probe)
	}
}

func BenchmarkSupersetsIn10_000(b *testing.B) {
	tree10_000 := mkTree(generateIvals(10_000))
	probe := generateIvals(1)[0]
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = tree10_000.Supersets(probe)
	}
}

func BenchmarkSupersetsIn100_000(b *testing.B) {
	tree100_000 := mkTree(generateIvals(100_000))
	probe := generateIvals(1)[0]
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = tree100_000.Supersets(probe)
	}
}

func BenchmarkSupersetsIn1_000_000(b *testing.B) {
	tree1_000_000 := mkTree(generateIvals(1_000_000))
	probe := generateIvals(1)[0]
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = tree1_000_000.Supersets(probe)
	}
}
