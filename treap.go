// Package interval provides fast lookups and various other methods for generic one-dimensional intervals.
//
// The author of the library uses the package for fast IP range lookups in access control lists (ACL)
// and in the author's own IP address management (IPAM) and network management software,
// see also the author's [iprange package].
//
// However, the interval package is useful for all one-dimensional intervals, e.g. time intervals.
//
// [iprange package]: https://github.com/gaissmai/iprange
package interval

import (
	"math/rand"
)

// Interface is the type constraint for generic interval items.
// Compare the lower and upper points of two intervals.
type Interface[T any] interface {
	CompareLower(T) int
	CompareUpper(T) int
}

// Node is the basic recursive data structure. This is a generic type,
// the implementation constraint is defined by the interval.Interface.
type Node[T Interface[T]] struct {
	//
	// augment the treap for interval lookups
	minUpper *Node[T] // pointer to node in subtree with min upper value, just needed for Subsets()
	maxUpper *Node[T] // pointer to node in subtree with max upper value, needed for all other lookups
	//
	// base treap fields, in memory efficient order
	left  *Node[T]
	right *Node[T]
	prio  float64 // random key for binary heap, balances the tree
	item  T       // generic key/value
}

// Tree is just a handle for the root node.
type Tree[T Interface[T]] struct {
	root *Node[T]
}

// NewTree initializes the interval tree.
func NewTree[T Interface[T]](items ...T) Tree[T] {
	var t Tree[T]

	for i := range items {
		t.root = t.root.insert(makeNode(items[i]), false)
	}
	return t
}

// makeNode, create new node with item and random priority.
func makeNode[T Interface[T]](item T) *Node[T] {
	n := new(Node[T])
	n.item = item
	n.prio = rand.Float64()
	n.recalc() // initial calculation of finger pointers...

	return n
}

// copyNode, make a shallow copy of the pointers and the item, no recalculation necessary.
func (n *Node[T]) copyNode() *Node[T] {
	if n == nil {
		return n
	}

	m := *n
	return &m
}

// Insert elements into the tree, if an element is a duplicate, it replaces the previous element.
func (t Tree[T]) Insert(items ...T) Tree[T] {
	// something to preserve?
	immutable := true
	if t.root == nil {
		immutable = false
	}

	for i := range items {
		t.root = t.root.insert(makeNode(items[i]), immutable)
	}
	return t
}

// insert into tree, changing nodes are copied, new treap is returned, old treap is modified if immutable is false.
func (n *Node[T]) insert(b *Node[T], immutable bool) *Node[T] {
	if n == nil {
		return b
	}

	// if b is the new root?
	if b.prio >= n.prio {
		//
		//          b
		//          | split t in ( <b | dupe? | >b )
		//          v
		//       t
		//      / \
		//    l     d(upe)
		//   / \   / \
		//  l   r l   r
		//           /
		//          l
		//
		l, dupe, r := n.split(b.item, immutable)

		// replace dupe with b. b has same key but different prio than dupe, a join() is required
		if dupe != nil {
			return join(l, join(b, r, immutable), immutable)
		}

		// no duplicate, take b as new root
		//
		//     b
		//   /  \
		//  <b   >b
		//
		b.left, b.right = l, r
		b.recalc()
		return b
	}

	cmp := compare(b.item, n.item)
	if cmp == 0 {
		// replace duplicate item with b, but b has different prio, a join() is required
		return join(n.left, join(b, n.right, immutable), immutable)
	}

	if immutable {
		n = n.copyNode()
	}

	switch {
	case cmp < 0: // rec-descent
		n.left = n.left.insert(b, immutable)
		//
		//       R
		// b    l r
		//     l   r
		//
	case cmp > 0: // rec-descent
		n.right = n.right.insert(b, immutable)
		//
		//   R
		//  l r    b
		// l   r
		//
	}

	n.recalc() // node has changed, recalc
	return n
}

// Delete removes an item if it exists, returns the new tree and true, false if not found.
func (t Tree[T]) Delete(item T) (Tree[T], bool) {
	immutable := true
	l, m, r := t.root.split(item, immutable)
	t.root = join(l, r, immutable)

	if m == nil {
		return t, false
	}
	return t, true
}

// Union combines any two trees. In case of duplicate items, the "overwrite" flag
// controls whether the union keeps the original or whether it is replaced by the item in the b treap.
//
// The immutable flag controls whether the old treaps are allowed to be modified.
//
// To create very large trees, it may be time-saving to split the input data into chunks,
// fan out for Insert and combine the generated subtrees with Union.
func (t Tree[T]) Union(b Tree[T], overwrite bool, immutable bool) Tree[T] {
	t.root = t.root.union(b.root, overwrite, immutable)
	return t
}

