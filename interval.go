// package interval provides fast lookups and various methods
// on generic one-dimensional intervals.
//
// The author of the library uses it for IP-Range lookups in
// Access-Control-Lists (ACL) and in the authors own
// IP-Address-Management (IPAM) and network-management software.
//
// But the library is also useful for all one-dimensional arrays,
// e.g. time intervals.
//
// Thanks to generics this could be abstracted with minimal constraints.
package interval

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

// parent index of all childs
const root = -1

// Interface for generic one dimensional intervals.
//
// Compare the lower and upper points of two intervals.
type Interface[T any] interface {
	CompareLower(T) int
	CompareUpper(T) int
}

// Tree is the handle to perform various methods on a slice of intervals.
// This is a generic type, the implementation constraint is defined by the interval.Interface.
type Tree[T Interface[T]] struct {
	// the sorted items, immutable, stored as slice, not as tree, duplicates removed
	items []T

	// item indices: parent -> []child
	idxTree map[int][]int
}

// NewTree takes a slice of intervals and returns the tree handle.
// The algorithm prohibits duplicates and these are therefore sorted out.
func NewTree[T Interface[T]](items []T) *Tree[T] {
	t := new(Tree[T])

	if len(items) == 0 {
		return t
	}

	// the underlying data structure is a sorted slice ...
	t.items = make([]T, len(items))

	// ... and a parent -> []child tree with the items indices
	t.idxTree = make(map[int][]int)

	// clone and sort the items
	copy(t.items, items)
	sortDefault(t.items)

	// skip/drop duplicates
	pos := 0
	for i := range t.items {
		// skip duplicates, don't incr pos
		if i > 1 && equal(t.items[i], t.items[i-1]) {
			continue
		}

		// move item in slice
		if i != pos {
			t.items[pos] = t.items[i]
		}

		pos++
	}
	// shrink len and cap after removing duplicates
	t.items = t.items[:pos:pos]

	// build parent -> child(s) relationship
	t.buildIndexTree()

	return t
}

// Size returns the number of intervals in tree.
func (t *Tree[T]) Size() int {
	return len(t.items)
}

// ######################################################################################
// mother's little helpers
// ######################################################################################

// equal reports whether a == b. Operator == not available.
func equal[T Interface[T]](a, b T) bool {
	return a.CompareLower(b) == 0 && a.CompareUpper(b) == 0
}

// covers reports whether a truly covers b (not equal).
func covers[T Interface[T]](a, b T) bool {
	if equal(a, b) {
		return false
	}
	return a.CompareLower(b) <= 0 && a.CompareUpper(b) >= 0
}

// compareDefault, compare the lower points,
// sort supersets to the left as tiebreaker
//
// a wrapper for CompareLower[T] with added functionality for superset sorting
func compareDefault[T Interface[T]](a, b T) int {
	if equal(a, b) {
		return 0
	}

	// cmp lower
	if cmp := a.CompareLower(b); cmp != 0 {
		return cmp
	}

	// lower is equal, sort supersets to the left
	return -(a.CompareUpper(b))
}

// ######################################################################################
// SORTING
// ######################################################################################

// sortDefault the slice in place, with lower points ascending.
// As tie breaker sort supersets to the left.
//
//	[2...9 3...5 3...4 7...9]
func sortDefault[T Interface[T]](items []T) {
	sort.Slice(items, func(i, j int) bool { return compareDefault(items[i], items[j]) < 0 })
}

// sortUpper sorts the slice in place, with upper point ascending
func sortUpper[T Interface[T]](items []T) {
	sort.Slice(items, func(i, j int) bool { return items[i].CompareUpper(items[j]) < 0 })
}

// ######################################################################################
// LOOKUP, use the provided interface methods
// ######################################################################################

// Shortest returns the shortest interval that covers item.
// ok is true on success.
//
// Returns the identical interval if it exists in the tree,
// or the interval at which the item would be inserted.
//
// If the interval tree consists of IP CIDRs, shortest is identical to the longest-prefix-match.
//
// The meaning of 'shortest' is best explained with an example
//
//	e.g. for this interval tree
//
//		 ▼
//		 ├─ 0...6
//		 │  └─ 0...5
//		 ├─ 1...8
//		 │  ├─ 1...7
//		 │  │  └─ 1...5
//		 │  │     └─ 1...4
//		 │  └─ 2...8
//		 │     ├─ 2...7
//		 │     └─ 4...8
//		 │        └─ 6...7
//		 └─ 7...9
//
//	 tree.Shortest(ival{0,5}) returns ival{0,5}, true
//	 tree.Shortest(ival{3,6}) returns ival{2,7}, true
//	 tree.Shortest(ival{6,9}) returns ival{},    false
//
// If the item would be inserted directly under root,
// the zero value and false is be returned.
func (t *Tree[T]) Shortest(item T) (match T, ok bool) {
	// rec-descent
	return t.lookup(root, item)
}

