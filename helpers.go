package interval

import (
	"fmt"
	"io"
	"strings"
)

type traverseOrder int

const (
	preorder traverseOrder = iota
	inorder
	reverse
)

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// compare is a wrapper for CompareLower, CompareUpper with added functionality for superset sorting
func compare[T Interval[T]](a, b T) int {
	cmpLower := a.CompareLower(b)
	cmpUpper := a.CompareUpper(b)

	// lower interval value is the primary sort key
	if cmpLower != 0 {
		return cmpLower
	}

	// if lower interval values are equal, sort supersets to the left
	if cmpUpper != 0 {
		return -cmpUpper
	}

	// both, lower and upper are equal
	return 0
}

// covers reports whether a truly covers b (not equal).
func covers[T Interval[T]](a, b T) bool {
	cmpLower := a.CompareLower(b)
	cmpUpper := a.CompareUpper(b)

	// equal
	if cmpLower == 0 && cmpUpper == 0 {
		return false
	}

	return cmpLower <= 0 && cmpUpper >= 0
}

// traverse the BST in some order, call the visitor function for each node.
// Prematurely stop traversion if visitor function returns false.
func (t *Tree[T]) traverse(order traverseOrder, visitFn Visitor[T]) {
	switch order {
	case inorder:
		// left
		if t.left != nil {
			t.left.traverse(order, visitFn)
		}
		// this
		if !visitFn(t) {
			return
		}
		// right
		if t.right != nil {
			t.right.traverse(order, visitFn)
		}
	case preorder:
		// this
		if !visitFn(t) {
			return
		}
		// left
		if t.left != nil {
			t.left.traverse(order, visitFn)
		}
		// right
		if t.right != nil {
			t.right.traverse(order, visitFn)
		}
	case reverse:
		// right
		if t.right != nil {
			t.right.traverse(order, visitFn)
		}
		// this
		if !visitFn(t) {
			return
		}
		// left
		if t.left != nil {
			t.left.traverse(order, visitFn)
		}
	}
}

// String returns the ordered tree as a directory graph.
//
// example: IP CIDRs as intervals
//
//     ▼
//     ├─ 0.0.0.0/0
//     │  ├─ 10.0.0.0/8
//     │  │  ├─ 10.0.0.0/24
//     │  │  └─ 10.0.1.0/24
//     │  ├─ 127.0.0.0/8
//     │  │  └─ 127.0.0.1/32
//     │  ├─ 169.254.0.0/16
//     │  ├─ 172.16.0.0/12
//     │  └─ 192.168.0.0/16
//     │     └─ 192.168.1.0/24
//     └─ ::/0
//        ├─ ::1/128
//        ├─ 2000::/3
//        │  └─ 2001:db8::/32
//        ├─ fc00::/7
//        ├─ fe80::/10
//        └─ ff00::/8
//
// If the interval items don't implement fmt.Stringer they are stringified with
// their default format %v.
//
func (t *Tree[T]) String() string {
	pcm := t.buildParentChildsMap()

	if len(pcm.pcMap) == 0 {
		return ""
	}

	w := new(strings.Builder)

	// start symbol
	w.WriteString("▼\n")

	// start recursion with root and empty padding
	walkAndStringify(pcm, nil, "", w)

	return w.String()
}

func walkAndStringify[T Interval[T]](pcm parentChildsMap[T], parent *Tree[T], pad string, w io.StringWriter) {
	// the prefix (pad + glyphe) is already printed on the line on upper level
	if parent != nil {
		w.WriteString(fmt.Sprintf("%v\n", parent.item)) //nolint:errcheck
	}

	glyphe := "├─ "
	spacer := "│  "

	// dereference child-slice for clearer code
	childs := pcm.pcMap[parent]

	// for all childs do, but ...
	for i, c := range childs {
		// ... treat last child special
		if i == len(childs)-1 {
			glyphe = "└─ "
			spacer = "   "
		}
		// print prefix for next item
		w.WriteString(pad + glyphe) //nolint:errcheck

		// recdescent down
		walkAndStringify(pcm, c, pad+spacer, w)
	}
}

