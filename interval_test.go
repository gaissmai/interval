package interval_test

import (
	"math/rand"
	"strings"
	"testing"

	"github.com/gaissmai/interval"
	"github.com/gaissmai/interval/internal/period"
)

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

var treap *interval.Tree[period.Ival]

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
	var zero *interval.Tree[period.Ival]

	w := new(strings.Builder)
	zero.Fprint(w)
	if w.String() != "" {
		t.Errorf("tree.Write(w) = %v, want \"\"", w.String())
	}

	if _, ok := zero.Shortest(period.Ival{}); ok {
		t.Errorf("tree.Shortest(), got: %v, want: false", ok)
	}

	if _, ok := zero.Largest(period.Ival{}); ok {
		t.Errorf("tree.Largest(), got: %v, want: false", ok)
	}

	if s := zero.Subsets(period.Ival{}); s != nil {
		t.Errorf("tree.Subsets(), got: %v, want: nil", s)
	}

	if s := zero.Supersets(period.Ival{}); s != nil {
		t.Errorf("tree.Supersets(), got: %v, want: nil", s)
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

func TestTreeLookup(t *testing.T) {
	t.Parallel()
	is := []period.Ival{
		{1, 100},
		{45, 60},
	}

	tree := treap.Insert(is...)

	item := period.Ival{0, 6}
	if got, ok := tree.Shortest(item); ok {
		t.Errorf("Shortest(%v) = %v, want %v", item, got, !ok)
	}

	item = period.Ival{47, 62}
	if got, _ := tree.Shortest(item); got != is[0] {
		t.Errorf("Shortest(%v) = %v, want %v", item, got, is[0])
	}

	item = period.Ival{45, 60}
	if got, _ := tree.Shortest(item); got != is[1] {
		t.Errorf("Shortest(%v) = %v, want %v", item, got, is[1])
	}
}

func TestTreeSuperset(t *testing.T) {
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

func TestTreeRandom(t *testing.T) {
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
