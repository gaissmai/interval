package interval

import "math/rand"

// Interval is the type constraint for generic interval items.
type Interval[T any] interface {
	CompareLower(T) int
	CompareUpper(T) int
}

// Tree is the basic recursive data structure, augmented for fast interval lookups.
// This is a generic type, the implementation constraint is defined by the interval.Interface.
type Tree[T Interval[T]] struct {
	//
	// augment the treap for interval lookups
	minUpper *Tree[T] // pointer to tree in subtree with min upper value
	maxUpper *Tree[T] // pointer to tree in subtree with max upper value
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
func NewTree[T Interval[T]](items ...T) *Tree[T] {
	var t *Tree[T]
	for i := range items {
		t = t.insert(makeNode(items[i]))
	}
	return t
}

// makeNode, create new node with item and random priority and augment it.
func makeNode[T Interval[T]](item T) *Tree[T] {
	n := new(Tree[T])
	n.item = item
	n.prio = rand.Float64()
	n.augment()

	return n
}

// copyNode, make a shallow copy.
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

func (t *Tree[T]) insert(other *Tree[T]) *Tree[T] {
	if t == nil {
		return other
	}

	// other is the new root node?
	if other.prio >= t.prio {
		left, dupe, right := t.split(other.item)
		if dupe != nil {
			// duplicate, drop other
			return t
		}
		other.left, other.right = left, right
		other.augment() // node has changed, augment
		return other
	}

	// immutable insert, copy node
	root := t.copyNode()

	cmp := compare(other.item, root.item)
	switch {
	case cmp < 0: // rec-descent
		root.left = root.left.insert(other)
	case cmp > 0: // rec-descent
		root.right = root.right.insert(other)
	default: // drop duplicate
	}

	root.augment() // node has changed, augment
	return root
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

// find the item in the tree, return the node.
func (t *Tree[T]) find(item T) *Tree[T] {
	for {
		if t == nil {
			return nil
		}
		cmp := compare(t.item, item)
		switch {
		case cmp == 0:
			return t
		case cmp < 0:
			t = t.right
		case cmp > 0:
			t = t.left
		}
	}
}

// split the treap into all nodes that compare less-than, equal
// and greater-than the provided key. The resulting nodes are
// properly formed treaps or nil.
func (t *Tree[T]) split(item T) (left, mid, right *Tree[T]) {
	if t == nil {
		return
	}

	// immutable split, copy node
	root := t.copyNode()

	cmp := compare(root.item, item)
	switch {
	case cmp < 0:
		l, m, r := root.right.split(item)
		root.right = l
		root.augment() // node has changed, augment
		return root, m, r
	case cmp > 0:
		l, m, r := root.left.split(item)
		root.left = r
		root.augment() // node has changed, augment
		return l, m, root
	default:
		l, r := root.left, root.right
		root.left, root.right = nil, nil
		root.augment() // node has changed, augment
		return l, root, r
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

	// shortcut
	if m := t.find(item); m != nil {
		return m.item, true
	}

	l, _, _ := t.split(item)
	return l.shortest(item)
}

func (t *Tree[T]) shortest(item T) (result T, ok bool) {
	if t == nil {
		return
	}

	// nope, subtree has too small max upper interval value
	if item.CompareUpper(t.maxUpper.item) > 0 {
		return
	}

	// reverse-order traversal for shortest
	// try right wing for smallest containing hull
	if t.right != nil && item.CompareUpper(t.right.maxUpper.item) <= 0 {
		if result, ok = t.right.shortest(item); ok {
			return result, ok
		}
	}

	// this item
	if item.CompareUpper(t.item) <= 0 {
		return t.item, true
	}

	// recursive call to left wing
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
	// try left wing for largest containing hull
	if t.left != nil && item.CompareUpper(t.left.maxUpper.item) <= 0 {
		if result, ok = t.left.largest(item); ok {
			return result, ok
		}
	}

	// this item
	if item.CompareUpper(t.item) <= 0 {
		return t.item, true
	}

	// recursive call to right wing
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

	// in-order traversal for supersets, recursive call to left wing
	if t.left != nil && item.CompareUpper(t.left.maxUpper.item) <= 0 {
		result = append(result, t.left.supersets(item)...)
	}

	// this item
	if item.CompareUpper(t.item) <= 0 {
		result = append(result, t.item)
	}

	// recursive call to right wing
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

	// in-order traversal for subsets, recursive call to left wing
	if t.left != nil && item.CompareUpper(t.left.minUpper.item) >= 0 {
		result = append(result, t.left.subsets(item)...)
	}

	// this item
	if item.CompareUpper(t.item) >= 0 {
		result = append(result, t.item)
	}

	// recursive call to right wing
	if t.right != nil && item.CompareUpper(t.right.minUpper.item) >= 0 {
		result = append(result, t.right.subsets(item)...)
	}

	return
}

// Visitor function, returning false stops the iteration.
type Visitor[T Interval[T]] func(t *Tree[T]) bool

// Ascend traverses the tree in ascencding order, calls the visitFn for every subtree until visitFn returns false.
func (t *Tree[T]) Ascend(visitFn Visitor[T]) {
	if t == nil {
		return
	}
	t.traverse(inorder, visitFn)
}

// Descend traverses the tree in descending order, calls the visitFn for every subtree until visitFn returns false.
func (t *Tree[T]) Descend(visitFn Visitor[T]) {
	if t == nil {
		return
	}
	t.traverse(reverse, visitFn)
}

// join combines two disjunct treaps. All nodes in treap a have keys <= that of trep b
// for this algorithm to work correctly.
func join[T Interval[T]](a, b *Tree[T]) *Tree[T] {
	// recursion stop condition
	if a == nil {
		return b
	}
	if b == nil {
		return a
	}

	if a.prio > b.prio {
		// immutable join, copy node
		a = a.copyNode()
		a.right = join(a.right, b)
		a.augment() // node has changed, augment
		return a
	} else {
		// immutable join, copy node
		b = b.copyNode()
		b.left = join(a, b.left)
		b.augment() // node has changed, augment
		return b
	}
}

// augment the treap node after each creation/modification as an interval tree with Min/Max upper value in descendants.
// Only one level deeper must be considered. The treap datastructure is very easy to augment.
func (t *Tree[T]) augment() {
	if t == nil {
		return
	}

	// augment the node for some statistics, not really needed for interval algo
	t.size = 1 + t.left.Size() + t.right.Size()
	t.height = 1 + max(t.left.Height(), t.right.Height())

	// start with upper min/max to self
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
