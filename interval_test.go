package interval_test

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/gaissmai/interval"
)

func generateIvals(n int) []period {
	is := make([]period, n)
	for i := 0; i < n; i++ {
		a := rand.Intn(n)
		b := rand.Intn(n)
		if a > b {
			a, b = b, a
		}
		is = append(is, period{a, b})
	}
	return is
}

func TestTreeNil(t *testing.T) {
	tree := interval.NewTree[period](nil)

	if s := tree.String(); s != "" {
		t.Errorf("tree.String() = %v, want \"\"", s)
	}

	if s := tree.Size(); s != 0 {
		t.Errorf("tree.Size() = %v, want 0", s)
	}

	if _, ok := tree.Shortest(period{}); ok {
		t.Errorf("tree.Shortest(), got: %v, want: false", ok)
	}

	if _, ok := tree.Largest(period{}); ok {
		t.Errorf("tree.Largest(), got: %v, want: false", ok)
	}

	if s := tree.Subsets(period{}); s != nil {
		t.Errorf("tree.Subsets(), got: %v, want: nil", s)
	}

	if s := tree.Supersets(period{}); s != nil {
		t.Errorf("tree.Supersets(), got: %v, want: nil", s)
	}
}

func TestTreeWithDups(t *testing.T) {
	tree := interval.NewTree([]period{{0, 100}, {41, 102}, {42, 67}, {42, 67}, {48, 50}, {3, 13}})
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
	is := []period{
		{1, 100},
		{45, 60},
	}

	tree := interval.NewTree(is)

	item := period{0, 6}
	if got, ok := tree.Shortest(item); ok {
		t.Errorf("Shortest(%v) = %v, want %v", item, got, !ok)
	}

	item = period{47, 62}
	if got, _ := tree.Shortest(item); got != is[0] {
		t.Errorf("Shortest(%v) = %v, want %v", item, got, is[0])
	}

	item = period{45, 60}
	if got, _ := tree.Shortest(item); got != is[1] {
		t.Errorf("Shortest(%v) = %v, want %v", item, got, is[1])
	}
}

func TestTreeSuperset(t *testing.T) {
	is := []period{
		{1, 100},
		{45, 120},
		{46, 80},
	}

	tree := interval.NewTree(is)

	item := period{0, 6}
	if got, ok := tree.Largest(item); ok {
		t.Errorf("Largest(%v) = %v, want %v", item, got, !ok)
	}

	item = period{99, 200}
	if got, ok := tree.Largest(item); ok {
		t.Errorf("Largest(%v) = %v, want %v", item, got, !ok)
	}

	item = period{1, 100}
	if got, ok := tree.Largest(item); got != is[0] || !ok {
		t.Errorf("Largest(%v) = %v, want %v", item, got, is[0])
	}

	item = period{46, 80}
	if got, ok := tree.Largest(item); got != is[0] || !ok {
		t.Errorf("Largest(%v) = %v, want %v", item, got, is[0])
	}

	item = period{47, 62}
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
			shortest  period
			largest   period
			subsets   []period
			supersets []period
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

func ExampleTree_Supersets() {
	periods := []period{
		{3, 4},
		{2, 9},
		{7, 9},
		{3, 5},
	}

	tree := interval.NewTree(periods)
	item := period{3, 4}
	supersets := tree.Supersets(item)

	fmt.Println(tree)
	fmt.Printf("Supersets for item: %v\n", item)
	for _, p := range supersets {
		fmt.Println(p)
	}

	// Output:
	// ▼
	// └─ 2...9
	//    ├─ 3...5
	//    │  └─ 3...4
	//    └─ 7...9
	//
	// Supersets for item: 3...4
	// 2...9
	// 3...5
	// 3...4
}

func ExampleTree_Subsets() {
	periods := []period{
		{3, 4},
		{2, 9},
		{7, 9},
		{3, 5},
	}

	tree := interval.NewTree(periods)
	item := period{3, 10}
	subsets := tree.Subsets(item)

	fmt.Println(tree)
	fmt.Printf("Subsets for item: %v\n", item)
	for _, p := range subsets {
		fmt.Println(p)
	}

	// Output:
	// ▼
	// └─ 2...9
	//    ├─ 3...5
	//    │  └─ 3...4
	//    └─ 7...9
	//
	// Subsets for item: 3...10
	// 3...5
	// 3...4
	// 7...9
}
