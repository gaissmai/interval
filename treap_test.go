package interval_test

import (
	"fmt"
	"math"
	"math/rand"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/gaissmai/interval"
	"github.com/gaissmai/interval/internal/period"
)

var treap *interval.Tree[period.Ival]

// test data
var ps = []period.Ival{
	{0, 6},
	{0, 5},
	{1, 8},
	{1, 7},
	{1, 5},
	{1, 4},
	{2, 8},
	{2, 7},
	{4, 8},
	{6, 7},
	{7, 9},
}

// random test data
func generateIvals(n int) []period.Ival {
	is := make([]period.Ival, n)
	for i := 0; i < n; i++ {
		a := rand.Intn(1000 * n)
		b := rand.Intn(1000 * n)
		if a > b {
			a, b = b, a
		}
		is[i] = period.Ival{a, b}
	}
	return is
}

func TestTreeZeroValue(t *testing.T) {
	t.Parallel()

	var zeroItem period.Ival
	var zeroTree *interval.Tree[period.Ival]

	w := new(strings.Builder)
	_ = zeroTree.Fprint(w)

	if w.String() != "" {
		t.Errorf("Write(w) = %v, want \"\"", w.String())
	}

	w.Reset()
	_ = zeroTree.FprintBST(w)

	if w.String() != "" {
		t.Errorf("Write(w) = %v, want \"\"", w.String())
	}

	if s := zeroTree.Insert(zeroItem); s == nil {
		t.Errorf("Insert(), got: %v, want: !nil", s)
	}

	if _, ok := zeroTree.Delete(zeroItem); ok {
		t.Errorf("Delete(), got: %v, want: false", ok)
	}

	if s := zeroTree.Clone(); s != nil {
		t.Errorf("Clone(), got: %v, want: nil", s)
	}

	if _, ok := zeroTree.Shortest(zeroItem); ok {
		t.Errorf("Shortest(), got: %v, want: false", ok)
	}

	if _, ok := zeroTree.Largest(zeroItem); ok {
		t.Errorf("Largest(), got: %v, want: false", ok)
	}

	if s := zeroTree.Subsets(zeroItem); s != nil {
		t.Errorf("Subsets(), got: %v, want: nil", s)
	}

	if s := zeroTree.Supersets(zeroItem); s != nil {
		t.Errorf("Supersets(), got: %v, want: nil", s)
	}

	if s := zeroTree.Min(); s != zeroItem {
		t.Errorf("Min(), got: %v, want: %v", s, zeroItem)
	}

	if s := zeroTree.Max(); s != zeroItem {
		t.Errorf("Max(), got: %v, want: %v", s, zeroItem)
	}

	var items []period.Ival
	zeroTree.Visit(zeroItem, zeroItem, func(item period.Ival) bool {
		items = append(items, item)
		return true
	})
	if len(items) != 0 {
		t.Errorf("Visit(), got: %v, want: 0", len(items))
	}
}

func TestTreeWithDups(t *testing.T) {
	t.Parallel()

	is := []period.Ival{
		{0, 100},
		{41, 102},
		{42, 67},
		{42, 67},
		{48, 50},
		{3, 13},
	}

	tree := treap.Insert(is...)
	if s := tree.Size(); s != 5 {
		t.Errorf("Size() = %v, want 5", s)
	}

	asStr := `▼
├─ 0...100
│  └─ 3...13
└─ 41...102
   └─ 42...67
      └─ 48...50
`
	w := new(strings.Builder)
	tree.Fprint(w)

	if w.String() != asStr {
		t.Errorf("Fprint()\nwant:\n%sgot:\n%s", asStr, w.String())
	}
}

func TestImmutable(t *testing.T) {
	t.Parallel()
	tree1 := treap.Insert(ps...)
	tree2 := tree1.Clone()

	if !reflect.DeepEqual(tree1, tree2) {
		t.Fatal("cloned tree is not deep equal to original")
	}

	if _, ok := tree1.Delete(tree2.Min()); !ok {
		t.Fatal("Delete, could not delete min item")
	}
	if !reflect.DeepEqual(tree1, tree2) {
		t.Fatal("Delete changed receiver")
	}

	item := period.Ival{-111, 666}
	tree1.Insert(item)
	if !reflect.DeepEqual(tree1, tree2) {
		t.Fatal("Insert changed receiver")
	}

	tree1.Shortest(item)
	if !reflect.DeepEqual(tree1, tree2) {
		t.Fatal("Shortest changed receiver")
	}

	tree1.Largest(item)
	if !reflect.DeepEqual(tree1, tree2) {
		t.Fatal("Largest changed receiver")
	}

	tree1.Subsets(item)
	if !reflect.DeepEqual(tree1, tree2) {
		t.Fatal("Subsets changed receiver")
	}

	tree1.Supersets(item)
	if !reflect.DeepEqual(tree1, tree2) {
		t.Fatal("Supersets changed receiver")
	}
}