func (t *Tree[T]) lookup(p int, item T) (match T, ok bool) {
	// dereference
	cs := t.idxTree[p]

	// find pos in slice on this level where t.items.lower > item.lower
	// item: 0...5
	// t.items:    [0...6 0...5 1...8 1...7 1...5 1...4 2...8 2...7 4...8 6...7 7...9]
	// idx: 2                   ^
	idx := sort.Search(len(cs), func(i int) bool { return t.items[cs[i]].CompareLower(item) > 0 })

	// child before idx may be equal or covers item
	if idx > 0 {
		idx--
		if equal(t.items[cs[idx]], item) {
			return t.items[cs[idx]], true
		}

		if covers(t.items[cs[idx]], item) {
			return t.lookup(cs[idx], item)
		}
	}

	// not on this level, return parent
	if p == root {
		// root is no legal value, just synthetic
		return
	}
	return t.items[p], true
}

// Largest returns the lower superset (top-down) that covers item.
// ok is true on success.
//
// The meaning of 'largest' is best explained with an example
//
//	e.g. for this interval tree
//
//		 ▼
//		 ├─ 0...6
//		 │  └─ 0...5
//		 ├─ 1...8
//		 │  ├─ 1...7
//		 │  │  └─ 1...5
//		 │  │     └─ 1...4
//		 │  └─ 2...8
//		 │     ├─ 2...7
//		 │     └─ 4...8
//		 │        └─ 6...7
//		 └─ 7...9
//
//	 tree.Largest(ival{0,6}) returns ival{0,6}, true
//	 tree.Largest(ival{0,5}) returns ival{0,6}, true
//	 tree.Largest(ival{3,7}) returns ival{1,8}, true
//	 tree.Largest(ival{6,9}) returns ival{},    false
//
// If the item is not covered by any interval in the tree,
// the zero value and false is returned.
func (t *Tree[T]) Largest(item T) (match T, ok bool) {
	// dereference root level slice
	rs := t.idxTree[root]

	// find pos in slice on root level where t.items.lower > item.lower
	// t.items: [0...6 1...8 7...9]
	// item:           2...5
	// idx:                  !
	idx := sort.Search(len(rs), func(i int) bool { return t.items[rs[i]].CompareLower(item) > 0 })

	if idx == 0 {
		return
	}

	// item before idx may be equal
	idx--
	if equal(t.items[rs[idx]], item) {
		// the items on any level are sorted and disjunct, maybe overlapping, BUT NOT covering each other
		// therefore we can return here, no element before can overlap this item
		return t.items[rs[idx]], true
	}

	// item isn't equal to any root level interval, find and return leftmost superset

	// some items before idx may cover item, find the leftmost
	for j := idx; j >= 0; j-- {

		// match, but continue to find next to the left also covering
		if covers(t.items[rs[j]], item) {
			match = t.items[rs[j]]
			ok = true
			continue
		}

		// remember: the items on any level are sorted and disjunct, maybe overlapping, BUT NOT covering each other
		// premature stop condition without item coverage, last match was superset
		break
	}

	return
}

// Supersets returns all intervals that covers the item in sorted order.
func (t *Tree[T]) Supersets(item T) []T {
	// idx is first interval where t.items[i].lower > item.lower
	idxLower := sort.Search(len(t.items), func(i int) bool { return t.items[i].CompareLower(item) > 0 })

	// resort remaining intervals [:idxFirst]
	// clone and sortUpper
	sl := make([]T, idxLower)
	copy(sl, t.items[:idxLower])
	sortUpper(sl)

	// idx is first interval where sl[i].upper is >= item.upper
	// lower limit of supersets
	idxUpper := sort.Search(len(sl), func(i int) bool { return sl[i].CompareUpper(item) >= 0 })

	// sort.Search: ... if there is no such index, Search returns n.
	if idxUpper == len(sl) {
		return nil
	}

	// maybe nil
	result := append([]T(nil), sl[idxUpper:]...)

	// sort before return
	sortDefault(result)

	return result
}

