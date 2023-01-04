package interval

import (
	"fmt"
	"io"
	"math"
)

type traverseOrder uint8

const (
	preorder traverseOrder = iota
	inorder
	reverse
)

// compare is a wrapper for CompareLower, CompareUpper with added functionality for superset sorting
func compare[T Interface[T]](a, b T) int {
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
func covers[T Interface[T]](a, b T) bool {
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
func (t *Tree[T]) traverse(order traverseOrder, depth int, visitFn func(n *Tree[T], depth int) bool) {
	if t == nil {
		return
	}

	switch order {
	case inorder:
		// left, do-it, right
		t.left.traverse(order, depth+1, visitFn)
		if !visitFn(t, depth) {
			return
		}
		t.right.traverse(order, depth+1, visitFn)
	case preorder:
		// do-it, left, right
		if !visitFn(t, depth) {
			return
		}
		t.left.traverse(order, depth+1, visitFn)
		t.right.traverse(order, depth+1, visitFn)
	case reverse:
		// right, do-it, left
		t.right.traverse(order, depth+1, visitFn)
		if !visitFn(t, depth) {
			return
		}
		t.left.traverse(order, depth+1, visitFn)
	}
}

// Fprint writes a hierarchical tree diagram of the ordered intervals to w.
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
func (t *Tree[T]) Fprint(w io.Writer) error {
	// pcm = parent-child-mapping
	var pcm parentChildsMap[T]

	// init map
	pcm.pcMap = make(map[*Tree[T]][]*Tree[T])

	pcm = t.buildParentChildsMap(pcm)

	if len(pcm.pcMap) == 0 {
		return nil
	}

	// start symbol
	if _, err := fmt.Fprint(w, "▼\n"); err != nil {
		return err
	}

	// start recursion with root and empty padding
	return walkAndStringify(w, pcm, nil, "")
}

func walkAndStringify[T Interface[T]](w io.Writer, pcm parentChildsMap[T], parent *Tree[T], pad string) error {
	// the prefix (pad + glyphe) is already printed on the line on upper level
	if parent != nil {
		if _, err := fmt.Fprintf(w, "%v\n", parent.item); err != nil {
			return err
		}
	}

	glyphe := "├─ "
	spacer := "│  "

	// dereference child-slice for clearer code
	childs := pcm.pcMap[parent]

	// for all childs do, but ...
	for i := range childs {
		// ... treat last child special
		if i == len(childs)-1 {
			glyphe = "└─ "
			spacer = "   "
		}
		// print prefix for next item
		if _, err := fmt.Fprint(w, pad+glyphe); err != nil {
			return err
		}

		// recdescent down
		if err := walkAndStringify(w, pcm, childs[i], pad+spacer); err != nil {
			return err
		}
	}

	return nil
}

// FprintBST writes a horizontal tree diagram of the binary search tree (BST) to w.
//
// Note: This is for debugging purposes only during development in semver
// 0.x.y. In future versions this will be removed without increasing the main
// semantic version, so please do not rely on it for now.
//
// e.g. with left/right, item and [height:size:priority]
//
//  R 0...5 [h:6|s:11|p:0.9405]
//  ├─l 0...6 [h:1|s:1|p:0.6047]
//  └─r 1...4 [h:5|s:9|p:0.6868]
//      ├─l 1...8 [h:3|s:3|p:0.6646]
//      │   └─r 1...7 [h:2|s:2|p:0.4377]
//      │       └─r 1...5 [h:1|s:1|p:0.4246]
//      └─r 7...9 [h:4|s:5|p:0.5152]
//          └─l 6...7 [h:3|s:4|p:0.3009]
//              └─l 2...7 [h:2|s:3|p:0.1565]
//                  ├─l 2...8 [h:1|s:1|p:0.06564]
//                  └─r 4...8 [h:1|s:1|p:0.09697]
//
//
func (t *Tree[T]) FprintBST(w io.Writer) error {
	if t == nil {
		return nil
	}

	if _, err := fmt.Fprint(w, "R "); err != nil {
		return err
	}

	// start recursion with empty padding
	return t.preorderStringify(w, "")
}

// preorderStringify, traverse the tree, stringify the nodes in preorder
func (t *Tree[T]) preorderStringify(w io.Writer, pad string) error {
	// stringify this node
	if _, err := fmt.Fprintf(w, "%v [p:%.4g] [p:%p|%p|%p]\n", t.item, t.prio, t, t.left, t.right); err != nil {
		return err
	}

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
		if _, err := fmt.Fprint(w, pad+glyphe); err != nil {
			return err
		}
		if err := t.left.preorderStringify(w, pad+spacer); err != nil {
			return err
		}
	}

	// right wing
	if t.right != nil {
		glyphe = "└─r "
		spacer = "    "
		if _, err := fmt.Fprint(w, pad+glyphe); err != nil {
			return err
		}
		if err := t.right.preorderStringify(w, pad+spacer); err != nil {
			return err
		}
	}

	return nil
}

// parentChildsMap, needed for interval tree printing, this is not BST printing!
//
// randomly balanced BST tree printed
//
//  R 0...5            [priority: 0.9405]
//  ├─l 0...6            [priority: 0.6047]
//  └─r 1...4            [priority: 0.6868]
//      ├─l 1...8            [priority: 0.6646]
//      │   └─r 1...7            [priority: 0.4377]
//      │       └─r 1...5            [priority: 0.4246]
//      └─r 7...9            [priority: 0.5152]
//          └─l 6...7            [priority: 0.3009]
//              └─l 2...7            [priority: 0.1565]
//                  ├─l 2...8            [priority: 0.06564]
//                  └─r 4...8            [priority: 0.09697]
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
type parentChildsMap[T Interface[T]] struct {
	pcMap map[*Tree[T]][]*Tree[T] // parent -> []child map
	stack []*Tree[T]              // just needed for the algo
}

// buildParentChildsMap, in-order traversal
func (t *Tree[T]) buildParentChildsMap(pcm parentChildsMap[T]) parentChildsMap[T] {
	if t == nil {
		return pcm
	}

	// in-order traversal, left tree
	pcm = t.left.buildParentChildsMap(pcm)

	// detect parent-child-mapping for this node
	pcm = t.pcmForNode(pcm)

	// in-order traversal, right tree
	return t.right.buildParentChildsMap(pcm)
}

// pcmForNode, find parent in stack, remove items from stack, put this item on stack.
func (t *Tree[T]) pcmForNode(pcm parentChildsMap[T]) parentChildsMap[T] {
	// if this item is covered by a prev item on stack
	for j := len(pcm.stack) - 1; j >= 0; j-- {

		that := pcm.stack[j]
		if covers(that.item, t.item) {
			// item in node j is parent to item
			pcm.pcMap[that] = append(pcm.pcMap[that], t)
			break
		}

		// Remember: sort order of intervals is lower-left, superset to the left:
		// if this item wasn't covered by j, remove node at j from stack
		pcm.stack = pcm.stack[:j]
	}

	// stack is emptied, no item on stack covers current item
	if len(pcm.stack) == 0 {
		// parent is root
		pcm.pcMap[nil] = append(pcm.pcMap[nil], t)
	}

	// put current neode on stack for next node
	pcm.stack = append(pcm.stack, t)

	return pcm
}

// Statistics, returns the maxDepth, average and standard deviation of the nodes.
//
// Note: This is for debugging purposes only during development in semver
// 0.x.y. In future versions this will be removed without increasing the main
// semantic version, so please do not rely on it for now.
//
func (t *Tree[T]) Statistics() (maxDepth int, average, deviation float64) {
	// key is depth, value is the sum of nodes with this depth
	depths := make(map[int]int)

	// get the depths
	t.traverse(inorder, 0, func(t *Tree[T], depth int) bool {
		depths[depth] += 1
		return true
	})

	var weightedSum, sum int
	for k, v := range depths {
		weightedSum += k * v
		sum += v
		if k > maxDepth {
			maxDepth = k
		}
	}

	average = float64(weightedSum) / float64(sum)

	var variance float64
	for k := range depths {
		variance += math.Pow(float64(k)-average, 2.0)
	}
	variance = variance / float64(sum)
	deviation = math.Sqrt(variance)

	return maxDepth, average, deviation
}

// Min returns the min item in tree.
func (t *Tree[T]) Min() (min T) {
	if t == nil {
		return
	}

	for t.left != nil {
		t = t.left
	}
	return t.item
}

// Max returns the node with max item in tree.
func (t *Tree[T]) Max() (max T) {
	if t == nil {
		return
	}

	for t.right != nil {
		t = t.right
	}
	return t.item
}

// Size returns the number of items in tree.
func (t *Tree[T]) Size() int {
	var size int
	t.traverse(inorder, 0, func(t *Tree[T], _ int) bool {
		size++
		return true
	})
	return size
}

// Visit traverses the tree with item >= start until item <= stop in ascending order,
// if start > stop, the order is reversed.
//
// The visit function is called for each item. The traversion stops prematurely if the visit function returns false.
func (t *Tree[T]) Visit(start, stop T, visitFn func(t T) bool) {
	if t == nil {
		return
	}

	order := inorder
	if compare(start, stop) > 0 {
		start, stop = stop, start
		order = reverse
	}

	// treaps are really cool datastructures!
	_, mid1, r := t.split(start)
	l, mid2, _ := r.split(stop)

	span := join(mid1, join(l, mid2))

	span.traverse(order, 0, func(t *Tree[T], dummy int) bool {
		return visitFn(t.item)
	})
}

// Clone the tree.
func (t *Tree[T]) Clone() *Tree[T] {
	if t == nil {
		return t
	}

	t.left = t.left.Clone()
	t.right = t.right.Clone()

	return t.copyNode()
}
