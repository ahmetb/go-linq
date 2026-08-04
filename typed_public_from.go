package linq

import (
	"context"
	"iter"
)

// FromSlice creates a typed query over a slice without copying it.
func FromSlice[S ~[]T, T any](source S) Query[T] {
	return fromSlice(source)
}

// FromMap creates a typed query over a map's key-value pairs.
func FromMap[M ~map[K]V, K comparable, V any](source M) Query[KeyValue[K, V]] {
	return fromMap(source)
}

// FromChannel creates a typed query that receives values until the channel closes.
func FromChannel[T any](source <-chan T) Query[T] {
	return fromChannel(source)
}

// FromChannelWithContext creates a typed channel query that also observes ctx.
func FromChannelWithContext[T any](ctx context.Context, source <-chan T) Query[T] {
	return fromChannelWithContext(ctx, source)
}

// FromString creates a rune query over a string or named string type.
func FromString[S ~string](source S) Query[rune] {
	return fromString(source)
}

// FromIterable creates a typed query from a custom iterable collection.
func FromIterable[T any](source Iterable[T]) Query[T] {
	return fromIterable(source)
}

// FromSeq creates a typed query directly from a range-over-function sequence.
func FromSeq[T any](source iter.Seq[T]) Query[T] {
	return Query[T]{iterate: source}
}

// Range generates count consecutive integers starting at start.
func Range(start, count int) Query[int] {
	return integerRange(start, count)
}

// Repeat generates count copies of value.
func Repeat[T any](value T, count int) Query[T] {
	return repeatValue(value, count)
}
