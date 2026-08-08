package linq

import (
	"context"
	"iter"
	"slices"
)

// Query is the type returned from query functions. It represents a lazy,
// strongly-typed sequence of elements of type T. It can be iterated manually
// as shown in the example.
type Query[T any] struct {
	Iterate iter.Seq[T]

	// size hints the exact number of elements the query yields, when that is
	// cheaply known; zero means unknown. Sources with a known length set it,
	// and only operators that emit exactly one element per source element may
	// propagate it. Operators that change cardinality need to do nothing:
	// they construct a fresh Query without the field, and the hint safely
	// zeroes out.
	//
	// The hint is consumed only as the capacity of preallocated result
	// slices, so it can never change what a query produces. A missing hint
	// forfeits the preallocation; a stale one (e.g., a source map mutated
	// after the query was built) merely mis-sizes it.
	size int
}

// collect gathers all elements into a slice, preallocating when the query
// carries a size hint. The hinted allocation is deferred until the first
// element arrives so that an empty query collects to nil even when a stale
// hint promises elements, matching slices.Collect on the unhinted path.
func (q Query[T]) collect() []T {
	if size := q.size; size > 0 {
		var out []T
		q.Iterate(func(item T) bool {
			if out == nil {
				out = make([]T, 0, size)
			}
			out = append(out, item)
			return true
		})
		return out
	}
	return slices.Collect(q.Iterate)
}

// KeyValue is a type used to iterate over a map. This type is also used by
// the ToMap function to output the result of a query into a map.
type KeyValue[TKey comparable, TValue any] struct {
	Key   TKey
	Value TValue
}

// FromSlice initializes a linq query with a passed slice.
func FromSlice[S ~[]T, T any](source S) Query[T] {
	return Query[T]{
		Iterate: func(yield func(T) bool) {
			for _, item := range source {
				if !yield(item) {
					return
				}
			}
		},
		size: len(source),
	}
}

// FromMap initializes a linq query with a passed map. Elements are yielded as
// KeyValue pairs, in unspecified order.
func FromMap[M ~map[TKey]TValue, TKey comparable, TValue any](source M) Query[KeyValue[TKey, TValue]] {
	return Query[KeyValue[TKey, TValue]]{
		Iterate: func(yield func(KeyValue[TKey, TValue]) bool) {
			for k, v := range source {
				if !yield(KeyValue[TKey, TValue]{
					Key:   k,
					Value: v,
				}) {
					return
				}
			}
		},
		size: len(source),
	}
}

// FromChannel initializes a linq query with a passed channel, linq iterates over
// the channel until it is closed.
func FromChannel[T any](source <-chan T) Query[T] {
	return Query[T]{
		Iterate: func(yield func(T) bool) {
			for item := range source {
				if !yield(item) {
					return
				}
			}
		},
	}
}

// FromChannelWithContext initializes a linq query with a passed channel
// and stops iterating either when the channel is closed or when the context is canceled.
func FromChannelWithContext[T any](ctx context.Context, source <-chan T) Query[T] {
	return Query[T]{
		Iterate: func(yield func(T) bool) {
			for {
				select {
				case <-ctx.Done():
					// Context canceled or deadline exceeded
					return
				case item, ok := <-source:
					if !ok || !yield(item) {
						// Channel closed or Consumer stopped early
						return
					}
				}
			}
		},
	}
}

// FromString initializes a query from a string, iterating over its runes.
func FromString[S ~string](source S) Query[rune] {
	return Query[rune]{
		Iterate: func(yield func(rune) bool) {
			for _, ch := range string(source) {
				if !yield(ch) {
					return
				}
			}
		},
	}
}

// FromSeq initializes a linq query from an iter.Seq. This allows any
// range-over-func iterator to be used as a query source, including custom
// collections that expose an iterator method.
func FromSeq[T any](source iter.Seq[T]) Query[T] {
	return Query[T]{
		Iterate: source,
	}
}

// Range generates a sequence of integral numbers within a specified range.
func Range(start, count int) Query[int] {
	return Query[int]{
		Iterate: func(yield func(int) bool) {
			end := start + count
			for i := start; i < end; i++ {
				if !yield(i) {
					return
				}
			}
		},
		size: max(count, 0),
	}
}

// Repeat generates a sequence that contains one repeated value.
func Repeat[T any](value T, count int) Query[T] {
	return Query[T]{
		Iterate: func(yield func(T) bool) {
			for i := 0; i < count; i++ {
				if !yield(value) {
					return
				}
			}
		},
		size: max(count, 0),
	}
}
