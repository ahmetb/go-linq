package linq

func (q Query[T]) Concat(other Query[T]) Query[T] {
	return Query[T]{iterate: func(yield func(T) bool) {
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
