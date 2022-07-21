package interval_test

import (
	"fmt"

	"github.com/gaissmai/interval"
)

func cmp(a, b int) int {
	if a == b {
		return 0
	}
	if a < b {
		return -1
	}
	return 1
}

// simple test interval
type simple struct {
	lo, hi int
}

// implementing interval.Interface

func (a simple) CompareFirst(b simple) int { return cmp(a.lo, b.lo) }
func (a simple) CompareLast(b simple) int  { return cmp(a.hi, b.hi) }

func ExampleSort() {
	ivals := []simple{
		{3, 4},
		{2, 9},
		{7, 9},
		{3, 5},
	}

	// sort in place
	interval.Sort(ivals)

	for _, iv := range ivals {
		fmt.Println(iv)
	}

	// Output:
	// {2 9}
	// {3 5}
	// {3 4}
	// {7 9}
}
