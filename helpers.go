package interval

import (
	"fmt"
	"io"
	"math"
	"strings"
)

type traverseOrder uint8

const (
	inorder traverseOrder = iota
	reverse
)

// traverse the BST in some order, call the visitor function for each node.
// Prematurely stop traversion if visitor function returns false.
func (t *Tree[T]) traverse(n *node[T], order traverseOrder, depth int, visitFn func(n *node[T], depth int) bool) bool {
	if n == nil {
		return true
	}

	switch order {
	case inorder:
		// left, do-it, right
		if !t.traverse(n.left, order, depth+1, visitFn) {
			return false
		}

		if !visitFn(n, depth) {
			return false
		}

		if !t.traverse(n.right, order, depth+1, visitFn) {
			return false
		}

		return true
	case reverse:
		// right, do-it, left
		if !t.traverse(n.right, order, depth+1, visitFn) {
			return false
		}

		if !visitFn(n, depth) {
			return false
		}

		if !t.traverse(n.left, order, depth+1, visitFn) {
			return false
		}

		return true
	default:
		panic("unreachable")
	}
}

// String returns a hierarchical tree diagram of the ordered intervals as string, just a wrapper for [Fprint].
func (t Tree[T]) String() string {
	w := new(strings.Builder)
	_ = t.Fprint(w)
	return w.String()
}

// Fprint writes an ordered interval tree diagram to w.
//
// The order from top to bottom is in ascending order of the left edges of the intervals
// and the subtree structure is determined by the intervals coverage.
//
// example: IP CIDRs as intervals
//
//	▼
//	├─ 0.0.0.0/0
//	│  ├─ 10.0.0.0/8
//	│  │  ├─ 10.0.0.0/24
//	│  │  └─ 10.0.1.0/24
//	│  ├─ 127.0.0.0/8
//	│  │  └─ 127.0.0.1/32
//	│  ├─ 169.254.0.0/16
//	│  ├─ 172.16.0.0/12
//	│  └─ 192.168.0.0/16
//	│     └─ 192.168.1.0/24
//	└─ ::/0
//	   ├─ ::1/128
//	   ├─ 2000::/3
//	   │  └─ 2001:db8::/32
//	   ├─ fc00::/7
//	   ├─ fe80::/10
//	   └─ ff00::/8
//
// If the interval items don't implement fmt.Stringer they are stringified with
// their default format %v.
func (t Tree[T]) Fprint(w io.Writer) error {
	// pcm = parent-child-mapping
	var pcm parentChildsMap[T]

	// init map
	pcm.pcMap = make(map[*node[T]][]*node[T])

	pcm = t.buildParentChildsMap(t.root, pcm)

	if len(pcm.pcMap) == 0 {
		return nil
	}

	// start symbol
	if _, err := fmt.Fprint(w, "▼\n"); err != nil {
		return err
	}

	// start recursion with nil parent and empty padding
	return t.hierarchyStringify(w, nil, pcm, "")
}

