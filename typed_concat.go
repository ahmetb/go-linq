package linq

func (q query[T]) Concat(other query[T]) query[T] {
	return query[T]{iterate: func(yield func(T) bool) {
		continuing := true
		q.iterate(func(value T) bool {
			continuing = yield(value)
			return continuing
		})
		if continuing {
			other.iterate(yield)
		}
	}}
}
