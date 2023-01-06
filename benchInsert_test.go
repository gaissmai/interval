package interval_test

import (
	"testing"
)

func BenchmarkInsertInto1(b *testing.B) {
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = tree1.Insert(probe)
	}
}

func BenchmarkInsertInto10(b *testing.B) {
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = tree10.Insert(probe)
	}
}

func BenchmarkInsertInto100(b *testing.B) {
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = tree100.Insert(probe)
	}
}

func BenchmarkInsertInto1_000(b *testing.B) {
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = tree1_000.Insert(probe)
	}
}

func BenchmarkInsertInto10_000(b *testing.B) {
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = tree10_000.Insert(probe)
	}
}

func BenchmarkInsertInto100_000(b *testing.B) {
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = tree100_000.Insert(probe)
	}
}

func BenchmarkInsertInto1_000_000(b *testing.B) {
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = tree1_000_000.Insert(probe)
	}
}
