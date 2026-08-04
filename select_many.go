package linq

// SelectMany projects each value to a query and lazily flattens the results.
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

// SelectManyIndexed projects each index and value to a query and flattens it.
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

// SelectManyBy flattens projected queries and combines inner and outer values.
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

// SelectManyByIndexed is SelectManyBy with each outer source index.
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
