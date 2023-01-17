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

The author is propably the first (in december 2022) using augmented treaps
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
type Interface[T any] interface {
	// Compare the left (l) and right (r) points of two intervals and returns four integers.
	Compare(T) (ll, rr, lr, rl int)
}
```

## API
```go
  import "github.com/gaissmai/interval"

  type Tree[T Interface[T]] struct{ ... }

  func NewTree[T Interface[T]](items ...T) Tree[T]

  func (t Tree[T]) Insert(items ...T) Tree[T]
  func (t Tree[T]) Delete(item T) (Tree[T], bool)

  func (t *Tree[T]) InsertMutable(items ...T)
  func (t *Tree[T]) DeleteMutable(item T) bool

  func (t Tree[T]) Find(item T) (result T, ok bool)
  func (t Tree[T]) Shortest(item T) (result T, ok bool)
  func (t Tree[T]) Largest(item T) (result T, ok bool)

  func (t Tree[T]) Subsets(item T) []T
  func (t Tree[T]) Supersets(item T) []T

  func (t Tree[T]) Clone() Tree[T]
  func (t Tree[T]) Union(other Tree[T], overwrite bool, immutable bool) Tree[T]

  func (t Tree[T]) Visit(start, stop T, visitFn func(item T) bool)
  func (t Tree[T]) Fprint(w io.Writer) error
  func (t Tree[T]) String() string
  func (t Tree[T]) Size() int
  func (t Tree[T]) Min() (min T)
  func (t Tree[T]) Max() (max T)
```

## Benchmarks

### Insert

The benchmark for `Insert()` shows the values for inserting an item into trees with increasing size.

The trees are randomly generated, as is the item to be inserted.

The trees are immutable, insertions and deletions generate new nodes on the path. The expected depth
of the trees is **O(log(n))** and the **allocs/op** represent this well.

The data structure is a randomized BST, the expected depth is determined with very
high probability (for large n) but not deterministic.

If the original tree is allowed to mutate during insert and delete because the old state is no longer needed,
then the values are correspondingly better.

```
$ go test -benchmem -bench='Insert'
goos: linux
goarch: amd64
pkg: github.com/gaissmai/interval
cpu: Intel(R) Core(TM) i5-8250U CPU @ 1.60GHz
BenchmarkInsert/Into1-8                       7488871           189 ns/op       128 B/op          2 allocs/op
BenchmarkInsert/Into10-8                      3276116           398 ns/op       256 B/op          4 allocs/op
BenchmarkInsert/Into100-8                     1496136           951 ns/op       512 B/op          8 allocs/op
BenchmarkInsert/Into1_000-8                    533299          1959 ns/op       960 B/op         15 allocs/op
BenchmarkInsert/Into10_000-8                   727627          1398 ns/op       832 B/op         13 allocs/op
BenchmarkInsert/Into100_000-8                  342596          3339 ns/op      1600 B/op         25 allocs/op
BenchmarkInsert/Into1_000_000-8                396890          3681 ns/op      1728 B/op         27 allocs/op

BenchmarkInsertMutable/Into1-8               10415164           128 ns/op        64 B/op          1 allocs/op
BenchmarkInsertMutable/Into10-8               5160836           245 ns/op        64 B/op          1 allocs/op
BenchmarkInsertMutable/Into100-8              2707705           443 ns/op        64 B/op          1 allocs/op
BenchmarkInsertMutable/Into1_000-8            1694250           723 ns/op        64 B/op          1 allocs/op
BenchmarkInsertMutable/Into10_000-8           1204220           963 ns/op        64 B/op          1 allocs/op
BenchmarkInsertMutable/Into100_000-8           966566          1249 ns/op        64 B/op          1 allocs/op
BenchmarkInsertMutable/Into1_000_000-8         645523          1863 ns/op        64 B/op          1 allocs/op
```

### Delete

The benchmark for `Delete()` shows the same asymptotic behavior:

```
$ go test -benchmem -bench='Delete'
goos: linux
goarch: amd64
pkg: github.com/gaissmai/interval
cpu: Intel(R) Core(TM) i5-8250U CPU @ 1.60GHz
BenchmarkDelete/DeleteFrom10-8                 8018518           150 ns/op      128 B/op          2 allocs/op
BenchmarkDelete/DeleteFrom100-8                2349710           708 ns/op      448 B/op          7 allocs/op
BenchmarkDelete/DeleteFrom1_000-8               616173          2051 ns/op     1536 B/op         24 allocs/op
BenchmarkDelete/DeleteFrom10_000-8              446180          2362 ns/op     1856 B/op         29 allocs/op
BenchmarkDelete/DeleteFrom100_000-8             272798          4224 ns/op     2816 B/op         44 allocs/op
BenchmarkDelete/DeleteFrom1_000_000-8           231808          5897 ns/op     3520 B/op         55 allocs/op

