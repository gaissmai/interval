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

func generateIvals(n int) []period.Ival {
	is := make([]period.Ival, n)
	for i := 0; i < n; i++ {
		a := rand.Intn(n)
		b := rand.Intn(n)
		if a > b {
			a, b = b, a
		}
		is = append(is, period.Ival{a, b})
	}
	return is
}

func TestTreeNullValue(t *testing.T) {
	t.Parallel()
	var tree *interval.Tree[period.Ival]

	w := new(strings.Builder)
	tree.Fprint(w)
	if w.String() != "" {
		t.Errorf("tree.Write(w) = %v, want \"\"", w.String())
	}

	if _, ok := tree.Shortest(period.Ival{}); ok {
		t.Errorf("tree.Shortest(), got: %v, want: false", ok)
	}

	if _, ok := tree.Largest(period.Ival{}); ok {
		t.Errorf("tree.Largest(), got: %v, want: false", ok)
	}

	if s := tree.Subsets(period.Ival{}); s != nil {
		t.Errorf("tree.Subsets(), got: %v, want: nil", s)
	}

	if s := tree.Supersets(period.Ival{}); s != nil {
		t.Errorf("tree.Supersets(), got: %v, want: nil", s)
	}
}

func TestTreeNil(t *testing.T) {
	t.Parallel()
	tree := interval.NewTree[period.Ival]()

	w := new(strings.Builder)
	tree.Fprint(w)
	if w.String() != "" {
		t.Errorf("tree.Write(w) = %v, want \"\"", w.String())
	}

	if _, ok := tree.Shortest(period.Ival{}); ok {
		t.Errorf("tree.Shortest(), got: %v, want: false", ok)
	}

	if _, ok := tree.Largest(period.Ival{}); ok {
		t.Errorf("tree.Largest(), got: %v, want: false", ok)
	}

	if s := tree.Subsets(period.Ival{}); s != nil {
		t.Errorf("tree.Subsets(), got: %v, want: nil", s)
	}

	if s := tree.Supersets(period.Ival{}); s != nil {
		t.Errorf("tree.Supersets(), got: %v, want: nil", s)
	}
}

func TestTreeWithDups(t *testing.T) {
	t.Parallel()
	tree := interval.NewTree([]period.Ival{{0, 100}, {41, 102}, {42, 67}, {42, 67}, {48, 50}, {3, 13}}...)
	if s := tree.Size(); s != 5 {
		t.Errorf("tree.Size() = %v, want 5", s)
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
		t.Errorf("tree.String()\nwant:\n%sgot:\n%s", asStr, w.String())
	}
}

func TestTreeLookup(t *testing.T) {
	t.Parallel()
	is := []period.Ival{
		{1, 100},
		{45, 60},
	}

	tree := interval.NewTree(is...)

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

	tree := interval.NewTree(is...)

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
	tree := interval.NewTree(is...)

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
	tree := interval.NewTree(ps...)
	want := period.Ival{0, 6}
	if tree.Min().Item() != want {
		t.Fatalf("Min(), want: %v, got: %v", want, tree.Min().Item())
	}

	want = period.Ival{7, 9}
	if tree.Max().Item() != want {
		t.Fatalf("Max(), want: %v, got: %v", want, tree.Max().Item())
	}
}
