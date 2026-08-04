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

func (q query[T]) SelectIndexed[R any](selector func(int, T) R) query[R] {
	return query[R]{
		iterate: func(yield func(R) bool) {
			index := 0
			q.iterate(func(value T) bool {
				selected := selector(index, value)
				index++
				return yield(selected)
			})
		},
	}
}