// PrintBST, returns the string representation of the balanced BST.
//
// Note: This is for debugging purposes only during development in semver
// 0.x.y. In future versions this will be removed without increasing the main
// semantic version, so please do not rely on it for now.
//
// e.g.
//
//  R 0...5 [h:6|s:11]
//  ├─l 0...6 [h:1|s:1]
//  └─r 1...4 [h:5|s:9]
//      ├─l 1...8 [h:3|s:3]
//      │   └─r 1...7 [h:2|s:2]
//      │       └─r 1...5 [h:1|s:1]
//      └─r 7...9 [h:4|s:5]
//          └─l 6...7 [h:3|s:4]
//              └─l 2...7 [h:2|s:3]
//                  ├─l 2...8 [h:1|s:1]
//                  └─r 4...8 [h:1|s:1]
//
func (t *Tree[T]) PrintBST() string {
	if t == nil {
		return ""
	}

	w := new(strings.Builder)

	// start recursion with root and empty padding
	w.WriteString("R ")
	t.preorderStringify("", w)

	return w.String()
}

// preorderStringify, traverse the tree, stringify the nodes in preorder
func (t *Tree[T]) preorderStringify(pad string, w io.StringWriter) {
	// stringify this
	w.WriteString(fmt.Sprintf("%v [h:%d|s:%d|p:%.4g]\n", t.item, t.height, t.size, t.prio)) //nolint:errcheck

	// prepare glyphe, spacer and padding for next level
	var glyphe string
	var spacer string

	// left wing
	if t.left != nil {
		if t.right != nil {
			glyphe = "├─l "
			spacer = "│   "
		} else {
			glyphe = "└─l "
			spacer = "    "
		}
		w.WriteString(pad + glyphe) //nolint:errcheck
		t.left.preorderStringify(pad+spacer, w)
	}

	// right wing
	if t.right != nil {
		glyphe = "└─r "
		spacer = "    "
		w.WriteString(pad + glyphe) //nolint:errcheck
		t.right.preorderStringify(pad+spacer, w)
	}
}

// parentChildsMap, needed for interval tree printing, this is not BST printing!
//
// randomly balanced BST tree printed
//
//  R 0...5
//  ├─l 0...6
//  └─r 1...4
//      ├─l 1...8
//      │   └─r 1...7
//      │       └─r 1...5
//      └─r 7...9
//          └─l 6...7
//              └─l 2...7
//                  ├─l 2...8
//                  └─r 4...8
//
// Interval tree, parent->child relation printed
//  ▼
//  ├─ 0...6
//  │  └─ 0...5
//  ├─ 1...8
//  │  ├─ 1...7
//  │  │  └─ 1...5
//  │  │     └─ 1...4
//  │  └─ 2...8
//  │     ├─ 2...7
//  │     └─ 4...8
//  │        └─ 6...7
//  └─ 7...9
//
type parentChildsMap[T Interval[T]] struct {
	pcMap map[*Tree[T]][]*Tree[T] // parent -> []child map
	stack []*Tree[T]              // just needed for the algo
}

func (t *Tree[T]) buildParentChildsMap() parentChildsMap[T] {
	var pcm parentChildsMap[T]

	if t == nil {
		return pcm
	}

	pcm = parentChildsMap[T]{pcMap: make(map[*Tree[T]][]*Tree[T])}

	// this function is called in-order for every node
	visitFn := func(n *Tree[T]) bool {
		// if this item is covered by a prev item on stack
		for j := len(pcm.stack) - 1; j >= 0; j-- {

			that := pcm.stack[j]
			if covers(that.item, n.item) {
				// item in node j is parent to item
				pcm.pcMap[that] = append(pcm.pcMap[that], n)
				break
			}

			// Remember: sort order of intervals is lower-left, superset to the left:
			// if this item wasn't covered by j, remove node at j from stack
			pcm.stack = pcm.stack[:j]
		}

		// stack is emptied, no item on stack covers current item
		if len(pcm.stack) == 0 {
			// parent is root
			pcm.pcMap[nil] = append(pcm.pcMap[nil], n)
		}

		// put current item on stack für next round
		pcm.stack = append(pcm.stack, n)

		return true
	}

	t.traverse(inorder, visitFn)

	return pcm
}
