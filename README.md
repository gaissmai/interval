# package interval
`package interval` provides fast lookups on generic one-dimensional intervals.

## INTERFACE

To apply this library to types of intervals, they must implement the following interface:

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