BenchmarkDeleteMutable/DeleteFrom10-8          7682869           156 ns/op        0 B/op          0 allocs/op
BenchmarkDeleteMutable/DeleteFrom100-8        13009023            92 ns/op        0 B/op          0 allocs/op
BenchmarkDeleteMutable/DeleteFrom1_000-8       1912417           627 ns/op        0 B/op          0 allocs/op
BenchmarkDeleteMutable/DeleteFrom10_000-8      1362752           889 ns/op        0 B/op          0 allocs/op
BenchmarkDeleteMutable/DeleteFrom100_000-8      893157          1334 ns/op        0 B/op          0 allocs/op
BenchmarkDeleteMutable/DeleteFrom1_000_000-8    647199          1828 ns/op        0 B/op          0 allocs/op
```

### Lookups

The benchmark for `Shortest()` (a.k.a. longest-prefix-match if the interval is an IP CIDR prefix) is very promising:

```
$ go test -benchmem -bench='Shortest'
goos: linux
goarch: amd64
pkg: github.com/gaissmai/interval
cpu: Intel(R) Core(TM) i5-8250U CPU @ 1.60GHz
BenchmarkShortest/In100-8            6201722       183 ns/op        0 B/op    0 allocs/op
BenchmarkShortest/In1_000-8          5257042       228 ns/op        0 B/op    0 allocs/op
BenchmarkShortest/In10_000-8         2372185       503 ns/op        0 B/op    0 allocs/op
BenchmarkShortest/In100_000-8        1509024       788 ns/op        0 B/op    0 allocs/op
BenchmarkShortest/In1_000_000-8      1531923       773 ns/op        0 B/op    0 allocs/op
```
The benchmark `Find()` for the exact match:

```
$ go test -benchmem -bench='Find'
goos: linux
goarch: amd64
pkg: github.com/gaissmai/interval
cpu: Intel(R) Core(TM) i5-8250U CPU @ 1.60GHz
BenchmarkFind/In100-8               17299449        63 ns/op      0 B/op    0 allocs/op
BenchmarkFind/In1_000-8             17327350        69 ns/op      0 B/op    0 allocs/op
BenchmarkFind/In10_000-8            12858908        90 ns/op      0 B/op    0 allocs/op
BenchmarkFind/In100_000-8            4696676       256 ns/op      0 B/op    0 allocs/op
BenchmarkFind/In1_000_000-8          7131028       163 ns/op      0 B/op    0 allocs/op
```

The lookup `Largest()` needs a split() of the treap, which means memory allocations,
the same is true for `Supersets()` and `Subsets()`:

```
$ go test -benchmem -bench='Largest'
goos: linux
goarch: amd64
pkg: github.com/gaissmai/interval
cpu: Intel(R) Core(TM) i5-8250U CPU @ 1.60GHz
BenchmarkLargest/In100-8             1184683         991 ns/op      768 B/op      12 allocs/op
BenchmarkLargest/In1_000-8           1235576        1032 ns/op      704 B/op      11 allocs/op
BenchmarkLargest/In10_000-8           574849        1965 ns/op     1216 B/op      19 allocs/op
BenchmarkLargest/In100_000-8          396973        3079 ns/op     1728 B/op      27 allocs/op
BenchmarkLargest/In1_000_000-8        540117        2488 ns/op     1344 B/op      21 allocs/op
```
