# package interval
[![Go Reference](https://pkg.go.dev/badge/github.com/gaissmai/interval.svg)](https://pkg.go.dev/github.com/gaissmai/interval#section-documentation)
[![CI](https://github.com/gaissmai/interval/actions/workflows/go.yml/badge.svg)](https://github.com/gaissmai/interval/actions/workflows/go.yml)
[![Coverage Status](https://coveralls.io/repos/github/gaissmai/interval/badge.svg)](https://coveralls.io/github/gaissmai/interval)
[![Go Report Card](https://goreportcard.com/badge/github.com/gaissmai/interval)](https://goreportcard.com/report/github.com/gaissmai/interval)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

`package interval` provides fast lookups and various methods on generic one-dimensional intervals.

The author of the library uses it for IP-Range lookups in Access-Control-Lists (ACL)
and in the authors own IP-Address-Management (IPAM) and network-management software.

But the library is also useful for all one-dimensional arrays, e.g. time intervals.
Thanks to generics this could be abstracted with minimal constraints.

## Interface

To apply this library to types of one-dimensional intervals, they must just implement the following small interface:

```go
type Interface[T any] interface {
	// CompareFirst must compare the first points of two intervals.
	CompareFirst(T) int

	// CompareLast must compare the last points of two intervals.
	CompareLast(T) int
}
```

## API
```golang
func Sort[T Interface[T]](items []T)

func NewTree[T Interface[T]](items []T) *Tree[T]

func (t *Tree[T]) Shortest(item T) (match T, ok bool)
func (t *Tree[T]) Largest(item T) (match T, ok bool)

func (t *Tree[T]) Subsets(item T) []T
func (t *Tree[T]) Supersets(item T) []T

func (t *Tree[T]) Size() int
func (t *Tree[T]) String() string

```
