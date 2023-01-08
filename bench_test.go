package interval_test

import (
	"github.com/gaissmai/interval"
)

func mkTree[T interval.Interface[T]](ivals []T) (tree *interval.Tree[T]) {
	return tree.Insert(ivals...)
	// return insertCC(runtime.NumCPU(), ivals)
}

/*
func insertCC[T interval.Interface[T]](workers int, is []T) *interval.Tree[T] {
	var wg sync.WaitGroup
	var chunk int = len(is) / workers
	var treap *interval.Tree[T]
	var cTree chan *interval.Tree[T] = make(chan *interval.Tree[T])

	// fan out
	for len(is) > 0 {
		if len(is) >= chunk {
			wg.Add(1)
			go func(s []T, c chan<- *interval.Tree[T]) {
				defer wg.Done()
				c <- treap.Insert(s...)
			}(is[:chunk], cTree)

			is = is[chunk:]
		} else {
			wg.Add(1)
			go func(s []T, c chan<- *interval.Tree[T]) {
				defer wg.Done()
				var treap *interval.Tree[T]
				c <- treap.Insert(s...)
			}(is, cTree)
			is = nil
		}
	}

	go func() {
		wg.Wait()
		close(cTree)
	}()

	// fan in
	var union *interval.Tree[T]
	for treap := range cTree {
		union = union.Union(treap, false, false)
	}
	return union
}
*/
