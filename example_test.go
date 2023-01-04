package interval_test

import (
	"fmt"
	"os"

	"github.com/gaissmai/interval"
	"github.com/gaissmai/interval/internal/period"
)

// example data
var periods = []period.Ival{
	{3, 4},
	{2, 9},
	{7, 9},
	{3, 5},
}

var tree *interval.Tree[period.Ival]

func ExampleInterface() {
	tree = tree.Insert(periods...)
	tree.Fprint(os.Stdout)

	// Output:
	// ▼
	// └─ 2...9
	//    ├─ 3...5
	//    │  └─ 3...4
	//    └─ 7...9
}

func ExampleTree_Max() {
	tree = tree.Insert(periods...)
	tree.Fprint(os.Stdout)

	fmt.Println("\nInterval with max lower value in tree:")
	fmt.Println(tree.Max())

	// Output:
	// ▼
	// └─ 2...9
	//    ├─ 3...5
	//    │  └─ 3...4
	//    └─ 7...9
	//
	//Interval with max lower value in tree:
	//7...9
}

func ExampleTree_Supersets() {
	tree = tree.Insert(periods...)
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
	tree = tree.Insert(periods...)
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
	tree = tree.Insert(periods...)
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