// Subsets returns all intervals in tree that are covered by item in sorted order.
func (t *Tree[T]) Subsets(item T) []T {
	// idx is first interval where t.items.lower >= item.lower
	// item: 3...8
	// t.items:    [0...6 0...5 1...8 1...7 1...5 1...4 2...8 2...7 4...8 6...7 7...9]
	// idx: 8                                                 ^
	idxLower := sort.Search(len(t.items), func(i int) bool { return t.items[i].CompareLower(item) >= 0 })

	// remaining intervals [idx:]
	// t.items: [2...7 4...8 6...7 7...9]

	// resort remaining intervals [idxLower:]
	// clone and sortUpper
	sl := make([]T, len(t.items)-idxLower)
	copy(sl, t.items[idxLower:])
	sortUpper(sl)

	// idx is first interval where sl[i].upper is > item.upper
	// item: 3...8
	// sl:         [6...7 4...8 7...9]
	// idx: 2                       ^
	idxUpper := sort.Search(len(sl), func(i int) bool { return sl[i].CompareUpper(item) > 0 })

	// [6...7 4...8]
	result := append([]T(nil), sl[:idxUpper]...) // maybe nil

	// sort before return
	sortDefault(result)

	return result
}

// String returns the ordered tree as a directory graph.
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
// If the interval items don't implement fmt.Stringer they are stringified with their default format %v.
func (t *Tree[T]) String() string {
	if len(t.items) == 0 {
		return ""
	}

	w := new(strings.Builder)

	// start symbol
	w.WriteString("▼\n")

	// start recursion with root and empty padding
	t.walkAndStringify(root, "", w)

	return w.String()
}

// walkAndStringify rec-descent, top-down
//
//	p:   parent index
//	pad: padding hrows and shrinks during recursion
//	w:   a StringWriter
func (t *Tree[T]) walkAndStringify(p int, pad string, w io.StringWriter) {
	// the prefix (pad + glyphe) is already printed on the line on upper level
	if p != root {
		w.WriteString(fmt.Sprintf("%v\n", t.items[p])) //nolint:errcheck
	}

	glyphe := "├─ "
	spacer := "│  "

	// dereference child-slice for clearer code
	cs := t.idxTree[p]

	// for all childs do, but ...
	for i, ii := range cs {
		// ... treat last child special
		if i == len(cs)-1 {
			glyphe = "└─ "
			spacer = "   "
		}
		// print prefix for next item
		w.WriteString(pad + glyphe) //nolint:errcheck

		// recdescent down
		t.walkAndStringify(ii, pad+spacer, w)
	}
}

// buildIndexTree, parent->child map, iterative algo with stack.
// Just building the tree with the slice indices, the items itself are not moved.
//
//		e.g.
//		 items in sort order lower-left:
//		  [0...300 0...100 9...18 13...18 15...19 200...400 201...230 203...300]
//
//	  map[parent][]child indices, root == -1
//
//		 map[int][]int:
//		  -1: [0 5]
//		   0: [1]
//		   1: [2 4]
//		   2: [3]
//		   5: [6 7]
//
//		 ▼
//		 ├─ 0...300
//		 │  └─ 0...100
//		 │     ├─ 9...18
//		 │     │  └─ 13...18
//		 │     └─ 15...19
//		 └─ 200...400
//		    ├─ 201...230
//		    └─ 203...300
func (t *Tree[T]) buildIndexTree() {
	// prev item on top of stack
	var stack []int

	// for all items ...
	for i := range t.items {

		// if this item is covered by a prev item on stack
		for j := len(stack) - 1; j >= 0; j-- {

			// de-reference, stack values are indices into t.items[stack[j]]
			k := stack[j]

			if covers(t.items[k], t.items[i]) {
				// item k is parent to item i
				t.idxTree[k] = append(t.idxTree[k], i)
				break
			}

			// sort order is lower-left: if next item wasn't covered, remove it from stack
			stack = stack[:j]
		}

		// stack is emptied, no item on stack covers current item
		if len(stack) == 0 {
			// parent is root
			t.idxTree[root] = append(t.idxTree[root], i)
		}

		// put current item on stack für next round
		stack = append(stack, i)
	}
}
