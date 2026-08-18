package linq

import "math"

// concatSize returns the exact size of a query that yields every element of a
// size-a sequence followed by every element of a size-b one. Both arguments are
// sizes, never deltas, so zero means unknown and the result is known only when
// both inputs are and their sum fits in an int.
//
// Only the concatenating operators need this. An operator that shifts a size by
// a constant (Skip, Take) clamps with min or max instead, which already treat
// the zero sentinel correctly; plain addition does not.
func concatSize(a, b int) int {
	if a <= 0 || b <= 0 || a > math.MaxInt-b {
		return 0
	}
	return a + b
}

// Append inserts an item to the end of a collection, so it becomes the last
// item.
func (q Query[T]) Append(item T) Query[T] {
	return Query[T]{
		Iterate: func(yield func(T) bool) {
			stopped := false

			q.Iterate(func(originalItem T) bool {
				if !yield(originalItem) {
					stopped = true
					return false
				}
				return true
			})

			if !stopped {
				yield(item)
			}
		},
		size: concatSize(q.size, 1),
	}
}

// Concat concatenates two collections.
//
// The Concat method differs from the Union method because the Concat method
// returns all the original elements in the input sequences. The Union method
// returns only unique elements.
func (q Query[T]) Concat(q2 Query[T]) Query[T] {
	return Query[T]{
		Iterate: func(yield func(T) bool) {
			stopped := false

			q.Iterate(func(item T) bool {
				if !yield(item) {
					stopped = true
					return false
				}
				return true
			})

			if !stopped {
				q2.Iterate(yield)
			}
		},
		size: concatSize(q.size, q2.size),
	}
}

// Prepend inserts an item to the beginning of a collection, so it becomes the
// first item.
func (q Query[T]) Prepend(item T) Query[T] {
	return Query[T]{
		Iterate: func(yield func(T) bool) {
			if !yield(item) {
				return
			}

			q.Iterate(yield)
		},
		size: concatSize(q.size, 1),
	}
}
