package interval_test

import (
	"fmt"

	"github.com/gaissmai/interval"
)

// little helper, compare two ints
func cmp(a, b int) int {
	if a == b {
		return 0
	}
	if a < b {
		return -1
	}
	return 1
}

// example interval
type period struct {
	start int
	stop  int
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

func Example() {
	tree := interval.NewTree(periods)
	fmt.Println(tree)

	// Output:
	// ▼
	// └─ 2...9
	//    ├─ 3...5
	//    │  └─ 3...4
	//    └─ 7...9
}
