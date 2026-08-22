package linq

import "math/rand/v2"

// Shuffle returns the elements of a collection in randomized order. It uses a
// non-cryptographically-secure random number generator and reshuffles on every
// iteration.
func (q Query[T]) Shuffle() Query[T] {
	return Query[T]{
		Iterate: func(yield func(T) bool) {
			items := q.collect()
			for i := len(items) - 1; i > 0; i-- {
				j := rand.IntN(i + 1)
				items[i], items[j] = items[j], items[i]
			}

			for _, item := range items {
				if !yield(item) {
					return
				}
			}
		},
		size: q.size,
	}
}
