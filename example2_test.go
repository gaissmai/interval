package interval_test

import (
	"fmt"

	"github.com/gaissmai/interval"
)

func ExampleSort() {
	// sort in place
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

func ExampleTree_Supersets() {
	tree := interval.NewTree(ivals)
	fmt.Println(tree)
	item := ival{3, 4}

	supersets := tree.Supersets(item)
	interval.Sort(supersets)

	fmt.Printf("Supersets for item: %v\n", item)
	for _, iv := range supersets {
		fmt.Println(iv)
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
	tree := interval.NewTree(ivals)
	fmt.Println(tree)
	item := ival{3, 10}

	subsets := tree.Subsets(item)
	interval.Sort(subsets)

	fmt.Printf("Subsets for item: %v\n", item)
	for _, iv := range subsets {
		fmt.Println(iv)
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
