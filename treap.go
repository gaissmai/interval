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

// node is the basic recursive data structure.
type node[T Interface[T]] struct {
	// augment the treap for interval lookups
	minUpper *node[T] // pointer to node in subtree with min upper value
	maxUpper *node[T] // pointer to node in subtree with max upper value
	//
	// base treap fields, in memory efficient order
	left  *node[T]
	right *node[T]
	prio  uint32 // random key for binary heap, balances the tree
	item  T      // generic key/value
}

// Tree is the public handle to the hidden implementation.
//
// The zero value is useful without initialization, but it may be clearer to use [NewTree]
// because of the possibility of type inference.
type Tree[T Interface[T]] struct {
	root *node[T]
}

// NewTree initializes the interval tree with zero or more items of generic type T.
// The type constraint is defined by the [interval.Interface].
func NewTree[T Interface[T]](items ...T) Tree[T] {
	var t Tree[T]
	t.InsertMutable(items...)
	return t
}

// makeNode, create new node with item and random priority.
func makeNode[T Interface[T]](item T) *node[T] {
	n := new(node[T])
	n.item = item
	n.prio = rand.Uint32()
	n.recalc() // initial calculation of finger pointers...

	return n
}

// copyNode, make a shallow copy of the pointers and the item, no recalculation necessary.
func (n *node[T]) copyNode() *node[T] {
	c := *n
	return &c
}

// Insert elements into the tree, returns the new Tree.
// If an element is a duplicate, it replaces the previous element.
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

// InsertMutable inserts items into the tree, changing the original tree.
// If the original tree does not need to be preserved then this is much faster than the immutable insert.
func (t *Tree[T]) InsertMutable(items ...T) {
	for i := range items {
		t.root = t.root.insert(makeNode(items[i]), false)
	}
}

// insert into tree, changing nodes are copied, new treap is returned, old treap is modified if immutable is false.
func (n *node[T]) insert(m *node[T], immutable bool) *node[T] {
	if n == nil {
		return m
	}

	// if m is the new root?
	if m.prio >= n.prio {
		//
		//          m
		//          | split t in ( <m | dupe? | >m )
		//          v
		//       t
		//      / \
		//    l     d(upe)
		//   / \   / \
		//  l   r l   r
		//           /
		//          l
		//
		l, dupe, r := n.split(m.item, immutable)

		// replace dupe with m. m has same key but different prio than dupe, a join() is required
		if dupe != nil {
			return join(l, join(m, r, immutable), immutable)
		}

		// no duplicate, take m as new root
		//
		//     m
		//   /  \
		//  <m   >m
		//
		m.left, m.right = l, r
		m.recalc()
		return m
	}

	cmp := compare(m.item, n.item)
	if cmp == 0 {
		// replace duplicate item with m, but m has different prio, a join() is required
		return join(n.left, join(m, n.right, immutable), immutable)
	}

	if immutable {
		n = n.copyNode()
	}

	switch {
	case cmp < 0: // rec-descent
		n.left = n.left.insert(m, immutable)
		//
		//       R
		// m    l r
		//     l   r
		//
	case cmp > 0: // rec-descent
		n.right = n.right.insert(m, immutable)
		//
		//   R
		//  l r    m
		// l   r
		//
	}

	n.recalc() // node has changed, recalc
	return n
}

// Delete removes an item if it exists, returns the new tree and true, false if not found.
func (t Tree[T]) Delete(item T) (Tree[T], bool) {
	// split/join must be immutable
	l, m, r := t.root.split(item, true)
	t.root = join(l, r, true)

	ok := m != nil
	return t, ok
}

// DeleteMutable removes an item from tree, returns true if it exists, false otherwise.
// If the original tree does not need to be preserved then this is much faster than the immutable delete.
func (t *Tree[T]) DeleteMutable(item T) bool {
	l, m, r := t.root.split(item, false)
	t.root = join(l, r, false)

	return m != nil
}

