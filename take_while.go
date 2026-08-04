package linq

// TakeWhile yields the leading run of values satisfying predicate.
func (q Query[T]) TakeWhile(predicate func(T) bool) Query[T] {
	return Query[T]{iterate: func(yield func(T) bool) {
		q.iterate(func(value T) bool {
			if !predicate(value) {
				return false
			}

			return yield(value)
		})
	}}
}
