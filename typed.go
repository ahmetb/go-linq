package linq

import (
	"iter"
	"slices"
)

// Query is a lazily evaluated, statically typed sequence.
type Query[T any] struct {
	iterate iter.Seq[T]
}

func fromSlice[S ~[]T, T any](source S) Query[T] {
	return Query[T]{iterate: slices.Values(source)}
}

func (q Query[T]) toSlice() []T {
	return slices.Collect(q.iterate)
}
