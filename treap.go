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
	minUpper *node[T] // pointer to node in subtree with min upper value, just needed for Subsets()
	maxUpper *node[T] // pointer to node in subtree with max upper value, needed for all other lookups
	//
	// base treap fields, in memory efficient order
	left  *node[T]
	right *node[T]
	prio  float64 // random key for binary heap, balances the tree
	item  T       // generic key/value
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
	n.prio = rand.Float64()
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

// Find, searches for the interval in the tree and returns it as well as true,
// otherwise the zero value for item and false.
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
func (t Tree[T]) Shortest(item T) (result T, ok bool) {
	return t.root.shortest(item)
}

// shortest can't use tree.split(key) because of allocations or mutations.
func (n *node[T]) shortest(item T) (result T, ok bool) {
	if n == nil {
		return
	}

	// fast exit, node has too small max upper interval value (augmented value)
	if cmpUpper(item, n.maxUpper.item) > 0 {
		return
	}

	cmp := compare(n.item, item)
	switch {
	case cmp > 0:
		return n.left.shortest(item)
	case cmp == 0:
		// equality is always the shortest containing hull
		return n.item, true
	}

	// now on proper depth in tree
	// first try right subtree for shortest containing hull
	if n.right != nil {

		// rec-descent with n.right
		if compare(n.right.item, item) <= 0 {
			result, ok = n.right.shortest(item)
			if ok {
				return result, ok
			}
		}

		// try n.right.left subtree for smallest containing hull
		// take this path only if n.right.left.item > t.item (this node)
		if n.right.left != nil && compare(n.right.left.item, n.item) > 0 {
			// rec-descent with n.right.left
			result, ok = n.right.left.shortest(item)
			if ok {
				return result, ok
			}
		}

	}

	// not found in right subtree, try this node
	if covers(n.item, item) {
		return n.item, true
	}

	// rec-descent with t.left
	return n.left.shortest(item)
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
func (t Tree[T]) Largest(item T) (result T, ok bool) {
	if t.root == nil {
		return
	}

	// algo with tree.split(), allocations allowed
	l, m, _ := t.root.split(item, true)
	result, ok = l.largest(item)

	// if key is in treap and no other largest found...
	if !ok && m != nil {
		return m.item, true
	}

	return
}

func (n *node[T]) largest(item T) (result T, ok bool) {
	if n == nil {
		return
	}

	// fast exit, node has too small max upper interval value (augmented value)
	if cmpUpper(item, n.maxUpper.item) > 0 {
		return
	}

	// rec-descent left subtree
	if result, ok = n.left.largest(item); ok {
		return result, ok
	}

	// this item
	if cmpUpper(item, n.item) <= 0 {
		return n.item, true
	}

	return n.right.largest(item)
}

// Supersets returns all intervals that covers the item in sorted order.
func (t Tree[T]) Supersets(item T) []T {
	if t.root == nil {
		return nil
	}
	var result []T

	// supersets algo with tree.split(), allocations allowed
	l, m, _ := t.root.split(item, true)
	result = l.supersets(item)

	// if key is in treap, add key to result set
	if m != nil {
		result = append(result, item)
	}

	return result
}

func (n *node[T]) supersets(item T) (result []T) {
	if n == nil {
		return
	}

	// nope, subtree has too small upper interval value
	if cmpUpper(item, n.maxUpper.item) > 0 {
		return
	}

	// in-order traversal for supersets, recursive call to left tree
	result = append(result, n.left.supersets(item)...)

	// this n.item
	if covers(n.item, item) {
		result = append(result, n.item)
	}

	// recursive call to right tree
	return append(result, n.right.supersets(item)...)
}

// Subsets returns all intervals in tree that are covered by item in sorted order.
func (t Tree[T]) Subsets(item T) []T {
	if t.root == nil {
		return nil
	}
	var result []T

	// subsets algo with tree.split(), allocations allowed
	_, m, r := t.root.split(item, true)

	// if key is in treap, start with key in result
	if m != nil {
		result = []T{item}
	}
	result = append(result, r.subsets(item)...)

	return result
}

func (n *node[T]) subsets(item T) (result []T) {
	if n == nil {
		return
	}

	// nope, subtree has too big upper interval value
	if cmpUpper(item, n.minUpper.item) < 0 {
		return
	}

	// in-order traversal for subsets, recursive call to left tree
	result = append(result, n.left.subsets(item)...)

	// this n.item
	if covers(n.item, item) {
		result = append(result, n.item)
	}

	// recursive call to right tree
	return append(result, n.right.subsets(item)...)
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
		if cmpUpper(n.minUpper.item, n.right.minUpper.item) > 0 {
			n.minUpper = n.right.minUpper
		}

		if cmpUpper(n.maxUpper.item, n.right.maxUpper.item) < 0 {
			n.maxUpper = n.right.maxUpper
		}
	}

	if n.left != nil {
		if cmpUpper(n.minUpper.item, n.left.minUpper.item) > 0 {
			n.minUpper = n.left.minUpper
		}

		if cmpUpper(n.maxUpper.item, n.left.maxUpper.item) < 0 {
			n.maxUpper = n.left.maxUpper
		}
	}
}
