package linq

// Take lazily yields at most count values from the start of the query.
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
