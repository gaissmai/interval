package interval_test

import (
	"math/rand"
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

func BenchmarkInsertMutable(b *testing.B) {
	for n := 1; n <= 1_000_000; n *= 10 {
		tree := interval.NewTree(generateIvals(n)...)
		probe := generateIvals(1)[0]
		name := "Into" + intMap[n]

		b.Run(name, func(b *testing.B) {
			b.ResetTimer()
			for n := 0; n < b.N; n++ {
				(&tree).InsertMutable(probe)
			}
		})
	}
}

func BenchmarkDelete(b *testing.B) {
	for n := 1; n <= 1_000_000; n *= 10 {
		ivals := generateIvals(n)
		probe := ivals[rand.Intn(len(ivals))]

		tree := interval.NewTree(ivals...)
		name := "DeleteFrom" + intMap[n]

		b.Run(name, func(b *testing.B) {
			b.ResetTimer()
			for n := 0; n < b.N; n++ {
				_, _ = tree.Delete(probe)
			}
		})
	}
}

func BenchmarkDeleteMutable(b *testing.B) {
	for n := 1; n <= 1_000_000; n *= 10 {
		ivals := generateIvals(n)
		probe := ivals[rand.Intn(len(ivals))]

		tree := interval.NewTree(generateIvals(n)...)
		name := "DeleteFrom" + intMap[n]

		b.Run(name, func(b *testing.B) {
			b.ResetTimer()
			for n := 0; n < b.N; n++ {
				_ = (&tree).DeleteMutable(probe)
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

func BenchmarkUnionMutable(b *testing.B) {
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

func BenchmarkFind(b *testing.B) {
	for n := 1; n <= 1_000_000; n *= 10 {
		ivals := generateIvals(n)
		tree := interval.NewTree(ivals...)
		probe := ivals[rand.Intn(len(ivals))]
		name := "In" + intMap[n]

		b.Run(name, func(b *testing.B) {
			b.ResetTimer()
			for n := 0; n < b.N; n++ {
				_, _ = tree.Find(probe)
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
	for n := 100; n <= 100_000; n *= 10 {
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
	for n := 100; n <= 100_000; n *= 10 {
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

func BenchmarkIntersections(b *testing.B) {
	for n := 100; n <= 100_000; n *= 10 {
		tree := interval.NewTree(generateIvals(n)...)
		probe := generateIvals(1)[0]
		name := "In" + intMap[n]

		b.Run(name, func(b *testing.B) {
			b.ResetTimer()
			for n := 0; n < b.N; n++ {
				_ = tree.Intersections(probe)
			}
		})
	}
}

func BenchmarkSize(b *testing.B) {
	for n := 100; n <= 1_000_000; n *= 10 {
		ivals := generateIvals(n)
		tree := interval.NewTree(ivals...)
		name := "In" + intMap[n]

		b.Run(name, func(b *testing.B) {
			b.ResetTimer()
			for n := 0; n < b.N; n++ {
				_ = tree.Size()
			}
		})
	}
}

func BenchmarkMin(b *testing.B) {
	for n := 100; n <= 1_000_000; n *= 10 {
		ivals := generateIvals(n)
		tree := interval.NewTree(ivals...)
		name := "In" + intMap[n]

		b.Run(name, func(b *testing.B) {
			b.ResetTimer()
			for n := 0; n < b.N; n++ {
				_ = tree.Min()
			}
		})
	}
}

func BenchmarkMax(b *testing.B) {
	for n := 100; n <= 1_000_000; n *= 10 {
		ivals := generateIvals(n)
		tree := interval.NewTree(ivals...)
		name := "In" + intMap[n]

		b.Run(name, func(b *testing.B) {
			b.ResetTimer()
			for n := 0; n < b.N; n++ {
				_ = tree.Max()
			}
		})
	}
}
