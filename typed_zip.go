package linq

import "iter"

func (q Query[T]) Zip[U, R any](other Query[U], resultSelector func(T, U) R) Query[R] {
	return Query[R]{
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
