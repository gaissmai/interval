package interval_test

var (
	ivals = generateIvals(3_000_000)
	probe = generateIvals(1)[0]
)

var (
	tree1         = treap.Insert(ivals[:1]...)
	tree10        = treap.Insert(ivals[:10]...)
	tree100       = treap.Insert(ivals[:100]...)
	tree1_000     = treap.Insert(ivals[:1_000]...)
	tree10_000    = treap.Insert(ivals[:10_000]...)
	tree100_000   = treap.Insert(ivals[:100_000]...)
	tree1_000_000 = treap.Insert(ivals[:1_000_000]...)
)
