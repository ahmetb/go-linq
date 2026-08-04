package linq

func (q query[T]) Select[R any](selector func(T) R) query[R] {
	return query[R]{
		iterate: func(yield func(R) bool) {
			q.iterate(func(value T) bool {
				return yield(selector(value))
			})
		},
	}
}
