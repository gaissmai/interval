package interval_test

import (
	"testing"
)

func BenchmarkUnionImmutable100_000with10(b *testing.B) {
	tree10 := treap.Insert(generateIvals(10)...)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = tree100_000.Union(tree10, false, true)
	}
}

func BenchmarkUnionImmutable100_000with100(b *testing.B) {
	tree100 := treap.Insert(generateIvals(100)...)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = tree100_000.Union(tree100, false, true)
	}
}

func BenchmarkUnionImmutable100_000with1_000(b *testing.B) {
	tree1_000 := treap.Insert(generateIvals(1_000)...)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = tree100_000.Union(tree1_000, false, true)
	}
}

func BenchmarkUnionImmutable100_000with10_000(b *testing.B) {
	tree10_000 := treap.Insert(generateIvals(10_000)...)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = tree100_000.Union(tree10_000, false, true)
	}
}

func BenchmarkUnionImmutable100_000with100_000(b *testing.B) {
	this100_000 := treap.Insert(generateIvals(100_000)...)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = tree100_000.Union(this100_000, false, true)
	}
}

func BenchmarkUnion100_000with10(b *testing.B) {
	tree10 := treap.Insert(generateIvals(10)...)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = tree100_000.Union(tree10, false, false)
	}
}

func BenchmarkUnion100_000with100(b *testing.B) {
	tree100 := treap.Insert(generateIvals(100)...)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = tree100_000.Union(tree100, false, false)
	}
}

func BenchmarkUnion100_000with1_000(b *testing.B) {
	tree1_000 := treap.Insert(generateIvals(1_000)...)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = tree100_000.Union(tree1_000, false, false)
	}
}

func BenchmarkUnion100_000with10_000(b *testing.B) {
	tree10_000 := treap.Insert(generateIvals(10_000)...)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = tree100_000.Union(tree10_000, false, false)
	}
}

func BenchmarkUnion100_000with100_000(b *testing.B) {
	this100_000 := treap.Insert(generateIvals(100_000)...)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = tree100_000.Union(this100_000, false, false)
	}
}
