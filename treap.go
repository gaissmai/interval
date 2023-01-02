package interval

import "math/rand"

// Interface is the type constraint for generic interval items.
type Interface[T any] interface {
	CompareLower(T) int
	CompareUpper(T) int
}

// Tree is the basic recursive data structure, usable without initialization.
//
// This is a generic type, the implementation constraint is defined by the interval.Interface.
type Tree[T Interface[T]] struct {
	//
	// augment the treap for interval lookups
	minUpper *Tree[T] // finger pointer to node in subtree with min upper value, just needed for Subsets()
	maxUpper *Tree[T] // finger pointer to node in subtree with max upper value
	//
	// augment the treap for some statistics
	size   int // descendents at this node
	height int // height at this node
	//
	// base treap fields, in memory efficient order
	left  *Tree[T]
	right *Tree[T]
	prio  float64 // automatic balance the tree, random key for binary heap
	item  T       // generic key/value
}

// NewTree takes zero or more intervals and returns the new tree.
// Duplicate items are silently dropped during insert.
func NewTree[T Interface[T]](items ...T) *Tree[T] {
	var t *Tree[T]
	for i := range items {
		t = t.insert(makeNode(items[i]))
	}
	return t
}

// makeNode, create new node with item and random priority.
func makeNode[T Interface[T]](item T) *Tree[T] {
	n := new(Tree[T])
	n.item = item
	n.prio = rand.Float64()
	n.recalc() // initial calculation of augmented fields, size, height, finger pointers...

	return n
}

// copyNode, make a shallow copy, no recalculation necessary.
func (t *Tree[T]) copyNode() *Tree[T] {
	n := *t
	return &n
}

// Item returns the item field.
func (t *Tree[T]) Item() (item T) {
	if t == nil {
		return
	}
	return t.item
}

// Size returns the number of descendents at this position in the tree.
func (t *Tree[T]) Size() int {
	if t == nil {
		return 0
	}
	return t.size
}

// Height returns the height at this position in the tree.
//
// Note:
// This is for statistical purposes only during development in semver 0.x.y.
// In future versions this may be removed without increasing the main semantic version,
// so please do not rely on it for now.
func (t *Tree[T]) Height() int {
	if t == nil {
		return 0
	}
	return t.height
}

// Insert items into the tree, returns the new tree.
// Duplicate items are silently dropped during insert.
func (t *Tree[T]) Insert(items ...T) *Tree[T] {
	for i := range items {
		t = t.insert(makeNode(items[i]))
	}
	return t
}

// insert into tree, changing nodes are copied, new treap is returned, old treap isn't modified.
func (t *Tree[T]) insert(b *Tree[T]) *Tree[T] {
	if t == nil {
		return b
	}
	//
	//           b
	//     a
	//    l r
	//
	if b.prio >= t.prio {
		left, dupe, right := t.split(b.item)
		if dupe != nil {
			// duplicate, drop b
			return t
		}

		b.left, b.right = left, right
		b.recalc() // node has changed, recalc
		return b
		//
		//     b
		//    l r
		//
	}

	// immutable insert, copy node
	root := t.copyNode()

	cmp := compare(b.item, root.item)
	switch {
	case cmp < 0: // rec-descent
		root.left = root.left.insert(b)
		//
		//       R
		// b    l r
		//     l   r
		//
	case cmp > 0: // rec-descent
		root.right = root.right.insert(b)
		//
		//   R
		//  l r    b
		// l   r
		//
	default: // equal, drop duplicate
	}

	root.recalc() // node has changed, recalc
	return root
}

// Upsert, replace/insert item in tree, returns the new tree.
func (t *Tree[T]) Upsert(item T) *Tree[T] {
	k := makeNode(item)
	if t == nil {
		return k
	}
	// don't use middle, replace it with k
	l, _, r := t.split(item)
	return join(l, join(k, r))
}

// Delete removes an item if it exists, returns the new tree and true, false if not found.
func (t *Tree[T]) Delete(item T) (*Tree[T], bool) {
	l, m, r := t.split(item)
	t = join(l, r)
	if m == nil {
		return t, false
	}
	return t, true
}

