package interval_test

import (
	"fmt"
	"os"

	"github.com/gaissmai/interval"
)

// fmt.Stringer for formattting, not required for interval.Interface
func (p Ival) String() string {
	return fmt.Sprintf("%d...%d", p[0], p[1])
}

// little helper
func cmp(a, b uint) int {
	switch {
	case a == b:
		return 0
	case a < b:
		return -1
	}
	return 1
}

// implement interval.Interface
func (p Ival) Compare(q Ival) (ll, rr, lr, rl int) {
	return cmp(p[0], q[0]), cmp(p[1], q[1]), cmp(p[0], q[1]), cmp(p[1], q[0])
}

// example interval
type Ival [2]uint

// example data
var periods = []Ival{
	{3, 4},
	{2, 9},
	{7, 9},
	{3, 5},
}

func ExampleInterface_period() {
	tree := interval.NewTree(periods...)
	tree.Fprint(os.Stdout)

	// Output:
	// ▼
	// └─ 2...9
	//    ├─ 3...5
	//    │  └─ 3...4
	//    └─ 7...9
}

func ExampleTree_Max() {
	tree := interval.NewTree(periods...)
	tree.Fprint(os.Stdout)

	fmt.Println("\nInterval with max value in tree:")
	fmt.Println(tree.Max())

	// Output:
	// ▼
	// └─ 2...9
	//    ├─ 3...5
	//    │  └─ 3...4
	//    └─ 7...9
	//
	//Interval with max value in tree:
	//7...9
}

/*
func ExampleTree_Supersets() {
	tree := interval.NewTree(periods...)
	tree.Fprint(os.Stdout)

	item := period.Ival{3, 4}
	fmt.Printf("\nSupersets for item: %v\n", item)
	for _, p := range tree.Supersets(item) {
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
	tree := interval.NewTree(periods...)
	tree.Fprint(os.Stdout)

	item := period.Ival{3, 10}
	fmt.Printf("\nSubsets for item: %v\n", item)
	for _, p := range tree.Subsets(item) {
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

func ExampleTree_Visit() {
	tree := interval.NewTree(periods...)
	fmt.Println("parent/child printing")
	tree.Fprint(os.Stdout)

	start := period.Ival{3, 5}
	stop := period.Ival{7, 9}
	visitFn := func(item period.Ival) bool {
		fmt.Printf("%v\n", item)
		return true
	}

	fmt.Println("visit ascending")
	tree.Visit(start, stop, visitFn)

	fmt.Println("visit descending")
	tree.Visit(stop, start, visitFn)

	// Output:
	// parent/child printing
	// ▼
	// └─ 2...9
	//    ├─ 3...5
	//    │  └─ 3...4
	//    └─ 7...9
	// visit ascending
	// 3...5
	// 3...4
	// 7...9
	// visit descending
	// 7...9
	// 3...4
	// 3...5
}
*/
