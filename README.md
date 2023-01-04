# package interval
[![Go Reference](https://pkg.go.dev/badge/github.com/gaissmai/interval.svg)](https://pkg.go.dev/github.com/gaissmai/interval#section-documentation)
![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/gaissmai/interval)
[![CI](https://github.com/gaissmai/interval/actions/workflows/go.yml/badge.svg)](https://github.com/gaissmai/interval/actions/workflows/go.yml)
[![Coverage Status](https://coveralls.io/repos/github/gaissmai/interval/badge.svg)](https://coveralls.io/github/gaissmai/interval)
[![Stand With Ukraine](https://raw.githubusercontent.com/vshymanskyy/StandWithUkraine/main/badges/StandWithUkraine.svg)](https://stand-with-ukraine.pp.ua)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## Overview

`package interval` is an immutable datastructure for fast lookups in one dimensional intervals.

The implementation is based on Treaps, augmented for fast interval lookups.

Immutability is achived because insert/delete will return a new Treap which will share some nodes with the original Treap.
All nodes are read-only after creation, allowing concurrent readers to operate safely with concurrent writers.

The time complexity is **O(log(n))** or **O(k*log(n))** where k is the number of returned items, the space complexity is **O(n)**.

```
Insert()    O(log(n))
Upsert()    O(log(n))
Delete()    O(log(n))
Shortest()  O(log(n))
Largest()   O(log(n))

Subsets()   O(k*log(n))
Supersets() O(k*log(n))
```

The author is propably the first (december 2022) using augmented Treaps
as a very promising [data structure] for the representation of dynamic IP address tables
for arbitrary ranges, that enables most and least specific range matching and even more lookup methods
returning sets of intervals.

The library can be used for all comparable one-dimensional intervals,
but the author of the library uses it mainly for fast [IP range lookups] in access control lists (ACL)
and in his own IP address management (IPAM) and network management software.

The augmented Treap is NOT limited to IP CIDR ranges unlike the prefix trie.
Arbitrary IP ranges and both IP versions can be handled together in this data structure.

Due to the nature of Treaps the lookups and updates can be concurrently decoupled, without delayed rebalancing.

To familiarize yourself with Treaps, see the extraordinarily good lectures from
Pavel Mravin about Algorithms and Datastructures e.g. "[Treaps, implicit keys]"
or follow [some links about Treaps] from one of the inventors.

Especially useful is the paper "[Fast Set Operations Using Treaps]" by Guy E. Blelloch and Margaret Reid-Miller.

[IP Range lookups]: https://github.com/gaissmai/iprange
[data structure]: https://ieeexplore.ieee.org/abstract/document/912716
[Treaps, implicit keys]: https://youtu.be/svAHk-FAQgM
[some links about Treaps]: http://faculty.washington.edu/aragon/treaps.html
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
  func (t *Tree[T]) Upsert(items ...T) *Tree[T]
  func (t *Tree[T]) Delete(item T) (*Tree[T], bool)
  func (t *Tree[T]) Clone() *Tree[T]

  func (t *Tree[T]) Min() (min T)
  func (t *Tree[T]) Max() (max T)
  func (t *Tree[T]) Size() int
  func (t *Tree[T]) Visit(start, stop T, visitFn func(t T) bool)
  func (t *Tree[T]) Fprint(w io.Writer) error

  func (t *Tree[T]) Shortest(item T) (result T, ok bool)
  func (t *Tree[T]) Largest(item T) (result T, ok bool)
  func (t *Tree[T]) Subsets(item T) []T
  func (t *Tree[T]) Supersets(item T) []T
```