// split the treap into all nodes that compare less-than, equal
// and greater-than the provided item (BST key). The resulting nodes are
// properly formed treaps or nil.
func (t *Tree[T]) split(key T) (left, mid, right *Tree[T]) {
	// recursion stop condition
	if t == nil {
		return nil, nil, nil
	}

	// immutable split, copy node
	root := t.copyNode()

	cmp := compare(root.item, key)
	switch {
	case cmp < 0:
		l, m, r := root.right.split(key)
		root.right = l
		root.recalc() // node has changed, recalc
		return root, m, r
		//
		//       (k)
		//      R
		//     l r   ==> (R.r, m, r) = R.r.split(k)
		//    l   r
		//
	case cmp > 0:
		l, m, r := root.left.split(key)
		root.left = r
		root.recalc() // node has changed, recalc
		return l, m, root
		//
		//   (k)
		//      R
		//     l r   ==> (l, m, R.l) = R.l.split(k)
		//    l   r
		//
	default:
		l, r := root.left, root.right
		root.left, root.right = nil, nil
		root.recalc() // node has changed, recalc
		return l, root, r
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
// If the item would be inserted directly under root, the zero value and false
// is returned.
//
func (t *Tree[T]) Shortest(item T) (result T, ok bool) {
	if t == nil {
		return
	}

	// the shortest interval covering item must have t.item <= item
	l, m, _ := t.split(item)

	// item is in tree, return it as shortest.
	if m != nil {
		return m.item, true
	}
	return l.shortest(item)
}

// shortest, find rec-descent, use augmented maxUpper finger pointer.
func (t *Tree[T]) shortest(item T) (result T, ok bool) {
	if t == nil {
		return
	}

	// nope, whole subtree has too small max upper interval value
	if item.CompareUpper(t.maxUpper.item) > 0 {
		return
	}

	// reverse-order traversal for shortest
	// try right tree for smallest containing hull
	if t.right != nil && item.CompareUpper(t.right.maxUpper.item) <= 0 {
		if result, ok = t.right.shortest(item); ok {
			return result, ok
		}
	}

	// no match in right tree, try this item
	if item.CompareUpper(t.item) <= 0 {
		return t.item, true
	}

	// recursive call to left tree
	if t.left != nil && item.CompareUpper(t.left.maxUpper.item) <= 0 {
		if result, ok = t.left.shortest(item); ok {
			return result, ok
		}
	}

	// nope
	return
}

// Largest returns the largest superset (top-down in tree) that covers item.
// ok is true on success, otherwise the interval isn't contained in the tree.
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
//
func (t *Tree[T]) Largest(item T) (result T, ok bool) {
	if t == nil {
		return
	}

	l, m, _ := t.split(item)
	result, ok = l.largest(item)

	// if key is in treap and no outer hull found
	if m != nil && !ok {
		result, ok = item, true
	}

	return result, ok
}

// largest is the recursive workhorse for Largest().
func (t *Tree[T]) largest(item T) (result T, ok bool) {
	if t == nil {
		return
	}

	// nope, subtree has too small upper interval value
	if item.CompareUpper(t.maxUpper.item) > 0 {
		return
	}

	// in-order traversal for longest
	// try left tree for largest containing hull
	if t.left != nil && item.CompareUpper(t.left.maxUpper.item) <= 0 {
		if result, ok = t.left.largest(item); ok {
			return result, ok
		}
	}

	// this item
	if item.CompareUpper(t.item) <= 0 {
		return t.item, true
	}

	// recursive call to right tree
	if t.right != nil && item.CompareUpper(t.right.maxUpper.item) <= 0 {
		if result, ok = t.right.largest(item); ok {
			return result, ok
		}
	}

	// nope
	return
}

// Supersets returns all intervals that covers the item in sorted order.
func (t *Tree[T]) Supersets(item T) []T {
	if t == nil {
		return nil
	}
	var result []T

	l, m, _ := t.split(item)
	result = l.supersets(item)

	// if key is in treap, add key to result set
	if m != nil {
		result = append(result, item)
	}

	return result
}

func (t *Tree[T]) supersets(item T) (result []T) {
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
func (t *Tree[T]) Subsets(item T) []T {
	if t == nil {
		return nil
	}
	var result []T

	_, m, r := t.split(item)

	// if key is in treap, start with key in result
	if m != nil {
		result = []T{item}
	}
	result = append(result, r.subsets(item)...)

	return result
}

func (t *Tree[T]) subsets(item T) (result []T) {
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
func join[T Interface[T]](a, b *Tree[T]) *Tree[T] {
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
		a = a.copyNode() // immutable join, copy node
		a.right = join(a.right, b)
		a.recalc()
		return a
	} else {
		//            b
		//      a    l r
		//     l r
		//
		b = b.copyNode() // immutable join, copy node
		b.left = join(a, b.left)
		b.recalc()
		return b
	}
}

// recalc the augmented fields in treap node after each creation/modification with values in descendants.
// Only one level deeper must be considered. The treap datastructure is very easy to augment.
func (t *Tree[T]) recalc() {
	if t == nil {
		return
	}

	// recalc some statistics, not really needed for interval algo
	t.size = 1 + t.left.Size() + t.right.Size()
	t.height = 1 + max(t.left.Height(), t.right.Height())

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

// Visit traverses the tree with item >= start until item <= stop in ascending order,
// if start > stop, the order is reversed.
//
// The visit function is called for each item. The traversion stops prematurely if the visit function returns false.
func (t *Tree[T]) Visit(start, stop T, visitFn func(t T) bool) {
	if t == nil {
		return
	}

	order := inorder
	if start.CompareLower(stop) > 0 {
		start, stop = stop, start
		order = reverse
	}

	// treaps are really cool datastructures!
	_, mid1, r := t.split(start)
	l, mid2, _ := r.split(stop)

	span := join(mid1, join(l, mid2))

	span.traverse(order, visitFn)
}

// Min returns the node with min item in tree.
func (t *Tree[T]) Min() *Tree[T] {
	if t == nil {
		return t
	}

	for t.left != nil {
		t = t.left
	}
	return t
}

// Max returns the node with max item in tree.
func (t *Tree[T]) Max() *Tree[T] {
	if t == nil {
		return t
	}

	for t.right != nil {
		t = t.right
	}
	return t
}

// MinUpper returns the node with item with min upper value.
func (t *Tree[T]) MinUpper() *Tree[T] {
	return t.minUpper
}

// MaxUpper returns the node with item with max upper value.
func (t *Tree[T]) MaxUpper() *Tree[T] {
	return t.maxUpper
}