func (t *Tree[T]) hierarchyStringify(w io.Writer, n *node[T], pcm parentChildsMap[T], pad string) error {
	// the prefix (pad + glyphe) is already printed on the line on upper level
	if n != nil {
		if _, err := fmt.Fprintf(w, "%v\n", n.item); err != nil {
			return err
		}
	}

	glyphe := "├─ "
	spacer := "│  "

	// dereference child-slice for clearer code
	childs := pcm.pcMap[n]

	// for all childs do, but ...
	for i, child := range childs {
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
		if err := t.hierarchyStringify(w, child, pcm, pad+spacer); err != nil {
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
// e.g. with left/right, item priority [prio] and pointers [this|left|right]
//
//	R 0...5 [prio:0.9405] [0xc000024940|l:0xc000024140|r:0xc000024980]
//	├─l 0...6 [prio:0.6047] [0xc000024140|l:0x0|r:0x0]
//	└─r 1...4 [prio:0.6868] [0xc000024980|l:0xc000024440|r:0xc000024900]
//	    ├─l 1...8 [prio:0.6646] [0xc000024440|l:0x0|r:0xc000024480]
//	    │   └─r 1...7 [prio:0.4377] [0xc000024480|l:0x0|r:0xc0000244c0]
//	    │       └─r 1...5 [prio:0.4246] [0xc0000244c0|l:0x0|r:0x0]
//	    └─r 7...9 [prio:0.5152] [0xc000024900|l:0xc0000249c0|r:0x0]
//	        └─l 6...7 [prio:0.3009] [0xc0000249c0|l:0xc000024880|r:0x0]
//	            └─l 2...7 [prio:0.1565] [0xc000024880|l:0xc000024680|r:0xc0000248c0]
//	                ├─l 2...8 [prio:0.06564] [0xc000024680|l:0x0|r:0x0]
//	                └─r 4...8 [prio:0.09697] [0xc0000248c0|l:0x0|r:0x0]
func (t Tree[T]) FprintBST(w io.Writer) error {
	if t.root == nil {
		return nil
	}

	if _, err := fmt.Fprint(w, "R "); err != nil {
		return err
	}

	// start recursion with empty padding
	return t.binarytreeStringify(w, t.root, "")
}

// binarytreeStringify, traverse the tree, stringify the nodes in preorder
func (t *Tree[T]) binarytreeStringify(w io.Writer, n *node[T], pad string) error {
	// stringify this node
	_, err := fmt.Fprintf(w, "%v [prio:%.4g] [%p|l:%p|r:%p]\n",
		n.item, float64(n.prio)/math.MaxUint32, n, n.left, n.right)
	if err != nil {
		return err
	}

	// prepare glyphe, spacer and padding for next level
	var glyphe string
	var spacer string

	// left wing
	if n.left != nil {
		if n.right != nil {
			glyphe = "├─l "
			spacer = "│   "
		} else {
			glyphe = "└─l "
			spacer = "    "
		}
		if _, err := fmt.Fprint(w, pad+glyphe); err != nil {
			return err
		}
		if err := t.binarytreeStringify(w, n.left, pad+spacer); err != nil {
			return err
		}
	}

	// right wing
	if n.right != nil {
		glyphe = "└─r "
		spacer = "    "
		if _, err := fmt.Fprint(w, pad+glyphe); err != nil {
			return err
		}
		if err := t.binarytreeStringify(w, n.right, pad+spacer); err != nil {
			return err
		}
	}

	return nil
}

// parentChildsMap, needed for interval tree printing, this is not BST printing!
//
// Interval tree, parent->childs relation printed. A parent interval covers a child interval.
//
//	▼
//	├─ 0...6
//	│  └─ 0...5
//	├─ 1...8
//	│  ├─ 1...7
//	│  │  └─ 1...5
//	│  │     └─ 1...4
//	│  └─ 2...8
//	│     ├─ 2...7
//	│     └─ 4...8
//	│        └─ 6...7
//	└─ 7...9
type parentChildsMap[T any] struct {
	pcMap map[*node[T]][]*node[T] // parent -> []child map
	stack []*node[T]              // just needed for the algo
}

// buildParentChildsMap, in-order traversal
func (t *Tree[T]) buildParentChildsMap(n *node[T], pcm parentChildsMap[T]) parentChildsMap[T] {
	if n == nil {
		return pcm
	}

	// in-order traversal, left tree
	pcm = t.buildParentChildsMap(n.left, pcm)

	// detect parent-child-mapping for this node
	pcm = t.pcmForNode(n, pcm)

	// in-order traversal, right tree
	return t.buildParentChildsMap(n.right, pcm)
}

// pcmForNode, find parent in stack, remove items from stack, put this item on stack.
func (t *Tree[T]) pcmForNode(n *node[T], pcm parentChildsMap[T]) parentChildsMap[T] {
	// if this item is covered by a prev item on stack
	for j := len(pcm.stack) - 1; j >= 0; j-- {

		that := pcm.stack[j]
		if t.cmpCovers(that.item, n.item) {
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

	// put current neode on stack for next node
	pcm.stack = append(pcm.stack, n)

	return pcm
}

// Statistics, returns the maxDepth, average and standard deviation of the nodes.
//
// Note: This is for debugging and testing purposes only during development in semver
// 0.x.y. In future versions this will be removed without increasing the main
// semantic version, so please do not rely on it for now.
func (t Tree[T]) Statistics() (size int, maxDepth int, average, deviation float64) {
	// key is depth, value is the sum of nodes with this depth
	depths := make(map[int]int)

	// get the depths, sum up the size
	t.traverse(t.root, inorder, 0, func(n *node[T], depth int) bool {
		depths[depth] += 1
		size += 1
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

	return size, maxDepth, math.Round(average*10000) / 10000, math.Round(deviation*10000) / 10000
}

// Min returns the min item in tree.
func (t Tree[T]) Min() (min T) {
	n := t.root
	if n == nil {
		return
	}

	for n.left != nil {
		n = n.left
	}
	return n.item
}

// Max returns the max item in tree.
func (t Tree[T]) Max() (max T) {
	n := t.root
	if n == nil {
		return
	}

	for n.right != nil {
		n = n.right
	}
	return n.item
}

// Visit traverses the tree with item >= start to item <= stop in ascending order,
// or if start > stop, then the order is reversed. The visit function is called for each item.
//
// For example, the entire tree can be traversed as follows
//
//	t.Visit(t.Min(), t.Max(), visitFn)
//
// or in reverse order by
//
//	t.Visit(t.Max(), t.Min(), visitFn).
//
// The traversion terminates prematurely if the visit function returns false.
func (t Tree[T]) Visit(start, stop T, visitFn func(item T) bool) {
	if t.root == nil {
		return
	}

	order := inorder
	if t.compare(start, stop) > 0 {
		start, stop = stop, start
		order = reverse
	}

	// treaps are really cool datastructures!!!
	_, mid1, r := t.split(t.root, start, true)
	l, mid2, _ := t.split(r, stop, true)

	span := (&t).join(mid1, (&t).join(l, mid2, true), true)

	t.traverse(span, order, 0, func(n *node[T], _ int) bool {
		return visitFn(n.item)
	})
}

// Clone, deep cloning of the tree structure.
func (t Tree[T]) Clone() *Tree[T] {
	c := t
	c.root = t.clone(t.root)
	return &c
}

// clone rec-descent
func (t *Tree[T]) clone(n *node[T]) *node[T] {
	if n == nil {
		return n
	}
	n = n.copyNode()

	n.left = t.clone(n.left)
	n.right = t.clone(n.right)
	t.recalc(n)

	return n
}
