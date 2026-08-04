package linq

func (q query[T]) SelectMany[R any](selector func(T) query[R]) query[R] {
	return query[R]{
		iterate: func(yield func(R) bool) {
			q.iterate(func(outer T) bool {
				keepGoing := true
				selector(outer).iterate(func(inner R) bool {
					keepGoing = yield(inner)
					return keepGoing
				})
				return keepGoing
			})
		},
	}
}

func (q query[T]) SelectManyIndexed[R any](selector func(int, T) query[R]) query[R] {
	return query[R]{
		iterate: func(yield func(R) bool) {
			index := 0
			q.iterate(func(outer T) bool {
				innerQuery := selector(index, outer)
				index++
				keepGoing := true
				innerQuery.iterate(func(inner R) bool {
					keepGoing = yield(inner)
					return keepGoing
				})
				return keepGoing
			})
		},
	}
}

func (q query[T]) SelectManyBy[U, R any](
	selector func(T) query[U],
	resultSelector func(U, T) R,
) query[R] {
	return query[R]{
		iterate: func(yield func(R) bool) {
			q.iterate(func(outer T) bool {
				keepGoing := true
				selector(outer).iterate(func(inner U) bool {
					keepGoing = yield(resultSelector(inner, outer))
					return keepGoing
				})
				return keepGoing
			})
		},
	}
}
