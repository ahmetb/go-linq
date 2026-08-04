package linq

func (q Query[T]) Select[R any](selector func(T) R) Query[R] {
	return Query[R]{
		iterate: func(yield func(R) bool) {
			q.iterate(func(value T) bool {
				return yield(selector(value))
			})
		},
	}
}

func (q Query[T]) SelectIndexed[R any](selector func(int, T) R) Query[R] {
	return Query[R]{
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
