package linq

func (q query[T]) SkipWhile(predicate func(T) bool) query[T] {
	return query[T]{iterate: func(yield func(T) bool) {
		skipping := true
		q.iterate(func(value T) bool {
			if skipping && predicate(value) {
				return true
			}

			skipping = false
			return yield(value)
		})
	}}
}
