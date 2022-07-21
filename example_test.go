package interval_test

import (
	"fmt"

	"github.com/gaissmai/interval"
)

// simple test interval
type ival struct {
	lo, hi int
}

func cmp(a, b int) int {
	if a == b {
		return 0
	}
	if a < b {
		return -1
	}
	return 1
}

// implementing interval.Interface

func (a ival) CompareFirst(b ival) int { return cmp(a.lo, b.lo) }
func (a ival) CompareLast(b ival) int  { return cmp(a.hi, b.hi) }

var ivals = []ival{
	{3, 4},
	{2, 9},
	{7, 9},
	{3, 5},
}

func (i ival) String() string {
	return fmt.Sprintf("%d...%d", i.lo, i.hi)
}

func Example() {
	tree := interval.NewTree(ivals)
	fmt.Println(tree.String())

	// Output:
	// ▼
	// └─ 2...9
	//    ├─ 3...5
	//    │  └─ 3...4
	//    └─ 7...9
}
