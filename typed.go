package linq

import (
	"iter"
	"slices"
)

// query is the typed query engine under development for v5. It remains
// unexported while operators are migrated alongside the v4 implementation.
type query[T any] struct {
	iterate iter.Seq[T]
}

func fromSlice[S ~[]T, T any](source S) query[T] {
	return query[T]{iterate: slices.Values(source)}
}

func (q query[T]) toSlice() []T {
	return slices.Collect(q.iterate)
}
