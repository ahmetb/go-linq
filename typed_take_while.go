package linq

func (q query[T]) TakeWhile(predicate func(T) bool) query[T] {
	return query[T]{iterate: func(yield func(T) bool) {
		q.iterate(func(value T) bool {
			if !predicate(value) {
				return false
			}

			return yield(value)
		})
	}}
}
