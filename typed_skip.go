package linq

func (q Query[T]) Skip(count int) Query[T] {
	return Query[T]{iterate: func(yield func(T) bool) {
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