func TestLookup(t *testing.T) {
	t.Parallel()

	// bring some variance into the Treap
	for i := 0; i < 100; i++ {
		tree := treap.Insert(ps...)

		//     	 ▼
		//     	 ├─ 0...6
		//     	 │  └─ 0...5
		//     	 ├─ 1...8
		//     	 │  ├─ 1...7
		//     	 │  │  └─ 1...5
		//     	 │  │     └─ 1...4
		//     	 │  └─ 2...8
		//     	 │     ├─ 2...7
		//     	 │     └─ 4...8
		//     	 │        └─ 6...7
		//     	 └─ 7...9

		item := period.Ival{0, 5}
		if got, _ := tree.Shortest(item); got != item {
			t.Errorf("Shortest(%v) = %v, want %v", item, got, item)
		}

		item = period.Ival{5, 5}
		want := period.Ival{4, 8}
		if got, _ := tree.Shortest(item); got != want {
			t.Errorf("Shortest(%v) = %v, want %v", item, got, want)
		}

		item = period.Ival{8, 9}
		want = period.Ival{7, 9}
		if got, _ := tree.Shortest(item); got != want {
			t.Errorf("Shortest(%v) = %v, want %v", item, got, want)
		}

		item = period.Ival{3, 8}
		want = period.Ival{2, 8}
		if got, _ := tree.Shortest(item); got != want {
			t.Errorf("Shortest(%v) = %v, want %v", item, got, want)
		}

		item = period.Ival{19, 55}
		if got, ok := tree.Shortest(item); ok {
			t.Errorf("Shortest(%v) = %v, want %v", item, got, !ok)
		}

		item = period.Ival{-19, 0}
		if got, ok := tree.Shortest(item); ok {
			t.Errorf("Shortest(%v) = %v, want %v", item, got, !ok)
		}

		item = period.Ival{8, 8}
		want = period.Ival{1, 8}
		if got, _ := tree.Largest(item); got != want {
			t.Errorf("Largest(%v) = %v, want %v", item, got, want)
		}

		item = period.Ival{3, 6}
		want = period.Ival{0, 6}
		if got, _ := tree.Largest(item); got != want {
			t.Errorf("Largest(%v) = %v, want %v", item, got, want)
		}

		item = period.Ival{3, 7}
		want = period.Ival{1, 8}
		if got, _ := tree.Largest(item); got != want {
			t.Errorf("Largest(%v) = %v, want %v", item, got, want)
		}

		item = period.Ival{0, 7}
		if _, ok := tree.Largest(item); ok {
			t.Errorf("Largest(%v) = %v, want %v", item, ok, false)
		}

	}
}

func TestSuperset(t *testing.T) {
	t.Parallel()

	is := []period.Ival{
		{1, 100},
		{45, 120},
		{46, 80},
	}

	tree := treap.Insert(is...)

	item := period.Ival{0, 6}
	if got, ok := tree.Largest(item); ok {
		t.Errorf("Largest(%v) = %v, want %v", item, got, !ok)
	}

	item = period.Ival{99, 200}
	if got, ok := tree.Largest(item); ok {
		t.Errorf("Largest(%v) = %v, want %v", item, got, !ok)
	}

	item = period.Ival{1, 100}
	if got, ok := tree.Largest(item); got != is[0] || !ok {
		t.Errorf("Largest(%v) = %v, want %v", item, got, is[0])
	}

	item = period.Ival{46, 80}
	if got, ok := tree.Largest(item); got != is[0] || !ok {
		t.Errorf("Largest(%v) = %v, want %v", item, got, is[0])
	}

	item = period.Ival{47, 62}
	if got, ok := tree.Largest(item); got != is[0] || !ok {
		t.Errorf("Largest(%v) = %v, want %v", item, got, is[0])
	}
}

func TestVisit(t *testing.T) {
	t.Parallel()
	tree := treap.Insert(ps...)

	var collect []period.Ival
	want := 4
	tree.Visit(tree.Min(), tree.Max(), func(item period.Ival) bool {
		collect = append(collect, item)
		return len(collect) != want
	})

	if len(collect) != want {
		t.Fatalf("Visit() ascending, want to stop after %v visits, got: %v, %v", want, len(collect), collect)
	}

	collect = nil
	want = 9
	tree.Visit(tree.Min(), tree.Max(), func(item period.Ival) bool {
		collect = append(collect, item)
		return len(collect) != want
	})

	if len(collect) != want {
		t.Fatalf("Visit() ascending, want to stop after %v visits, got: %v, %v", want, len(collect), collect)
	}

	collect = nil
	want = 2
	tree.Visit(tree.Max(), tree.Min(), func(item period.Ival) bool {
		collect = append(collect, item)
		return len(collect) != want
	})

	collect = nil
	want = 5
	tree.Visit(tree.Max(), tree.Min(), func(item period.Ival) bool {
		collect = append(collect, item)
		return len(collect) != want
	})
}

