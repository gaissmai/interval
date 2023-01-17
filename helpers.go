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

// compare is a wrapper for Compare with added functionality for superset sorting
func compare[T Interface[T]](a, b T) int {
	ll, rr, _, _ := a.Compare(b)
	switch {
	case ll == 0:
		return -rr
	default:
		return ll
	}
}

func cmpUpper[T Interface[T]](a, b T) int {
	_, rr, _, _ := a.Compare(b)
	return rr
}

// traverse the BST in some order, call the visitor function for each node.
// covers reports whether a truly covers b (not equal).
func covers[T Interface[T]](a, b T) bool {
	ll, rr, _, _ := a.Compare(b)

	// equal
	if ll == 0 && rr == 0 {
		return false
	}

	return ll <= 0 && rr >= 0
}

// traverse the BST in some order, call the visitor function for each node.
// Prematurely stop traversion if visitor function returns false.
func (n *node[T]) traverse(order traverseOrder, depth int, visitFn func(n *node[T], depth int) bool) bool {
	if n == nil {
		return true
	}

	switch order {
	case inorder:
		// left, do-it, right
		if !n.left.traverse(order, depth+1, visitFn) {
			return false
		}

		if !visitFn(n, depth) {
			return false
		}

		if !n.right.traverse(order, depth+1, visitFn) {
			return false
		}

		return true
	case reverse:
		// right, do-it, left
		if !n.right.traverse(order, depth+1, visitFn) {
			return false
		}

		if !visitFn(n, depth) {
			return false
		}

		if !n.left.traverse(order, depth+1, visitFn) {
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
func (t Tree[T]) Fprint(w io.Writer) error {
	// pcm = parent-child-mapping
	var pcm parentChildsMap[T]

	// init map
	pcm.pcMap = make(map[*node[T]][]*node[T])

	pcm = t.root.buildParentChildsMap(pcm)

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

func walkAndStringify[T Interface[T]](w io.Writer, pcm parentChildsMap[T], parent *node[T], pad string) error {
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
// e.g. with left/right, item priority [prio] and pointers [this|left|right]
//
//  R 0...5 [prio:0.9405] [0xc000024940|l:0xc000024140|r:0xc000024980]
//  ├─l 0...6 [prio:0.6047] [0xc000024140|l:0x0|r:0x0]
//  └─r 1...4 [prio:0.6868] [0xc000024980|l:0xc000024440|r:0xc000024900]
//      ├─l 1...8 [prio:0.6646] [0xc000024440|l:0x0|r:0xc000024480]
//      │   └─r 1...7 [prio:0.4377] [0xc000024480|l:0x0|r:0xc0000244c0]
//      │       └─r 1...5 [prio:0.4246] [0xc0000244c0|l:0x0|r:0x0]
//      └─r 7...9 [prio:0.5152] [0xc000024900|l:0xc0000249c0|r:0x0]
//          └─l 6...7 [prio:0.3009] [0xc0000249c0|l:0xc000024880|r:0x0]
//              └─l 2...7 [prio:0.1565] [0xc000024880|l:0xc000024680|r:0xc0000248c0]
//                  ├─l 2...8 [prio:0.06564] [0xc000024680|l:0x0|r:0x0]
//                  └─r 4...8 [prio:0.09697] [0xc0000248c0|l:0x0|r:0x0]
//
func (t Tree[T]) FprintBST(w io.Writer) error {
	if t.root == nil {
		return nil
	}

	if _, err := fmt.Fprint(w, "R "); err != nil {
		return err
	}

	// start recursion with empty padding
	return t.root.preorderStringify(w, "")
}

// preorderStringify, traverse the tree, stringify the nodes in preorder
func (n *node[T]) preorderStringify(w io.Writer, pad string) error {
	// stringify this node
	if _, err := fmt.Fprintf(w, "%v [prio:%.4g] [%p|l:%p|r:%p]\n", n.item, n.prio, n, n.left, n.right); err != nil {
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
		if err := n.left.preorderStringify(w, pad+spacer); err != nil {
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
		if err := n.right.preorderStringify(w, pad+spacer); err != nil {
			return err
		}
	}

	return nil
}

// parentChildsMap, needed for interval tree printing, this is not BST printing!
//
// Interval tree, parent->childs relation printed. A parent interval covers a child interval.
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
	pcMap map[*node[T]][]*node[T] // parent -> []child map
	stack []*node[T]              // just needed for the algo
}

// buildParentChildsMap, in-order traversal
func (n *node[T]) buildParentChildsMap(pcm parentChildsMap[T]) parentChildsMap[T] {
	if n == nil {
		return pcm
	}

	// in-order traversal, left tree
	pcm = n.left.buildParentChildsMap(pcm)

	// detect parent-child-mapping for this node
	pcm = n.pcmForNode(pcm)

	// in-order traversal, right tree
	return n.right.buildParentChildsMap(pcm)
}

// pcmForNode, find parent in stack, remove items from stack, put this item on stack.
func (n *node[T]) pcmForNode(pcm parentChildsMap[T]) parentChildsMap[T] {
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

	// put current neode on stack for next node
	pcm.stack = append(pcm.stack, n)

	return pcm
}

// Statistics, returns the maxDepth, average and standard deviation of the nodes.
//
// Note: This is for debugging and testing purposes only during development in semver
// 0.x.y. In future versions this will be removed without increasing the main
// semantic version, so please do not rely on it for now.
//
func (t Tree[T]) Statistics() (maxDepth int, average, deviation float64) {
	// key is depth, value is the sum of nodes with this depth
	depths := make(map[int]int)

	// get the depths
	t.root.traverse(inorder, 0, func(n *node[T], depth int) bool {
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

// Max returns the node with max item in tree.
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

// Size returns the number of items in tree.
func (t Tree[T]) Size() int {
	size := 0
	t.root.traverse(inorder, 0, func(k *node[T], _ int) bool {
		size++
		return true
	})
	return size
}

// Visit traverses the tree with item >= start to item <= stop in ascending order,
// or if start > stop, then the order is reversed. The visit function is called for each item.
//
// For example, the entire tree can be traversed as follows
//  t.Visit(t.Min(), t.Max(), visitFn)
//
// or in reverse order by
//  t.Visit(t.Max(), t.Min(), visitFn).
//
// The traversion terminates prematurely if the visit function returns false.
//
func (t Tree[T]) Visit(start, stop T, visitFn func(item T) bool) {
	if t.root == nil {
		return
	}

	order := inorder
	if compare(start, stop) > 0 {
		start, stop = stop, start
		order = reverse
	}

	// treaps are really cool datastructures!!!
	_, mid1, r := t.root.split(start, true)
	l, mid2, _ := r.split(stop, true)

	span := join(mid1, join(l, mid2, true), true)

	span.traverse(order, 0, func(n *node[T], _ int) bool {
		return visitFn(n.item)
	})
}

// Clone, deep cloning of the tree structure, the items are copied.
func (t Tree[T]) Clone() Tree[T] {
	t.root = t.root.clone()
	return t
}

func (n *node[T]) clone() *node[T] {
	if n == nil {
		return n
	}
	n = n.copyNode()

	n.left = n.left.clone()
	n.right = n.right.clone()

	return n
}