func (n *Node[T]) union(b *Node[T], overwrite bool, immutable bool) *Node[T] {
	if n == nil {
		return b
	}
	if b == nil {
		return n
	}

	// swap treaps if needed, treap with higher prio remains as new root
	if n.prio < b.prio {
		n, b = b, n
		overwrite = !overwrite
	}

	// immutable union, copy remaining root
	if immutable {
		n = n.copyNode()
	}

	// the treap with the lower priority is split with the root key in the treap with the higher priority
	l, dupe, r := b.split(n.item, immutable)

	// the treaps may have duplicate items
	if overwrite && dupe != nil {
		n.item = dupe.item
	}

	// rec-descent
	n.left = n.left.union(l, overwrite, immutable)
	n.right = n.right.union(r, overwrite, immutable)
	n.recalc()

	return n
}

// split the treap into all nodes that compare less-than, equal
// and greater-than the provided item (BST key). The resulting nodes are
// properly formed treaps or nil.
func (t *Node[T]) split(key T, immutable bool) (left, mid, right *Node[T]) {
	// recursion stop condition
	if t == nil {
		return nil, nil, nil
	}

	if immutable {
		t = t.copyNode()
	}

	cmp := compare(t.item, key)
	switch {
	case cmp < 0:
		l, m, r := t.right.split(key, immutable)
		t.right = l
		t.recalc() // node has changed, recalc
		return t, m, r
		//
		//       (k)
		//      R
		//     l r   ==> (R.r, m, r) = R.r.split(k)
		//    l   r
		//
	case cmp > 0:
		l, m, r := t.left.split(key, immutable)
		t.left = r
		t.recalc() // node has changed, recalc
		return l, m, t
		//
		//   (k)
		//      R
		//     l r   ==> (l, m, R.l) = R.l.split(k)
		//    l   r
		//
	default:
		l, r := t.left, t.right
		t.left, t.right = nil, nil
		t.recalc() // node has changed, recalc
		return l, t, r
		//
		//     (k)
		//      R
		//     l r   ==> (R.l, R, R.r)
		//    l   r
		//
	}
}

// Shortest returns the most specific interval that covers item. ok is true on
// success.
//
// Returns the identical interval if it exists in the tree, or the interval at
// which the item would be inserted.
//
// If the item would be inserted directly under root, the zero value and false
// is returned.
//
// If the interval tree consists of IP CIDRs, shortest is identical to the
// longest-prefix-match.
//
// The meaning of 'shortest' is best explained with an example
//
//     e.g. for this interval tree
//
//     	 ▼
//     	 ├─ 0...6
//     	 │  └─ 0...5
//     	 ├─ 1...8
//     	 │  ├─ 1...7
//     	 │  │  └─ 1...5
//     	 │  │     └─ 1...4
//     	 │  └─ 2...8
//     	 │     ├─ 2...7
//     	 │     └─ 4...8
//     	 │        └─ 6...7
//     	 └─ 7...9
//
//      tree.Shortest(ival{0,5}) returns ival{0,5}, true
//      tree.Shortest(ival{3,6}) returns ival{2,7}, true
//      tree.Shortest(ival{6,9}) returns ival{},    false
//
func (t *Node[T]) Shortest(item T) (result T, ok bool) {
	if t == nil {
		return
	}

	// fast exit, node has too small max upper interval value (augmented value)
	if item.CompareUpper(t.maxUpper.item) > 0 {
		return
	}

	cmp := compare(t.item, item)
	switch {
	case cmp > 0:
		return t.left.Shortest(item)
	case cmp == 0:
		// equality is always the shortest containing hull
		return t.item, true
	}

	// now on proper depth in tree
	// first try right subtree for shortest containing hull
	if t.right != nil {

		// rec-descent with t.right
		if compare(t.right.item, item) <= 0 {
			result, ok = t.right.Shortest(item)
			if ok {
				return result, ok
			}
		}

		// try t.right.left subtree for smallest containing hull
		// take this path only if t.right.left.item > t.item (this node)
		if t.right.left != nil && compare(t.right.left.item, t.item) > 0 {
			// rec-descent with t.right.left
			result, ok = t.right.left.Shortest(item)
			if ok {
				return result, ok
			}
		}

	}

	// not found in right subtree, try this node
	if covers(t.item, item) {
		return t.item, true
	}

	// rec-descent with t.left
	return t.left.Shortest(item)
}

