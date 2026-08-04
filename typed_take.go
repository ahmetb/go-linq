package linq

func (q query[T]) Take(count int) query[T] {
	return query[T]{iterate: func(yield func(T) bool) {
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
