package interval_test

import (
	"testing"

	"github.com/gaissmai/interval"
)

func FuzzLCPandSCP(f *testing.F) {
	ivals := genUintIvals(100_000)
	tree := interval.NewTree(cmpUintInterval, ivals...)

	for i := 0; i < 100; i++ {
		probe := genUintIvals(1)[0]
		a := probe[0]
		b := probe[1]
		f.Add(a, b)
	}

	f.Fuzz(func(t *testing.T, a, b uint) {
		probe := makeUintIval(a, b)

		_, okLCP := tree.CoverLCP(probe)
		_, okSCP := tree.CoverSCP(probe)

		if okLCP != okSCP {
			// okLCP and okSCP must be both true or both false
			t.Fatalf("CoverLCP(%v) and CoverSCP(%v) mismatch", probe, probe)
		}
	})
}

func FuzzIntersects(f *testing.F) {
	ivals := genUintIvals(10_000)
	tree := interval.NewTree(cmpUintInterval, ivals...)

	for i := 0; i < 10; i++ {
		a := ivals[i][0]
		b := ivals[i][1]
		f.Add(a, b)
	}

	f.Fuzz(func(t *testing.T, a, b uint) {
		probe := makeUintIval(a, b)

		gotBool := tree.Intersects(probe)
		gotSlice := tree.Intersections(probe)

		if gotBool && gotSlice == nil || !gotBool && gotSlice != nil {
			t.Fatalf("Intersects(%v) and Intersections(%v) mismatch", probe, probe)
		}
	})
}
