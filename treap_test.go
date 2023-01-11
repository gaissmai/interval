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
)

// test data
var ps = []Ival{
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
func generateIvals(n int) []Ival {
	is := make([]Ival, n)
	for i := 0; i < n; i++ {
		a := rand.Int()
		b := rand.Int()
		if a > b {
			a, b = b, a
		}
		is[i] = Ival{a, b}
	}
	return is
}

func TestNewTree(t *testing.T) {
	t.Parallel()

	var zeroItem Ival
	var zeroTree interval.Tree[Ival]

	if zeroTree.String() != "" {
		t.Errorf("String() = %v, want \"\"", "")
	}

	w := new(strings.Builder)
	if err := zeroTree.Fprint(w); err != nil {
		t.Fatal(err)
	}

	if w.String() != "" {
		t.Errorf("Fprint(w) = %v, want \"\"", w.String())
	}

	w.Reset()
	if err := zeroTree.FprintBST(w); err != nil {
		t.Fatal(err)
	}

	if w.String() != "" {
		t.Errorf("FprintBST(w) = %v, want \"\"", w.String())
	}

	if _, ok := zeroTree.Delete(zeroItem); ok {
		t.Errorf("Delete(), got: %v, want: false", ok)
	}

	if _, ok := zeroTree.Shortest(zeroItem); ok {
		t.Errorf("Shortest(), got: %v, want: false", ok)
	}

	if _, ok := zeroTree.Largest(zeroItem); ok {
		t.Errorf("Largest(), got: %v, want: false", ok)
	}

	if s := zeroTree.Insert(zeroItem); s.Size() != 1 {
		t.Errorf("Insert(), got: %v, want: 1", s.Size())
	}

	if s := zeroTree.Clone(); s.Size() != 0 {
		t.Errorf("Clone(), got: %v, want: 0", s.Size())
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

	var items []Ival
	zeroTree.Visit(zeroItem, zeroItem, func(item Ival) bool {
		items = append(items, item)
		return true
	})
	if len(items) != 0 {
		t.Errorf("Visit(), got: %v, want: 0", len(items))
	}
}

func TestTreeWithDups(t *testing.T) {
	t.Parallel()

	is := []Ival{
		{0, 100},
		{41, 102},
		{41, 102},
		{41, 102},
		{41, 102},
		{41, 102},
		{41, 102},
		{41, 102},
		{42, 67},
		{42, 67},
		{42, 67},
		{42, 67},
		{42, 67},
		{42, 67},
		{42, 67},
		{42, 67},
		{48, 50},
		{3, 13},
		{3, 13},
		{3, 13},
		{3, 13},
		{3, 13},
		{3, 13},
	}

	tree := interval.NewTree(is...)
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
	if tree.String() != asStr {
		t.Errorf("Fprint()\nwant:\n%sgot:\n%s", asStr, tree.String())
	}
}

func TestImmutable(t *testing.T) {
	t.Parallel()
	tree1 := interval.NewTree(ps...)
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

	item := Ival{-111, 666}
	_ = tree1.Insert(item)
	if !reflect.DeepEqual(tree1, tree2) {
		t.Fatal("Insert changed receiver")
	}

	_, _ = tree1.Shortest(item)
	if !reflect.DeepEqual(tree1, tree2) {
		t.Fatal("Shortest changed receiver")
	}

	_, _ = tree1.Largest(item)
	if !reflect.DeepEqual(tree1, tree2) {
		t.Fatal("Largest changed receiver")
	}

	_ = tree1.Subsets(item)
	if !reflect.DeepEqual(tree1, tree2) {
		t.Fatal("Subsets changed receiver")
	}

	_ = tree1.Supersets(item)
	if !reflect.DeepEqual(tree1, tree2) {
		t.Fatal("Supersets changed receiver")
	}
}

func TestMutable(t *testing.T) {
	tree1 := interval.NewTree(ps...)
	tree2 := tree1.Clone()

	min := tree1.Min()

	var ok bool
	if ok = (&tree1).DeleteMutable(min); !ok {
		t.Fatal("DeleteMutable, could not delete min item")
	}
	if reflect.DeepEqual(tree1, tree2) {
		t.Fatal("DeleteMutable didn't change receiver")
	}

	// reset tree1, tree2
	tree1 = interval.NewTree(ps...)
	tree2 = tree1.Clone()

	item := Ival{-111, 666}
	(&tree1).InsertMutable(item)

	if reflect.DeepEqual(tree1, tree2) {
		t.Fatal("InsertMutable didn't change receiver")
	}
	if _, ok := tree1.Delete(item); !ok {
		t.Fatal("InsertMutable didn't change receiver")
	}
}

func TestFind(t *testing.T) {
	t.Parallel()

	ivals := generateIvals(100_00)
	tree := interval.NewTree(ivals...)

	for _, ival := range ivals {
		item, ok := tree.Find(ival)
		if ok != true {
			t.Errorf("Find(%v) = %v, want %v", item, ok, true)
		}
		if item.CompareLower(ival) != 0 || item.CompareUpper(ival) != 0 {
			t.Errorf("Find(%v) = %v, want %v", ival, item, ival)
		}
	}
}

func TestLookup(t *testing.T) {
	t.Parallel()

	for i := 0; i < 100; i++ {
		// bring some variance into the Treap due to the prio randomness
		tree := interval.NewTree(ps...)

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

		item := Ival{0, 5}
		if got, _ := tree.Shortest(item); got != item {
			t.Errorf("Shortest(%v) = %v, want %v", item, got, item)
		}

		item = Ival{5, 5}
		want := Ival{4, 8}
		if got, _ := tree.Shortest(item); got != want {
			t.Errorf("Shortest(%v) = %v, want %v", item, got, want)
		}

		item = Ival{8, 9}
		want = Ival{7, 9}
		if got, _ := tree.Shortest(item); got != want {
			t.Errorf("Shortest(%v) = %v, want %v", item, got, want)
		}

		item = Ival{3, 8}
		want = Ival{2, 8}
		if got, _ := tree.Shortest(item); got != want {
			t.Errorf("Shortest(%v) = %v, want %v", item, got, want)
		}

		item = Ival{19, 55}
		if got, ok := tree.Shortest(item); ok {
			t.Errorf("Shortest(%v) = %v, want %v", item, got, !ok)
		}

		item = Ival{-19, 0}
		if got, ok := tree.Shortest(item); ok {
			t.Errorf("Shortest(%v) = %v, want %v", item, got, !ok)
		}

		item = Ival{8, 8}
		want = Ival{1, 8}
		if got, _ := tree.Largest(item); got != want {
			t.Errorf("Largest(%v) = %v, want %v", item, got, want)
		}

		item = Ival{3, 6}
		want = Ival{0, 6}
		if got, _ := tree.Largest(item); got != want {
			t.Errorf("Largest(%v) = %v, want %v", item, got, want)
		}

		item = Ival{3, 7}
		want = Ival{1, 8}
		if got, _ := tree.Largest(item); got != want {
			t.Errorf("Largest(%v) = %v, want %v", item, got, want)
		}

		item = Ival{0, 7}
		if _, ok := tree.Largest(item); ok {
			t.Errorf("Largest(%v) = %v, want %v", item, ok, false)
		}

	}
}

func TestSuperset(t *testing.T) {
	t.Parallel()

	is := []Ival{
		{1, 100},
		{45, 120},
		{46, 80},
	}

	tree := interval.NewTree(is...)

	item := Ival{0, 6}
	if got, ok := tree.Largest(item); ok {
		t.Errorf("Largest(%v) = %v, want %v", item, got, !ok)
	}

	item = Ival{99, 200}
	if got, ok := tree.Largest(item); ok {
		t.Errorf("Largest(%v) = %v, want %v", item, got, !ok)
	}

	item = Ival{1, 100}
	if got, ok := tree.Largest(item); got != is[0] || !ok {
		t.Errorf("Largest(%v) = %v, want %v", item, got, is[0])
	}

	item = Ival{46, 80}
	if got, ok := tree.Largest(item); got != is[0] || !ok {
		t.Errorf("Largest(%v) = %v, want %v", item, got, is[0])
	}

	item = Ival{47, 62}
	if got, ok := tree.Largest(item); got != is[0] || !ok {
		t.Errorf("Largest(%v) = %v, want %v", item, got, is[0])
	}
}

func TestVisit(t *testing.T) {
	t.Parallel()
	tree := interval.NewTree(ps...)

	var collect []Ival
	want := 4
	tree.Visit(tree.Min(), tree.Max(), func(item Ival) bool {
		collect = append(collect, item)
		return len(collect) != want
	})

	if len(collect) != want {
		t.Fatalf("Visit() ascending, want to stop after %v visits, got: %v, %v", want, len(collect), collect)
	}

	collect = nil
	want = 9
	tree.Visit(tree.Min(), tree.Max(), func(item Ival) bool {
		collect = append(collect, item)
		return len(collect) != want
	})

	if len(collect) != want {
		t.Fatalf("Visit() ascending, want to stop after %v visits, got: %v, %v", want, len(collect), collect)
	}

	collect = nil
	want = 2
	tree.Visit(tree.Max(), tree.Min(), func(item Ival) bool {
		collect = append(collect, item)
		return len(collect) != want
	})

	collect = nil
	want = 5
	tree.Visit(tree.Max(), tree.Min(), func(item Ival) bool {
		collect = append(collect, item)
		return len(collect) != want
	})
}

func TestMinMax(t *testing.T) {
	t.Parallel()
	tree := interval.NewTree(ps...)
	want := Ival{0, 6}
	if tree.Min() != want {
		t.Fatalf("Min(), want: %v, got: %v", want, tree.Min())
	}

	want = Ival{7, 9}
	if tree.Max() != want {
		t.Fatalf("Max(), want: %v, got: %v", want, tree.Max())
	}
}

func TestUnion(t *testing.T) {
	t.Parallel()
	tree := interval.NewTree[Ival]()

	for i := range ps {
		b := interval.NewTree(ps[i])
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

	if tree.String() != asStr {
		t.Errorf("Fprint()\nwant:\n%sgot:\n%s", asStr, tree.String())
	}

	// now with dupe overwrite
	for i := range ps {
		b := interval.NewTree(ps[i])
		tree = tree.Union(b, true, true)
	}

	w := new(strings.Builder)
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
			tree := interval.NewTree(generateIvals(n)...)

			_, averageDepth, deviation := tree.Statistics()
			t.Logf("stats: n=%d, averageDepth=%.4g, deviation=%.4g\n", n, averageDepth, deviation)

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
	tree := interval.NewTree(is...)

	rand.Shuffle(len(is), func(i, j int) { is[i], is[j] = is[j], is[i] })

	for _, item := range is {
		tname := fmt.Sprintf("%v", item)
		t.Run(tname, func(t *testing.T) {
			t.Parallel()
			var (
				shortest  Ival
				largest   Ival
				subsets   []Ival
				supersets []Ival
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
	tree := interval.NewTree(ps...)

	w := new(strings.Builder)
	_ = tree.FprintBST(w)

	lc := len(strings.Split(w.String(), "\n"))
	want := 12
	if lc != want {
		t.Fatalf("FprintBST(), want line count: %d, got: %d", want, lc)
	}
}