// Union combines any two trees. In case of duplicate items, the "overwrite" flag
// controls whether the union keeps the original or whether it is replaced by the item in the other treap.
//
// The "immutable" flag controls whether the two trees are allowed to be modified.
//
// To create very large trees, it may be time-saving to slice the input data into chunks,
// fan out for creation and combine the generated subtrees with non-immutable unions.
func (t Tree[T]) Union(other Tree[T], overwrite bool, immutable bool) Tree[T] {
	t.root = t.root.union(other.root, overwrite, immutable)
	return t
}

func (n *node[T]) union(b *node[T], overwrite bool, immutable bool) *node[T] {
	// recursion stop condition
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
// If the split must be immutable, first copy concerned nodes.
func (n *node[T]) split(key T, immutable bool) (left, mid, right *node[T]) {
	// recursion stop condition
	if n == nil {
		return nil, nil, nil
	}

	if immutable {
		n = n.copyNode()
	}

	cmp := compare(n.item, key)
	switch {
	case cmp < 0:
		l, m, r := n.right.split(key, immutable)
		n.right = l
		n.recalc() // node has changed, recalc
		return n, m, r
		//
		//       (k)
		//      R
		//     l r   ==> (R.r, m, r) = R.r.split(k)
		//    l   r
		//
	case cmp > 0:
		l, m, r := n.left.split(key, immutable)
		n.left = r
		n.recalc() // node has changed, recalc
		return l, m, n
		//
		//   (k)
		//      R
		//     l r   ==> (l, m, R.l) = R.l.split(k)
		//    l   r
		//
	default:
		l, r := n.left, n.right
		n.left, n.right = nil, nil
		n.recalc() // node has changed, recalc
		return l, n, r
		//
		//     (k)
		//      R
		//     l r   ==> (R.l, R, R.r)
		//    l   r
		//
	}
}

// Find, searches for the exact interval in the tree and returns it as well as true,
// otherwise the zero value for item is returned and false.
func (t Tree[T]) Find(item T) (result T, ok bool) {
	n := t.root
	for {
		if n == nil {
			return
		}

		cmp := compare(item, n.item)
		switch {
		case cmp == 0:
			return n.item, true
		case cmp < 0:
			n = n.left
		case cmp > 0:
			n = n.right
		}
	}
}

// CoverLCP returns the interval with the longest-common-prefix that covers the item.
// If the item isn't covered by any interval, the zero value and false is returned.
//
// The meaning of 'LCP' is best explained with examples:
//
//   A, B and C covers the item, but B has longest-common-prefix (LCP) with item.
//
//   ------LCP--->|
//
//   Item            |----|
//
//   A |------------------------|
//   B            |---------------------------|
//   C     |---------------|
//   D              |--|
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
//      tree.CoverLCP(ival{0,5}) returns ival{0,5}, true
//      tree.CoverLCP(ival{3,6}) returns ival{2,7}, true
//      tree.CoverLCP(ival{6,9}) returns ival{},    false
//
// If the interval tree consists of IP CIDRs, CoverLCP is identical to the
// longest-prefix-match.
//
//  example: IP CIDRs as intervals
//
//     ▼
//     ├─ 0.0.0.0/0
//     │  ├─ 10.0.0.0/8
//     │  │  ├─ 10.0.0.0/24
//     │  │  └─ 10.0.1.0/24
//     │  └─ 127.0.0.0/8
//     │     └─ 127.0.0.1/32
//     └─ ::/0
//        ├─ ::1/128
//        ├─ 2000::/3
//        │  └─ 2001:db8::/32
//        ├─ fc00::/7
//        ├─ fe80::/10
//        └─ ff00::/8
//
//      tree.CoverLCP("10.0.1.17/32")       returns "10.0.1.0/24", true
//      tree.CoverLCP("2001:7c0:3100::/40") returns "2000::/3",    true
//
func (t Tree[T]) CoverLCP(item T) (result T, ok bool) {
	return t.root.lcp(item)
}

// lcp
func (n *node[T]) lcp(item T) (result T, ok bool) {
	if n == nil {
		return
	}

	// fast exit, node has too small max upper interval value (augmented value)
	if cmpRR(item, n.maxUpper.item) > 0 {
		return
	}

	cmp := compare(n.item, item)
	switch {
	case cmp > 0:
		// left rec-descent
		return n.left.lcp(item)
	case cmp == 0:
		// equality is always the shortest containing hull
		return n.item, true
	}

	// right backtracking
	result, ok = n.right.lcp(item)
	if ok {
		return result, ok
	}

	// not found in right subtree, try this node
	if covers(n.item, item) {
		return n.item, true
	}

	// left rec-descent
	return n.left.lcp(item)
}

// CoverSCP returns the interval with the shortest-common-prefix that covers the item.
// If the item isn't covered by any interval, the zero value and false is returned.
//
// The meaning of 'SCP' is best explained with examples:
//
//   A, B and C covers the item, but A has shortest-common-prefix (SCP) with item.
//
//   --SCP-->|
//
//   Item                  |----|
//
//   A       |------------------------|
//   B                  |---------------------------|
//   C           |---------------|
//   D     |-----------------|
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
//	 tree.CoverSCP(ival{0,6}) returns ival{0,6}, true
//	 tree.CoverSCP(ival{0,5}) returns ival{0,6}, true
//	 tree.CoverSCP(ival{3,7}) returns ival{1,8}, true
//	 tree.CoverSCP(ival{6,9}) returns ival{},    false
//
func (t Tree[T]) CoverSCP(item T) (result T, ok bool) {
	return t.root.scp(item)
}

// scp
func (n *node[T]) scp(item T) (result T, ok bool) {
	if n == nil {
		return
	}

	// fast exit, node has too small max upper interval value (augmented value)
	if cmpRR(item, n.maxUpper.item) > 0 {
		return
	}

	// left backtracking
	if result, ok = n.left.scp(item); ok {
		return result, ok
	}

	// this item
	if covers(n.item, item) {
		return n.item, true
	}

	// right rec-descent
	return n.right.scp(item)
}

// Covers returns all intervals that cover the item.
// The returned intervals are in sorted order.
func (t Tree[T]) Covers(item T) []T {
	return t.root.covers(item)
}

// covers
func (n *node[T]) covers(item T) (result []T) {
	if n == nil {
		return
	}

	// nope, subtree has too small upper interval value
	if cmpRR(item, n.maxUpper.item) > 0 {
		return
	}

	// in-order traversal for supersets, recursive call to left tree
	result = append(result, n.left.covers(item)...)

	// n.item covers item
	if covers(n.item, item) {
		result = append(result, n.item)
	}

	// recursive call to right tree
	return append(result, n.right.covers(item)...)
}

// CoveredBy returns all intervals that are covered by item.
// The returned intervals are in sorted order.
func (t Tree[T]) CoveredBy(item T) []T {
	return t.root.coveredBy(item)
}

// coveredBy
func (n *node[T]) coveredBy(item T) (result []T) {
	if n == nil {
		return
	}

	// nope, subtree has too big upper interval value
	if cmpRR(item, n.minUpper.item) < 0 {
		return
	}

	// in-order traversal for subsets, recursive call to left tree
	result = append(result, n.left.coveredBy(item)...)

	// item covers n.item
	if covers(item, n.item) {
		result = append(result, n.item)
	}

	// recursive call to right tree
	return append(result, n.right.coveredBy(item)...)
}

// Intersects returns true if any interval intersects item.
func (t Tree[T]) Intersects(item T) bool {
	return t.root.intersects(item)
}

func (n *node[T]) intersects(item T) bool {
	if n == nil {
		return false
	}

	// nope, subtree has too small upper value for intersection
	if cmpLR(item, n.maxUpper.item) > 0 {
		return false
	}

	// recursive call to left tree
	if n.left.intersects(item) {
		return true
	}

	// this n.item
	if intersects(n.item, item) {
		return true
	}

	// recursive call to right tree
	return n.right.intersects(item)
}

// Intersections returns all intervals that intersect with item.
// The returned intervals are in sorted order.
func (t Tree[T]) Intersections(item T) []T {
	return t.root.isections(item)
}

func (n *node[T]) isections(item T) (result []T) {
	if n == nil {
		return
	}

	// nope, subtree has too small upper value for intersection
	if cmpLR(item, n.maxUpper.item) > 0 {
		return
	}

	// in-order traversal for intersections, recursive call to left tree
	result = append(result, n.left.isections(item)...)

	// this n.item
	if intersects(n.item, item) {
		result = append(result, n.item)
	}

	// recursive call to right tree
	return append(result, n.right.isections(item)...)
}

// Precedes returns all intervals that precedes the item.
// The returned intervals are in sorted order.
//
//  example:
//
//   Item                       |-----------------|
//
//   A       |---------------------------------------|
//   B                  |-----|
//   C           |-----------------|
//   D     |-----------------|
//
//  Precedes(item) => [D, B]
//
func (t Tree[T]) Precedes(item T) []T {
	l, _, _ := t.root.split(item, true)
	return l.precedes(item)
}

func (n *node[T]) precedes(item T) (result []T) {
	if n == nil {
		return
	}

	// nope, all intervals in this subtree intersects with item
	if cmpLR(item, n.minUpper.item) <= 0 {
		return
	}

	// recursive call to ...
	result = append(result, n.left.precedes(item)...)

	// this n.item
	if !intersects(n.item, item) {
		result = append(result, n.item)
	}

	// recursive call to right tree
	return append(result, n.right.precedes(item)...)
}

// PrecededBy returns all intervals that are preceded by the item.
// The returned intervals are in sorted order.
//
//  example:
//
//   Item       |-----|
//
//   A       |---------------------------------------|
//   B                  |-----|
//   C           |-----------------|
//   D                    |-----------------|
//
//  PrecededBy(item) => [B, D]
//
func (t Tree[T]) PrecededBy(item T) []T {
	_, _, r := t.root.split(item, true)
	return r.precededby(item)
}

func (n *node[T]) precededby(item T) (result []T) {
	if n == nil {
		return
	}

	// skip some left wings
	if n.left != nil {
		if intersects(item, n.item) {
			// skip left, proceed instead with left.right
			result = append(result, n.left.right.precededby(item)...)
		} else {
			result = append(result, n.left.precededby(item)...)
		}
	}

	// this n.item
	if !intersects(n.item, item) {
		result = append(result, n.item)
	}

	// recursive call to right
	return append(result, n.right.precededby(item)...)
}

// join combines two disjunct treaps. All nodes in treap n have keys <= that of treap m
// for this algorithm to work correctly. If the join must be immutable, first copy concerned nodes.
func join[T Interface[T]](n, m *node[T], immutable bool) *node[T] {
	// recursion stop condition
	if n == nil {
		return m
	}
	if m == nil {
		return n
	}

	if n.prio > m.prio {
		//     n
		//    l r    m
		//          l r
		//
		if immutable {
			n = n.copyNode()
		}
		n.right = join(n.right, m, immutable)
		n.recalc()
		return n
	} else {
		//            m
		//      n    l r
		//     l r
		//
		if immutable {
			m = m.copyNode()
		}
		m.left = join(n, m.left, immutable)
		m.recalc()
		return m
	}
}

// recalc the augmented fields in treap node after each creation/modification with values in descendants.
// Only one level deeper must be considered. The treap datastructure is very easy to augment.
func (n *node[T]) recalc() {
	if n == nil {
		return
	}

	// start with upper min/max pointing to self
	n.minUpper = n
	n.maxUpper = n

	if n.right != nil {
		if cmpRR(n.minUpper.item, n.right.minUpper.item) > 0 {
			n.minUpper = n.right.minUpper
		}

		if cmpRR(n.maxUpper.item, n.right.maxUpper.item) < 0 {
			n.maxUpper = n.right.maxUpper
		}
	}

	if n.left != nil {
		if cmpRR(n.minUpper.item, n.left.minUpper.item) > 0 {
			n.minUpper = n.left.minUpper
		}

		if cmpRR(n.maxUpper.item, n.left.maxUpper.item) < 0 {
			n.maxUpper = n.left.maxUpper
		}
	}
}
