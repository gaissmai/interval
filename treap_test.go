package interval_test

import (
	"math"
	"math/rand"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"testing"

	"github.com/gaissmai/interval"
)

// test data
var ps = []uintInterval{
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

func makeUintIval(a, b uint) uintInterval {
	if a > b {
		a, b = b, a
	}
	return uintInterval{a, b}
}

// random test data, end is random to start
func genUintIvals(n int) []uintInterval {
	is := make([]uintInterval, n)
	for i := 0; i < n; i++ {
		a := rand.Int()
		b := rand.Int()
		is[i] = makeUintIval(uint(a), uint(b))
	}
	return is
}

// random test data, end depends on start
func gen2UintIvals(n int) []uintInterval {
	is := make([]uintInterval, n)
	for i := 0; i < n; i++ {
		a := rand.Int()
		b := a + 100
		is[i] = makeUintIval(uint(a), uint(b))
	}
	return is
}

func equals(a, b uintInterval) bool {
	return a[0] == b[0] && a[1] == b[1]
}

func equalStatistics(t1, t2 *interval.Tree[uintInterval]) bool {
	a1, b1, c1, d1 := t1.Statistics()
	a2, b2, c2, d2 := t2.Statistics()
	return a1 == a2 && b1 == b2 && c1 == c2 && d1 == d2
}

func equalsSizeAndOrder[T any](t1, t2 *interval.Tree[T]) bool {
	var t1InOrder []T
	t1.Visit(t1.Min(), t1.Max(), func(item T) bool {
		t1InOrder = append(t1InOrder, item)
		return true
	})

	var t2InOrder []T
	t2.Visit(t2.Min(), t2.Max(), func(item T) bool {
		t2InOrder = append(t2InOrder, item)
		return true
	})

	return reflect.DeepEqual(t1InOrder, t2InOrder)
}

func TestNewTree(t *testing.T) {
	t.Parallel()

	var zeroItem uintInterval
	tree := interval.NewTree(cmpUintInterval)

	if tree.String() != "" {
		t.Errorf("String() = %v, want \"\"", "")
	}

	tree = interval.NewTreeConcurrent(0, cmpUintInterval)

	if tree.String() != "" {
		t.Errorf("String() = %v, want \"\"", "")
	}

	w := new(strings.Builder)
	if err := tree.Fprint(w); err != nil {
		t.Fatal(err)
	}

	if w.String() != "" {
		t.Errorf("Fprint(w) = %v, want \"\"", w.String())
	}

	w.Reset()
	if err := tree.FprintBST(w); err != nil {
		t.Fatal(err)
	}

	if w.String() != "" {
		t.Errorf("FprintBST(w) = %v, want \"\"", w.String())
	}

	if _, ok := tree.Find(zeroItem); ok {
		t.Errorf("Find(), got: %v, want: false", ok)
	}

	if _, ok := tree.DeleteImmutable(zeroItem); ok {
		t.Errorf("Delete(), got: %v, want: false", ok)
	}

	if _, ok := tree.CoverLCP(zeroItem); ok {
		t.Errorf("CoverLCP(), got: %v, want: false", ok)
	}

	if _, ok := tree.CoverSCP(zeroItem); ok {
		t.Errorf("CoverSCP(), got: %v, want: false", ok)
	}

	if size, _, _, _ := tree.InsertImmutable(zeroItem).Statistics(); size != 1 {
		t.Errorf("Insert(), got: %v, want: 1", size)
	}

	if s := tree.CoveredBy(zeroItem); s != nil {
		t.Errorf("CoveredBy(), got: %v, want: nil", s)
	}

	if s := tree.Covers(zeroItem); s != nil {
		t.Errorf("Covers(), got: %v, want: nil", s)
	}

	if s := tree.Intersects(zeroItem); s != false {
		t.Errorf("Intersectons(), got: %v, want: false", s)
	}

	if s := tree.Intersections(zeroItem); s != nil {
		t.Errorf("Intersections(), got: %v, want: nil", s)
	}

	if s := tree.Precedes(zeroItem); s != nil {
		t.Errorf("Precedes(), got: %v, want: nil", s)
	}

	if s := tree.PrecededBy(zeroItem); s != nil {
		t.Errorf("PrecededBy(), got: %v, want: nil", s)
	}

	if s := tree.Min(); s != zeroItem {
		t.Errorf("Min(), got: %v, want: %v", s, zeroItem)
	}

	if s := tree.Max(); s != zeroItem {
		t.Errorf("Max(), got: %v, want: %v", s, zeroItem)
	}

	var items []uintInterval
	tree.Visit(zeroItem, zeroItem, func(item uintInterval) bool {
		items = append(items, item)
		return true
	})
	if len(items) != 0 {
		t.Errorf("Visit(), got: %v, want: 0", len(items))
	}
}

func TestNewTreeConcurrent(t *testing.T) {
	t.Parallel()

	ivals := genUintIvals(100_000)

	tree1 := interval.NewTree(cmpUintInterval, ivals[0])
	tree2 := interval.NewTreeConcurrent(1, cmpUintInterval, ivals[0])

	if !equalsSizeAndOrder(tree1, tree2) {
		t.Fatal("New() differs with NewConcurrent()")
	}

	tree1 = interval.NewTree(cmpUintInterval, ivals[:2]...)
	tree2 = interval.NewTreeConcurrent(2, cmpUintInterval, ivals[:2]...)

	if !equalsSizeAndOrder(tree1, tree2) {
		t.Fatal("New() differs with NewConcurrent()")
	}

	tree1 = interval.NewTree(cmpUintInterval, ivals[:30_000]...)
	tree2 = interval.NewTreeConcurrent(3, cmpUintInterval, ivals[:30_000]...)

	if !equalsSizeAndOrder(tree1, tree2) {
		t.Fatal("New() differs with NewConcurrent()")
	}

	tree1 = interval.NewTree(cmpUintInterval, ivals...)
	tree2 = interval.NewTreeConcurrent(runtime.NumCPU(), cmpUintInterval, ivals...)

	if !equalsSizeAndOrder(tree1, tree2) {
		t.Fatal("New() differs with NewConcurrent()")
	}
}

func TestTreeWithDups(t *testing.T) {
	t.Parallel()

	is := []uintInterval{
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

	tree1 := interval.NewTree(cmpUintInterval, is...)
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
	tree1 := interval.NewTree(cmpUintInterval, ps...)

	if _, ok := tree1.DeleteImmutable(tree1.Min()); !ok {
		t.Fatal("Delete, could not delete min item")
	}
	if _, ok := tree1.DeleteImmutable(tree1.Min()); !ok {
		t.Fatal("Delete changed receiver")
	}

	item := uintInterval{111, 666}
	_ = tree1.InsertImmutable(item)

	if _, ok := tree1.Find(item); ok {
		t.Fatal("Insert changed receiver")
	}
}

func TestMutable(t *testing.T) {
	tree1 := interval.NewTree(cmpUintInterval, ps...)
	clone := tree1.Clone()

	if !equalStatistics(tree1, clone) {
		t.Error("Clone, something wrong, statistics differs")
	}

	min := tree1.Min()

	var ok bool
	if ok = tree1.Delete(min); !ok {
		t.Fatal("Delete, could not delete min item")
	}

	if equalStatistics(tree1, clone) {
		t.Fatal("Delete didn't change receiver")
	}

	if ok = tree1.Delete(min); ok {
		t.Fatal("Delete didn't change receiver")
	}

	// reset
	tree1 = interval.NewTree(cmpUintInterval, ps...)
	clone = tree1.Clone()

	if !equalStatistics(tree1, clone) {
		t.Fatalf("Clone, something wrong, statistics differs")
	}

	item := uintInterval{111, 666}
	tree1.Insert(item)

	if _, ok := tree1.Find(item); !ok {
		t.Fatal("Insert didn't changed receiver")
	}
}

func TestFind(t *testing.T) {
	t.Parallel()

	ivals := genUintIvals(100_00)
	tree1 := interval.NewTree(cmpUintInterval, ivals...)

	for _, ival := range ivals {
		ival := ival
		t.Run(ival.String(), func(t *testing.T) {
			t.Parallel()
			item, ok := tree1.Find(ival)
			if ok != true {
				t.Errorf("Find(%v) = %v, want %v", item, ok, true)
			}
			ll, rr, _, _ := cmpUintInterval(item, ival)
			if ll != 0 || rr != 0 {
				t.Errorf("Find(%v) = %v, want %v", ival, item, ival)
			}
		})
	}
}

func TestCoverLCP(t *testing.T) {
	t.Parallel()

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

	testcases := []struct {
		in     uintInterval
		want   uintInterval
		wantOK bool
	}{
		{
			in:     uintInterval{0, 6},
			want:   uintInterval{0, 6},
			wantOK: true,
		},

		{
			in:     uintInterval{5, 5},
			want:   uintInterval{4, 8},
			wantOK: true,
		},
		{
			in:     uintInterval{8, 9},
			want:   uintInterval{7, 9},
			wantOK: true,
		},
		{
			in:     uintInterval{3, 5},
			want:   uintInterval{2, 7},
			wantOK: true,
		},
	}

	for i := 0; i < 100; i++ {
		// bring some variance into the Treap due to the prio randomness
		tree1 := interval.NewTree(cmpUintInterval, ps...)

		for _, tc := range testcases {
			tc := tc // capture range variable
			t.Run("", func(t *testing.T) {
				t.Parallel()
				if got, ok := tree1.CoverLCP(tc.in); got != tc.want || ok != tc.wantOK {
					t.Errorf("CoverLCP(%v) = (%v, %v) want (%v, %v)", tc.in, got, ok, tc.want, tc.wantOK)
				}
			})
		}
	}
}

func TestCoverSCP(t *testing.T) {
	t.Parallel()

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

	testcases := []struct {
		in     uintInterval
		want   uintInterval
		wantOK bool
	}{
		{
			in:     uintInterval{7, 7},
			want:   uintInterval{1, 8},
			wantOK: true,
		},

		{
			in:     uintInterval{1, 8},
			want:   uintInterval{1, 8},
			wantOK: true,
		},
		{
			in:     uintInterval{7, 9},
			want:   uintInterval{7, 9},
			wantOK: true,
		},
		{
			in:     uintInterval{1, 8},
			want:   uintInterval{1, 8},
			wantOK: true,
		},
		{
			in:     uintInterval{7, 9},
			want:   uintInterval{7, 9},
			wantOK: true,
		},
		{
			in:     uintInterval{8, 9},
			want:   uintInterval{7, 9},
			wantOK: true,
		},
		{
			in:     uintInterval{0, 6},
			want:   uintInterval{0, 6},
			wantOK: true,
		},
		{
			in:     uintInterval{3, 6},
			want:   uintInterval{0, 6},
			wantOK: true,
		},
		{
			in:     uintInterval{3, 7},
			want:   uintInterval{1, 8},
			wantOK: true,
		},
		{
			in:     uintInterval{0, 7},
			want:   uintInterval{},
			wantOK: false,
		},
		{
			in:     uintInterval{6, 10},
			want:   uintInterval{},
			wantOK: false,
		},
	}

	for i := 0; i < 100; i++ {
		// bring some variance into the Treap due to the prio randomness
		tree1 := interval.NewTree(cmpUintInterval, ps...)

		for _, tc := range testcases {
			tc := tc
			t.Run("", func(t *testing.T) {
				t.Parallel()
				if got, ok := tree1.CoverSCP(tc.in); got != tc.want || ok != tc.wantOK {
					t.Errorf("CoverSCP(%v) = (%v, %v) want (%v, %v)", tc.in, got, ok, tc.want, tc.wantOK)
				}
			})
		}
	}
}

func TestCoveredBy(t *testing.T) {
	t.Parallel()

	tree1 := interval.NewTree(cmpUintInterval, ps...)
	var want []uintInterval

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

	item := uintInterval{0, 6}
	want = []uintInterval{{0, 6}, {0, 5}, {1, 5}, {1, 4}}
	covered := tree1.CoveredBy(item)

	if !reflect.DeepEqual(covered, want) {
		t.Fatalf("Covered, got: %v, want: %v", covered, want)
	}

	// ###
	item = uintInterval{3, 6}
	want = nil
	covered = tree1.CoveredBy(item)

	if !reflect.DeepEqual(covered, want) {
		t.Fatalf("Covered, got: %v, want: %v", covered, want)
	}

	// ###
	item = uintInterval{3, 11}
	want = []uintInterval{{4, 8}, {6, 7}, {7, 9}}
	covered = tree1.CoveredBy(item)

	if !reflect.DeepEqual(covered, want) {
		t.Fatalf("Covered(%v), got: %+v, want: %+v", item, covered, want)
	}
}

func TestCovers(t *testing.T) {
	t.Parallel()

	tree1 := interval.NewTree(cmpUintInterval, ps...)
	var want []uintInterval

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

	item := uintInterval{0, 6}
	want = []uintInterval{{0, 6}}
	covers := tree1.Covers(item)

	if !reflect.DeepEqual(covers, want) {
		t.Fatalf("Covers(%v), got: %v, want: %v", item, covers, want)
	}

	// ###
	item = uintInterval{3, 7}
	want = []uintInterval{{1, 8}, {1, 7}, {2, 8}, {2, 7}}
	covers = tree1.Covers(item)

	if !reflect.DeepEqual(covers, want) {
		t.Fatalf("Covers(%v), got: %v, want: %v", item, covers, want)
	}

	// ###
	item = uintInterval{3, 11}
	want = nil
	covers = tree1.Covers(item)

	if !reflect.DeepEqual(covers, want) {
		t.Fatalf("Covers(%v), got: %+v, want: %+v", item, covers, want)
	}
}

func TestIntersects(t *testing.T) {
	t.Parallel()

	tree1 := interval.NewTree(cmpUintInterval, ps...)

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

	testcases := []struct {
		in   uintInterval
		want bool
	}{
		{
			in:   uintInterval{0, 1},
			want: true,
		},
		{
			in:   uintInterval{1, 1},
			want: true,
		},
		{
			in:   uintInterval{7, 7},
			want: true,
		},
		{
			in:   uintInterval{9, 17},
			want: true,
		},
		{
			in:   uintInterval{10, 12},
			want: false,
		},
	}
	for _, tc := range testcases {
		tc := tc
		t.Run("", func(t *testing.T) {
			t.Parallel()
			if ok := tree1.Intersects(tc.in); ok != tc.want {
				t.Errorf("Intersects(%v) = %v, want: %v", tc.in, ok, tc.want)
			}
		})
	}
}

func TestIntersections(t *testing.T) {
	t.Parallel()

	tree1 := interval.NewTree(cmpUintInterval, ps...)
	var want []uintInterval

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

	item := uintInterval{7, 7}
	want = []uintInterval{{1, 8}, {1, 7}, {2, 8}, {2, 7}, {4, 8}, {6, 7}, {7, 9}}
	intersections := tree1.Intersections(item)

	if !reflect.DeepEqual(intersections, want) {
		t.Fatalf("Intersections(%v), got: %v, want: %v", item, intersections, want)
	}

	// ###
	item = uintInterval{8, 10}
	want = []uintInterval{{1, 8}, {2, 8}, {4, 8}, {7, 9}}
	intersections = tree1.Intersections(item)

	if !reflect.DeepEqual(intersections, want) {
		t.Fatalf("Intersections(%v), got: %v, want: %v", item, intersections, want)
	}

	// ###
	item = uintInterval{10, 15}
	want = nil
	intersections = tree1.Intersections(item)

	if !reflect.DeepEqual(intersections, want) {
		t.Fatalf("Intersections(%v), got: %+v, want: %+v", item, intersections, want)
	}
}

func TestPrecedes(t *testing.T) {
	t.Parallel()

	tree1 := interval.NewTree(cmpUintInterval, ps...)
	var want []uintInterval

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

	item := uintInterval{7, 7}
	want = []uintInterval{{0, 6}, {0, 5}, {1, 5}, {1, 4}}
	precedes := tree1.Precedes(item)

	if !reflect.DeepEqual(precedes, want) {
		t.Fatalf("Precedes(%v), got: %v, want: %v", item, precedes, want)
	}

	// ###
	item = uintInterval{5, 10}
	want = []uintInterval{{1, 4}}
	precedes = tree1.Precedes(item)

	if !reflect.DeepEqual(precedes, want) {
		t.Fatalf("Precedes(%v), got: %v, want: %v", item, precedes, want)
	}

	// ###
	item = uintInterval{0, 9}
	want = nil
	precedes = tree1.Precedes(item)

	if !reflect.DeepEqual(precedes, want) {
		t.Fatalf("Precedes(%v), got: %+v, want: %+v", item, precedes, want)
	}
}

func TestPrecededBy(t *testing.T) {
	t.Parallel()

	tree1 := interval.NewTree(cmpUintInterval, ps...)
	var want []uintInterval

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

	item := uintInterval{4, 4}
	want = []uintInterval{{6, 7}, {7, 9}}
	precedes := tree1.PrecededBy(item)

	if !reflect.DeepEqual(precedes, want) {
		t.Fatalf("PrecededBy(%v), got: %v, want: %v", item, precedes, want)
	}

	// ###
	item = uintInterval{1, 2}
	want = []uintInterval{{4, 8}, {6, 7}, {7, 9}}
	precedes = tree1.PrecededBy(item)

	if !reflect.DeepEqual(precedes, want) {
		t.Fatalf("PrecededBy(%v), got: %v, want: %v", item, precedes, want)
	}

	// ###
	item = uintInterval{0, 7}
	want = nil
	precedes = tree1.PrecededBy(item)

	if !reflect.DeepEqual(precedes, want) {
		t.Fatalf("PrecededBy(%v), got: %+v, want: %+v", item, precedes, want)
	}
}

func TestVisit(t *testing.T) {
	t.Parallel()
	tree1 := interval.NewTree(cmpUintInterval, ps...)

	var collect []uintInterval
	want := 4
	tree1.Visit(tree1.Min(), tree1.Max(), func(item uintInterval) bool {
		collect = append(collect, item)
		return len(collect) != want
	})

	if len(collect) != want {
		t.Fatalf("Visit() ascending, want to stop after %v visits, got: %v, %v", want, len(collect), collect)
	}

	collect = nil
	want = 9
	tree1.Visit(tree1.Max(), tree1.Min(), func(item uintInterval) bool {
		collect = append(collect, item)
		return true
	})

	want, _, _, _ = tree1.Statistics()
	if len(collect) != want {
		t.Fatalf("Visit() descending, want: %d  got: %v, %v", want, len(collect), collect)
	}

	collect = nil
	want = 2
	tree1.Visit(tree1.Max(), tree1.Min(), func(item uintInterval) bool {
		collect = append(collect, item)
		return len(collect) != want
	})
}

func TestMinMax(t *testing.T) {
	t.Parallel()
	tree1 := interval.NewTree(cmpUintInterval, ps...)
	want := uintInterval{0, 6}
	if tree1.Min() != want {
		t.Fatalf("Min(), want: %v, got: %v", want, tree1.Min())
	}

	want = uintInterval{7, 9}
	if tree1.Max() != want {
		t.Fatalf("Max(), want: %v, got: %v", want, tree1.Max())
	}
}

func TestUnion(t *testing.T) {
	t.Parallel()
	tree1 := interval.NewTree(cmpUintInterval)

	for i := range ps {
		b := tree1.InsertImmutable(ps[i])
		tree1 = tree1.UnionImmutable(b, false)
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
		b := tree1.InsertImmutable(ps[i])
		tree1 = tree1.UnionImmutable(b, true)
	}

	if tree1.String() != asStr {
		t.Errorf("String()\nwant:\n%sgot:\n%s", asStr, tree1.String())
	}

	ps2 := []uintInterval{
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

	tree2 := tree1.InsertImmutable(ps2...)
	tree1.Union(tree2, false)

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
		n := n
		count := strconv.Itoa(n)
		t.Run(count, func(t *testing.T) {
			t.Parallel()
			tree1 := interval.NewTree(cmpUintInterval, genUintIvals(n)...)

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
	tree1 := interval.NewTree(cmpUintInterval, ps...)

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
	tree1 := interval.NewTree(cmpUintInterval, genUintIvals(10_000)...)

	n := 100
	for i := 0; i < n; i++ {
		probe := genUintIvals(10_000)[0]

		t.Run(probe.String(), func(t *testing.T) {
			t.Parallel()
			tree1 := tree1.InsertImmutable(probe)

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
	tree1 := interval.NewTree(cmpUintInterval, genUintIvals(10_000)...)

	n := 100
	for i := 0; i < n; i++ {
		probe := genUintIvals(1)[0]

		t.Run(probe.String(), func(t *testing.T) {
			t.Parallel()
			tree1 := tree1.InsertImmutable(probe)
			var ok bool

			if _, ok = tree1.Find(probe); !ok {
				t.Fatalf("inserted item not found in tree: %v", probe)
			}

			if tree1, ok = tree1.DeleteImmutable(probe); !ok {
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
