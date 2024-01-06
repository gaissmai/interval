# package interval
[![Go Reference](https://pkg.go.dev/badge/github.com/gaissmai/interval.svg)](https://pkg.go.dev/github.com/gaissmai/interval#section-documentation)
![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/gaissmai/interval)
[![CI](https://github.com/gaissmai/interval/actions/workflows/go.yml/badge.svg)](https://github.com/gaissmai/interval/actions/workflows/go.yml)
[![Coverage Status](https://coveralls.io/repos/github/gaissmai/interval/badge.svg)](https://coveralls.io/github/gaissmai/interval)
[![Stand With Ukraine](https://raw.githubusercontent.com/vshymanskyy/StandWithUkraine/main/badges/StandWithUkraine.svg)](https://stand-with-ukraine.pp.ua)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## API CHANGE !!!

The API has changed from v0.9.1 to v0.10.0

## Overview

`package interval` is an immutable datastructure for fast lookups in one dimensional intervals.

The implementation is based on treaps, augmented for intervals. Treaps are randomized self balancing binary search trees.

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
Pavel Mavrin about Algorithms and Datastructures e.g. "[Treaps, implicit keys]"
or follow [some links about treaps] from one of the inventors.

Especially useful is the paper "[Fast Set Operations Using Treaps]" by Guy E. Blelloch and Margaret Reid-Miller.

[IP Range lookups]: https://github.com/gaissmai/iprange
[data structure]: https://ieeexplore.ieee.org/abstract/document/912716
[Treaps, implicit keys]: https://youtu.be/svAHk-FAQgM
[some links about treaps]: http://faculty.washington.edu/aragon/treaps.html
[Fast Set Operations Using Treaps]: https://www.cs.cmu.edu/~scandal/papers/treaps-spaa98.pdf

## Compare function

To apply this library to types of one-dimensional intervals, you must provide a compare function:

```go
  // cmp must return four int values:
  //
  //  ll: left  point interval a compared with left  point interval b (-1, 0, +1)
  //  rr: right point interval a compared with right point interval b (-1, 0, +1)
  //  lr: left  point interval a compared with right point interval b (-1, 0, +1)
  //  rl: right point interval a compared with left  point interval b (-1, 0, +1)
  //
  func[T any] cmp(a, b T) (ll, rr, lr, rl int)
```

## API
```go
  import "github.com/gaissmai/interval"

  type Tree[T any] struct{ ... }
  func NewTree[T any](cmp func(a, b T) (ll, rr, lr, rl int), items ...T) *Tree[T]

  func NewTreeConcurrent[T any](jobs int, cmp func(a, b T) (ll, rr, lr, rl int), items ...T) *Tree[T]

  func (t *Tree[T]) Insert(items ...T)
  func (t *Tree[T]) Delete(item T) bool
  func (t *Tree[T]) Union(other *Tree[T], overwrite bool)

  func (t Tree[T]) InsertImmutable(items ...T) *Tree[T]
  func (t Tree[T]) DeleteImmutable(item T) (*Tree[T], bool)
  func (t Tree[T]) UnionImmutable(other *Tree[T], overwrite bool) *Tree[T]
  func (t Tree[T]) Clone() *Tree[T]

  func (t Tree[T]) Find(item T) (result T, ok bool)
  func (t Tree[T]) CoverLCP(item T) (result T, ok bool)
  func (t Tree[T]) CoverSCP(item T) (result T, ok bool)
  func (t Tree[T]) Intersects(item T) bool

  func (t Tree[T]) Covers(item T) []T
  func (t Tree[T]) Precedes(item T) []T

  func (t Tree[T]) CoveredBy(item T) []T
  func (t Tree[T]) PrecededBy(item T) []T

  func (t Tree[T]) Intersections(item T) []T

  func (t Tree[T]) Visit(start, stop T, visitFn func(item T) bool)
  func (t Tree[T]) Fprint(w io.Writer) error
  func (t Tree[T]) String() string
  func (t Tree[T]) Min() (min T)
  func (t Tree[T]) Max() (max T)
```

## Benchmarks

### Insert

The benchmark for `Insert()` shows the values for inserting an item into trees with increasing size.

The trees are randomly generated, as is the item to be inserted.

```
$ go test -benchmem -bench='Insert'
goos: linux
goarch: amd64
pkg: github.com/gaissmai/interval
cpu: Intel(R) Core(TM) i5-8250U CPU @ 1.60GHz

BenchmarkInsert/Into1-8               10415164           128 ns/op        64 B/op          1 allocs/op
BenchmarkInsert/Into10-8               5160836           245 ns/op        64 B/op          1 allocs/op
BenchmarkInsert/Into100-8              2707705           443 ns/op        64 B/op          1 allocs/op
BenchmarkInsert/Into1_000-8            1694250           723 ns/op        64 B/op          1 allocs/op
BenchmarkInsert/Into10_000-8           1204220           963 ns/op        64 B/op          1 allocs/op
BenchmarkInsert/Into100_000-8           966566          1249 ns/op        64 B/op          1 allocs/op
BenchmarkInsert/Into1_000_000-8         645523          1863 ns/op        64 B/op          1 allocs/op
```

### Delete

```
$ go test -benchmem -bench='Delete'
goos: linux
goarch: amd64
pkg: github.com/gaissmai/interval
cpu: Intel(R) Core(TM) i5-8250U CPU @ 1.60GHz

BenchmarkDelete/DeleteFrom10-8          7682869           156 ns/op        0 B/op          0 allocs/op
BenchmarkDelete/DeleteFrom100-8        13009023            92 ns/op        0 B/op          0 allocs/op
BenchmarkDelete/DeleteFrom1_000-8       1912417           627 ns/op        0 B/op          0 allocs/op
BenchmarkDelete/DeleteFrom10_000-8      1362752           889 ns/op        0 B/op          0 allocs/op
BenchmarkDelete/DeleteFrom100_000-8      893157          1334 ns/op        0 B/op          0 allocs/op
BenchmarkDelete/DeleteFrom1_000_000-8    647199          1828 ns/op        0 B/op          0 allocs/op
```

### Lookups

The benchmark for `CoverLCP()` (also known as longest-prefix-match when the interval is an IP-CIDR prefix) is very
promising, as it is a generalized algorithm that is not specifically optimized only for IP address lookup:

```
$ go test -benchmem -bench='CoverLCP'
goos: linux
goarch: amd64
pkg: github.com/gaissmai/interval
cpu: Intel(R) Core(TM) i5-8250U CPU @ 1.60GHz

BenchmarkCoverLCP/In100-8            4833073       239 ns/op        0 B/op    0 allocs/op
BenchmarkCoverLCP/In1_000-8          3714054       323 ns/op        0 B/op    0 allocs/op
BenchmarkCoverLCP/In10_000-8         1897353       630 ns/op        0 B/op    0 allocs/op
BenchmarkCoverLCP/In100_000-8        1533745       782 ns/op        0 B/op    0 allocs/op
BenchmarkCoverLCP/In1_000_000-8      1732171       693 ns/op        0 B/op    0 allocs/op
```

... and for `Intersects()`:


```
$ go test -benchmem -bench='Intersects'
goos: linux
goarch: amd64
pkg: github.com/gaissmai/interval
cpu: Intel(R) Core(TM) i5-7500T CPU @ 2.70GHz

BenchmarkIntersects/In1-4           58757540        19 ns/op        0 B/op    0 allocs/op
BenchmarkIntersects/In10-4          27267051        40 ns/op        0 B/op    0 allocs/op
BenchmarkIntersects/In100-4         17559418        60 ns/op        0 B/op    0 allocs/op
BenchmarkIntersects/In1_000-4       10471032       120 ns/op        0 B/op    0 allocs/op
BenchmarkIntersects/In10_000-4       6546705       175 ns/op        0 B/op    0 allocs/op
BenchmarkIntersects/In100_000-4      7483621       191 ns/op        0 B/op    0 allocs/op
BenchmarkIntersects/In1_000_000-4    8134428       149 ns/op        0 B/op    0 allocs/op
```

... or Find(), the exact match:

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

With all other methods, the treap is split before the search, which allocates memory
but minimises the search space.
