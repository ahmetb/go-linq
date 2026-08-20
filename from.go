package linq

import (
	"context"
	"iter"
	"math"
	"slices"
)

// Query is the type returned from query functions. It represents a lazy,
// strongly-typed sequence of elements of type T. It can be iterated manually
// as shown in the example.
type Query[T any] struct {
	Iterate iter.Seq[T]

	// size hints the exact number of elements the query yields, when that is
	// cheaply known; zero means unknown. Sources with a known length set it,
	// and operators may propagate a size that can be derived exactly from their
	// inputs. Operators whose cardinality is unknown construct a fresh Query
	// without the field, and the hint safely zeroes out.
	//
	// The hint is consumed only as the capacity of preallocated result
	// slices, so it can never change what a query produces. A missing hint
	// forfeits the preallocation; a stale one (e.g., a source map mutated
	// after the query was built) merely mis-sizes it.
	size int
}

func emptyQuery[T any]() Query[T] {
	return Query[T]{Iterate: func(func(T) bool) {}}
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

// KeyValue pairs a key with a value. Map queries and keyed operators yield
// KeyValue elements, which ToMap can collect into a map.
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

// Sequence generates a numeric sequence from start towards endInclusive by
// repeatedly adding step. It panics if any argument is NaN, if step is zero
// while the bounds differ, or if step points away from endInclusive.
//
// endInclusive is yielded only when the sum lands on it exactly, so
// Sequence(0, 5, 2) stops at 4. A step that would overflow T instead of
// reaching the bound also ends the sequence: Sequence(int8(126), 127, 2)
// yields 126 alone.
func Sequence[T Number](start, endInclusive, step T) Query[T] {
	if start != start {
		panic("linq: sequence start must not be NaN")
	}
	if endInclusive != endInclusive {
		panic("linq: sequence end must not be NaN")
	}
	if step != step {
		panic("linq: sequence step must not be NaN")
	}

	if start == endInclusive {
		return Repeat(start, 1)
	}
	if step == 0 {
		panic("linq: sequence step must not be zero unless the bounds are equal")
	}

	increasing := step > 0
	if increasing != (endInclusive > start) {
		panic("linq: sequence step points away from end")
	}

	return Query[T]{
		Iterate: func(yield func(T) bool) {
			for current := start; yield(current); {
				next := current + step
				// Stop once next passes the bound, or once overflow wrapped it
				// back the way it came.
				if increasing {
					if next > endInclusive || next <= current {
						return
					}
				} else if next < endInclusive || next >= current {
					return
				}

				current = next
			}
		},
		size: sequenceSize(start, endInclusive, step),
	}
}

// sequenceSize returns the exact number of elements Sequence yields, or zero
// when that count is not cheaply derivable: a floating point step accumulates
// rounding, and a count wider than an int cannot be a capacity anyway.
//
// The span is measured in uint64 rather than in T so that it survives bounds a
// signed T cannot hold (int8 spans -128..127 as 255), and so that negating the
// most negative step stays correct.
//
// It expects what Sequence has already validated: a nonzero step pointing from
// start towards endInclusive. A zero step would divide by zero below.
func sequenceSize[T Number](start, endInclusive, step T) int {
	// Integer division truncates to zero; only a floating point T keeps a half.
	if half := T(1) / 2; half != 0 {
		return 0
	}

	var span, magnitude uint64
	if step > 0 {
		span, magnitude = uint64(endInclusive)-uint64(start), uint64(step)
	} else {
		span, magnitude = uint64(start)-uint64(endInclusive), -uint64(step)
	}

	steps := span / magnitude
	if steps >= uint64(math.MaxInt) {
		return 0
	}
	return int(steps) + 1
}

// InfiniteSequence generates an unbounded numeric sequence by repeatedly
// adding step to start. Iteration ends only when the consumer stops it.
func InfiniteSequence[T Number](start, step T) Query[T] {
	return Query[T]{
		Iterate: func(yield func(T) bool) {
			for current := start; yield(current); current += step {
			}
		},
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
