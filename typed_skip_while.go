package linq

func (q Query[T]) SkipWhile(predicate func(T) bool) Query[T] {
	return Query[T]{iterate: func(yield func(T) bool) {
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
