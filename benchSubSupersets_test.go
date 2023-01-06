package interval_test

import (
	"testing"
)

func BenchmarkSubsetsIn1(b *testing.B) {
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = tree1.Subsets(probe)
	}
}

func BenchmarkSubsetsIn10(b *testing.B) {
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = tree10.Subsets(probe)
	}
}

func BenchmarkSubsetsIn100(b *testing.B) {
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = tree100.Subsets(probe)
	}
}

func BenchmarkSubsetsIn1_000(b *testing.B) {
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = tree1_000.Subsets(probe)
	}
}

func BenchmarkSubsetsIn10_000(b *testing.B) {
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = tree10_000.Subsets(probe)
	}
}

func BenchmarkSubsetsIn100_000(b *testing.B) {
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = tree100_000.Subsets(probe)
	}
}

func BenchmarkSubsetsIn1_000_000(b *testing.B) {
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = tree1_000_000.Subsets(probe)
	}
}

// #################################################

func BenchmarkSupersetsIn1(b *testing.B) {
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = tree1.Supersets(probe)
	}
}

func BenchmarkSupersetsIn10(b *testing.B) {
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = tree10.Supersets(probe)
	}
}

func BenchmarkSupersetsIn100(b *testing.B) {
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = tree100.Supersets(probe)
	}
}

func BenchmarkSupersetsIn1_000(b *testing.B) {
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = tree1_000.Supersets(probe)
	}
}

func BenchmarkSupersetsIn10_000(b *testing.B) {
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = tree10_000.Supersets(probe)
	}
}

func BenchmarkSupersetsIn100_000(b *testing.B) {
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = tree100_000.Supersets(probe)
	}
}

func BenchmarkSupersetsIn1_000_000(b *testing.B) {
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = tree1_000_000.Supersets(probe)
	}
}