// Largest returns the largest interval (top-down in tree) that covers item.
// ok is true on success, otherwise the item isn't contained in the tree.
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
//
func (t *Node[T]) Largest(item T) (result T, ok bool) {
	// find node.item < item
	for {
		if t == nil {
			return
		}

		if compare(t.item, item) > 0 {
			t = t.left
			continue
		}
		break
	}

	// start recursion
	return t.largest(item)
}

func (t *Node[T]) largest(item T) (result T, ok bool) {
	if t == nil {
		return
	}

	// fast exit, node has too small max upper interval value (augmented value)
	if item.CompareUpper(t.maxUpper.item) > 0 {
		return
	}

	// rec-descent left subtree
	if result, ok = t.left.largest(item); ok {
		return result, ok
	}

	// this item
	if item.CompareUpper(t.item) <= 0 {
		return t.item, true
	}

	// rec-descent right
	return t.right.largest(item)
}

// Supersets returns all intervals that covers the item in sorted order.
func (t *Node[T]) Supersets(item T) []T {
	if t == nil {
		return nil
	}
	var result []T

	l, m, _ := t.split(item, true)
	result = l.supersets(item)

	// if key is in treap, add key to result set
	if m != nil {
		result = append(result, item)
	}

	return result
}

func (t *Node[T]) supersets(item T) (result []T) {
	if t == nil {
		return
	}

	// nope, subtree has too small upper interval value
	if item.CompareUpper(t.maxUpper.item) > 0 {
		return
	}

	// in-order traversal for supersets, recursive call to left tree
	if t.left != nil && item.CompareUpper(t.left.maxUpper.item) <= 0 {
		result = append(result, t.left.supersets(item)...)
	}

	// this item
	if item.CompareUpper(t.item) <= 0 {
		result = append(result, t.item)
	}

	// recursive call to right tree
	if t.right != nil && item.CompareUpper(t.right.maxUpper.item) <= 0 {
		result = append(result, t.right.supersets(item)...)
	}

	return
}

// Subsets returns all intervals in tree that are covered by item in sorted order.
func (t *Node[T]) Subsets(item T) []T {
	if t == nil {
		return nil
	}
	var result []T

	_, m, r := t.split(item, true)

	// if key is in treap, start with key in result
	if m != nil {
		result = []T{item}
	}
	result = append(result, r.subsets(item)...)

	return result
}

func (t *Node[T]) subsets(item T) (result []T) {
	if t == nil {
		return
	}

	// nope, subtree has too big upper interval value
	if item.CompareUpper(t.minUpper.item) < 0 {
		return
	}

	// in-order traversal for subsets, recursive call to left tree
	if t.left != nil && item.CompareUpper(t.left.minUpper.item) >= 0 {
		result = append(result, t.left.subsets(item)...)
	}

	// this item
	if item.CompareUpper(t.item) >= 0 {
		result = append(result, t.item)
	}

	// recursive call to right tree
	if t.right != nil && item.CompareUpper(t.right.minUpper.item) >= 0 {
		result = append(result, t.right.subsets(item)...)
	}

	return
}

// join combines two disjunct treaps. All nodes in treap a have keys <= that of treap b
// for this algorithm to work correctly. The join is immutable, first copy concerned nodes.
func join[T Interface[T]](a, b *Node[T], immutable bool) *Node[T] {
	// recursion stop condition
	if a == nil {
		return b
	}
	if b == nil {
		return a
	}

	if a.prio > b.prio {
		//     a
		//    l r    b
		//          l r
		//
		if immutable {
			a = a.copyNode()
		}
		a.right = join(a.right, b, immutable)
		a.recalc()
		return a
	} else {
		//            b
		//      a    l r
		//     l r
		//
		if immutable {
			b = b.copyNode()
		}
		b.left = join(a, b.left, immutable)
		b.recalc()
		return b
	}
}

// recalc the augmented fields in treap node after each creation/modification with values in descendants.
// Only one level deeper must be considered. The treap datastructure is very easy to augment.
func (t *Node[T]) recalc() {
	if t == nil {
		return
	}

	// start with upper min/max pointing to self
	t.minUpper = t
	t.maxUpper = t

	if t.right != nil {
		if t.minUpper.item.CompareUpper(t.right.minUpper.item) > 0 {
			t.minUpper = t.right.minUpper
		}

		if t.maxUpper.item.CompareUpper(t.right.maxUpper.item) < 0 {
			t.maxUpper = t.right.maxUpper
		}
	}

	if t.left != nil {
		if t.minUpper.item.CompareUpper(t.left.minUpper.item) > 0 {
			t.minUpper = t.left.minUpper
		}

		if t.maxUpper.item.CompareUpper(t.left.maxUpper.item) < 0 {
			t.maxUpper = t.left.maxUpper
		}
	}
}
