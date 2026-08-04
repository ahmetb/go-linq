package linq

import "iter"

func (q query[T]) Zip[U, R any](other query[U], resultSelector func(T, U) R) query[R] {
	return query[R]{
		iterate: func(yield func(R) bool) {
			nextLeft, stopLeft := iter.Pull(q.iterate)
			defer stopLeft()
			nextRight, stopRight := iter.Pull(other.iterate)
			defer stopRight()

			for {
				left, leftOK := nextLeft()
				right, rightOK := nextRight()
				if !leftOK || !rightOK || !yield(resultSelector(left, right)) {
					return
				}
			}
		},
	}
}
