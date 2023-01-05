package interval_test

import (
	"math/rand"
	"reflect"
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
		a := rand.Intn(n)
		b := rand.Intn(n)
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
	zeroTree.Fprint(w)

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

func TestRandom(t *testing.T) {
	t.Parallel()
	is := generateIvals(1000)
	tree := treap.Insert(is...)

	rand.Shuffle(len(is), func(i, j int) { is[i], is[j] = is[j], is[i] })

	for _, item := range is {
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
	}
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
		tree = tree.Union(b, false)
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
		tree = tree.Union(b, true)
	}

	w.Reset()
	tree.Fprint(w)
	if w.String() != asStr {
		t.Errorf("Fprint()\nwant:\n%sgot:\n%s", asStr, w.String())
	}
}
