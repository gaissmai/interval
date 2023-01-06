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

/*
type Item = period.Ival

// ccInsert is a fan-out/fan-in version of interval.Insert
func ccInsert(is []Item) *interval.Tree[Item] {
	var treap *interval.Tree[Item]
	var chunk int = len(is) / runtime.NumCPU()

	trees := make(chan *interval.Tree[Item])
	var wg sync.WaitGroup

	// fan out
	for len(is) > 0 {
		if len(is) >= chunk {
			wg.Add(1)
			go func(items []Item, c chan<- *interval.Tree[Item]) {
				defer wg.Done()
				c <- treap.Insert(items...)
			}(is[:chunk], trees)
			is = is[chunk:]
		} else {
			wg.Add(1)
			go func(s []Item, c chan<- *interval.Tree[Item]) {
				defer wg.Done()
				c <- treap.Insert(s...)
			}(is, trees)
			is = nil
		}
	}

	go func() {
		wg.Wait()
		close(trees)
	}()

	// fan in
	var union *interval.Tree[Item]
	for treap := range trees {
		union = union.Union(treap, false)
	}
	return union
}
*/
