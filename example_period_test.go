package interval_test

import (
	"fmt"
	"os"

	"github.com/gaissmai/interval"
)

// example interval
type Ival [2]uint

// example data
var periods = []Ival{
	{3, 4},
	{2, 9},
	{7, 9},
	{3, 5},
}

// fmt.Stringer for formattting, not required
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

// cmp function for type Ival
func compareIval(p, q Ival) (ll, rr, lr, rl int) {
	return cmp(p[0], q[0]), cmp(p[1], q[1]), cmp(p[0], q[1]), cmp(p[1], q[0])
}

func ExampleNewTree() {
	tree1 := interval.NewTree(compareIval, periods...)

	tree1.Fprint(os.Stdout)
	// Output:
	// ▼
	// └─ 2...9
	//    ├─ 3...5
	//    │  └─ 3...4
	//    └─ 7...9
}

func ExampleTree_Max() {
	tree1 := interval.NewTree(compareIval, periods...)
	tree1.Fprint(os.Stdout)

	fmt.Println("\nInterval with max value in tree:")
	fmt.Println(tree1.Max())

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

func ExampleTree_Covers() {
	tree1 := interval.NewTree(compareIval, periods...)
	tree1.Fprint(os.Stdout)

	item := Ival{3, 4}
	fmt.Printf("\nCovers for item: %v\n", item)
	for _, p := range tree1.Covers(item) {
		fmt.Println(p)
	}

	// Output:
	// ▼
	// └─ 2...9
	//    ├─ 3...5
	//    │  └─ 3...4
	//    └─ 7...9
	//
	// Covers for item: 3...4
	// 2...9
	// 3...5
	// 3...4
}

func ExampleTree_CoveredBy() {
	tree1 := interval.NewTree(compareIval, periods...)
	tree1.Fprint(os.Stdout)

	item := Ival{3, 10}
	fmt.Printf("\nCoveredBy item: %v\n", item)
	for _, p := range tree1.CoveredBy(item) {
		fmt.Println(p)
	}

	// Output:
	// ▼
	// └─ 2...9
	//    ├─ 3...5
	//    │  └─ 3...4
	//    └─ 7...9
	//
	// CoveredBy item: 3...10
	// 3...5
	// 3...4
	// 7...9
}

func ExampleTree_Precedes_period() {
	tree1 := interval.NewTree(compareIval, periods...)
	tree1.Fprint(os.Stdout)

	item := Ival{6, 6}
	fmt.Printf("\nPrecedes item: %v\n", item)
	for _, p := range tree1.Precedes(item) {
		fmt.Println(p)
	}

	// Output:
	// ▼
	// └─ 2...9
	//    ├─ 3...5
	//    │  └─ 3...4
	//    └─ 7...9
	//
	// Precedes item: 6...6
	// 3...5
	// 3...4
}

func ExampleTree_Visit() {
	tree1 := interval.NewTree(compareIval, periods...)
	fmt.Println("parent/child printing")
	tree1.Fprint(os.Stdout)

	start := Ival{3, 5}
	stop := Ival{7, 9}
	visitFn := func(item Ival) bool {
		fmt.Printf("%v\n", item)
		return true
	}

	fmt.Println("visit ascending")
	tree1.Visit(start, stop, visitFn)

	fmt.Println("visit descending")
	tree1.Visit(stop, start, visitFn)

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
