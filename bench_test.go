package interval_test

import (
	"testing"

	"github.com/gaissmai/interval"
)

var intMap = map[int]string{
	1:         "1",
	10:        "10",
	100:       "100",
	1_000:     "1_000",
	10_000:    "10_000",
	100_000:   "100_000",
	1_000_000: "1_000_000",
}

func BenchmarkInsert(b *testing.B) {
	for n := 1; n <= 1_000_000; n *= 10 {
		tree := interval.NewTree(generateIvals(n)...)
		probe := generateIvals(1)[0]
		name := "Into" + intMap[n]

		b.Run(name, func(b *testing.B) {
			b.ResetTimer()
			for n := 0; n < b.N; n++ {
				_ = tree.Insert(probe)
			}
		})
	}
}

func BenchmarkDelete(b *testing.B) {
	for n := 10; n <= 1_000_000; n *= 10 {
		tree := interval.NewTree(generateIvals(n)...)
		probe := generateIvals(1)[0]
		name := "DeleteFrom" + intMap[n]

		b.Run(name, func(b *testing.B) {
			b.ResetTimer()
			for n := 0; n < b.N; n++ {
				_, _ = tree.Delete(probe)
			}
		})
	}
}

func BenchmarkUnionImmutable(b *testing.B) {
	this100_000 := interval.NewTree(generateIvals(100_000)...)
	for n := 10; n <= 100_000; n *= 10 {
		tree := interval.NewTree(generateIvals(n)...)
		name := "size100_000with" + intMap[n]

		b.Run(name, func(b *testing.B) {
			b.ResetTimer()
			for n := 0; n < b.N; n++ {
				_ = this100_000.Union(tree, false, true)
			}
		})
	}
}

func BenchmarkUnionNonImmutable(b *testing.B) {
	for n := 10; n <= 100_000; n *= 10 {
		this100_000 := interval.NewTree(generateIvals(100_000)...)
		tree := interval.NewTree(generateIvals(n)...)
		name := "size100_000with" + intMap[n]

		b.Run(name, func(b *testing.B) {
			b.ResetTimer()
			for n := 0; n < b.N; n++ {
				_ = this100_000.Union(tree, false, false)
			}
		})
	}
}

func BenchmarkShortest(b *testing.B) {
	for n := 100; n <= 1_000_000; n *= 10 {
		tree := interval.NewTree(generateIvals(n)...)
		probe := generateIvals(1)[0]
		name := "In" + intMap[n]

		b.Run(name, func(b *testing.B) {
			b.ResetTimer()
			for n := 0; n < b.N; n++ {
				_, _ = tree.Shortest(probe)
			}
		})
	}
}

func BenchmarkLargest(b *testing.B) {
	for n := 100; n <= 1_000_000; n *= 10 {
		tree := interval.NewTree(generateIvals(n)...)
		probe := generateIvals(1)[0]
		name := "In" + intMap[n]

		b.Run(name, func(b *testing.B) {
			b.ResetTimer()
			for n := 0; n < b.N; n++ {
				_, _ = tree.Largest(probe)
			}
		})
	}
}

func BenchmarkSubsets(b *testing.B) {
	for n := 100; n <= 1_000_000; n *= 10 {
		tree := interval.NewTree(generateIvals(n)...)
		probe := generateIvals(1)[0]
		name := "In" + intMap[n]

		b.Run(name, func(b *testing.B) {
			b.ResetTimer()
			for n := 0; n < b.N; n++ {
				_ = tree.Subsets(probe)
			}
		})
	}
}

func BenchmarkSupersets(b *testing.B) {
	for n := 100; n <= 1_000_000; n *= 10 {
		tree := interval.NewTree(generateIvals(n)...)
		probe := generateIvals(1)[0]
		name := "In" + intMap[n]

		b.Run(name, func(b *testing.B) {
			b.ResetTimer()
			for n := 0; n < b.N; n++ {
				_ = tree.Supersets(probe)
			}
		})
	}
}