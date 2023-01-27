package interval_test

import (
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

func mkIval(a, b uint) Ival {
	if a > b {
		a, b = b, a
	}
	return Ival{a, b}
}

// random test data
func generateIvals(n int) []Ival {
	is := make([]Ival, n)
	for i := 0; i < n; i++ {
		a := rand.Int()
		b := rand.Int()
		is[i] = mkIval(uint(a), uint(b))
	}
	return is
}

func equals(a, b Ival) bool {
	return a[0] == b[0] && a[1] == b[1]
}

func equalStatistics(t1, t2 interval.Tree[Ival]) bool {
	a1, b1, c1, d1 := t1.Statistics()
	a2, b2, c2, d2 := t2.Statistics()
	return a1 == a2 && b1 == b2 && c1 == c2 && d1 == d2
}

func TestNewTree(t *testing.T) {
	t.Parallel()

	var zeroItem Ival
	zeroTree := interval.NewTree(compareIval)

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

	if _, ok := zeroTree.Find(zeroItem); ok {
		t.Errorf("Find(), got: %v, want: false", ok)
	}

	if _, ok := zeroTree.Delete(zeroItem); ok {
		t.Errorf("Delete(), got: %v, want: false", ok)
	}

	if _, ok := zeroTree.CoverLCP(zeroItem); ok {
		t.Errorf("CoverLCP(), got: %v, want: false", ok)
	}

	if _, ok := zeroTree.CoverSCP(zeroItem); ok {
		t.Errorf("CoverSCP(), got: %v, want: false", ok)
	}

	if size, _, _, _ := zeroTree.Insert(zeroItem).Statistics(); size != 1 {
		t.Errorf("Insert(), got: %v, want: 1", size)
	}

	if s := zeroTree.CoveredBy(zeroItem); s != nil {
		t.Errorf("CoveredBy(), got: %v, want: nil", s)
	}

	if s := zeroTree.Covers(zeroItem); s != nil {
		t.Errorf("Covers(), got: %v, want: nil", s)
	}

	if s := zeroTree.Intersects(zeroItem); s != false {
		t.Errorf("Intersectons(), got: %v, want: false", s)
	}

	if s := zeroTree.Intersections(zeroItem); s != nil {
		t.Errorf("Intersections(), got: %v, want: nil", s)
	}

	if s := zeroTree.Precedes(zeroItem); s != nil {
		t.Errorf("Precedes(), got: %v, want: nil", s)
	}

	if s := zeroTree.PrecededBy(zeroItem); s != nil {
		t.Errorf("PrecededBy(), got: %v, want: nil", s)
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

	tree1 := interval.NewTree(compareIval, is...)
	if size, _, _, _ := tree1.Statistics(); size != 5 {
		t.Errorf("Size() = %v, want 5", size)
	}

	asStr := `▼
├─ 0...100
│  └─ 3...13
└─ 41...102
   └─ 42...67
      └─ 48...50
`
	if tree1.String() != asStr {
		t.Errorf("Fprint()\nwant:\n%sgot:\n%s", asStr, tree1.String())
	}
}

func TestImmutable(t *testing.T) {
	t.Parallel()
	tree1 := interval.NewTree(compareIval, ps...)

	if _, ok := tree1.Delete(tree1.Min()); !ok {
		t.Fatal("Delete, could not delete min item")
	}
	if _, ok := tree1.Delete(tree1.Min()); !ok {
		t.Fatal("Delete changed receiver")
	}

	item := Ival{111, 666}
	_ = tree1.Insert(item)

	if _, ok := tree1.Find(item); ok {
		t.Fatal("Insert changed receiver")
	}
}

func TestMutable(t *testing.T) {
	tree1 := interval.NewTree(compareIval, ps...)
	clone := tree1.Clone()

	if !equalStatistics(tree1, clone) {
		t.Fatalf("Clone, something wrong, statistics differs")
	}

	min := tree1.Min()

	var ok bool
	if ok = (&tree1).DeleteMutable(min); !ok {
		t.Fatal("DeleteMutable, could not delete min item")
	}

	if equalStatistics(tree1, clone) {
		t.Fatal("DeleteMutable didn't change receiver")
	}

	if ok = (&tree1).DeleteMutable(min); ok {
		t.Fatal("DeleteMutable didn't change receiver")
	}

	// reset
	tree1 = interval.NewTree(compareIval, ps...)
	clone = tree1.Clone()

	if !equalStatistics(tree1, clone) {
		t.Fatalf("Clone, something wrong, statistics differs")
	}

	item := Ival{111, 666}
	(&tree1).InsertMutable(item)

	if _, ok := tree1.Find(item); !ok {
		t.Fatal("InsertMutable didn't changed receiver")
	}
}

func TestFind(t *testing.T) {
	t.Parallel()

	ivals := generateIvals(100_00)
	tree1 := interval.NewTree(compareIval, ivals...)

	for _, ival := range ivals {
		item, ok := tree1.Find(ival)
		if ok != true {
			t.Errorf("Find(%v) = %v, want %v", item, ok, true)
		}
		ll, rr, _, _ := compareIval(item, ival)
		if ll != 0 || rr != 0 {
			t.Errorf("Find(%v) = %v, want %v", ival, item, ival)
		}
	}
}

func TestLookup(t *testing.T) {
	t.Parallel()

	for i := 0; i < 100; i++ {
		// bring some variance into the Treap due to the prio randomness
		tree1 := interval.NewTree(compareIval, ps...)

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
		if got, _ := tree1.CoverLCP(item); got != item {
			t.Errorf("CoverLCP(%v) = %v, want %v", item, got, item)
		}

		item = Ival{5, 5}
		want := Ival{4, 8}
		if got, _ := tree1.CoverLCP(item); got != want {
			t.Errorf("CoverLCP(%v) = %v, want %v", item, got, want)
		}

		item = Ival{8, 9}
		want = Ival{7, 9}
		if got, _ := tree1.CoverLCP(item); got != want {
			t.Errorf("CoverLCP(%v) = %v, want %v", item, got, want)
		}

		item = Ival{3, 8}
		want = Ival{2, 8}
		if got, _ := tree1.CoverLCP(item); got != want {
			t.Errorf("CoverLCP(%v) = %v, want %v", item, got, want)
		}

		item = Ival{19, 55}
		if got, ok := tree1.CoverLCP(item); ok {
			t.Errorf("CoverLCP(%v) = %v, want %v", item, got, !ok)
		}

		item = Ival{0, 19}
		if got, ok := tree1.CoverLCP(item); ok {
			t.Errorf("CoverLCP(%v) = %v, want %v", item, got, !ok)
		}

		item = Ival{7, 7}
		want = Ival{1, 8}
		if got, _ := tree1.CoverSCP(item); got != want {
			t.Errorf("CoverSCP(%v) = %v, want %v", item, got, want)
		}

		item = Ival{3, 6}
		want = Ival{0, 6}
		if got, _ := tree1.CoverSCP(item); got != want {
			t.Errorf("CoverSCP(%v) = %v, want %v", item, got, want)
		}

		item = Ival{3, 7}
		want = Ival{1, 8}
		if got, _ := tree1.CoverSCP(item); got != want {
			t.Errorf("CoverSCP(%v) = %v, want %v", item, got, want)
		}

		item = Ival{0, 7}
		if _, ok := tree1.CoverSCP(item); ok {
			t.Errorf("CoverSCP(%v) = %v, want %v", item, ok, false)
		}

	}
}

func TestCoveredBy(t *testing.T) {
	t.Parallel()

	tree1 := interval.NewTree(compareIval, ps...)
	var want []Ival

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

	item := Ival{0, 6}
	want = []Ival{{0, 6}, {0, 5}, {1, 5}, {1, 4}}
	covered := tree1.CoveredBy(item)

	if !reflect.DeepEqual(covered, want) {
		t.Fatalf("Covered, got: %v, want: %v", covered, want)
	}

	// ###
	item = Ival{3, 6}
	want = nil
	covered = tree1.CoveredBy(item)

	if !reflect.DeepEqual(covered, want) {
		t.Fatalf("Covered, got: %v, want: %v", covered, want)
	}

	// ###
	item = Ival{3, 11}
	want = []Ival{{4, 8}, {6, 7}, {7, 9}}
	covered = tree1.CoveredBy(item)

	if !reflect.DeepEqual(covered, want) {
		t.Fatalf("Covered(%v), got: %+v, want: %+v", item, covered, want)
	}
}

func TestCovers(t *testing.T) {
	t.Parallel()

	tree1 := interval.NewTree(compareIval, ps...)
	var want []Ival

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

	item := Ival{0, 6}
	want = []Ival{{0, 6}}
	covers := tree1.Covers(item)

	if !reflect.DeepEqual(covers, want) {
		t.Fatalf("Covers(%v), got: %v, want: %v", item, covers, want)
	}

	// ###
	item = Ival{3, 7}
	want = []Ival{{1, 8}, {1, 7}, {2, 8}, {2, 7}}
	covers = tree1.Covers(item)

	if !reflect.DeepEqual(covers, want) {
		t.Fatalf("Covers(%v), got: %v, want: %v", item, covers, want)
	}

	// ###
	item = Ival{3, 11}
	want = nil
	covers = tree1.Covers(item)

	if !reflect.DeepEqual(covers, want) {
		t.Fatalf("Covers(%v), got: %+v, want: %+v", item, covers, want)
	}
}

func TestIntersects(t *testing.T) {
	t.Parallel()

	tree1 := interval.NewTree(compareIval, ps...)

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

	item := Ival{7, 7}
	want := true
	got := tree1.Intersects(item)

	if got != want {
		t.Fatalf("Intersects(%v), got: %v, want: %v", item, got, want)
	}

	item = Ival{9, 17}
	want = true
	got = tree1.Intersects(item)

	if got != want {
		t.Fatalf("Intersects(%v), got: %v, want: %v", item, got, want)
	}

	item = Ival{1, 1}
	want = true
	got = tree1.Intersects(item)

	if got != want {
		t.Fatalf("Intersects(%v), got: %v, want: %v", item, got, want)
	}

	item = Ival{10, 12}
	want = false
	got = tree1.Intersects(item)

	if got != want {
		t.Fatalf("Intersects(%v), got: %v, want: %v", item, got, want)
	}
}

func TestIntersections(t *testing.T) {
	t.Parallel()

	tree1 := interval.NewTree(compareIval, ps...)
	var want []Ival

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

	item := Ival{7, 7}
	want = []Ival{{1, 8}, {1, 7}, {2, 8}, {2, 7}, {4, 8}, {6, 7}, {7, 9}}
	intersections := tree1.Intersections(item)

	if !reflect.DeepEqual(intersections, want) {
		t.Fatalf("Intersections(%v), got: %v, want: %v", item, intersections, want)
	}

	// ###
	item = Ival{8, 10}
	want = []Ival{{1, 8}, {2, 8}, {4, 8}, {7, 9}}
	intersections = tree1.Intersections(item)

	if !reflect.DeepEqual(intersections, want) {
		t.Fatalf("Intersections(%v), got: %v, want: %v", item, intersections, want)
	}

	// ###
	item = Ival{10, 15}
	want = nil
	intersections = tree1.Intersections(item)

	if !reflect.DeepEqual(intersections, want) {
		t.Fatalf("Intersections(%v), got: %+v, want: %+v", item, intersections, want)
	}
}

func TestPrecedes(t *testing.T) {
	t.Parallel()

	tree1 := interval.NewTree(compareIval, ps...)
	var want []Ival

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

	item := Ival{7, 7}
	want = []Ival{{0, 6}, {0, 5}, {1, 5}, {1, 4}}
	precedes := tree1.Precedes(item)

	if !reflect.DeepEqual(precedes, want) {
		t.Fatalf("Precedes(%v), got: %v, want: %v", item, precedes, want)
	}

	// ###
	item = Ival{5, 10}
	want = []Ival{{1, 4}}
	precedes = tree1.Precedes(item)

	if !reflect.DeepEqual(precedes, want) {
		t.Fatalf("Precedes(%v), got: %v, want: %v", item, precedes, want)
	}

	// ###
	item = Ival{0, 9}
	want = nil
	precedes = tree1.Precedes(item)

	if !reflect.DeepEqual(precedes, want) {
		t.Fatalf("Precedes(%v), got: %+v, want: %+v", item, precedes, want)
	}
}

func TestPrecededBy(t *testing.T) {
	t.Parallel()

	tree1 := interval.NewTree(compareIval, ps...)
	var want []Ival

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

	item := Ival{4, 4}
	want = []Ival{{6, 7}, {7, 9}}
	precedes := tree1.PrecededBy(item)

	if !reflect.DeepEqual(precedes, want) {
		t.Fatalf("PrecededBy(%v), got: %v, want: %v", item, precedes, want)
	}

	// ###
	item = Ival{1, 2}
	want = []Ival{{4, 8}, {6, 7}, {7, 9}}
	precedes = tree1.PrecededBy(item)

	if !reflect.DeepEqual(precedes, want) {
		t.Fatalf("PrecededBy(%v), got: %v, want: %v", item, precedes, want)
	}

	// ###
	item = Ival{0, 7}
	want = nil
	precedes = tree1.PrecededBy(item)

	if !reflect.DeepEqual(precedes, want) {
		t.Fatalf("PrecededBy(%v), got: %+v, want: %+v", item, precedes, want)
	}
}

func TestVisit(t *testing.T) {
	t.Parallel()
	tree1 := interval.NewTree(compareIval, ps...)

	var collect []Ival
	want := 4
	tree1.Visit(tree1.Min(), tree1.Max(), func(item Ival) bool {
		collect = append(collect, item)
		return len(collect) != want
	})

	if len(collect) != want {
		t.Fatalf("Visit() ascending, want to stop after %v visits, got: %v, %v", want, len(collect), collect)
	}

	collect = nil
	want = 9
	tree1.Visit(tree1.Max(), tree1.Min(), func(item Ival) bool {
		collect = append(collect, item)
		return true
	})

	want, _, _, _ = tree1.Statistics()
	if len(collect) != want {
		t.Fatalf("Visit() descending, want: %d  got: %v, %v", want, len(collect), collect)
	}

	collect = nil
	want = 2
	tree1.Visit(tree1.Max(), tree1.Min(), func(item Ival) bool {
		collect = append(collect, item)
		return len(collect) != want
	})
}

func TestMinMax(t *testing.T) {
	t.Parallel()
	tree1 := interval.NewTree(compareIval, ps...)
	want := Ival{0, 6}
	if tree1.Min() != want {
		t.Fatalf("Min(), want: %v, got: %v", want, tree1.Min())
	}

	want = Ival{7, 9}
	if tree1.Max() != want {
		t.Fatalf("Max(), want: %v, got: %v", want, tree1.Max())
	}
}

func TestUnion(t *testing.T) {
	t.Parallel()
	tree1 := interval.NewTree(compareIval)

	for i := range ps {
		b := tree1.Insert(ps[i])
		tree1 = tree1.Union(b, false, true)
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

	if tree1.String() != asStr {
		t.Errorf("Fprint()\nwant:\n%sgot:\n%s", asStr, tree1.String())
	}

	// now with dupe overwrite
	for i := range ps {
		b := tree1.Insert(ps[i])
		tree1 = tree1.Union(b, true, true)
	}

	if tree1.String() != asStr {
		t.Errorf("String()\nwant:\n%sgot:\n%s", asStr, tree1.String())
	}

	ps2 := []Ival{
		{7, 60},
		{8, 50},
		{9, 80},
		{9, 70},
		{9, 50},
		{9, 40},
		{2, 8},
		{2, 7},
		{4, 8},
		{6, 7},
		{7, 9},
	}

	tree2 := tree1.Insert(ps2...)
	tree1 = tree1.Union(tree2, false, false)

	asStr = `▼
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
├─ 7...60
│  ├─ 7...9
│  └─ 8...50
└─ 9...80
   └─ 9...70
      └─ 9...50
         └─ 9...40
`

	if tree1.String() != asStr {
		t.Errorf("String()\nwant:\n%sgot:\n%s", asStr, tree1.String())
	}
}

func TestStatistics(t *testing.T) {
	t.Parallel()

	for n := 10_000; n <= 1_000_000; n *= 10 {
		count := strconv.Itoa(n)
		t.Run(count, func(t *testing.T) {
			tree1 := interval.NewTree(compareIval, generateIvals(n)...)

			size, _, averageDepth, deviation := tree1.Statistics()
			if size != n {
				t.Fatalf("size, got: %d, want: %d", size, n)
			}

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

func TestPrintBST(t *testing.T) {
	t.Parallel()
	tree1 := interval.NewTree(compareIval, ps...)

	w := new(strings.Builder)
	_ = tree1.FprintBST(w)

	lc := len(strings.Split(w.String(), "\n"))
	want := 12
	if lc != want {
		t.Fatalf("FprintBST(), want line count: %d, got: %d", want, lc)
	}
}

func TestMatch(t *testing.T) {
	t.Parallel()
	tree1 := interval.NewTree(compareIval, generateIvals(100_000)...)

	n := 100
	for i := 0; i < n; i++ {
		probe := generateIvals(100_000)[0]

		t.Run(probe.String(), func(t *testing.T) {
			tree1 := tree1.Insert(probe)

			if _, ok := tree1.Find(probe); !ok {
				t.Fatalf("inserted item not found in tree: %v", probe)
			}

			shortest, short_ok := tree1.CoverLCP(probe)
			largest, large_ok := tree1.CoverSCP(probe)

			covers := tree1.Covers(probe)
			covered := tree1.CoveredBy(probe)
			intersections := tree1.Intersections(probe)

			// either both or neither
			if short_ok && !large_ok || large_ok && !short_ok {
				t.Fatalf("logic error: short_ok: %v, large_ok: %v", short_ok, large_ok)
			}

			lenCovers := len(covers)
			lenCovered := len(covered)
			lenIntersects := len(intersections)

			if short_ok && lenCovers == 0 {
				t.Fatalf("logic error: shortest: %v, len(covered): %v, len(covers): %v", shortest, lenCovered, lenCovers)
			}

			if short_ok && !equals(covers[lenCovers-1], shortest) {
				t.Fatalf("logic error: covers[last]: %v IS NOT shortest: %v", covers[lenCovers-1], shortest)
			}

			if large_ok && !equals(covers[0], largest) {
				t.Fatalf("logic error: covers[0]: %v IS NOT largest: %v", covers[0], largest)
			}

			if lenIntersects < lenCovered+lenCovers {
				t.Fatalf("logic error: len(intersections) MUST BE >= len(covered) + len(covers): %d IS NOT > %d + %d",
					lenIntersects, lenCovered, lenCovers)
			}
		})
	}
}

func TestMissing(t *testing.T) {
	t.Parallel()
	tree1 := interval.NewTree(compareIval, generateIvals(100_000)...)

	n := 100
	for i := 0; i < n; i++ {
		probe := generateIvals(100_000)[0]

		t.Run(probe.String(), func(t *testing.T) {
			tree1 := tree1.Insert(probe)
			var ok bool

			if _, ok = tree1.Find(probe); !ok {
				t.Fatalf("inserted item not found in tree: %v", probe)
			}

			if tree1, ok = tree1.Delete(probe); !ok {
				t.Fatalf("delete, inserted item not found in tree: %v", probe)
			}

			if _, ok = tree1.Find(probe); ok {
				t.Fatalf("deleted item still found in tree: %v", probe)
			}

			shortest, short_ok := tree1.CoverLCP(probe)
			largest, large_ok := tree1.CoverSCP(probe)

			covers := tree1.Covers(probe)
			covered := tree1.CoveredBy(probe)
			intersections := tree1.Intersections(probe)

			// either both or neither
			if short_ok && !large_ok || large_ok && !short_ok {
				t.Fatalf("logic error: short_ok: %v, large_ok: %v", short_ok, large_ok)
			}

			lenCovers := len(covers)
			lenCovered := len(covered)
			lenIntersects := len(intersections)

			if short_ok && lenCovers == 0 {
				t.Fatalf("logic error: shortest: %v, len(covered): %v, len(covers): %v", shortest, lenCovered, lenCovers)
			}

			if short_ok && !equals(covers[lenCovers-1], shortest) {
				t.Fatalf("logic error: covers[last]: %v IS NOT shortest: %v", covers[lenCovers-1], shortest)
			}

			if large_ok && !equals(covers[0], largest) {
				t.Fatalf("logic error: covers[0]: %v IS NOT largest: %v", covers[0], largest)
			}

			if lenIntersects < lenCovered+lenCovers {
				t.Fatalf("logic error: len(intersections) MUST BE >= len(covered) + len(covers): %d IS NOT > %d + %d",
					lenIntersects, lenCovered, lenCovers)
			}
		})
	}
}
