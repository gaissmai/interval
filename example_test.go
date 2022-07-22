package interval_test

import (
	"fmt"

	"github.com/gaissmai/interval"
)

// example interval
type period struct {
	start int
	stop  int
}

// little helper, compare two ints
func cmp(a, b int) int {
	switch {
	case a == b:
		return 0
	case a < b:
		return -1
	}
	return 1
}

// implementing the interval.Interface
func (p period) CompareLower(q period) int { return cmp(p.start, q.start) }
func (p period) CompareUpper(q period) int { return cmp(p.stop, q.stop) }

// example data
var periods = []period{
	{3, 4},
	{2, 9},
	{7, 9},
	{3, 5},
}

// fmt.Stringer for formattting, not required for interval.Interface
func (p period) String() string {
	return fmt.Sprintf("%d...%d", p.start, p.stop)
}

func ExampleInterface() {
	tree := interval.NewTree(periods)
	fmt.Println(tree)

	// Output:
	// ▼
	// └─ 2...9
	//    ├─ 3...5
	//    │  └─ 3...4
	//    └─ 7...9
}

func ExampleTree_Supersets() {
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
