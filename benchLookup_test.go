package interval_test

import (
	"testing"
)

func BenchmarkShortestIn1(b *testing.B) {
	tree1 := mkTree(generateIvals(1))
	probe := generateIvals(1)[0]
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, _ = tree1.Shortest(probe)
	}
}

func BenchmarkShortestIn10(b *testing.B) {
	tree10 := mkTree(generateIvals(10))
	probe := generateIvals(1)[0]
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, _ = tree10.Shortest(probe)
	}
}

func BenchmarkShortestIn100(b *testing.B) {
	tree100 := mkTree(generateIvals(100))
	probe := generateIvals(1)[0]
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, _ = tree100.Shortest(probe)
	}
}

func BenchmarkShortestIn1_000(b *testing.B) {
	tree1_000 := mkTree(generateIvals(1_000))
	probe := generateIvals(1)[0]
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, _ = tree1_000.Shortest(probe)
	}
}

func BenchmarkShortestIn10_000(b *testing.B) {
	tree10_000 := mkTree(generateIvals(10_000))
	probe := generateIvals(1)[0]
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, _ = tree10_000.Shortest(probe)
	}
}

func BenchmarkShortestIn100_000(b *testing.B) {
	tree100_000 := mkTree(generateIvals(100_000))
	probe := generateIvals(1)[0]
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, _ = tree100_000.Shortest(probe)
	}
}

func BenchmarkShortestIn1_000_000(b *testing.B) {
	tree1_000_000 := mkTree(generateIvals(1_000_000))
	probe := generateIvals(1)[0]
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, _ = tree1_000_000.Shortest(probe)
	}
}

// #################################################

func BenchmarkLargestIn1(b *testing.B) {
	tree1 := mkTree(generateIvals(1))
	probe := generateIvals(1)[0]
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, _ = tree1.Largest(probe)
	}
}

func BenchmarkLargestIn10(b *testing.B) {
	tree10 := mkTree(generateIvals(10))
	probe := generateIvals(1)[0]
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, _ = tree10.Largest(probe)
	}
}

func BenchmarkLargestIn100(b *testing.B) {
	tree100 := mkTree(generateIvals(100))
	probe := generateIvals(1)[0]
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, _ = tree100.Largest(probe)
	}
}

func BenchmarkLargestIn1_000(b *testing.B) {
	tree1_000 := mkTree(generateIvals(1_000))
	probe := generateIvals(1)[0]
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, _ = tree1_000.Largest(probe)
	}
}

func BenchmarkLargestIn10_000(b *testing.B) {
	tree10_000 := mkTree(generateIvals(10_000))
	probe := generateIvals(1)[0]
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, _ = tree10_000.Largest(probe)
	}
}

func BenchmarkLargestIn100_000(b *testing.B) {
	tree100_000 := mkTree(generateIvals(100_000))
	probe := generateIvals(1)[0]
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, _ = tree100_000.Largest(probe)
	}
}

func BenchmarkLargestIn1_000_000(b *testing.B) {
	tree1_000_000 := mkTree(generateIvals(1_000_000))
	probe := generateIvals(1)[0]
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, _ = tree1_000_000.Largest(probe)
	}
}
