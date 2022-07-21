package interval_test

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/gaissmai/interval"
)

// simple test interval
type ival struct {
	lo, hi int
}

// implementing interval.Interface

func (a ival) CompareFirst(b ival) int {
	if a.lo == b.lo {
		return 0
	}
	if a.lo < b.lo {
		return -1
	}
	return 1
}

func (a ival) CompareLast(b ival) int {
	if a.hi == b.hi {
		return 0
	}
	if a.hi < b.hi {
		return -1
	}
	return 1
}

// fmt.Stringer
func (a ival) String() string {
	return fmt.Sprintf("%d...%d", a.lo, a.hi)
}

func generateIvals(n int) []ival {
	is := make([]ival, n)
	for i := 0; i < n; i++ {
		a := rand.Intn(n)
		b := rand.Intn(n)
		if a > b {
			a, b = b, a
		}
		is = append(is, ival{a, b})
	}
	return is
}

func ExampleSort() {
	ivals := []ival{
		{2, 9},
		{3, 5},
		{3, 4},
		{7, 9},
	}
	interval.Sort(ivals)
	for _, iv := range ivals {
		fmt.Println(iv)
	}
	// Output:
	// 2...9
	// 3...5
	// 3...4
	// 7...9
}

func TestTreeNil(t *testing.T) {
	tree := interval.NewTree[ival](nil)

	if s := tree.String(); s != "" {
		t.Errorf("tree.String() = %v, want \"\"", s)
	}

	if s := tree.Size(); s != 0 {
		t.Errorf("tree.Size() = %v, want 0", s)
	}

	if _, ok := tree.Shortest(ival{}); ok {
		t.Errorf("tree.Shortest(), got: %v, want: false", ok)
	}

	if _, ok := tree.Largest(ival{}); ok {
		t.Errorf("tree.Largest(), got: %v, want: false", ok)
	}

	if s := tree.Subsets(ival{}); s != nil {
		t.Errorf("tree.Subsets(), got: %v, want: nil", s)
	}

	if s := tree.Supersets(ival{}); s != nil {
		t.Errorf("tree.Supersets(), got: %v, want: nil", s)
	}
}

func TestTreeWithDups(t *testing.T) {
	tree := interval.NewTree([]ival{{0, 100}, {41, 102}, {42, 67}, {42, 67}, {48, 50}, {3, 13}})
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
	if s := tree.String(); s != asStr {
		t.Errorf("tree.String()\nwant:\n%sgot:\n%s", asStr, s)
	}
}

func TestTreeLookup(t *testing.T) {
	is := []ival{
		{1, 100},
		{45, 60},
	}

	tree := interval.NewTree(is)

	item := ival{0, 6}
	if got, ok := tree.Shortest(item); ok {
		t.Errorf("Shortest(%v) = %v, want %v", item, got, !ok)
	}

	item = ival{47, 62}
	if got, _ := tree.Shortest(item); got != is[0] {
		t.Errorf("Shortest(%v) = %v, want %v", item, got, is[0])
	}

	item = ival{45, 60}
	if got, _ := tree.Shortest(item); got != is[1] {
		t.Errorf("Shortest(%v) = %v, want %v", item, got, is[1])
	}
}

func TestTreeSuperset(t *testing.T) {
	is := []ival{
		{1, 100},
		{45, 120},
		{46, 80},
	}

	tree := interval.NewTree(is)

	item := ival{0, 6}
	if got, ok := tree.Largest(item); ok {
		t.Errorf("Largest(%v) = %v, want %v", item, got, !ok)
	}

	item = ival{99, 200}
	if got, ok := tree.Largest(item); ok {
		t.Errorf("Largest(%v) = %v, want %v", item, got, !ok)
	}

	item = ival{1, 100}
	if got, ok := tree.Largest(item); got != is[0] || !ok {
		t.Errorf("Largest(%v) = %v, want %v", item, got, is[0])
	}

	item = ival{46, 80}
	if got, ok := tree.Largest(item); got != is[0] || !ok {
		t.Errorf("Largest(%v) = %v, want %v", item, got, is[0])
	}

	item = ival{47, 62}
	if got, ok := tree.Largest(item); got != is[0] || !ok {
		t.Errorf("Largest(%v) = %v, want %v", item, got, is[0])
	}
}

func TestTreeRandom(t *testing.T) {
	is := generateIvals(10)
	tree := interval.NewTree(is)

	rand.Shuffle(len(is), func(i, j int) { is[i], is[j] = is[j], is[i] })

	for _, item := range is {
		var (
			shortest  ival
			largest   ival
			subsets   []ival
			supersets []ival
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
		interval.Sort(subsets)
		if subsets[0] != shortest {
			t.Errorf("Subsets(%v).[0], want %v, got %v", item, shortest, subsets[0])
		}
		if supersets = tree.Supersets(item); supersets == nil {
			t.Errorf("Supersets(%v), got %v", item, supersets)
		}
		interval.Sort(supersets)
		if supersets[0] != largest {
			t.Errorf("Supersets(%v).[0], want %v, got %v", item, largest, supersets[0])
		}
	}
}
