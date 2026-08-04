package linq

func (q Query[T]) SelectMany[R any](selector func(T) Query[R]) Query[R] {
	return Query[R]{
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

func (q Query[T]) SelectManyIndexed[R any](selector func(int, T) Query[R]) Query[R] {
	return Query[R]{
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

func (q Query[T]) SelectManyBy[U, R any](
	selector func(T) Query[U],
	resultSelector func(U, T) R,
) Query[R] {
	return Query[R]{
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

func (q Query[T]) SelectManyByIndexed[U, R any](
	selector func(int, T) Query[U],
	resultSelector func(U, T) R,
) Query[R] {
	return Query[R]{
		iterate: func(yield func(R) bool) {
			index := 0
			q.iterate(func(outer T) bool {
				innerQuery := selector(index, outer)
				index++
				keepGoing := true
				innerQuery.iterate(func(inner U) bool {
					keepGoing = yield(resultSelector(inner, outer))
					return keepGoing
				})
				return keepGoing
			})
		},
	}
}
