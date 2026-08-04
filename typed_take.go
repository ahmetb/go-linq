package linq

func (q Query[T]) Take(count int) Query[T] {
	return Query[T]{iterate: func(yield func(T) bool) {
		remaining := count
		if remaining <= 0 {
			return
		}

		q.iterate(func(value T) bool {
			remaining--
			return yield(value) && remaining > 0
		})
	}}
}
