package interval_test

import (
	"fmt"

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

func ExampleInterface() {
	tree := interval.NewTree(periods...)
	fmt.Println(tree)

	// Output:
	// ▼
	// └─ 2...9
	//    ├─ 3...5
	//    │  └─ 3...4
	//    └─ 7...9
}

func ExampleTree_Supersets() {
	tree := interval.NewTree(periods...)
	item := period.Ival{3, 4}
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
	tree := interval.NewTree(periods...)
	item := period.Ival{3, 10}
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
