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
