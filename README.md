# package interval
[![Go Reference](https://pkg.go.dev/badge/github.com/gaissmai/interval.svg)](https://pkg.go.dev/github.com/gaissmai/interval#section-documentation)
![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/gaissmai/interval)
[![CI](https://github.com/gaissmai/interval/actions/workflows/go.yml/badge.svg)](https://github.com/gaissmai/interval/actions/workflows/go.yml)
[![Coverage Status](https://coveralls.io/repos/github/gaissmai/interval/badge.svg)](https://coveralls.io/github/gaissmai/interval)
[![Stand With Ukraine](https://raw.githubusercontent.com/vshymanskyy/StandWithUkraine/main/badges/StandWithUkraine.svg)](https://stand-with-ukraine.pp.ua)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## Overview

`package interval` is an immutable datastructure for fast lookups in one dimensional intervals.

The implementation is based on treaps, augmented for intervals. Treaps are randomized self balancing binary search trees.

Immutability is achieved because insert/delete will return a new treap which will share some nodes with the original treap.
All nodes are read-only after creation, allowing concurrent readers to operate safely with concurrent writers.

The time complexity is **O(log(n))** or **O(k*log(n))** where k is the number of returned items, the space complexity is **O(n)**.

```
Insert()    O(log(n))
Delete()    O(log(n))
Shortest()  O(log(n))
Largest()   O(log(n))

Subsets()   O(k*log(n))
Supersets() O(k*log(n))
```

The author is propably the first (december 2022) using augmented treaps
as a very promising [data structure] for the representation of dynamic IP address tables
for arbitrary ranges, that enables most and least specific range matching and even more lookup methods
returning sets of intervals.

The library can be used for all comparable one-dimensional intervals,
but the author of the library uses it mainly for fast [IP range lookups] in access control lists (ACL)
and in his own IP address management (IPAM) and network management software.

The augmented treap is NOT limited to IP CIDR ranges unlike the prefix trie.
Arbitrary IP ranges and both IP versions can be handled together in this data structure.

Due to the nature of treaps the lookups and updates can be concurrently decoupled,
without delayed rebalancing, promising to be a perfect match for a software-router or firewall.

To familiarize yourself with treaps, see the extraordinarily good lectures from
Pavel Mravin about Algorithms and Datastructures e.g. "[Treaps, implicit keys]"
or follow [some links about treaps] from one of the inventors.

Especially useful is the paper "[Fast Set Operations Using Treaps]" by Guy E. Blelloch and Margaret Reid-Miller.

[IP Range lookups]: https://github.com/gaissmai/iprange
[data structure]: https://ieeexplore.ieee.org/abstract/document/912716
[Treaps, implicit keys]: https://youtu.be/svAHk-FAQgM
[some links about treaps]: http://faculty.washington.edu/aragon/treaps.html
[Fast Set Operations Using Treaps]: https://www.cs.cmu.edu/~scandal/papers/treaps-spaa98.pdf

## Interface

To apply this library to types of one-dimensional intervals, they must just implement the following small interface:

```go
// Compare the lower and upper points of two intervals.
type Interface[T any] interface {
	CompareLower(T) int
	CompareUpper(T) int
}
```

## API
```go
import "github.com/gaissmai/interval"

type Tree[T Interface[T]] struct{ ... }

  func (t *Tree[T]) Insert(items ...T) *Tree[T]
  func (t *Tree[T]) Delete(item T) (*Tree[T], bool)

  func (t *Tree[T]) Shortest(item T) (result T, ok bool)
  func (t *Tree[T]) Largest(item T) (result T, ok bool)

  func (t *Tree[T]) Subsets(item T) []T
  func (t *Tree[T]) Supersets(item T) []T

  func (t *Tree[T]) Clone() *Tree[T]
  func (t *Tree[T]) Union(b *Tree[T], overwrite bool) *Tree[T]

  func (t *Tree[T]) Visit(start, stop T, visitFn func(t T) bool)
  func (t *Tree[T]) Fprint(w io.Writer) error
  func (t *Tree[T]) Size() int
  func (t *Tree[T]) Min() (min T)
  func (t *Tree[T]) Max() (max T)

```

## Benchmarks

### Insert

The benchmark for `Insert()` shows the values for inserting an item into trees with increasing size.

The trees are randomly generated, as is the item to be inserted.

The trees are immutable, insertions and deletions generate new nodes on the path. The expected depth
of the trees is **O(log(n))** and the **allocs/op** represent this well.

The data structure is a randomized BST, the expected depth is determined with very
high probability (for large n) but not deterministic.

```
$ go test -benchmem -bench='Insert' -cpu=1
goos: linux
goarch: amd64
pkg: github.com/gaissmai/interval
cpu: Intel(R) Core(TM) i5-8250U CPU @ 1.60GHz
BenchmarkInsertInto1         	 3780463	       295.8 ns/op	     128 B/op	       2 allocs/op
BenchmarkInsertInto10        	 2949376	       477.6 ns/op	     192 B/op	       3 allocs/op
BenchmarkInsertInto100       	 1000000	      1270 ns/op	     512 B/op	       8 allocs/op
BenchmarkInsertInto1_000     	 1335027	       941.0 ns/op	     384 B/op	       6 allocs/op
BenchmarkInsertInto10_000    	 1000000	      2117 ns/op	     768 B/op	      12 allocs/op
BenchmarkInsertInto100_000   	 1000000	      2590 ns/op	     960 B/op	      15 allocs/op
BenchmarkInsertInto1_000_000 	 1000000	      2682 ns/op	    1024 B/op	      16 allocs/op
```

### Delete

The benchmark for `Delete()` shows the same asymptotic behavior:

```
$ go test -benchmem -bench='Delete' -cpu=1
goos: linux
goarch: amd64
pkg: github.com/gaissmai/interval
cpu: Intel(R) Core(TM) i5-8250U CPU @ 1.60GHz
BenchmarkDeleteFrom1         	20716108	       137.2 ns/op	      64 B/op	       1 allocs/op
BenchmarkDeleteFrom10        	 3499713	       338.8 ns/op	     128 B/op	       2 allocs/op
BenchmarkDeleteFrom100       	 1000000	      1207 ns/op	     448 B/op	       7 allocs/op
BenchmarkDeleteFrom1_000     	 1461568	       821.2 ns/op	     320 B/op	       5 allocs/op
BenchmarkDeleteFrom10_000    	 1000000	      1795 ns/op	     704 B/op	      11 allocs/op
BenchmarkDeleteFrom100_000   	 1000000	      2358 ns/op	     896 B/op	      14 allocs/op
BenchmarkDeleteFrom1_000_000 	 1000000	      2561 ns/op	     960 B/op	      15 allocs/op
```

### Lookup

The benchmark for `Shortest()` (a.k.a. longest-prefix-match if the interval is an IP CIDR prefix) is very promising:

```
goos: linux
goarch: amd64
pkg: github.com/gaissmai/interval
cpu: Intel(R) Core(TM) i5-8250U CPU @ 1.60GHz
BenchmarkShortestIn1-8           	64952856	        18.36 ns/op	       0 B/op	       0 allocs/op
BenchmarkShortestIn10-8          	22210141	       122.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkShortestIn100-8         	33358066	       113.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkShortestIn1_000-8       	 7419712	       174.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkShortestIn10_000-8      	 8476441	       335.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkShortestIn100_000-8     	 4940514	       262.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkShortestIn1_000_000-8   	 4152391	       298.0 ns/op	       0 B/op	       0 allocs/op
```
