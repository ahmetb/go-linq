package linq

func (q query[T]) Skip(count int) query[T] {
	return query[T]{iterate: func(yield func(T) bool) {
		remaining := count
		q.iterate(func(value T) bool {
			if remaining > 0 {
				remaining--
				return true
			}
			return yield(value)
		})
	}}
}