func TestMinMax(t *testing.T) {
	t.Parallel()
	tree := treap.Insert(ps...)
	want := period.Ival{0, 6}
	if tree.Min() != want {
		t.Fatalf("Min(), want: %v, got: %v", want, tree.Min())
	}

	want = period.Ival{7, 9}
	if tree.Max() != want {
		t.Fatalf("Max(), want: %v, got: %v", want, tree.Max())
	}
}

func TestUnion(t *testing.T) {
	t.Parallel()
	tree := treap

	for i := range ps {
		b := treap.Insert(ps[i])
		tree = tree.Union(b, false, true)
	}

	asStr := `▼
├─ 0...6
│  └─ 0...5
├─ 1...8
│  ├─ 1...7
│  │  └─ 1...5
│  │     └─ 1...4
│  └─ 2...8
│     ├─ 2...7
│     └─ 4...8
│        └─ 6...7
└─ 7...9
`

	w := new(strings.Builder)
	tree.Fprint(w)

	if w.String() != asStr {
		t.Errorf("Fprint()\nwant:\n%sgot:\n%s", asStr, w.String())
	}

	// now with dupe overwrite
	for i := range ps {
		b := treap.Insert(ps[i])
		tree = tree.Union(b, true, true)
	}

	w.Reset()
	tree.Fprint(w)
	if w.String() != asStr {
		t.Errorf("Fprint()\nwant:\n%sgot:\n%s", asStr, w.String())
	}
}

func TestStatistics(t *testing.T) {
	t.Parallel()

	for n := 10_000; n <= 1_000_000; n *= 10 {
		count := strconv.Itoa(n)
		t.Run(count, func(t *testing.T) {
			is := generateIvals(n)
			treap = nil

			tree := treap.Insert(is...)

			_, averageDepth, deviation := tree.Statistics()

			maxAverageDepth := 2 * math.Log2(float64(n))
			if averageDepth > maxAverageDepth {
				t.Fatalf("n: %d, average > max expected average, got: %.4g, want: < %.4g", n, averageDepth, maxAverageDepth)
			}

			maxDeviation := 1.0
			if deviation > maxDeviation {
				t.Fatalf("n: %d, deviation > max expected deviation, got: %.4g, want: < %.4g", n, deviation, maxDeviation)
			}
		})
	}
}

func TestRandom(t *testing.T) {
	t.Parallel()
	is := generateIvals(100)
	tree := treap.Insert(is...)

	rand.Shuffle(len(is), func(i, j int) { is[i], is[j] = is[j], is[i] })

	for _, item := range is {
		tname := fmt.Sprintf("%v", item)
		t.Run(tname, func(t *testing.T) {
			t.Parallel()
			var (
				shortest  period.Ival
				largest   period.Ival
				subsets   []period.Ival
				supersets []period.Ival
				ok        bool
			)
			if shortest, ok = tree.Shortest(item); !ok {
				t.Errorf("Shortest(%v), got %v, %v", item, shortest, ok)
			}
			if largest, ok = tree.Largest(item); !ok {
				t.Errorf("Largest(%v), got %v, %v", item, largest, ok)
			}
			if subsets = tree.Subsets(item); subsets == nil {
				t.Errorf("Subsets(%v), got %v", item, subsets)
			}
			if subsets[0] != shortest {
				t.Errorf("Subsets(%v).[0], want %v, got %v", item, shortest, subsets[0])
			}
			if supersets = tree.Supersets(item); supersets == nil {
				t.Errorf("Supersets(%v), got %v", item, supersets)
			}
			if supersets[0] != largest {
				t.Errorf("Supersets(%v).[0], want %v, got %v", item, largest, supersets[0])
			}
		})
	}
}

func TestPrintBST(t *testing.T) {
	t.Parallel()
	tree := treap.Insert(ps...)

	w := new(strings.Builder)
	_ = tree.FprintBST(w)

	lc := len(strings.Split(w.String(), "\n"))
	want := 12
	if lc != want {
		t.Fatalf("FprintBST(), want line count: %d, got: %d", want, lc)
	}
}
