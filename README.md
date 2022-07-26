# package interval
[![Go Reference](https://pkg.go.dev/badge/github.com/gaissmai/interval.svg)](https://pkg.go.dev/github.com/gaissmai/interval#section-documentation)
![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/gaissmai/interval)
[![CI](https://github.com/gaissmai/interval/actions/workflows/go.yml/badge.svg)](https://github.com/gaissmai/interval/actions/workflows/go.yml)
[![Coverage Status](https://coveralls.io/repos/github/gaissmai/interval/badge.svg)](https://coveralls.io/github/gaissmai/interval)
[![Stand With Ukraine](https://raw.githubusercontent.com/vshymanskyy/StandWithUkraine/main/badges/StandWithUkraine.svg)](https://stand-with-ukraine.pp.ua)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

`package interval` provides fast lookups and various methods on generic one-dimensional intervals.

The author of the library uses it for IP-Range lookups in Access-Control-Lists (ACL)
and in the authors own IP-Address-Management (IPAM) and network-management software,
see also: https://github.com/gaissmai/iprange

But the library is also useful for all one-dimensional intervals, e.g. time intervals.
Thanks to generics this could be abstracted with minimal constraints.

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

  func NewTree[T Interface[T]](items []T) *Tree[T]

  func (t *Tree[T]) Shortest(item T) (match T, ok bool)
  func (t *Tree[T]) Largest(item T) (match T, ok bool)

  func (t *Tree[T]) Subsets(item T) []T
  func (t *Tree[T]) Supersets(item T) []T

  func (t *Tree[T]) Size() int
  func (t *Tree[T]) String() string

```
